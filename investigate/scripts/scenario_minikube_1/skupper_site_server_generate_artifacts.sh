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

# Configuration string for CSR
CSR_CONFIG="
[ req ]
distinguished_name = req_distinguished_name
req_extensions = req_ext
prompt = no

[ req_distinguished_name ]
CN = 192.168.49.240

[ req_ext ]
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
basicConstraints = critical, CA:false
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = skupper-router.west
DNS.2 = skupper-router.west.svc.cluster.local
DNS.3 = 192.168.49.240
IP.1 = 192.168.49.240
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
DNS.3 = 192.168.49.240
IP.1 = 192.168.49.240
"

# Generate RSA private key in PKCS#8 format
openssl genpkey -algorithm RSA -out $SKUPPER_SITE_SERVER_DIR/server-pkcs8.key -pkeyopt rsa_keygen_bits:2048

# Convert to PKCS#1 format using the traditional flag
openssl rsa -in $SKUPPER_SITE_SERVER_DIR/server-pkcs8.key -out $SKUPPER_SITE_SERVER_DIR/server-pkcs1.key -traditional
cp $SKUPPER_SITE_SERVER_DIR/server-pkcs1.key $SKUPPER_SITE_SERVER_DIR/tls.key

# Generate the certificate signing request (CSR)
openssl req -new -key $SKUPPER_SITE_SERVER_DIR/tls.key -out $SKUPPER_SITE_SERVER_DIR/skupper-site-server.csr -config <(echo "$CSR_CONFIG")

# Sign the CSR with the skupper-site-ca certificate and key
openssl x509 -req -in $SKUPPER_SITE_SERVER_DIR/skupper-site-server.csr -CA $SKUPPER_CA_DIR/tls.crt -CAkey $SKUPPER_CA_DIR/tls.key -CAcreateserial -out $SKUPPER_SITE_SERVER_DIR/tls.crt -days 1825 -extensions req_ext -extfile <(echo "$CERT_CONFIG")

