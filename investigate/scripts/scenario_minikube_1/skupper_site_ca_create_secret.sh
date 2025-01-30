#!/bin/bash

kubectl create secret tls skupper-site-ca \
  --cert=tls.crt \
  --key=tls.key \
  -n west

