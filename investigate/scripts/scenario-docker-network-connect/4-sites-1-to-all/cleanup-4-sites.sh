#!/bin/bash
set -e # Exit immediately if a command fails

NUM_SITES=4

echo "========================================================"
echo "CLEANUP: Removing ${NUM_SITES} k3d clusters and Docker networks."
echo "========================================================"

# --- CLUSTER CLEANUP ---
echo "[1/2] Cleaning up k3d clusters (1 to ${NUM_SITES})..."
for i in $(seq 1 $NUM_SITES); do
    CLUSTER_NAME="cluster$i"
    echo "    -> Deleting cluster: ${CLUSTER_NAME}"
    # Suppress errors if cluster does not exist, and delete
    k3d cluster delete "$CLUSTER_NAME" 2>/dev/null || true
done

# --- NETWORK CLEANUP ---
echo "[2/2] Cleaning up old Docker networks (1 to ${NUM_SITES})..."
for i in $(seq 1 $NUM_SITES); do
    NETWORK_NAME="cluster$i-net"
    echo "    -> Removing network: ${NETWORK_NAME}"
    # Suppress errors if network does not exist, and delete
    docker network rm "$NETWORK_NAME" 2>/dev/null || true
done

echo "========================================================"
echo "CLEANUP COMPLETE."
echo "========================================================"
