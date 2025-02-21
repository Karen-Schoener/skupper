#!/bin/bash

# Set SKUPPER_SITE_SERVER_CN to default if not provided
SKUPPER_SITE_SERVER_CN=${1:-mytest-skupper-router-west.local}

# Check if SKUPPER_CA_INTERMEDIATE_DIR is set
if [ -z "$SKUPPER_CA_INTERMEDIATE_DIR" ]; then
  echo "Error: SKUPPER_CA_INTERMEDIATE_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_INTERMEDIATE_DIR to the path of the skupper-site-ca-intermediate directory."
  exit 1
fi

# Check if SKUPPER_CA_DIR is set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the skupper-site-ca directory."
  exit 1
fi

# Check if SKUPPER_LINK_DIR is set
if [ -z "$SKUPPER_LINK_DIR" ]; then
  echo "Error: SKUPPER_LINK_DIR environment variable is not set."
  echo "Please set SKUPPER_LINK_DIR to the path of the skupper-site-server directory."
  exit 1
fi

# Check if SKUPPER_LINK_NAME is set
if [ -z "$SKUPPER_LINK_NAME" ]; then
  echo "Error: SKUPPER_LINK_NAME environment variable is not set."
  echo "Please set SKUPPER_LINK_NAME to the name of the link."
  exit 1
fi

echo "Using SKUPPER_CA_DIR: $SKUPPER_CA_DIR"
echo "Using SKUPPER_CA_INTERMEDIATE_DIR: $SKUPPER_CA_INTERMEDIATE_DIR"
echo "Using SKUPPER_LINK_DIR: $SKUPPER_LINK_DIR"
echo "Using SKUPPER_LINK_NAME: $SKUPPER_LINK_NAME"
echo "Using SKUPPER_SITE_SERVER_CN: $SKUPPER_SITE_SERVER_CN"

# Concatenate the link certificate and the intermediate certificate into a single file
cat $SKUPPER_LINK_DIR/tls.crt $SKUPPER_CA_INTERMEDIATE_DIR/tls.crt > $SKUPPER_LINK_DIR/combined-tls.crt

# Create the initial secret with tls.crt and tls.key
kubectl -n east create secret tls $SKUPPER_LINK_NAME --cert=$SKUPPER_LINK_DIR/combined-tls.crt --key=$SKUPPER_LINK_DIR/tls.key

# Patch the secret to add the root CA certificate
kubectl -n east patch secret $SKUPPER_LINK_NAME -p "$(printf '{"data":{"ca.crt":"%s"}}' $(base64 -w 0 $SKUPPER_CA_DIR/tls.crt))"

# Patch the secret to add annotations and labels
kubectl -n east patch secret $SKUPPER_LINK_NAME -p "{
    \"metadata\": {
        \"annotations\": {
            \"edge-host\": \"$SKUPPER_SITE_SERVER_CN\",
            \"edge-port\": \"45671\",
            \"inter-router-host\": \"$SKUPPER_SITE_SERVER_CN\",
            \"inter-router-port\": \"55671\",
            \"skupper.io/generated-by\": \"0f663c41-2974-45b9-aa05-ffaf9259a60c\",
            \"skupper.io/site-version\": \"v1-dev-release-1-g5834403-modified\"
        },
        \"labels\": {
            \"skupper.io/type\": \"connection-token\"
        }
    }
}"
