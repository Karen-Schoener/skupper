#!/bin/bash

# Check if SKUPPER_CA_DIR and SKUPPER_CA_INTERMEDIATE_DIR are set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the output directory."
  exit 1
fi

if [ -z "$SKUPPER_CA_INTERMEDIATE_DIR" ]; then
  echo "Error: SKUPPER_CA_INTERMEDIATE_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_INTERMEDIATE_DIR to the path of the Intermediate CA directory."
  exit 1
fi

echo "Using SKUPPER_CA_DIR: $SKUPPER_CA_DIR"
echo "Using SKUPPER_CA_INTERMEDIATE_DIR: $SKUPPER_CA_INTERMEDIATE_DIR"

# Ensure the output directories exist
mkdir -p "$SKUPPER_CA_DIR"
mkdir -p "$SKUPPER_CA_INTERMEDIATE_DIR"

# Default values (can be overridden by arguments)
CN=${1:-Intermediate-CA}
DAYS=365

# Configuration string for CSR
CSR_CONFIG="
[ req ]
distinguished_name = req_distinguished_name
prompt = no

[ req_distinguished_name ]
CN = $CN
"

# Set a default PEM passphrase
PEM_PASS="test"

# Generate Intermediate CA private key in PKCS#8 format
openssl genpkey -algorithm RSA -out "$SKUPPER_CA_INTERMEDIATE_DIR/tls.key" -aes256 -pass pass:$PEM_PASS 2>&1 | tee "$SKUPPER_CA_INTERMEDIATE_DIR/tls_key_debug.log"

# Check if private key generation was successful
if [ ! -f "$SKUPPER_CA_INTERMEDIATE_DIR/tls.key" ]; then
  echo "Error: Failed to generate Intermediate CA private key" | tee -a "$SKUPPER_CA_INTERMEDIATE_DIR/tls_key_debug.log"
  exit 1
fi

# Generate the certificate signing request (CSR) for Intermediate CA
openssl req -new -key "$SKUPPER_CA_INTERMEDIATE_DIR/tls.key" -out "$SKUPPER_CA_INTERMEDIATE_DIR/tls.csr" -config <(echo "$CSR_CONFIG") -subj "/CN=$CN" -passin pass:$PEM_PASS 2>&1 | tee "$SKUPPER_CA_INTERMEDIATE_DIR/tls_csr_debug.log"

# Check if CSR generation was successful
if [ ! -f "$SKUPPER_CA_INTERMEDIATE_DIR/tls.csr" ]; then
  echo "Error: Failed to generate Intermediate CA CSR" | tee -a "$SKUPPER_CA_INTERMEDIATE_DIR/tls_csr_debug.log"
  exit 1
fi

# Inspect the CSR to verify its contents
openssl req -text -noout -verify -in "$SKUPPER_CA_INTERMEDIATE_DIR/tls.csr" 2>&1 | tee -a "$SKUPPER_CA_INTERMEDIATE_DIR/tls_csr_debug.log"

# Generate the Intermediate CA certificate signed by the Root CA
openssl x509 -req -in "$SKUPPER_CA_INTERMEDIATE_DIR/tls.csr" -CA "$SKUPPER_CA_DIR/tls.crt" -CAkey "$SKUPPER_CA_DIR/tls.key" -CAcreateserial -out "$SKUPPER_CA_INTERMEDIATE_DIR/tls.crt" -days "$DAYS" -sha256 -extfile <(echo "basicConstraints=CA:TRUE,pathlen:0") 2>&1 | tee "$SKUPPER_CA_INTERMEDIATE_DIR/tls_crt_debug.log"

# Check if the certificate generation was successful
if [ ! -f "$SKUPPER_CA_INTERMEDIATE_DIR/tls.crt" ]; then
  echo "Error: Failed to generate Intermediate CA certificate" | tee -a "$SKUPPER_CA_INTERMEDIATE_DIR/tls_crt_debug.log"
  exit 1
fi

# Output the filenames for tls.crt and tls.key
echo "Intermediate CA certificate: $SKUPPER_CA_INTERMEDIATE_DIR/tls.crt"
echo "Intermediate CA private key: $SKUPPER_CA_INTERMEDIATE_DIR/tls.key"

