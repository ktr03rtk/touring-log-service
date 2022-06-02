#!/bin/bash

cd "$(dirname "$0")" || exit 1

IP="example.com"

SUBJECT_CA="/C=JP/OU=ca/CN=$IP"
SUBJECT_SERVER="/C=JP/OU=server/CN=$IP"
SUBJECT_CLIENT="/C=JP/OU=client/CN=$IP"

# Certificate Authority
# Generate a certificate authority certificate and key.
openssl req -newkey rsa:2048 -x509 -sha256 -nodes -subj "$SUBJECT_CA" -days 3650 -extensions v3_ca -keyout ca.key -out ca.crt

# Server
# Generate a server key.
openssl genrsa -out server.key 2048

# Generate a certificate signing request to send to the CA.
openssl req -subj "$SUBJECT_SERVER" -out server.csr -key server.key -new

# Send the CSR to the CA, or sign it with your CA key:
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650 -extfile subjectnames.txt

# Client
# Generate a client key.
openssl genrsa -out client.key 2048

# Generate a certificate signing request to send to the CA.
openssl req -subj "$SUBJECT_CLIENT" -out client.csr -key client.key -new

# Send the CSR to the CA, or sign it with your CA key:
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 3650 -extfile subjectnames.txt

mkdir -p ../cert
mv *.srl *.csr *.crt *.key ../cert
