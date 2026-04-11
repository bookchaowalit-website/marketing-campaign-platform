#!/bin/bash
# Generate Python gRPC stubs from protobuf definitions

set -e

PROTO_DIR="../../packages/grpc-protos"
OUTPUT_DIR="./app/grpc_client"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "Generating Python gRPC stubs..."

# Generate Python gRPC stubs
python -m grpc_tools.protoc \
    -I="$PROTO_DIR" \
    --python_out="$OUTPUT_DIR" \
    --grpc_python_out="$OUTPUT_DIR" \
    --pyi_out="$OUTPUT_DIR" \
    "$PROTO_DIR/marketing.proto"

# Rename generated files to match imports
cd "$OUTPUT_DIR"
if [ -f "marketing_pb2.py" ]; then
    mv marketing_pb2.py marketing_pb2.pyi.bak 2>/dev/null || true
fi

echo "gRPC stubs generated successfully in $OUTPUT_DIR"
echo "Generated files:"
ls -la "$OUTPUT_DIR"/*.py 2>/dev/null || echo "No .py files found"
