#!/bin/bash
kubectl -n west patch deployment skupper-router --type='json' -p='[
  {
    "op": "add",
    "path": "/spec/template/spec/volumes/-",
    "value": {
      "name": "ca-bundle",
      "secret": {
        "secretName": "ca-bundle",
        "defaultMode": 420
      }
    }
  },
  {
    "op": "add",
    "path": "/spec/template/spec/containers/0/volumeMounts/-",
    "value": {
      "mountPath": "/etc/skupper-router-certs/ca-bundle/",
      "name": "ca-bundle"
    }
  },
  {
    "op": "add",
    "path": "/spec/template/spec/containers/0/env/-",
    "value": {
      "name": "SSL_CERT_FILE",
      "value": "/etc/skupper-router-certs/ca-bundle/ca_bundle.pem"
    }
  }
]'

