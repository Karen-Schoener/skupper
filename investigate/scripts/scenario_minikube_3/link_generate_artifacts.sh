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
  echo "Please set SKUPPER_LINK_NAME to the name of the link."
  exit 1
fi

echo "Using SKUPPER_CA_DIR: $SKUPPER_CA_DIR"
echo "Using SKUPPER_LINK_DIR: $SKUPPER_LINK_DIR"
echo "Using SKUPPER_LINK_NAME: $SKUPPER_LINK_NAME"

# Ensure the output directory exists
mkdir -p "$SKUPPER_LINK_DIR"

# Common Name (CN) based on SKUPPER_LINK_NAME
CN="testlink-${SKUPPER_LINK_NAME}"
# Hardcoded validity period in days
DAYS=365

# Configuration string for CSR
CSR_CONFIG="
[ req ]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no

[ req_distinguished_name ]
CN = $CN

[ req_ext ]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
basicConstraints = critical, CA:false
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
"

# Configuration string for certificate extensions
CERT_CONFIG="
[ req_ext ]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
basicConstraints = critical, CA:false
authorityKeyIdentifier = keyid,issuer
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
"

# Generate RSA private key in PKCS#8 format
openssl genpkey -algorithm RSA -out "$SKUPPER_LINK_DIR/tls.key" -pkeyopt rsa_keygen_bits:2048 2>&1 | tee "$SKUPPER_LINK_DIR/tls_key_debug.log"

# Check if private key generation was successful
if [ ! -f "$SKUPPER_LINK_DIR/tls.key" ]; then
  echo "Error: Failed to generate private key" | tee -a "$SKUPPER_LINK_DIR/tls_key_debug.log"
  exit 1
fi

# Generate the certificate signing request (CSR)
openssl req -new -key "$SKUPPER_LINK_DIR/tls.key" -out "$SKUPPER_LINK_DIR/link.csr" -config <(echo "$CSR_CONFIG") 2>&1 | tee "$SKUPPER_LINK_DIR/csr_debug.log"

# Check if CSR generation was successful
if [ ! -f "$SKUPPER_LINK_DIR/link.csr" ]; then
  echo "Error: Failed to generate CSR" | tee -a "$SKUPPER_LINK_DIR/csr_debug.log"
  exit 1
fi

# Inspect the CSR to verify its contents
openssl req -text -noout -verify -in "$SKUPPER_LINK_DIR/link.csr" 2>&1 | tee -a "$SKUPPER_LINK_DIR/csr_debug.log"

# Sign the CSR with the skupper-site-ca certificate and key
openssl x509 -req -in "$SKUPPER_LINK_DIR/link.csr" -CA "$SKUPPER_CA_DIR/tls.crt" -CAkey "$SKUPPER_CA_DIR/tls.key" -CAcreateserial -out "$SKUPPER_LINK_DIR/tls.crt" -days "$DAYS" -extensions req_ext -extfile <(echo "$CERT_CONFIG") 2>&1 | tee "$SKUPPER_LINK_DIR/tls_crt_debug.log"

# Check if the certificate generation was successful
if [ ! -f "$SKUPPER_LINK_DIR/tls.crt" ]; then
  echo "Error: Failed to generate tls.crt" | tee -a "$SKUPPER_LINK_DIR/tls_crt_debug.log"
  exit 1
fi

echo "Link artifacts generated successfully in $SKUPPER_LINK_DIR" | tee -a "$SKUPPER_LINK_DIR/tls_crt_debug.log"
