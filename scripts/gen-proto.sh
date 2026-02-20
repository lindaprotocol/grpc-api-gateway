#!/bin/bash

PROJECT_ROOT=$(pwd)
PROTO_DIR="${PROJECT_ROOT}/proto"
OUT_DIR="${PROJECT_ROOT}/pkg/api"
MODULE_PATH="github.com/lindaprotocol/grpc-api-gateway/pkg/api"

echo "Generating protos from: ${PROTO_DIR}"

# Create output directories
mkdir -p ${OUT_DIR}/protocol
mkdir -p ${OUT_DIR}/protocol/core

# Download googleapis if not present
if [ ! -d "${PROJECT_ROOT}/third_party/googleapis" ]; then
    echo "Downloading googleapis..."
    git clone https://github.com/googleapis/googleapis.git ${PROJECT_ROOT}/third_party/googleapis
fi

# Define module mapping for well-known types
MAPPING="Mgoogle/protobuf/any.proto=google.golang.org/protobuf/types/known/anypb,\
Mgoogle/protobuf/duration.proto=google.golang.org/protobuf/types/known/durationpb,\
Mgoogle/protobuf/struct.proto=google.golang.org/protobuf/types/known/structpb,\
Mgoogle/protobuf/timestamp.proto=google.golang.org/protobuf/types/known/timestamppb,\
Mgoogle/protobuf/wrappers.proto=google.golang.org/protobuf/types/known/wrapperspb,\
Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations"

echo "Generating core protos..."
# Generate core protos first (no mappings needed for core)
protoc -I${PROTO_DIR} \
    -I${PROJECT_ROOT}/third_party/googleapis \
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
# Generate main api.proto with mappings
protoc -I${PROTO_DIR} \
    -I${PROJECT_ROOT}/third_party/googleapis \
    --go_out=${OUT_DIR} \
    --go_opt=module=${MODULE_PATH} \
    --go_opt=${MAPPING} \
    --go-grpc_out=${OUT_DIR} \
    --go-grpc_opt=module=${MODULE_PATH} \
    --go-grpc_opt=${MAPPING} \
    --grpc-gateway_out=${OUT_DIR} \
    --grpc-gateway_opt=module=${MODULE_PATH} \
    --grpc-gateway_opt=${MAPPING} \
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
    exit 0  # Exit with success since proto generation worked
fi

echo "All tasks completed successfully!"