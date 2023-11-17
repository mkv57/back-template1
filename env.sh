#!/bin/sh

# Database settings.
export BACK_TEMPLATE_COCKROACH_NODE1_ADDR_TCP_PORT="26257"
export BACK_TEMPLATE_COCKROACH_NODE1_ADDR_HTTP_PORT="8080"
export BACK_TEMPLATE_COCKROACH_NODE2_ADDR_TCP_PORT="26258"
export BACK_TEMPLATE_COCKROACH_NODE2_ADDR_HTTP_PORT="8081"
export BACK_TEMPLATE_COCKROACH_NODE3_ADDR_TCP_PORT="26259"
export BACK_TEMPLATE_COCKROACH_NODE3_ADDR_HTTP_PORT="8082"
export COCKROACH_TLS_CERT="./certs/cockroach/ca.crt"
export COCKROACH_TLS_KEY="./certs/cockroach/ca.key"
export COCKROACH_CLIENT_TLS_CERT="./certs/cockroach/client.root.crt"
export COCKROACH_CLIENT_TLS_KEY="./certs/cockroach/client.root.key"
export COCKROACH_NODE1_CRT="./certs/cockroach/nodes/node1/node.crt"
export COCKROACH_NODE1_KEY="./certs/cockroach/nodes/node1/node.key"
export COCKROACH_NODE2_CRT="./certs/cockroach/nodes/node2/node.crt"
export COCKROACH_NODE2_KEY="./certs/cockroach/nodes/node2/node.key"
export COCKROACH_NODE3_CRT="./certs/cockroach/nodes/node3/node.crt"
export COCKROACH_NODE3_KEY="./certs/cockroach/nodes/node3/node.key"

# Queue settings.
export BACK_TEMPLATE_NATS_NODE1_ADDR_TCP_PORT="4222"
export BACK_TEMPLATE_NATS_NODE1_ADDR_HTTP_PORT="8222"
export BACK_TEMPLATE_NATS_NODE2_ADDR_TCP_PORT="4223"
export BACK_TEMPLATE_NATS_NODE2_ADDR_HTTP_PORT="8223"
export BACK_TEMPLATE_NATS_NODE3_ADDR_TCP_PORT="4224"
export BACK_TEMPLATE_NATS_NODE3_ADDR_HTTP_PORT="8224"

# File storage settings.
export BACK_TEMPLATE_MINIO_NODE1_ADDR_HTTP_PORT="9100"
export BACK_TEMPLATE_MINIO_NODE1_ADDR_ADMIN_PORT="9101"

# user svc
export USER_SVC_CONFIG="./cmd/user/config.yml"
export USER_SVC_GRPC_PORT=10000
export USER_SVC_METRIC_PORT=10001
export USER_SVC_GRPC_GW_PORT=10002
export USER_SVC_GRPC_FILES_PORT=10003

# session svc
export SESSION_SVC_CONFIG="./cmd/session/config.yml"
export SESSION_SVC_GRPC_PORT=10100
export SESSION_SVC_METRIC_PORT=10101
