#!/bin/bash
set -e # Exit immediately if a command fails

echo "[1/4] Cleaning up old clusters..."
for i in {1..2}; do
    # Suppress errors if cluster does not exist, and delete
    k3d cluster delete cluster$i 2>/dev/null || true
done

echo "[2/4] Cleaning up old Docker networks..."
for i in {1..2}; do
    # Suppress errors if network does not exist, and delete
    docker network rm cluster$i-net 2>/dev/null || true
done
