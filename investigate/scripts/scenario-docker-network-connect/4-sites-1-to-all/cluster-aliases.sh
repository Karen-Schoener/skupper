#!/bin/bash
set -e # Exit immediately if a command fails

# --- DEDICATED SHELL INSTRUCTIONS ---
echo "--- POST-SETUP: DEDICATED SHELL WORKFLOW ---"
echo ""
echo ""
echo ""
echo "    unset KUBECONFIG "
echo "    echo "KUBECONFIG unset. Ready for dedication." "
echo ""
echo ""
# --- DEDICATED SHELL INSTRUCTIONS (Using $HOME for robustness) ---
echo "--- POST-SETUP: DEDICATED SHELL WORKFLOW ---"
echo ""
echo "1. ACTIVATE DEDICATED SHELL FUNCTIONS (Run these aliases):"
echo "   # Copy/paste ALL of these lines into your CURRENT shell."
echo "   # They define two new, robust aliases: 'dedicate-c1' and 'dedicate-c2'."
echo "   # They use --flatten and $HOME to prevent certificate/permission errors."
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

