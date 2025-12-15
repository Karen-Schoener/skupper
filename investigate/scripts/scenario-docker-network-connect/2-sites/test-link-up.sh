#!/bin/bash
set -euo pipefail # set -e: Exit on error; set -u: Exit on unset variables; set -o pipefail: Fail if any command in a pipeline fails

# --- Configuration (Defaults) ---
N_ITERATIONS=15      # Number of checks to perform for restoration
SLEEP_TIME=4         # Seconds to sleep between checks (Total time: 60s)
LINK_STATUS_FILE="/tmp/skupper_link_status.txt"

# --- Function to display usage ---
usage() {
    echo "Usage: $0 <Source_Cluster_A_Num> <Target_Cluster_B_Num> [Link_Name]"
    echo ""
    echo "Restores link connectivity by reconnecting the Docker network bridge."
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

# Define the status string to look for
if [ -n "$LINK_NAME" ]; then
    STATUS_STRING="${LINK_NAME} is connected"
    TARGET_CHECK="specific link: ${LINK_NAME}"
else
    STATUS_STRING="is connected"
    TARGET_CHECK="at least one link"
fi


# --- Script Start ---
echo "========================================================"
echo "PHASE 2: RESTORING LINK (Reconnecting ${HUB_NODE} to ${LINK_NETWORK})"
echo "========================================================"

# Check for dedication
if [[ -z "$KUBECONFIG" ]] || [[ "$KUBECONFIG" != *k8s-c${SOURCE_CLUSTER_NUM}.yaml* ]]; then
    echo "ERROR: KUBECONFIG is not correctly dedicated for Cluster ${SOURCE_CLUSTER_NUM}."
    echo "Please run 'dedicate-c${SOURCE_CLUSTER_NUM}' before executing this script."
    exit 1
fi

echo "[1/2] RECONNECTING network bridge: ${HUB_NODE} to ${LINK_NETWORK}..."
docker network connect "${LINK_NETWORK}" "${HUB_NODE}"
echo "Docker network reconnected."
echo "--------------------------------------------------------"

echo "[2/2] MONITORING: Waiting for ${TARGET_CHECK} to come back UP ($N_ITERATIONS checks, Total max time: $((N_ITERATIONS * SLEEP_TIME))s)..."
sleep 5 # Initial delay for Skupper to detect the connection

RESTORED="false"

# Use a fixed loop for monitoring restoration
for i in $(seq 1 $N_ITERATIONS); do
    skupper link status > "$LINK_STATUS_FILE" 2>&1

    # Accurate Check: Look for the specific status string
    if grep -q "${STATUS_STRING}" "$LINK_STATUS_FILE"; then
        echo "   [$(date +%T)] Status: UP! Link restored successfully on check $i."
        RESTORED="true"
        break
    else
        echo "   [$(date +%T)] Status: DOWN. Check $i/$N_ITERATIONS. Waiting ${SLEEP_TIME}s..."
        sleep $SLEEP_TIME
    fi
done

echo "--------------------------------------------------------"

if [ "$RESTORED" == "true" ]; then
    echo "TEST COMPLETE: Link successfully restored."
else
    echo "TEST FAILURE: Link did not restore after $N_ITERATIONS checks."
    cat "$LINK_STATUS_FILE"
fi

echo "========================================================"
