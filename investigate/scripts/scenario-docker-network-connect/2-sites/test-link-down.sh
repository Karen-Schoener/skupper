#!/bin/bash
set -euo pipefail # set -e: Exit on error; set -u: Exit on unset variables; set -o pipefail: Fail if any command in a pipeline fails

# --- Configuration (Defaults) ---
N_ITERATIONS=10      # Number of checks to perform
SLEEP_TIME=3         # Seconds to sleep between checks (Total time: 30s)
LINK_STATUS_FILE="/tmp/skupper_link_status.txt"

# --- Function to display usage ---
usage() {
    echo "Usage: $0 <Source_Cluster_A_Num> <Target_Cluster_B_Num> [Link_Name]"
    echo ""
    echo "Tests link failure by disconnecting the Docker network bridge used for the link."
    echo ""
    echo "Arguments:"
    echo "  <Source_Cluster_A_Num> (e.g., 1): The cluster running the link command (k3d-clusterA-server-0)."
    echo "  <Target_Cluster_B_Num> (e.g., 2): The cluster whose network is being linked to (clusterB-net)."
    echo "  [Link_Name] (e.g., link1): The specific Skupper link name to monitor (OPTIONAL, but recommended)."
    echo ""
    echo "Pre-requisite: Must be run in the shell dedicated to Cluster A (i.e., after running 'dedicate-cA')."
    exit 1
}

# --- Parameter Parsing ---

# Check for required arguments
if [ "$#" -lt 2 ]; then
    echo "ERROR: At least Source and Target cluster numbers are required."
    usage
fi

SOURCE_CLUSTER_NUM="$1"
TARGET_CLUSTER_NUM="$2"
LINK_NAME="${3:-}" # If $3 is empty, set LINK_NAME to empty string

# Derive K3d and Docker names from arguments
HUB_NODE="k3d-cluster${SOURCE_CLUSTER_NUM}-server-0"
LINK_NETWORK="cluster${TARGET_CLUSTER_NUM}-net"

# Define the status string to look for (for specific link or generic)
if [ -n "$LINK_NAME" ]; then
    STATUS_STRING="${LINK_NAME} is connected"
    TARGET_CHECK="specific link: ${LINK_NAME}"
else
    STATUS_STRING="is connected"
    TARGET_CHECK="at least one link"
fi

# --- Script Start ---
echo "========================================================"
echo "PHASE 1: LINK DOWN TEST (Disconnecting ${HUB_NODE} from ${LINK_NETWORK})"
echo "========================================================"

# Check for dedication
if [[ -z "$KUBECONFIG" ]] || [[ "$KUBECONFIG" != *k8s-c${SOURCE_CLUSTER_NUM}.yaml* ]]; then
    echo "ERROR: KUBECONFIG is not correctly dedicated for Cluster ${SOURCE_CLUSTER_NUM}."
    echo "Please run 'dedicate-c${SOURCE_CLUSTER_NUM}' before executing this script."
    exit 1
fi

echo "[1/3] Current Skupper link status on Site ${SOURCE_CLUSTER_NUM} (Checking for: \"${STATUS_STRING}\")."
skupper link status

echo "--------------------------------------------------------"
echo "[2/3] DISCONNECTING network bridge: ${HUB_NODE} from ${LINK_NETWORK}..."
docker network disconnect "${LINK_NETWORK}" "${HUB_NODE}" || {
    echo "WARNING: Disconnect command failed, proceeding with monitoring (already disconnected?)."
}
echo "Docker network disconnected."
echo "--------------------------------------------------------"

echo "[3/3] MONITORING: Waiting for ${TARGET_CHECK} to go DOWN ($N_ITERATIONS checks, Total max time: $((N_ITERATIONS * SLEEP_TIME))s)..."

LINK_IS_UP=true # Assume UP until proven otherwise

# Use a fixed loop for monitoring
for i in $(seq 1 $N_ITERATIONS); do
    skupper link status > "$LINK_STATUS_FILE" 2>&1
    
    # Accurate Check: Look for the specific status string
    if grep -q "${STATUS_STRING}" "$LINK_STATUS_FILE"; then
        echo "   [$(date +%T)] Status: UP. Check $i/$N_ITERATIONS. Waiting ${SLEEP_TIME}s..."
        sleep $SLEEP_TIME
    else
        echo "   [$(date +%T)] Status: DOWN. Link failure detected on check $i."
        LINK_IS_UP=false
        break
    fi
done

echo "--------------------------------------------------------"

if [ "$LINK_IS_UP" == false ]; then
    echo "TEST STEP COMPLETE: Skupper link is confirmed DOWN."
    echo "Current status of the link on Site ${SOURCE_CLUSTER_NUM}:"
    cat "$LINK_STATUS_FILE"
    echo ""
    echo "ACTION: Run test-link-up.sh ${SOURCE_CLUSTER_NUM} ${TARGET_CLUSTER_NUM} ${LINK_NAME} to restore."
else
    echo "TEST STEP WARNING: Link did not go down after $N_ITERATIONS checks. Manual check required."
    cat "$LINK_STATUS_FILE"
fi

echo "========================================================"
