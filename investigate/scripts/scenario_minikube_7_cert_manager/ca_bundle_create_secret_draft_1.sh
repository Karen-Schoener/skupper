
# Check if CA_BUNDLE_DIR is set
if [ -z "$CA_BUNDLE_DIR" ]; then
  echo "Error: CA_BUNDLE_DIR environment variable is not set."
  echo "Please set CA_BUNDLE_DIR to the path of the directory containing ca_bundle.pem."
  exit 1
fi

# Check if CA_BUNDLE_FILE is set
if [ -z "$CA_BUNDLE_FILE" ]; then
  echo "Error: CA_BUNDLE_FILE environment variable is not set."
  echo "Please set CA_BUNDLE_FILE to the name of the CA bundle file (e.g., ca_bundle_test1.pem)."
  exit 1
fi

# Base64 encode the CA bundle file and save it to a variable
CA_BUNDLE_CONTENT=$(base64 -w 0 $CA_BUNDLE_DIR/$CA_BUNDLE_FILE)

# Create a YAML manifest for the ca-bundle secret and save it to a variable
SECRET_MANIFEST=$(cat <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ca-bundle
  namespace: west
type: Opaque
data:
  ca_bundle.pem: $CA_BUNDLE_CONTENT
EOF
)

# Apply the secret using the YAML manifest stored in the variable
echo "$SECRET_MANIFEST" | kubectl apply -f -

