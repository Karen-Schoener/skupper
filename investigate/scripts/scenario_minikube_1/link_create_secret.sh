#!/bin/bash

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
  echo "Please set SKUPPER_LINK_NAME to the path of the skupper-site-server directory."
  exit 1
fi

kubectl -n east create secret tls $SKUPPER_LINK_NAME --cert=$SKUPPER_LINK_DIR/tls.crt --key=$SKUPPER_LINK_DIR/tls.key 

kubectl -n east patch secret $SKUPPER_LINK_NAME -p "$(printf '{"data":{"ca.crt":"%s"}}' $(base64 -w 0 $SKUPPER_CA_DIR/tls.crt))"

kubectl -n east patch secret $SKUPPER_LINK_NAME -p '{
    "metadata": {
        "annotations": {
            "edge-host": "192.168.49.240",
            "edge-port": "45671",
            "inter-router-host": "192.168.49.240",
            "inter-router-port": "55671",
            "skupper.io/generated-by": "0f663c41-2974-45b9-aa05-ffaf9259a60c",
            "skupper.io/site-version": "v1-dev-release-1-g5834403-modified"
        },
        "labels": {
            "skupper.io/type": "connection-token"
        }
    }
}'

