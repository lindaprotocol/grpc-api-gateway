#!/bin/bash

set -e

PROJECT_ROOT=$(pwd)

# Check protoc version (3.19+ recommended for protoc-gen-go; 3.6.1 may fail)
PROTOC_VER=$(protoc --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
if [ -n "$PROTOC_VER" ]; then
    MAJOR=$(echo "$PROTOC_VER" | cut -d. -f1)
    MINOR=$(echo "$PROTOC_VER" | cut -d. -f2)
    if [ "$MAJOR" -lt 3 ] || { [ "$MAJOR" -eq 3 ] && [ "$MINOR" -lt 19 ]; }; then
        echo "WARNING: protoc $PROTOC_VER detected. protoc-gen-go and grpc-gateway recommend protoc 3.19+."
        echo "If generation fails, upgrade: https://github.com/protocolbuffers/protobuf/releases"
    fi
fi
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
    git clone --depth 1 https://github.com/googleapis/googleapis.git ${THIRD_PARTY_DIR}/googleapis
fi

# Download grpc-gateway for openapiv2 options
if [ ! -d "${THIRD_PARTY_DIR}/grpc-gateway" ]; then
    echo "Downloading grpc-gateway..."
    git clone --depth 1 https://github.com/grpc-ecosystem/grpc-gateway.git ${THIRD_PARTY_DIR}/grpc-gateway
fi

# Ensure required protoc plugins are installed (must be in PATH, typically $GOPATH/bin or $HOME/go/bin)
for plugin in protoc-gen-go protoc-gen-go-grpc protoc-gen-grpc-gateway protoc-gen-openapiv2; do
    if ! command -v $plugin &> /dev/null; then
        echo "Installing $plugin..."
        case $plugin in
            protoc-gen-go)          go install google.golang.org/protobuf/cmd/protoc-gen-go@latest ;;
            protoc-gen-go-grpc)     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest ;;
            protoc-gen-grpc-gateway) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest ;;
            protoc-gen-openapiv2)   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest ;;
        esac
    fi
done
# Ensure Go bin is in PATH for this session
export PATH="${PATH}:$(go env GOPATH)/bin"

# Define proto include paths (-I with no space between -I and path)
PROTO_INCLUDES="-I${PROTO_DIR} -I${THIRD_PARTY_DIR}/googleapis -I${THIRD_PARTY_DIR}/grpc-gateway"

# Define module mapping for well-known types
MAPPING="\
Mgoogle/protobuf/any.proto=google.golang.org/protobuf/types/known/anypb,\
Mgoogle/protobuf/duration.proto=google.golang.org/protobuf/types/known/durationpb,\
Mgoogle/protobuf/struct.proto=google.golang.org/protobuf/types/known/structpb,\
Mgoogle/protobuf/timestamp.proto=google.golang.org/protobuf/types/known/timestamppb,\
Mgoogle/protobuf/wrappers.proto=google.golang.org/protobuf/types/known/wrapperspb,\
Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations"

echo "Generating core protos..."
# Generate core protos
protoc ${PROTO_INCLUDES} \
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
protoc ${PROTO_INCLUDES} \
    --go_out=${OUT_DIR} \
    --go_opt=module=${MODULE_PATH} \
    --go_opt=${MAPPING} \
    --go-grpc_out=${OUT_DIR} \
    --go-grpc_opt=module=${MODULE_PATH} \
    --go-grpc_opt=${MAPPING} \
    --grpc-gateway_out=${OUT_DIR} \
    --grpc-gateway_opt=module=${MODULE_PATH} \
    --grpc-gateway_opt=${MAPPING} \
    --grpc-gateway_opt=logtostderr=true \
    --grpc-gateway_opt=generate_unbound_methods=true \
    --openapiv2_out=${PROJECT_ROOT}/api_docs \
    --openapiv2_opt=logtostderr=true \
    --openapiv2_opt=json_names_for_fields=true \
    --openapiv2_opt=allow_merge=true \
    --openapiv2_opt=merge_file_name=api \
    ${PROTO_DIR}/api.proto

if [ $? -ne 0 ]; then
    echo "Failed to generate api.proto"
    exit 1
fi

echo "Proto generation complete!"

# Generate swagger.json if needed
if command -v swagger &> /dev/null; then
    echo "Generating swagger.json..."
    swagger generate spec -o ${PROJECT_ROOT}/api_docs/swagger.json
fi

# Run go mod tidy to update dependencies
echo "Running go mod tidy to update dependencies..."
cd ${PROJECT_ROOT}
go mod tidy

if [ $? -ne 0 ]; then
    echo "Warning: go mod tidy encountered issues, but proto generation succeeded"
    exit 0
fi

echo "All tasks completed successfully!"