#!/bin/bash

# Check if SKUPPER_CA_DIR is set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the output directory."
  exit 1
fi

echo "Using SKUPPER_CA_DIR: $SKUPPER_CA_DIR"

# Ensure the output directory exists
mkdir -p "$SKUPPER_CA_DIR"

# Default values (can be overridden by arguments)
CN=${1:-my-root-ca-cert}
DAYS=365

# Configuration string for CSR
CSR_CONFIG="
[ req ]
distinguished_name = req_distinguished_name
prompt = no

[ req_distinguished_name ]
CN = $CN
"

# Configuration string for the certificate
CERT_CONFIG="
[ x509_exts ]
basicConstraints = critical,CA:TRUE
keyUsage = critical,keyCertSign,cRLSign
"

# Generate RSA private key in PKCS#8 format
openssl genpkey -algorithm RSA -out "$SKUPPER_CA_DIR/tls.key" -pkeyopt rsa_keygen_bits:2048 2>&1 | tee "$SKUPPER_CA_DIR/tls_key_debug.log"

# Check if private key generation was successful
if [ ! -f "$SKUPPER_CA_DIR/tls.key" ]; then
  echo "Error: Failed to generate private key" | tee -a "$SKUPPER_CA_DIR/tls_key_debug.log"
  exit 1
fi

# Generate the certificate signing request (CSR)
openssl req -new -key "$SKUPPER_CA_DIR/tls.key" -out "$SKUPPER_CA_DIR/tls.csr" -config <(echo "$CSR_CONFIG") 2>&1 | tee "$SKUPPER_CA_DIR/csr_debug.log"

# Check if CSR generation was successful
if [ ! -f "$SKUPPER_CA_DIR/tls.csr" ]; then
  echo "Error: Failed to generate CSR" | tee -a "$SKUPPER_CA_DIR/csr_debug.log"
  exit 1
fi

# Inspect the CSR to verify its contents
openssl req -text -noout -verify -in "$SKUPPER_CA_DIR/tls.csr" 2>&1 | tee -a "$SKUPPER_CA_DIR/csr_debug.log"

# Generate the self-signed certificate with CA:TRUE
openssl x509 -req -days "$DAYS" -in "$SKUPPER_CA_DIR/tls.csr" -signkey "$SKUPPER_CA_DIR/tls.key" -out "$SKUPPER_CA_DIR/tls.crt" -extfile <(echo "$CERT_CONFIG") -extensions x509_exts 2>&1 | tee "$SKUPPER_CA_DIR/tls_crt_debug.log"

# Check if the certificate generation was successful
if [ ! -f "$SKUPPER_CA_DIR/tls.crt" ]; then
  echo "Error: Failed to generate tls.crt" | tee -a "$SKUPPER_CA_DIR/tls_crt_debug.log"
  exit 1
fi

echo "CA artifacts generated successfully in $SKUPPER_CA_DIR" | tee -a "$SKUPPER_CA_DIR/tls_crt_debug.log"
