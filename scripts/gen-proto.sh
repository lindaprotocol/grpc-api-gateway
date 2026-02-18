#!/bin/bash

PROJECT_ROOT=$(pwd)
PROTO_DIR="${PROJECT_ROOT}/proto"
OUT_DIR="${PROJECT_ROOT}/pkg/api"

# Create output directory
mkdir -p ${OUT_DIR}

# Download googleapis if not present
if [ ! -d "${PROJECT_ROOT}/third_party/googleapis" ]; then
    git clone https://github.com/googleapis/googleapis.git ${PROJECT_ROOT}/third_party/googleapis
fi

# Generate protos
protoc -I${PROTO_DIR} \
    -I${PROJECT_ROOT}/third_party/googleapis \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${OUT_DIR} \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=${OUT_DIR} \
    --grpc-gateway_opt=paths=source_relative \
    ${PROTO_DIR}/core/*.proto \
    ${PROTO_DIR}/api.proto

echo "Proto generation complete!"
