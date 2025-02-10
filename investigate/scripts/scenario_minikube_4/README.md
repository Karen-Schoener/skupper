# Test case: minikube, metallb; link sites east to west

The scripts in this directory have been tested in a minikube cluster, configured as follows:
* metallb is configured.  This allows the skupper-router in west to be assigned the same IP, assuming that the skupper-site in west is created first.

Goal is: 
* establish link from east to west via DNS name.
* Place user-supplied CA in trusted directory.
* When creating secrets skupper-site-server, connection-token, do not populate field ca.crt.
  This should result in skupper using CAs from the trusted CA store of the host OS.

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

## Step: copy resulting CA certificate to trusted store on host OS

Testing with minikube running on fedora VM.

The trusted CA certificates are stored in the /etc/pki/ca-trust/ directory. 

To ensure the CA certificate is trusted by the system, add it to the trusted store and update the CA trust settings.
```
    sudo cp your-ca.crt /etc/pki/ca-trust/source/anchors/
```

Update the CA Trust: Run the following command to update the CA trust settings:
```
    sudo update-ca-trust
```

This will ensure that your CA certificate is recognized and trusted by the system. 
After doing this, any certificates signed by your CA should be trusted by the Fedora host.

### Example
```
    [kschoener@fedora scenario_minikube_4]$ printenv SKUPPER_CA_DIR
    ./test1/west/secrets/skupper-site-ca
    [kschoener@fedora scenario_minikube_4]$ ls $SKUPPER_CA_DIR
    csr_debug.log  tls.crt  tls_crt_debug.log  tls.csr  tls.key  tls_key_debug.log
    [kschoener@fedora scenario_minikube_4]$ ls /etc/pki/ca-trust/source/anchors/
    [kschoener@fedora scenario_minikube_4]$ cp $SKUPPER_CA_DIR/tls.crt $SKUPPER_CA_DIR/20250210_my_ca_tls.crt
    [kschoener@fedora scenario_minikube_4]$ ls $SKUPPER_CA_DIR
    20250210_my_ca_tls.crt  csr_debug.log  tls.crt  tls_crt_debug.log  tls.csr  tls.key  tls_key_debug.log
    [kschoener@fedora scenario_minikube_4]$ sudo cp  $SKUPPER_CA_DIR/20250210_my_ca_tls.crt /etc/pki/ca-trust/source/anchors/
    [kschoener@fedora scenario_minikube_4]$ ls -alc /etc/pki/ca-trust/source/anchors/
    total 4
    drwxr-xr-x. 1 root root  44 Feb 10 12:00 .
    drwxr-xr-x. 1 root root 102 Dec 17 09:17 ..
    -rw-r--r--. 1 root root 985 Feb 10 12:00 20250210_my_ca_tls.crt
    [kschoener@fedora scenario_minikube_4]$ sudo update-ca-trust
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

### Example: verify skupper-site-server tls.crt
```
$ openssl verify $SKUPPER_SITE_SERVER_DIR/tls.crt
./test1/west/secrets/skupper-site-server/tls.crt: OK
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
