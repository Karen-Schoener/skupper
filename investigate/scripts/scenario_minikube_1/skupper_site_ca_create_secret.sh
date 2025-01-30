#!/bin/bash


# Check if SKUPPER_CA_DIR is set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the skupper-site-ca directory."
  exit 1
fi

kubectl create secret tls skupper-site-ca \
  --cert=$SKUPPER_CA_DIR/tls.crt \
  --key=$SKUPPER_CA_DIR/tls.key \
  -n west

