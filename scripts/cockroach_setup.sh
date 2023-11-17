#!/bin/sh

chmod 700 insecure/cockroach/certs/ca.crt
chmod 700 insecure/cockroach/certs/ca.key

chmod 700 insecure/cockroach/certs/nodes/node1/node.key
chmod 700 insecure/cockroach/certs/nodes/node1/node.crt
chmod 700 insecure/cockroach/certs/nodes/node2/node.key
chmod 700 insecure/cockroach/certs/nodes/node2/node.crt
chmod 700 insecure/cockroach/certs/nodes/node3/node.key
chmod 700 insecure/cockroach/certs/nodes/node3/node.crt

cockroach init --certs-dir=${CERTS_DIR} --host=${HOST}

# create admin
create_admin_user_query="create user if not exists root_user with password 'root'"
cockroach sql --certs-dir=${CERTS_DIR} --host=${HOST} --execute="${create_admin_user_query}"

for var in "$@"; do
  #  create database
  create_db_query="create database if not exists ${var}_db"
  cockroach sql --certs-dir=${CERTS_DIR} --host=${HOST} --execute="${create_db_query}"
  # create user
  create_user_query="create user if not exists ${var}_svc with password '${var}_pass'"
  cockroach sql --certs-dir=${CERTS_DIR} --host=${HOST} --execute="${create_user_query}"
  # grant access
  grant_access_query="grant all on database ${var}_db to ${var}_svc"
  cockroach sql --certs-dir=${CERTS_DIR} --host=${HOST} --execute="${grant_access_query}"
  # grant access for admin user
  grant_admin_access_query="grant all on database ${var}_db to root_user"
  cockroach sql --certs-dir=${CERTS_DIR} --host=${HOST} --execute="${grant_admin_access_query}"
done

echo "Setup finished"
