#!/bin/bash

# Check if necessary environment variables are set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the directory containing the updated secrets."
  exit 1
fi

# Check if SKUPPER_SITE_SERVER_DIR is set
if [ -z "$SKUPPER_SITE_SERVER_DIR" ]; then
  echo "Error: SKUPPER_SITE_SERVER_DIR environment variable is not set."
  echo "Please set SKUPPER_SITE_SERVER_DIR to the path of the skupper-site-server directory."
  exit 1
fi

echo "Using SKUPPER_CA_DIR: $SKUPPER_CA_DIR"
echo "Using SKUPPER_SITE_SERVER_DIR: $SKUPPER_SITE_SERVER_DIR"

# Check if necessary files exist in SKUPPER_CA_DIR
if [ ! -f "$SKUPPER_CA_DIR/tls.crt" ]; then
  echo "Error: tls.crt file not found in $SKUPPER_CA_DIR"
  exit 1
fi

if [ ! -f "$SKUPPER_SITE_SERVER_DIR/tls.crt" ]; then
  echo "Error: tls.crt file not found in $SKUPPER_SITE_SERVER_DIR"
  exit 1
fi

if [ ! -f "$SKUPPER_SITE_SERVER_DIR/tls.key" ]; then
  echo "Error: tls.key file not found in $SKUPPER_SITE_SERVER_DIR"
  exit 1
fi

# Encode files in base64
TLS_CRT_BASE64=$(base64 -w 0 "$SKUPPER_CA_DIR/tls.crt")
TLS_KEY_BASE64=$(base64 -w 0 "$SKUPPER_CA_DIR/tls.key")
CA_CRT_BASE64=$(base64 -w 0 "$SKUPPER_CA_DIR/tls.crt")

# Patch the existing skupper-site-server secret with the new values
kubectl -n west patch secret skupper-site-server -p "{\"data\":{\"tls.crt\":\"$TLS_CRT_BASE64\",\"tls.key\":\"$TLS_KEY_BASE64\",\"ca.crt\":\"$CA_CRT_BASE64\"}}"
