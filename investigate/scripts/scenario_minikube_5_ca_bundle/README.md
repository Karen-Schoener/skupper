# Test case: minikube, metallb; link sites east to west

The scripts in this directory have been tested in a minikube cluster, configured as follows:
* metallb is configured.  This allows the skupper-router in west to be assigned the same IP, assuming that the skupper-site in west is created first.

Goal is: 
  - create 2 sets of credentials: set 1 and set 2
    Each set of credentails contains: 
      1. user-supplied CA credentials: ca.crt
      2. user-supplied skupper-site-server credentials: tls.crt, tls.key
      3. user-supplied connection-token credentials: tls.crt, tls.key
  - configure a secret/ca-bundle
  - populate secret/ca-bundle with 2 sets of CAs: from set 1 and set 2.
  - patch skupper-router deployment to mount volume to secret/ca-bundle.
    patch skupper-router deployment to set env var: SSL_CERT_FILE
  - in west, configure skupper-site-server with user-supplied credentials.
  - in east, configure connection-token with user-supplied credentials.
    verify sites link.

Assumptions:
* west's skupper-router service will run with external-ip: 192.168.49.240

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

## Step: generate credentials: set 1
```
   export TEST_OUTPUT_DIR="./test1"
   export SKUPPER_CA_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca"
   export SKUPPER_SITE_SERVER_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-server"
   export SKUPPER_LINK_NAME="link1"
   export SKUPPER_LINK_DIR="${TEST_OUTPUT_DIR}/east/secrets/${SKUPPER_LINK_NAME}"

   ./user_supplied_ca_generate_artifacts.sh

   # ./skupper_site_server_generate_artifacts.sh 192.168.49.240
     ./skupper_site_server_generate_artifacts.sh mytest-skupper-router-west.local

   ./link_generate_artifacts.sh
```

## Step: generate credentials: set 2
```
   export TEST_OUTPUT_DIR="./test2"
   export SKUPPER_CA_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca"
   export SKUPPER_SITE_SERVER_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-server"
   export SKUPPER_LINK_NAME="link1"
   export SKUPPER_LINK_DIR="${TEST_OUTPUT_DIR}/east/secrets/${SKUPPER_LINK_NAME}"

   ./user_supplied_ca_generate_artifacts.sh

   # ./skupper_site_server_generate_artifacts.sh 192.168.49.240
     ./skupper_site_server_generate_artifacts.sh mytest-skupper-router-west.local

   ./link_generate_artifacts.sh
```

## Step: create ca-bundle files
```
mkdir ca_bundle

cat ./test1/west/secrets/skupper-site-ca/tls.crt > ca_bundle/ca_bundle_test1.pem
cat ./test2/west/secrets/skupper-site-ca/tls.crt > ca_bundle/ca_bundle_test2.pem
cat ./test1/west/secrets/skupper-site-ca/tls.crt ./test2/west/secrets/skupper-site-ca/tls.crt > ca_bundle/ca_bundle_test1_test2.pem
```

## Step: create skupper site west
```
skupper -n west init
```

## Step: create skupper site east
```
skupper -n east init
```

## Step: verify no CA errors in skupper-router logs
```
kubectl -n west logs $(kubectl get -n west pods -l app.kubernetes.io/name=skupper-router -o name) | grep CA
kubectl -n east logs $(kubectl get -n east pods -l app.kubernetes.io/name=skupper-router -o name) | grep CA
```

## Step: create secret ca-bundle
```
./ca_bundle_create_secret_test1.sh
```

## Step: patch skupper-router deployments to mount ca-bundle, set env var SSL_CERT_FILE
```
./deployment_skupper_router_west_patch.sh
./deployment_skupper_router_east_patch.sh
```

## Step: patch west/secret/skupper-site-server with set-1 user-supplied certs
OBSERVATION Note do not delete secret skupper-site-server.  Seems to cause skupper-router error logs if there is no ca in the skupper-site-server secret
```
./skupper_site_server_create_secret_NO_CA.sh
```

## Step: optional: from west skupper-router bash shell, run openssl commands
```
kubectl -n west exec -it $(kubectl get -n west pods -l app.kubernetes.io/name=skupper-router -o name) -- /bin/bash
openssl verify -CAfile /etc/skupper-router-certs/skupper-internal/ca.crt /etc/skupper-router-certs/skupper-internal/tls.crt
(expected to fail)
openssl verify -CAfile /etc/skupper-router-certs/ca-bundle/ca_bundle.pem  /etc/skupper-router-certs/skupper-internal/tls.crt
(expected OK)
openssl verify /etc/skupper-router-certs/skupper-internal/tls.crt
(expected OK)
```

## Step: in east, configure set-1 connection-token
OBSERVATION Note create connect-token with ca.crt specified.  At the moment, skupper will not configure east/configmap/skupper-ineternal/skrouter.json if ca.crt not set in connection token.
```
# ./link_create_secret_NO_CA.sh
./link_create_secret.sh
```

## Step: verify sites are linked
```
# Verify that sites are connected
skupper -n east link status
skupper -n east status -v
skupper -n west status -v
```
