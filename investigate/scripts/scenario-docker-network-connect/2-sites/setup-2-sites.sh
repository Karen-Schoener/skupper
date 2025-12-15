#!/bin/bash
set -e # Exit immediately if a command fails

echo "========================================================"
echo "PHASE 1: INFRASTRUCTURE SETUP (2 SITES)"
echo "========================================================"

# --- ROBUST TWO-PASS CLEANUP ---
echo "[1/4] Cleaning up old clusters..."
for i in {1..2}; do
    k3d cluster delete cluster$i 2>/dev/null || true
done

echo "[2/4] Cleaning up old Docker networks..."
for i in {1..2}; do
    docker network rm cluster$i-net 2>/dev/null || true
done
# -------------------------------

echo "[3/4] Creating Clusters and Networks..."
sleep 2

for i in {1..2}; do
    echo "   -> Processing Cluster $i..."
    docker network create cluster$i-net
    
    LB_PORT=$((9080 + i))
    
    k3d cluster create cluster$i \
      --network "cluster$i-net" \
      -p "$LB_PORT:80@loadbalancer" \
      --wait > /dev/null
      
    k3d kubeconfig merge cluster$i --kubeconfig-merge-default
done

# 4. NETWORKING (Connect Cluster 1 to Cluster 2's network for Skupper inter-router linking)
echo "[4/4] Bridging Cluster 1 to cluster2-net for direct communication..."
sleep 5
HUB_NODE="k3d-cluster1-server-0"
docker network connect cluster2-net $HUB_NODE

echo "========================================================"
echo "PHASE 2: SKUPPER INSTALLATION (ALL SITES)"
echo "========================================================"

for i in {1..2}; do
    CLUSTER_NAME="k3d-cluster$i"
    SITE_NAME="site-$i"
    LB_HOST="k3d-cluster$i-serverlb" 
    
    echo "[Site $i] Initializing Skupper on $CLUSTER_NAME..."
    
    # Switch context for the kubectl/skupper commands
    kubectl config use-context $CLUSTER_NAME
    
    skupper init --site-name "$SITE_NAME" \
        --enable-console \
        --enable-flow-collector \
        --ingress loadbalancer \
        --router-ingress-host "$LB_HOST"
    
    echo "   -> Checking Skupper site-$i status..."
    skupper status
done

echo "========================================================"
echo "SETUP COMPLETE"
echo "========================================================"

# --- DEDICATED SHELL INSTRUCTIONS (Using $HOME for robustness) ---
echo "--- POST-SETUP: DEDICATED SHELL WORKFLOW ---"
echo ""
echo "1. ACTIVATE DEDICATED SHELL FUNCTIONS (Run these aliases):"
echo "   # Copy/paste ALL of these lines into your CURRENT shell."
echo "   # They define two new, robust aliases: 'dedicate-c1' and 'dedicate-c2'."
echo "   # They use --flatten and $HOME to prevent certificate/permission errors."
echo ""
echo "    unset KUBECONFIG "
echo "    echo "KUBECONFIG unset. Ready for dedication." "
echo ""
echo "   alias dedicate-c1='kubectl config view --minify --flatten --context k3d-cluster1 > $HOME/k8s-c1.yaml && export KUBECONFIG=$HOME/k8s-c1.yaml && echo \"\n*** Shell 1 is now dedicated to Cluster 1. Run commands without --context. ***\n\"'"
echo "   alias dedicate-c2='kubectl config view --minify --flatten --context k3d-cluster2 > $HOME/k8s-c2.yaml && export KUBECONFIG=$HOME/k8s-c2.yaml && echo \"\n*** Shell 2 is now dedicated to Cluster 2. Run commands without --context. ***\n\"'"
echo ""
echo "2. OPEN TWO SHELLS & DEDICATE THEM:"
echo "   - **Shell 1 (Current):** Run:  dedicate-c1"
echo "   - **Shell 2 (New Terminal):** Paste aliases, then Run:  dedicate-c2"
echo "   (This prevents conflicts by isolating the KUBECONFIG environment variable in each shell.)"
echo ""
echo "3. GENERATE TOKEN (In Shell 2, dedicated to Cluster 2):"
echo "   skupper token create ~/work/scratch/token-to-cluster2.yaml --token-type cert"
echo ""
echo "4. CREATE LINK (In Shell 1, dedicated to Cluster 1):"
echo "   skupper link create ~/work/scratch/token-to-cluster2.yaml"
# ------------------------------------

echo ""
echo "Script finished. Please follow the 4 POST-SETUP steps above to complete the Skupper link."
