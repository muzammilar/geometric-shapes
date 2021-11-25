
# For now, we use relative paths
SRC_DIR := .
IMP_DIR := ..
PROJECT_ROOT:= ../../
GO_DST_DIR := $(PROJECT_ROOT)/go/protos/shapestore

.PHONY: all clean protos

all: clean protos

# Include both current directory and parent directory
# All paths are relative to source
protos:
	mkdir -p $(GO_DST_DIR)
	protoc --proto_path=$(SRC_DIR) \
			--proto_path=$(IMP_DIR) \
			--go_opt=paths=source_relative \
			--go_out=$(GO_DST_DIR) \
			--go-grpc_opt=paths=source_relative	\
			--go-grpc_out=$(GO_DST_DIR) \
		$(SRC_DIR)/shapestore.proto

clean:
	-rm -f $(GO_DST_DIR)/*.pb.go
	-rm -f $(SRC_DIR)/*.pb.go
