.PHONY: all clean protos

# For now, we use relative paths
SRC_DIR := .
PROJECT_ROOT:= ../../
GO_DST_DIR := $(PROJECT_ROOT)/go/protos/serviceinfo


all: clean protos

# All paths are relative to source
protos:
	mkdir -p $(GO_DST_DIR)
	protoc --proto_path=$(SRC_DIR) \
			--go_opt=paths=source_relative \
			--go_out=$(GO_DST_DIR) \
		$(SRC_DIR)/serviceinfo.proto

clean:
	-rm -f $(GO_DST_DIR)/*.pb.go
	-rm -f $(SRC_DIR)/*.pb.go
