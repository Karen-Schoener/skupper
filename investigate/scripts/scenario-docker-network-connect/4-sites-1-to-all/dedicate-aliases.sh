#!/bin/bash
# connect-hub.sh
set -e # Exit immediately if a command fails

NUM_SPOKES=3

echo "========================================================"
echo "CONNECTING HUB (CLUSTER 1) TO SPOKE NETWORKS (2, 3, 4)"
echo "========================================================"

HUB_NODE="k3d-cluster1-server-0"

# Connect Cluster 1 node to networks of Clusters 2, 3, and 4
for i in $(seq 2 $((NUM_SPOKES + 1))); do
    TARGET_NETWORK="cluster${i}-net"
    echo "    -> Connecting ${HUB_NODE} to ${TARGET_NETWORK}..."
    docker network connect "$TARGET_NETWORK" "$HUB_NODE"
done

echo "HUB CONNECTION COMPLETE. Cluster 1 is now bridged."
echo "========================================================"
