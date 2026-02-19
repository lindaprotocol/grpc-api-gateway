#!/bin/bash

PROJECT_ROOT=$(pwd)
PROTO_DIR="${PROJECT_ROOT}/proto"
OUT_DIR="${PROJECT_ROOT}/pkg/api"

# Create output directories
mkdir -p ${OUT_DIR}/protocol
mkdir -p ${OUT_DIR}/protocol/core

# Download googleapis if not present
if [ ! -d "${PROJECT_ROOT}/third_party/googleapis" ]; then
    git clone https://github.com/googleapis/googleapis.git ${PROJECT_ROOT}/third_party/googleapis
fi

# Define module mapping to avoid import cycles
MAPPING="Mgoogle/protobuf/any.proto=google.golang.org/protobuf/types/known/anypb,\
Mgoogle/protobuf/duration.proto=google.golang.org/protobuf/types/known/durationpb,\
Mgoogle/protobuf/struct.proto=google.golang.org/protobuf/types/known/structpb,\
Mgoogle/protobuf/timestamp.proto=google.golang.org/protobuf/types/known/timestamppb,\
Mgoogle/protobuf/wrappers.proto=google.golang.org/protobuf/types/known/wrapperspb,\
Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations"

# Generate core protos first
protoc -I${PROTO_DIR} \
    -I${PROJECT_ROOT}/third_party/googleapis \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${OUT_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/core/*.proto

# Generate main api.proto with proper mappings
protoc -I${PROTO_DIR} \
    -I${PROJECT_ROOT}/third_party/googleapis \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    --go_opt=${MAPPING} \
    --go-grpc_out=${OUT_DIR} \
    --go-grpc_opt=paths=source_relative \
    --go-grpc_opt=${MAPPING} \
    --grpc-gateway_out=${OUT_DIR} \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=${MAPPING} \
    ${PROTO_DIR}/api.proto

echo "Proto generation complete!"