#!/usr/bin/env bash

source .envrc

mkdir -p .ssl

# check if certificate expired or is going to expiry in less than 2 days
openssl x509 -in .ssl/ca.pem -noout -checkend 172800 > /dev/null 2>&1
if [ $? -eq 0 ]; then
  echo "Existing certificates are valid"
else
  echo "Creating new certificates"

  cat << EOF > .ssl/req.cnf
  [req]
  req_extensions = v3_req
  distinguished_name = req_distinguished_name

  [req_distinguished_name]

  [ v3_req ]
  basicConstraints = CA:FALSE
  keyUsage = nonRepudiation, digitalSignature, keyEncipherment
  subjectAltName = @alt_names

  [alt_names]
  DNS.1 = *.kube.local
EOF

  openssl genrsa -out .ssl/ca-key.pem 2048
  openssl req -x509 -new -nodes -key .ssl/ca-key.pem -days "${KKA_CERT_VALIDITY_DAYS}" -out .ssl/ca.pem -subj "/CN=root-ca"

  openssl genrsa -out .ssl/key-ingress-gateway.pem 2048
  openssl req -new -key .ssl/key-ingress-gateway.pem -out .ssl/csr.pem -subj "/CN=kube-local" -config .ssl/req.cnf
  openssl x509 -req -in .ssl/csr.pem -CA .ssl/ca.pem -CAkey .ssl/ca-key.pem -CAcreateserial -out .ssl/cert-ingress-gateway.pem -days "${KKA_CERT_VALIDITY_DAYS}" -extensions v3_req -extfile .ssl/req.cnf
fi
