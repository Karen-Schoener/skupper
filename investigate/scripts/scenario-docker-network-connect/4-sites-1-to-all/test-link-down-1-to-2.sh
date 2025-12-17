#!/bin/bash
set -e # Exit immediately if a command fails

# --- Configuration ---
HUB_NODE="k3d-cluster1-server-0"
LINK_NETWORK="cluster2-net"
LINK_STATUS_FILE="/tmp/skupper_link_status.txt"
N_ITERATIONS=10      # Number of checks to perform
SLEEP_TIME=3         # Seconds to sleep between checks (Total time: 30s)

echo "========================================================"
echo "PHASE 1: INITIATING LINK DOWN TEST (Disconnect)"
echo "========================================================"

# Check for dedication
if [ -z "$KUBECONFIG" ]; then
    echo "ERROR: KUBECONFIG variable is not set. Please run 'dedicate-c1' first."
    exit 1
fi

echo "[1/3] Current Skupper link status on Site 1 (k3d-cluster1):"
skupper link status

echo "--------------------------------------------------------"
echo "[2/3] DISCONNECTING network bridge between $HUB_NODE and $LINK_NETWORK..."
echo "Simulating link failure..."
docker network disconnect $LINK_NETWORK $HUB_NODE || { 
    echo "WARNING: Docker network was possibly already disconnected, proceeding with monitoring."
}
echo "Docker network disconnected."
echo "--------------------------------------------------------"

echo "[3/3] MONITORING: Checking status $N_ITERATIONS times (Total max time: $((N_ITERATIONS * SLEEP_TIME))s)..."

LINK_IS_UP=true # Assume UP until proven otherwise

# Use a fixed loop for monitoring
for i in $(seq 1 $N_ITERATIONS); do
    skupper link status > "$LINK_STATUS_FILE" 2>&1
    
    # Accurate Check: Look for the specific "is connected" string
    if grep -q "is connected" "$LINK_STATUS_FILE"; then
        echo "   [$(date +%T)] Status: UP (is connected). Check $i/$N_ITERATIONS. Waiting ${SLEEP_TIME}s..."
        sleep $SLEEP_TIME
    else
        echo "   [$(date +%T)] Status: DOWN (not connected). Link failure detected on check $i."
        LINK_IS_UP=false
        break
    fi
done

echo "--------------------------------------------------------"

if [ "$LINK_IS_UP" == false ]; then
    echo "TEST STEP COMPLETE: Skupper link is confirmed DOWN."
    echo "Current status of the link on Site 1:"
    cat "$LINK_STATUS_FILE"
    echo ""
    echo "Run test-link-up-1-to-2.sh to restore."
else
    echo "TEST STEP WARNING: Link did not go down after $N_ITERATIONS checks. Manual check required."
    cat "$LINK_STATUS_FILE"
fi

echo "========================================================"
