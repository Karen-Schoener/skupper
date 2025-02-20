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
   export SKUPPER_CA_INTERMEDIATE_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca-intermediate"
   export SKUPPER_SITE_SERVER_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-server"
   export SKUPPER_LINK_NAME="link1"
   export SKUPPER_LINK_DIR="${TEST_OUTPUT_DIR}/east/secrets/${SKUPPER_LINK_NAME}"

   ./user_supplied_ca_generate_artifacts.sh
   ./user_supplied_ca_intermediate_generate_artifacts.sh

   # ./skupper_site_server_generate_artifacts.sh 192.168.49.240
     ./skupper_site_server_generate_artifacts.sh mytest-skupper-router-west.local # defaults to signing with intermediate CA

     ./skupper_site_server_generate_artifacts.sh mytest-skupper-router-west.local intermediate
     ./skupper_site_server_generate_artifacts.sh mytest-skupper-router-west.local root



   ./link_generate_artifacts.sh # defaults to signing with intermediate CA
   ./link_generate_artifacts.sh intermediate
   ./link_generate_artifacts.sh root
```

## Step: debug commands to verify that intermediate ca was created

```
$ openssl x509 -in ./test1/west/secrets/skupper-site-ca-intermediate/tls.crt -text -noout

What to to look For:
    Issuer:     The issuer should be the Root CA.

    Subject:    The subject should indicate that it is the Intermediate CA.

    Extensions: basicConstraints extension should show CA:TRUE and 
                a pathlen constraint indicating it’s an intermediate certificate.

    Extensions notes: CA:TRUE: Indicates that the certificate is a Certificate Authority (CA) certificate.

                      pathlen:0: Specifies that the certificate is an intermediate certificate and
                      it cannot issue certificates to other intermediate CAs
                      (it can only issue end-entity certificates).
```

```
$ openssl x509 -in ./test1/west/secrets/skupper-site-ca-intermediate/tls.crt -text -noout
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            16:b3:09:c4:f9:b7:94:da:11:23:a2:d3:6b:11:65:dd:6f:a9:3d:c7
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: CN = my-cert
        Validity
            Not Before: Feb 19 15:05:39 2025 GMT
            Not After : Feb 19 15:05:39 2026 GMT
        Subject: CN = Intermediate-CA
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:dc:7e:4b:f3:df:4f:36:01:7b:02:5e:1f:87:fc:
                    71:86:13:f3:23:1c:87:d6:92:3d:40:19:f3:e4:22:
                    58:73:e5:a7:a3:d0:5c:0c:d1:a9:fd:81:b3:5e:01:
                    c6:15:3b:fa:fc:51:21:59:5d:fc:9a:9d:e2:f0:fe:
                    0b:99:8d:0d:a5:9d:d4:f8:fa:b3:ad:df:4d:07:90:
                    e9:ab:c7:15:f9:1c:17:a7:88:a1:67:c7:23:40:79:
                    19:d2:82:0a:9b:23:23:9e:d0:32:1b:d0:ec:bc:19:
                    62:ac:9d:9a:25:7c:38:80:e2:eb:a4:94:d0:77:32:
                    3c:db:21:2e:2e:2f:db:ac:4a:2a:b8:07:0e:a7:07:
                    fb:a9:94:75:01:47:58:a5:45:db:95:30:b8:97:d7:
                    43:45:e1:ae:51:06:78:10:33:f0:b5:28:94:70:15:
                    74:5d:bd:da:78:93:9c:23:8f:02:9e:12:f0:0c:b1:
                    6b:d6:0d:cf:c8:ed:b7:23:e9:8f:3c:61:d1:ea:7e:
                    95:52:8f:78:e2:76:da:00:91:88:18:70:42:52:97:
                    be:3f:ae:9e:18:69:46:35:3b:7f:12:17:57:c8:f1:
                    49:af:44:f4:21:b7:8e:f8:c7:a7:07:12:bb:44:a3:
                    88:11:2e:43:99:86:46:72:98:ee:a7:2b:eb:be:e2:
                    47:df
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints:
                CA:TRUE, pathlen:0
            X509v3 Subject Key Identifier:
                24:96:E1:71:DD:82:30:C3:27:08:A3:D7:94:B7:04:C6:56:0E:0E:42
            X509v3 Authority Key Identifier:
                DirName:/CN=my-cert
                serial:05:CE:3E:B2:F6:9F:1D:60:7E:FF:44:0F:89:24:1B:C9:F1:E8:84:FF
    Signature Algorithm: sha256WithRSAEncryption
    Signature Value:
        a2:80:dd:25:c9:7d:0a:f2:64:be:3e:fa:bb:52:69:4c:3e:81:
        d4:6e:75:9a:4e:21:90:4c:5c:2b:c1:9b:42:48:a5:c2:07:dd:
        8d:db:a9:0c:10:d9:61:f4:8a:b6:55:95:d1:2e:03:09:ce:ca:
        f4:b9:0d:79:5c:cf:a6:02:8e:ce:05:bf:46:db:5b:95:4b:71:
        3d:95:07:7c:2c:66:05:ed:3d:31:7f:8c:23:d9:8d:60:ac:37:
        36:39:fd:cf:77:cb:43:69:93:cd:3e:3c:56:63:e9:97:e9:9f:
        e0:c7:7e:0c:57:7c:f7:49:d6:0a:04:a1:f5:0c:1f:5e:2d:79:
        a3:3f:ca:e1:3f:81:b1:65:1b:80:a9:98:33:1f:51:85:13:9b:
        c2:1e:bc:fd:07:2e:eb:f9:19:5a:08:bc:29:88:ca:70:c0:fb:
        5d:78:47:4a:86:d3:80:c8:cc:d4:c6:7d:75:63:1e:1e:12:80:
        38:e7:d7:19:aa:0d:0a:ff:a6:b3:ef:78:15:d3:97:a5:6a:a9:
        ab:b1:4e:1e:da:e8:03:46:00:a3:22:29:bb:84:0b:d5:35:1f:
        5e:19:56:43:8b:fe:fe:da:4f:ac:1a:31:1c:da:a8:0f:5a:44:
        77:38:fd:1c:a2:4b:bd:62:93:39:f0:c9:f7:ec:23:45:2e:72:
        bd:6a:93:29

```

## Step: debug commands to verify that skupper-site-server tls.crt was signed by intermediate ca
```
$ openssl x509 -in $SKUPPER_SITE_SERVER_DIR/tls.crt -text -noout

What to to look For:
    Issuer:                   The issuer should be the intermediate CA.

    Authority Key Identifier: This field should match the key identifier of the Intermediate CA.

    Signature Algorithm:      This should show that the certificate was signed using the algorithm specified by the Intermediate CA.
```

```
Notes on checking: Authority Key Identifier:

$ openssl x509 -in $SKUPPER_CA_INTERMEDIATE_DIR/tls.crt -text -noout

In the intermediate CA:
    Issuer: This confirms that the Intermediate CA certificate was issued by "my-cert".
    Subject: The subject of the Intermediate CA certificate is "Intermediate-CA".
    Subject Key Identifier: This value is 24:96:E1:71:DD:82:30:C3:27:08:A3:D7:94:B7:04:C6:56:0E:0E:42.
    Authority Key Identifier: This confirms that the Intermediate CA certificate was issued by "my-cert".

$ openssl x509 -in $SKUPPER_SITE_SERVER_DIR/tls.crt -text -noout              

In the resulting tls.crt:
    Issuer: The server certificate should list "CN = Intermediate-CA" as the issuer.
    Authority Key Identifier: The server certificate should have the Authority Key Identifier 
    that which matches the Subject Key Identifier of the Intermediate CA.
```

```
$ openssl x509 -in $SKUPPER_SITE_SERVER_DIR/tls.crt -text -noout
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            6b:5d:77:45:74:c7:32:50:65:65:2a:17:5f:29:e9:ab:ac:a7:c7:48
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: CN = Intermediate-CA
        Validity
            Not Before: Feb 19 15:43:39 2025 GMT
            Not After : Feb 19 15:43:39 2026 GMT
        Subject: CN = mytest-skupper-router-west.local
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:b6:17:5e:6c:2f:77:95:9b:11:f2:14:c9:bf:50:
                    55:9d:bb:ac:d1:23:b3:e9:83:c0:80:10:54:df:de:
                    a5:7d:98:cb:6b:28:24:cc:7c:76:57:92:51:c7:b5:
                    e9:b2:24:cc:ec:d0:67:37:78:48:94:f5:5c:4b:85:
                    0d:be:73:44:23:12:fe:10:8c:e7:ea:30:38:af:ef:
                    e3:4d:2c:b0:82:db:d4:f2:b0:0f:bf:9d:07:fb:73:
                    3a:16:ef:07:fc:0d:33:18:54:fa:9b:dc:7b:94:67:
                    22:ca:40:43:5b:4b:73:17:b1:14:a9:8e:96:94:7d:
                    03:ba:fe:94:40:99:7f:1f:b9:5a:58:81:7a:be:01:
                    e1:62:5c:77:f1:12:78:f4:d6:fa:74:a0:06:7c:2d:
                    91:5d:e0:56:c1:e8:eb:f7:52:56:58:bc:86:b7:01:
                    78:6e:0f:0e:77:a4:0a:7e:78:b5:8a:60:16:94:49:
                    0f:6e:b2:cf:11:58:33:9f:ad:4a:be:bd:a6:df:95:
                    3b:87:7a:ef:b2:22:48:45:f6:2d:ca:0f:72:af:df:
                    9e:9e:18:a2:18:c8:c4:0b:6f:7d:0e:3b:81:87:dd:
                    46:1b:41:95:15:a5:38:5e:4f:98:22:ac:ff:22:90:
                    13:1b:09:78:85:ae:c0:a1:3f:f2:73:d7:0a:df:8e:
                    4c:e3
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage:
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Basic Constraints: critical
                CA:FALSE
            X509v3 Authority Key Identifier:
                24:96:E1:71:DD:82:30:C3:27:08:A3:D7:94:B7:04:C6:56:0E:0E:42
            X509v3 Subject Alternative Name:
                DNS:skupper-router.west.svc.cluster.local, DNS:mytest-skupper-router-west.local
            X509v3 Subject Key Identifier:
                AE:5D:71:06:DC:61:64:12:5A:5C:33:34:E2:F5:83:FB:98:9F:90:96
    Signature Algorithm: sha256WithRSAEncryption
    Signature Value:
        b5:db:e4:6a:aa:fa:dc:ea:8d:f4:69:b3:d0:ae:0f:3e:c5:e5:
        36:9e:e6:d6:5d:ee:ac:e1:54:55:f2:d9:45:00:80:65:ca:59:
        1b:e8:0e:a6:90:e8:2a:1c:6d:cf:3d:86:cd:d4:ff:e1:66:96:
        7d:a2:b2:ef:bf:12:51:90:18:a9:ed:7a:ae:6d:97:e6:70:26:
        42:8c:a6:74:2f:5f:9a:34:1a:1c:72:2f:05:2a:60:ab:d8:aa:
        25:3c:36:e8:e0:6a:e6:1d:5f:bc:a7:15:96:81:c7:51:57:12:
        f6:5d:f9:34:1e:c2:7f:0d:fd:38:3b:c8:03:db:a1:0d:49:bf:
        51:dd:c4:b5:4e:a5:df:81:ee:b3:fe:12:1e:13:b4:e2:df:c9:
        49:fb:0b:eb:2e:b1:a7:e5:63:f0:b8:39:da:fb:8b:46:4c:1f:
        74:ac:82:5d:61:e5:d7:02:27:8a:c5:96:88:d9:34:11:3a:b0:
        89:07:9d:ac:a7:e2:7a:3a:4c:8f:74:6a:16:9d:c4:02:ee:49:
        c3:e2:43:6c:7a:2a:be:a2:62:70:08:73:cb:df:c1:fe:4b:04:
        5b:f2:b0:28:1c:fe:c6:4d:4e:63:45:fe:a5:a6:38:d2:02:ea:
        6c:41:5c:f9:ba:a2:a3:70:47:a7:1f:49:cf:5b:11:07:77:bb:
        82:37:e9:b4
```

## Step: debug commands to verify that link credentials were signed by intermediate ca
```
openssl x509 -in test1/east/secrets/link1/tls.crt -text -noout

What to look for:
    Issuer: The issuer of the certificate should be the Intermediate CA.
    Authority Key Identifier: <matches intermediate CA: Subject Key Identifier>
```

## Step: generate credentials: set 2
```
   export TEST_OUTPUT_DIR="./test2"
   export SKUPPER_CA_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-ca"
   export SKUPPER_SITE_SERVER_DIR="${TEST_OUTPUT_DIR}/west/secrets/skupper-site-server"
   export SKUPPER_LINK_NAME="link1"
   export SKUPPER_LINK_DIR="${TEST_OUTPUT_DIR}/east/secrets/${SKUPPER_LINK_NAME}"

   ./user_supplied_ca_generate_artifacts.sh
   ./user_supplied_ca_intermediate_generate_artifacts.sh

   # ./skupper_site_server_generate_artifacts.sh 192.168.49.240
     ./skupper_site_server_generate_artifacts.sh mytest-skupper-router-west.local

   ./link_generate_artifacts.sh
```

## Step: create ca-bundle files
```
mkdir ca_bundle

# NOTE: the intermediate CA certificate must be listed before the root CA certificate in the ca-bundle.

cat ./test1/west/secrets/skupper-site-ca/tls.crt > ca_bundle/ca_bundle_test1.pem
cat ./test1/west/secrets/skupper-site-ca-intermediate/tls.crt ./test1/west/secrets/skupper-site-ca/tls.crt > ca_bundle/ca_bundle_test1_with_intermediate.pem
# cat ./test2/west/secrets/skupper-site-ca/tls.crt > ca_bundle/ca_bundle_test2.pem
# cat ./test1/west/secrets/skupper-site-ca/tls.crt ./test2/west/secrets/skupper-site-ca/tls.crt > ca_bundle/ca_bundle_test1_test2.pem
```

## Step: verify the link certificate against the CA bundle
```
$ openssl verify -CAfile ca_bundle/ca_bundle_test1.pem test1/east/secrets/link1/tls.crt
test1/east/secrets/link1/tls.crt: OK
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

```
skupper -n west status -v
skupper -n east status -v
```

## Step: patch west/secret/skupper-site-server with set-1 user-supplied certs
OBSERVATION Note do not delete secret skupper-site-server.  Seems to cause skupper-router error logs if there is no ca in the skupper-site-server secret
```
# ./skupper_site_server_create_secret_NO_CA.sh
./skupper_site_server_create_secret_NO_CA_TLS_CRT_INCLUDES_INTERMEDIATE_CA.sh
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

```
kubectl -n east exec -it $(kubectl get -n east pods -l app.kubernetes.io/name=skupper-router -o name) -- /bin/bash
```

## Step: in east, configure set-1 connection-token
OBSERVATION Note create connect-token with ca.crt specified.  At the moment, skupper will not configure east/configmap/skupper-ineternal/skrouter.json if ca.crt not set in connection token.
```
# ./link_create_secret_NO_CA.sh
# ./link_create_secret.sh
./link_create_secret_TLS_CRT_INCLUDES_INTERMEDIATE_CA.sh
```

## Step: verify sites are linked
```
# Verify that sites are connected
skupper -n east link status
skupper -n east status -v
skupper -n west status -v
```
