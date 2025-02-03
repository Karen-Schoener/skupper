#!/bin/bash

# Check if SKUPPER_CA_DIR is set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the CA directory."
  exit 1
fi

# Check if SKUPPER_SITE_SERVER_DIR is set
if [ -z "$SKUPPER_SITE_SERVER_DIR" ]; then
  echo "Error: SKUPPER_SITE_SERVER_DIR environment variable is not set."
  echo "Please set SKUPPER_SITE_SERVER_DIR to the path of the output directory."
  exit 1
fi

# Check if CN (Common Name) is provided
if [ -z "$1" ]; then
  echo "Error: Common Name (CN) must be provided as the first argument."
  echo "Usage: $0 <CN>"
  echo "Example: $0 192.168.49.240"
  exit 1
fi

echo "Using SKUPPER_CA_DIR: $SKUPPER_CA_DIR"
echo "Using SKUPPER_SITE_SERVER_DIR: $SKUPPER_SITE_SERVER_DIR"

# Ensure the output directory exists
mkdir -p "$SKUPPER_SITE_SERVER_DIR"

# Mandatory Common Name (CN)
CN=$1
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
DNS.1 = skupper-router.west
DNS.2 = skupper-router.west.svc.cluster.local
DNS.3 = $CN
IP.1 = $CN
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
DNS.1 = skupper-router.west
DNS.2 = skupper-router.west.svc.cluster.local
DNS.3 = $CN
IP.1 = $CN
"

# Generate RSA private key in PKCS#8 format
openssl genpkey -algorithm RSA -out "$SKUPPER_SITE_SERVER_DIR/tls.key" -pkeyopt rsa_keygen_bits:2048 2>&1 | tee "$SKUPPER_SITE_SERVER_DIR/tls_key_debug.log"

# Check if private key generation was successful
if [ ! -f "$SKUPPER_SITE_SERVER_DIR/tls.key" ]; then
  echo "Error: Failed to generate private key" | tee -a "$SKUPPER_SITE_SERVER_DIR/tls_key_debug.log"
  exit 1
fi

# Generate the certificate signing request (CSR)
openssl req -new -key "$SKUPPER_SITE_SERVER_DIR/tls.key" -out "$SKUPPER_SITE_SERVER_DIR/skupper-site-server.csr" -config <(echo "$CSR_CONFIG") 2>&1 | tee "$SKUPPER_SITE_SERVER_DIR/csr_debug.log"

# Check if CSR generation was successful
if [ ! -f "$SKUPPER_SITE_SERVER_DIR/skupper-site-server.csr" ]; then
  echo "Error: Failed to generate CSR" | tee -a "$SKUPPER_SITE_SERVER_DIR/csr_debug.log"
  exit 1
fi

# Inspect the CSR to verify its contents
openssl req -text -noout -verify -in "$SKUPPER_SITE_SERVER_DIR/skupper-site-server.csr" 2>&1 | tee -a "$SKUPPER_SITE_SERVER_DIR/csr_debug.log"

# Sign the CSR with the CA certificate and key
openssl x509 -req -in "$SKUPPER_SITE_SERVER_DIR/skupper-site-server.csr" -CA "$SKUPPER_CA_DIR/tls.crt" -CAkey "$SKUPPER_CA_DIR/tls.key" -CAcreateserial -out "$SKUPPER_SITE_SERVER_DIR/tls.crt" -days "$DAYS" -extensions req_ext -extfile <(echo "$CERT_CONFIG") 2>&1 | tee "$SKUPPER_SITE_SERVER_DIR/tls_crt_debug.log"

# Check if the certificate generation was successful
if [ ! -f "$SKUPPER_SITE_SERVER_DIR/tls.crt" ]; then
  echo "Error: Failed to generate tls.crt" | tee -a "$SKUPPER_SITE_SERVER_DIR/tls_crt_debug.log"
  exit 1
fi

echo "Skupper site server artifacts generated successfully in $SKUPPER_SITE_SERVER_DIR" | tee -a "$SKUPPER_SITE_SERVER_DIR/tls_crt_debug.log"
