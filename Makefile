#!/bin/bash

set -e

PROJECT_ROOT=$(pwd)
PROTO_DIR="${PROJECT_ROOT}/proto"
OUT_DIR="${PROJECT_ROOT}/pkg/lindapb"
THIRD_PARTY_DIR="${PROJECT_ROOT}/third_party"
MODULE_PATH="github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"

echo "Generating protos from: ${PROTO_DIR}"

# Create output directories
mkdir -p ${OUT_DIR}
mkdir -p ${OUT_DIR}/core

# Download googleapis if not present
if [ ! -d "${THIRD_PARTY_DIR}/googleapis" ]; then
    echo "Downloading googleapis..."
    git clone https://github.com/googleapis/googleapis.git ${THIRD_PARTY_DIR}/googleapis
fi

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "protoc not found. Please install protobuf compiler:"
    echo "  apt-get install -y protobuf-compiler"
    exit 1
fi

# Get the path to the locally installed plugins
GOBIN=$(go env GOPATH)/bin
export PATH=$PATH:$GOBIN

# Define proto paths
PROTO_PATH="\
-I${PROTO_DIR} \
-I${THIRD_PARTY_DIR}/googleapis"

echo "Generating core protos..."
# Generate core protos
protoc ${PROTO_PATH} \
    --go_out=${OUT_DIR} \
    --go_opt=module=${MODULE_PATH} \
    --go-grpc_out=${OUT_DIR} \
    --go-grpc_opt=module=${MODULE_PATH} \
    ${PROTO_DIR}/core/*.proto

if [ $? -ne 0 ]; then
    echo "Failed to generate core protos"
    exit 1
fi

echo "Generating api.proto..."
# Generate main api.proto with gateway
protoc ${PROTO_PATH} \
    --go_out=${OUT_DIR} \
    --go_opt=module=${MODULE_PATH} \
    --go-grpc_out=${OUT_DIR} \
    --go-grpc_opt=module=${MODULE_PATH} \
    --grpc-gateway_out=${OUT_DIR} \
    --grpc-gateway_opt=module=${MODULE_PATH} \
    --grpc-gateway_opt=logtostderr=true \
    --grpc-gateway_opt=generate_unbound_methods=true \
    ${PROTO_DIR}/api.proto

if [ $? -ne 0 ]; then
    echo "Failed to generate api.proto"
    exit 1
fi

echo "Proto generation complete!"

# Run go mod tidy to update dependencies
echo "Running go mod tidy to update dependencies..."
cd ${PROJECT_ROOT}
go mod tidy

if [ $? -ne 0 ]; then
    echo "Warning: go mod tidy encountered issues, but proto generation succeeded"
    exit 0
fi

echo "All tasks completed successfully!"