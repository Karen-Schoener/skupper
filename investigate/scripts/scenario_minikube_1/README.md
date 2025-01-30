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
mkdir west
cd west
mkdir secrets
cd secrets

mkdir skupper-site-ca
cd skupper-site-ca

./skupper_site_ca_generate_artifacts.sh

./skupper_site_ca_create_secret.sh

kubectl -n west get secrets

export SKUPPER_CA_DIR=$(pwd)
```

## Step: create west/secret skupper-site-service
```
cd west/secrets

mkdir skupper-site-server
cd skupper-site-server

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_CA_DIR

./skupper_site_server_generate_artifacts.sh

./skupper_site_server_create_secret.sh

kubectl -n west get secrets
```

## Step: create east/secret link1
```
mkdir east
cd east
mkdir secrets
cd secrets

# Verify that SKUPPER_CA_DIR is set
printenv SKUPPER_CA_DIR

./link1_generate_artifacts.sh

./link1_create_secret.sh

kubectl -n east get secrets

# Verify that sites are connected
skupper -n east link status
skupper -n east status -v
skupper -n west status -v
```
