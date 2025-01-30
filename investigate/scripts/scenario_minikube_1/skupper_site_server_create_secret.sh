#!/bin/bash

# Check if SKUPPER_CA_DIR is set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the skupper-site-ca directory."
  exit 1
fi

# Check if SKUPPER_SITE_SERVER_DIR is set
if [ -z "$SKUPPER_SITE_SERVER_DIR" ]; then
  echo "Error: SKUPPER_SITE_SERVER_DIR environment variable is not set."
  echo "Please set SKUPPER_SITE_SERVER_DIR to the path of the skupper-site-server directory."
  exit 1
fi

kubectl create secret tls skupper-site-server \
  --cert=$SKUPPER_SITE_SERVER_DIR/tls.crt \
  --key=$SKUPPER_SITE_SERVER_DIR/tls.key \
  -n west

kubectl -n west patch secret skupper-site-server -p "$(printf '{"data":{"ca.crt":"%s"}}' $(base64 -w 0 $SKUPPER_CA_DIR/tls.crt))"

