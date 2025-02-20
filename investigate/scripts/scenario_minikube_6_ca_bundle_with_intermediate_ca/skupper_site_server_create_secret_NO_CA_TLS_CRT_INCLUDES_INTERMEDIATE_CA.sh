#!/bin/bash

# Check if SKUPPER_CA_DIR is set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the skupper-site-ca-intermediate directory."
  exit 1
fi

# Check if SKUPPER_CA_INTERMEDIATE_DIR is set
if [ -z "$SKUPPER_CA_INTERMEDIATE_DIR" ]; then
  echo "Error: SKUPPER_CA_INTERMEDIATE_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_INTERMEDIATE_DIR to the path of the skupper-site-ca-intermediate directory."
  exit 1
fi

# Check if SKUPPER_SITE_SERVER_DIR is set
if [ -z "$SKUPPER_SITE_SERVER_DIR" ]; then
  echo "Error: SKUPPER_SITE_SERVER_DIR environment variable is not set."
  echo "Please set SKUPPER_SITE_SERVER_DIR to the path of the skupper-site-server directory."
  exit 1
fi

# Concatenate the server certificate and the intermediate certificate into a single file
cat $SKUPPER_SITE_SERVER_DIR/tls.crt $SKUPPER_CA_INTERMEDIATE_DIR/tls.crt > $SKUPPER_SITE_SERVER_DIR/combined-tls.crt

# Create a YAML manifest for the secret and place it in SKUPPER_SITE_SERVER_DIR
cat <<EOF > $SKUPPER_SITE_SERVER_DIR/skupper-site-server-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: skupper-site-server
type: kubernetes.io/tls
data:
  tls.crt: $(base64 -w 0 $SKUPPER_SITE_SERVER_DIR/combined-tls.crt)
  tls.key: $(base64 -w 0 $SKUPPER_SITE_SERVER_DIR/tls.key)
EOF

# Apply the secret using the YAML manifest from SKUPPER_SITE_SERVER_DIR
kubectl -n west apply -f $SKUPPER_SITE_SERVER_DIR/skupper-site-server-secret.yaml
kubectl -n east apply -f $SKUPPER_SITE_SERVER_DIR/skupper-site-server-secret.yaml

kubectl -n west get pods -l skupper.io/component=router -o jsonpath="{.items[0].metadata.name}" | xargs -I{} kubectl -n west annotate pod {} testing-secret-updated/force-reconcile=$(date +%s) --overwrite
kubectl -n east get pods -l skupper.io/component=router -o jsonpath="{.items[0].metadata.name}" | xargs -I{} kubectl -n east annotate pod {} testing-secret-updated/force-reconcile=$(date +%s) --overwrite

