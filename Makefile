
# includes
include ./Makefile.variable

# Applications
GODIR=go

.PHONY: all clean certs clean_certs protos clean_protos test lint gomodule gomodinit go

all: clean protos go ${GO_DATASERVER} ${GO_GEOMSERVER} ${GO_GOCLIENT}

clean:
	-rm -f ${BUILD_DIR}/${GO_DATASERVER}
	-rm -f ${BUILD_DIR}/${GO_GEOMSERVER}
	-rm -f ${BUILD_DIR}/${GO_CLIENT}

clean_protos:
	$(MAKE) clean -C protodata

protos: clean_protos
	$(MAKE) $@ -C protodata

# run make all on the go directory
go:
	$(MAKE) -C ${GODIR}

# language agnostic commands
test lint:
	$(MAKE) $@ -C ${GODIR}

# go specific commands
gomodule gomodinit ${GO_DATASERVER} ${GO_GEOMSERVER} ${GO_GOCLIENT}:
	$(MAKE) $@ -C ${GODIR}

clean_certs:
	-rm -f certs/*

# source: https://github.com/denji/golang-tls
certs: clean_certs
	mkdir -p certs
# Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
	openssl ecparam -genkey -name secp384r1 -out certs/server.grpc.key
# Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
	openssl req -new -x509 -sha256 \
	-key certs/server.grpc.key -out certs/server.grpc.crt -days 3650 \
	-subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=server.grpc"
