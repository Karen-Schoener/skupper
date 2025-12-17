#!/bin/bash
# setup-cluster-N.sh <Cluster_Num>
set -e # Exit immediately if a command fails

# --- INITIAL ENVIRONMENT CHECK (KUBECONFIG) ---
if [ -n "$KUBECONFIG" ]; then
    echo "🛑 ERROR: KUBECONFIG Environment Variable is Set."
    echo "Please run 'unset KUBECONFIG' before running any setup script."
    exit 1
fi
# --- End KUBECONFIG Check ---

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <Cluster_Num>"
    exit 1
fi

CLUSTER_NUM="$1"
CLUSTER_NAME="k3d-cluster${CLUSTER_NUM}"
SITE_NAME="site-${CLUSTER_NUM}"
LB_HOST="${CLUSTER_NAME}-serverlb"
LB_PORT=$((9080 + CLUSTER_NUM))
NETWORK_NAME="cluster${CLUSTER_NUM}-net"

echo "========================================================"
echo "SETTING UP CLUSTER ${CLUSTER_NUM} (${SITE_NAME})"
echo "========================================================"

echo "[1/3] Creating Docker Network: ${NETWORK_NAME}..."
docker network create "$NETWORK_NAME"

echo "[2/3] Creating k3d Cluster: ${CLUSTER_NAME} (LB Port: ${LB_PORT})..."
k3d cluster create cluster${CLUSTER_NUM} \
  --network "$NETWORK_NAME" \
  -p "$LB_PORT:80@loadbalancer" \
  --wait # wait for the cluster to be ready

# Merge the new context into the system default (~/.kube/config)
k3d kubeconfig merge cluster${CLUSTER_NUM} --kubeconfig-merge-default

echo "[3/3] Initializing Skupper on ${CLUSTER_NAME}..."
# Switch context and initialize Skupper
kubectl config use-context "$CLUSTER_NAME"

skupper init --site-name "$SITE_NAME" \
    --enable-console \
    --enable-flow-collector \
    --ingress loadbalancer \
    --router-ingress-host "$LB_HOST"

echo "    -> Checking Skupper status..."
skupper status
kubectl get pods

echo "SETUP FOR CLUSTER ${CLUSTER_NUM} COMPLETE."
echo "========================================================"
