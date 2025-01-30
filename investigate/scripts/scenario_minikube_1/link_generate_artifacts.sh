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

# Configuration string for CSR
CSR_CONFIG="
[ req ]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no

[ req_distinguished_name ]
CN = testnamecert

[ req_ext ]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
basicConstraints = critical, CA:false
subjectAltName = @alt_names

[ alt_names ]
DNS.1 =
"

# Configuration string for certificate extensions
CERT_CONFIG="
[ req ]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no

[ req_distinguished_name ]
CN = testnamecert

[ req_ext ]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
basicConstraints = critical, CA:false
subjectAltName = @alt_names

[ alt_names ]
DNS.1 =
"

# Generate RSA private key in PKCS#8 format
openssl genpkey -algorithm RSA -out $SKUPPER_LINK_DIR/pkcs8.key -pkeyopt rsa_keygen_bits:2048

# Convert to PKCS#1 format using the traditional flag
openssl rsa -in $SKUPPER_LINK_DIR/pkcs8.key -out $SKUPPER_LINK_DIR/pkcs1.key -traditional
cp $SKUPPER_LINK_DIR/pkcs1.key $SKUPPER_LINK_DIR/tls.key

# Generate the certificate signing request (CSR)
openssl req -new -key $SKUPPER_LINK_DIR/tls.key -out $SKUPPER_LINK_DIR/link.csr -config <(echo "$CSR_CONFIG")

# Sign the CSR with the skupper-site-ca certificate and key
openssl x509 -req -in $SKUPPER_LINK_DIR/link.csr -CA $SKUPPER_CA_DIR/tls.crt -CAkey $SKUPPER_CA_DIR/tls.key -CAcreateserial -out $SKUPPER_LINK_DIR/tls.crt -days 1825 -extensions req_ext -extfile <(echo "$CERT_CONFIG")
