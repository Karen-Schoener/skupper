#!/bin/bash
set -e # Exit immediately if a command fails

# --- Configuration ---
NETWORK_NAME="$1"

# --- Script Logic ---

if [ -z "$NETWORK_NAME" ]; then
    echo "ERROR: Please provide a Docker network name as the first argument."
    echo "Usage: $0 <network_name>"
    echo "Examples: "
    echo "     # Check Cluster 1 network (should show k3d-cluster1-server-0) "
    echo "     ./list-connected-containers.sh cluster1-net "
    echo "     # Check Cluster 2 network (should show k3d-cluster2-server-0 and potentially the local-registry) "
    echo "     ./list-connected-containers.sh cluster2-netUsage: $0 <network_name>"
    exit 1
fi

echo "========================================================"
echo "Containers Attached to Network: $NETWORK_NAME"
echo "========================================================"

# Check if the network exists before proceeding
if ! docker network inspect "$NETWORK_NAME" > /dev/null 2>&1; then
    echo "ERROR: Docker network '$NETWORK_NAME' does not exist."
    exit 1
fi

# Use docker network inspect with jq to extract the Name and IPv4Address
# We use a combined filter to select only the container objects and then map them
# into a clean array structure for jq to process and output.

# We check for an empty object `{}` as this indicates no containers are attached.
CONTAINER_DATA=$(docker network inspect "$NETWORK_NAME" --format '{{json .Containers}}')

if [ "$CONTAINER_DATA" == "{}" ]; then
    echo "No containers found attached to '$NETWORK_NAME'."
else
    # Process the JSON data to list names and IPs
    echo "$CONTAINER_DATA" | jq -r '
        .[] | {
            Name: .Name, 
            IP: .IPv4Address
        } | "\(.Name) (\(.IP))"
    ' | sort
fi

echo "========================================================"
