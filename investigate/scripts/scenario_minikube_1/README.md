# Test case: minikube, metallb; link sites east to west

The scripts in this directory have been tested in a minikube cluster, configured as follows:
* metallb is configured.  This allows the skupper-router in west to be assigned the same IP, assuming that the skupper-site in west is created first.

## Step: create namespaces: east, west
```
kubectl delete namespace east
kubectl delete namespace west

kubectl create namespace east
kubectl create namespace west
```

## Step: create west/secret skupper-site-ca
```
export TEST_OUTPUT_DIR="$(pwd)/test"

export SKUPPER_CA_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca"
mkdir -p "${SKUPPER_CA_DIR}"

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_CA_DIR

./skupper_site_ca_generate_artifacts.sh

./skupper_site_ca_create_secret.sh

kubectl -n west get secrets
```

## Step: create west/secret skupper-site-service
```
export TEST_OUTPUT_DIR="$(pwd)/test"
export SKUPPER_CA_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca"

export SKUPPER_SITE_SERVER_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-server"
mkdir -p "${SKUPPER_SITE_SERVER_DIR}"

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_CA_DIR
printenv SKUPPER_SITE_SERVER_DIR

./skupper_site_server_generate_artifacts.sh

./skupper_site_server_create_secret.sh

kubectl -n west get secrets
```

## Step: skupper -n west init
```
kubectl -n west get secret

skupper -n west init

kubectl -n west get secret

kubectl -n west get service

# Confirm that skupper-router service gets expected external-ip
```

## Step: skupper -n east init
```
skupper -n east init
```

## Step: create east/secret link
```
export TEST_OUTPUT_DIR="$(pwd)/test"
export SKUPPER_CA_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca"

export SKUPPER_LINK_NAME="link1"

export SKUPPER_LINK_DIR="${TEST_OUTPUT_DIR}/east/secrets/${SKUPPER_LINK_NAME}"
mkdir -p "${SKUPPER_LINK_DIR}"

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_CA_DIR
printenv SKUPPER_LINK_DIR
printenv SKUPPER_LINK_NAME

./link_generate_artifacts.sh

./link_create_secret.sh

kubectl -n east get secrets

# Verify that sites are connected
skupper -n east link status
skupper -n east status -v
skupper -n west status -v
```
