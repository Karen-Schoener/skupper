#!/bin/bash

# Configuration string
CONFIG="
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
openssl genpkey -algorithm RSA -out pkcs8.key -pkeyopt rsa_keygen_bits:2048

# Convert to PKCS#1 format using the traditional flag
openssl rsa -in pkcs8.key -out pkcs1.key -traditional
cp pkcs1.key tls.key

# Generate the certificate signing request (CSR)
openssl req -new -key tls.key -out tls.csr -config <(echo "$CONFIG")

# Generate the self-signed certificate with extensions
openssl x509 -req -days 1825 -in tls.csr -signkey tls.key -out tls.crt -extensions req_ext -extfile <(echo "$CONFIG")

