#!/usr/bin/env bash

set -xe

openssl genrsa -out ca.key 2048

openssl req -new -x509 -days 365 -key ca.key \
  -subj "/C=US/O=apex/OU=dev/CN=ca" \
  -out ca.cert

openssl req -newkey rsa:2048 -nodes -keyout server.key \
  -subj "/C=US/O=apex/OU=dev/CN=localhost" \
  -out server.csr

openssl x509 -req \
  -extfile <(printf "subjectAltName=DNS:localhost") \
  -days 365 \
  -in server.csr \
  -CA ca.cert \
  -CAkey ca.key \
  -CAcreateserial \
  -out server.crt

rm -f bundle.crt
cat server.crt ca.cert > bundle.crt
