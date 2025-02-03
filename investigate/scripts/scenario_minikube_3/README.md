# Test case: minikube, metallb; link sites east to west

The scripts in this directory have been tested in a minikube cluster, configured as follows:
* metallb is configured.  This allows the skupper-router in west to be assigned the same IP, assuming that the skupper-site in west is created first.

Goal is: establish link from east to west via DNS name.

Assumptions:
* west's skupper-router service will run with external-ip: 192.168.49.240

These scripts perform the following:
* create user-supplied, self-signed CA artifacts
* create user-supplied skupper-site-server artifacts
* create secret in west: skupper-site-server
* create user-supplied connection-token artifacts
* create secret in east: connection-token

## Step: create test DNS name to west's skupper-router service external-ip
```
echo "192.168.49.240 mytest-skupper-router-west.local" | sudo tee -a /etc/hosts
```

## Step: create namespaces: east, west
```
kubectl delete namespace east
kubectl delete namespace west

kubectl create namespace east
kubectl create namespace west
```

## Step: create west/secret skupper-site-ca
```
export TEST_OUTPUT_DIR="./test1"

export SKUPPER_CA_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca"

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_CA_DIR

./user_supplied_ca_generate_artifacts.sh

```

## Step: create west/secret skupper-site-service
```
printenv TEST_OUTPUT_DIR
printenv SKUPPER_CA_DIR
printenv SKUPPER_SITE_SERVER_DIR

export SKUPPER_SITE_SERVER_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-server"

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_SITE_SERVER_DIR

# ./skupper_site_server_generate_artifacts.sh 192.168.49.240
  ./skupper_site_server_generate_artifacts.sh mytest-skupper-router-west.local

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
printenv TEST_OUTPUT_DIR
printenv SKUPPER_CA_DIR
printenv SKUPPER_SITE_SERVER_DIR
printenv SKUPPER_LINK_NAME
printenv SKUPPER_LINK_DIR

export SKUPPER_LINK_NAME="link1"
export SKUPPER_LINK_DIR="${TEST_OUTPUT_DIR}/east/secrets/${SKUPPER_LINK_NAME}"

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_CA_DIR
printenv SKUPPER_LINK_DIR
printenv SKUPPER_LINK_NAME

./link_generate_artifacts.sh

./link_create_secret.sh mytest-skupper-router-west.local

kubectl -n east get secrets

# Verify that sites are connected
skupper -n east link status
skupper -n east status -v
skupper -n west status -v
```
