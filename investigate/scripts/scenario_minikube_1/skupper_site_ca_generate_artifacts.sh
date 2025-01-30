#!/bin/bash

# Check if SKUPPER_CA_DIR is set
if [ -z "$SKUPPER_CA_DIR" ]; then
  echo "Error: SKUPPER_CA_DIR environment variable is not set."
  echo "Please set SKUPPER_CA_DIR to the path of the skupper-site-ca directory."
  exit 1
fi

echo "Using SKUPPER_CA_DIR: $SKUPPER_CA_DIR"

# Ensure the output directory exists
mkdir -p $SKUPPER_CA_DIR

# Configuration string for CSR
CSR_CONFIG="
[ req ]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no

[ req_distinguished_name ]
CN = skupper-site-ca

[ req_ext ]
keyUsage = critical, digitalSignature, keyEncipherment, keyCertSign
extendedKeyUsage = serverAuth, clientAuth
basicConstraints = critical, CA:true
subjectKeyIdentifier = hash
subjectAltName = @alt_names

[ alt_names ]
DNS.1 =
"

# Generate RSA private key in PKCS#8 format
openssl genpkey -algorithm RSA -out $SKUPPER_CA_DIR/pkcs8.key -pkeyopt rsa_keygen_bits:2048

# Convert to PKCS#1 format using the traditional flag
openssl rsa -in $SKUPPER_CA_DIR/pkcs8.key -out $SKUPPER_CA_DIR/pkcs1.key -traditional
cp $SKUPPER_CA_DIR/pkcs1.key $SKUPPER_CA_DIR/tls.key
cp $SKUPPER_CA_DIR/pkcs8.key $SKUPPER_CA_DIR/tls.key

# Generate the certificate signing request (CSR)
openssl req -new -key $SKUPPER_CA_DIR/tls.key -out $SKUPPER_CA_DIR/tls.csr -config <(echo "$CSR_CONFIG")

# Generate the self-signed certificate with extensions
openssl x509 -req -days 1825 -in $SKUPPER_CA_DIR/tls.csr -signkey $SKUPPER_CA_DIR/tls.key -out $SKUPPER_CA_DIR/tls.crt -extensions req_ext -extfile <(echo "$CSR_CONFIG")
