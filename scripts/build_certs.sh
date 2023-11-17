#!/bin/bash

set -x -e -o pipefail

# init ca for certs
cockroach cert create-ca \
  --certs-dir . \
  --ca-key ca.key \
  --allow-ca-key-reuse

# init node certs
for i in {1..3}; do
  dir=nodes/node${i}
  mkdir -p ${dir}/

  cockroach cert create-node \
    0.0.0.0 \
    localhost \
    127.0.0.1 \
    cockroach-balancer \
    cockroach-node${i} \
    cockroach \
    --certs-dir . \
    --ca-key ca.key

  mv node.crt ${dir}
  mv node.key ${dir}
  chmod 700 ${dir}/node.key
  chmod 700 ${dir}/node.crt

done

# init client certs
cockroach cert create-client root --certs-dir . --ca-key=ca.key
chmod 700 ca.crt
chmod 700 ca.key
