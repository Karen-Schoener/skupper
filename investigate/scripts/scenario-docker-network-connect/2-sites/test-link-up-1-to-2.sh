#!/bin/bash
set -e # Exit immediately if a command fails

# --- Configuration ---
HUB_NODE="k3d-cluster1-server-0"
LINK_NETWORK="cluster2-net"
LINK_STATUS_FILE="/tmp/skupper_link_status.txt"
N_ITERATIONS=15      # Number of checks to perform for restoration
SLEEP_TIME=4         # Seconds to sleep between checks (Total time: 60s)

echo "========================================================"
echo "PHASE 2: RESTORING LINK (Reconnect)"
echo "========================================================"

# Check for dedication
if [ -z "$KUBECONFIG" ]; then
    echo "ERROR: KUBECONFIG variable is not set. Please run 'dedicate-c1' first."
    exit 1
fi

echo "[1/2] RECONNECTING network bridge between $HUB_NODE and $LINK_NETWORK..."
docker network connect $LINK_NETWORK $HUB_NODE
echo "Docker network reconnected."
echo "--------------------------------------------------------"

echo "[2/2] MONITORING: Waiting for Skupper link to come back UP ($N_ITERATIONS checks, Total max time: $((N_ITERATIONS * SLEEP_TIME))s)..."
sleep 5 # Initial delay for Skupper to detect the connection

RESTORED="false"

# Use a fixed loop for monitoring restoration
for i in $(seq 1 $N_ITERATIONS); do
    skupper link status > "$LINK_STATUS_FILE" 2>&1

    # Accurate Check: Look for the specific "is connected" string
    if grep -q "is connected" "$LINK_STATUS_FILE"; then
        echo "   [$(date +%T)] Status: UP (is connected)! Link restored successfully on check $i."
        RESTORED="true"
        break
    else
        echo "   [$(date +%T)] Status: DOWN (not connected). Check $i/$N_ITERATIONS. Waiting ${SLEEP_TIME}s..."
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
