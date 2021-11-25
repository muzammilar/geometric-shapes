

# Variables
COMMIT := $(shell /usr/bin/git describe --always)
DEFAULT_VERSION := 0.0.1 # the default application version
BENCH_CPUS := 1          # number of cpus for benchmark testing
BENCH_ITERATIONS := 1000 # number of iterations for benchmark testing
GOMODULENAME := "github.com/muzammilar/geometric-shapes"

SERVER_VERSION := ${DEFAULT_VERSION}
CLIENT_VERSION := ${DEFAULT_VERSION}

# Applications
GEOMSERVER=geomserver
DATASERVER=dataserver
GOCLIENT=goclient
GODIR=go

.PHONY: all clean certs clean_certs protos clean_protos test lint gomodule gomodinit go

all: clean protos go certs ${DATASERVER} ${GEOMSERVER} ${GOCLIENT}

clean:
	-rm -f build/${DATASERVER}
	-rm -f build/${GEOMSERVER}
	-rm -f build/${GOCLIENT}

clean_protos:
	$(MAKE) clean -C protodata

protos: clean_protos
	$(MAKE) $@ -C protodata

${DATASERVER}:
	go build -ldflags "-X main.version=${SERVER_VERSION}" -o build/${DATASERVER} ./${GODIR}/cmd/dataserver/dataserver.go

${GEOMSERVER}:
	go build -ldflags "-X main.version=${SERVER_VERSION}" -o build/${GEOMSERVER} ./${GODIR}/cmd/geomserver/geomserver.go

${GOCLIENT}:
	go build -ldflags "-X main.version=${CLIENT_VERSION}" -o build/${GOCLIENT} ./${GODIR}/cmd/client/client.go

# run make all on the go directory
go:
	$(MAKE) -C ${GODIR}

# language agnostic commands
test lint:
	$(MAKE) $@ -C ${GODIR}

# go specific commands
gomodule gomodinit:
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
