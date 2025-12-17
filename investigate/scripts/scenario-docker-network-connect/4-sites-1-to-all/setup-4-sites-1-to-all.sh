#!/bin/bash
set -e # Exit immediately if a command fails

NUM_SITES=4

# --- New: INITIAL ENVIRONMENT CHECK (KUBECONFIG) ---
if [ -n "$KUBECONFIG" ]; then
    echo "========================================================"
    echo "🛑 WARNING: KUBECONFIG Environment Variable is Set."
    echo "========================================================"
    echo "The setup script requires KUBECONFIG to be unset to ensure"
    echo "that k3d correctly merges all four cluster contexts into"
    echo "your default system configuration (~/.kube/config)."
    echo ""
    echo "Currently set to: $KUBECONFIG"
    echo ""
    echo "Please unset the variable and restart the script:"
    echo ""
    echo "    unset KUBECONFIG"
    echo "    ./setup-4-sites.sh"
    echo ""
    exit 1
fi
# --- End KUBECONFIG Check ---


echo "========================================================"
echo "PHASE 1: INFRASTRUCTURE SETUP (${NUM_SITES} SITES)"
echo "========================================================"

# --- ROBUST TWO-PASS CLEANUP ---
echo "[1/4] Cleaning up old clusters (1 to ${NUM_SITES})..."
for i in $(seq 1 $NUM_SITES); do
    k3d cluster delete cluster$i 2>/dev/null || true
done

echo "[2/4] Cleaning up old Docker networks (1 to ${NUM_SITES})..."
for i in $(seq 1 $NUM_SITES); do
    docker network rm cluster$i-net 2>/dev/null || true
done
# -------------------------------

echo "[3/4] Creating Clusters and Networks (1 to ${NUM_SITES})..."
sleep 2

for i in $(seq 1 $NUM_SITES); do
    echo "    -> Processing Cluster $i..."
    docker network create cluster$i-net

    # Unique Load Balancer Port: 9081, 9082, 9083, 9084...
    LB_PORT=$((9080 + i))

    k3d cluster create cluster$i \
      --network "cluster$i-net" \
      -p "$LB_PORT:80@loadbalancer" \
      --wait > /dev/null

    # k3d will now use the system default location because KUBECONFIG is unset.
    k3d kubeconfig merge cluster$i --kubeconfig-merge-default
done

# 4. NETWORKING (Connect Cluster 1 to the network of all other clusters)
echo "[4/4] Bridging Cluster 1 to cluster2-net, cluster3-net, and cluster4-net for hub-and-spoke communication..."
sleep 5
HUB_NODE="k3d-cluster1-server-0"

for i in $(seq 2 $NUM_SITES); do
    TARGET_NETWORK="cluster${i}-net"
    echo "    -> Connecting ${HUB_NODE} to ${TARGET_NETWORK}..."
    docker network connect "$TARGET_NETWORK" "$HUB_NODE"
done

echo "========================================================"
echo "PHASE 2: SKUPPER INSTALLATION (ALL SITES)"
echo "========================================================"

for i in $(seq 1 $NUM_SITES); do
    CLUSTER_NAME="k3d-cluster$i"
    SITE_NAME="site-$i"
    LB_HOST="k3d-cluster$i-serverlb"

    echo "[Site $i] Initializing Skupper on $CLUSTER_NAME..."

    # Switch context for the kubectl/skupper commands
    kubectl config use-context "$CLUSTER_NAME"

    skupper init --site-name "$SITE_NAME" \
        --enable-console \
        --enable-flow-collector \
        --ingress loadbalancer \
        --router-ingress-host "$LB_HOST"

    echo "    -> Checking Skupper site-$i status..."
    skupper status
done

echo "========================================================"
echo "SETUP COMPLETE"
echo "========================================================"

# --- DEDICATED SHELL INSTRUCTIONS (Using $HOME for robustness) ---
echo "--- POST-SETUP: DEDICATED SHELL WORKFLOW ---"
# ... (The dedication instructions remain the same) ...
# ... (Note: The dedication aliases intentionally set KUBECONFIG later.)

# 1. ACTIVATE DEDICATED SHELL FUNCTIONS (Copy/paste ALL aliases below):
echo ""
echo "1. ACTIVATE DEDICATED SHELL FUNCTIONS (Copy/paste ALL aliases below):"
echo "    # The aliases now include C3 and C4."
echo ""
echo "    unset KUBECONFIG "
echo "    echo \"KUBECONFIG unset. Ready for dedication.\""
echo ""
echo "    alias dedicate-c1='kubectl config view --minify --flatten --context k3d-cluster1 > \$HOME/k8s-c1.yaml && export KUBECONFIG=\$HOME/k8s-c1.yaml && echo \"\n*** Shell 1 is now dedicated to Cluster 1. Run commands without --context. ***\n\"'"
echo "    alias dedicate-c2='kubectl config view --minify --flatten --context k3d-cluster2 > \$HOME/k8s-c2.yaml && export KUBECONFIG=\$HOME/k8s-c2.yaml && echo \"\n*** Shell 2 is now dedicated to Cluster 2. Run commands without --context. ***\n\"'"
echo "    alias dedicate-c3='kubectl config view --minify --flatten --context k3d-cluster3 > \$HOME/k8s-c3.yaml && export KUBECONFIG=\$HOME/k8s-c3.yaml && echo \"\n*** Shell 3 is now dedicated to Cluster 3. Run commands without --context. ***\n\"'"
echo "    alias dedicate-c4='kubectl config view --minify --flatten --context k3d-cluster4 > \$HOME/k8s-c4.yaml && export KUBECONFIG=\$HOME/k8s-c4.yaml && echo \"\n*** Shell 4 is now dedicated to Cluster 4. Run commands without --context. ***\n\"'"
echo ""
echo "2. OPEN FOUR SHELLS & DEDICATE THEM:"
echo "    - **Shell 1 (Hub):** Run:  dedicate-c1"
echo "    - **Shell 2 (Spoke 1):** Paste aliases, then Run:  dedicate-c2"
echo "    - **Shell 3 (Spoke 2):** Paste aliases, then Run:  dedicate-c3"
echo "    - **Shell 4 (Spoke 3):** Paste aliases, then Run:  dedicate-c4"
echo ""
echo "3. GENERATE TOKENS (In Shells 2, 3, and 4) - Using token-to-c[X].yaml:"
echo "    - **Shell 2:** skupper token create ~/tokens/token-to-c2.yaml --token-type cert"
echo "    - **Shell 3:** skupper token create ~/tokens/token-to-c3.yaml --token-type cert"
echo "    - **Shell 4:** skupper token create ~/tokens/token-to-c4.yaml --token-type cert"
echo ""
echo "4. CREATE LINKS (In Shell 1, dedicated to Cluster 1 - Hub):"
echo "    - skupper link create ~/tokens/token-to-c2.yaml"
echo "    - skupper link create ~/tokens/token-to-c3.yaml"
echo "    - skupper link create ~/tokens/token-to-c4.yaml"

echo ""
echo "Script finished. Please follow the 4 POST-SETUP steps above to complete the Skupper mesh link."
