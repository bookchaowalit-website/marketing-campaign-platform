.PHONY: generate clean

PROTO_DIR=.
GO_OUTPUT=../../api-go/pkg/grpc
PYTHON_OUTPUT=../../api-python/app/grpc

generate:
	@echo "Generating Go gRPC code..."
	mkdir -p $(GO_OUTPUT)
	protoc --go_out=$(GO_OUTPUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUTPUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/marketing.proto
	
	@echo "Generating Python gRPC code..."
	mkdir -p $(PYTHON_OUTPUT)
	python3 -m grpc_tools.protoc -I$(PROTO_DIR) \
		--python_out=$(PYTHON_OUTPUT) \
		--grpc_python_out=$(PYTHON_OUTPUT) \
		$(PROTO_DIR)/marketing.proto
	
	@echo "gRPC code generation complete!"

clean:
	rm -rf $(GO_OUTPUT)/*.pb.go
	rm -rf $(PYTHON_OUTPUT)/*_pb2.py
	rm -rf $(PYTHON_OUTPUT)/*_pb2_grpc.py
