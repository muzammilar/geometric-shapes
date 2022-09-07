
# includes
include ./Makefile.variable

# Applications
GODIR=go

.PHONY: all clean certs clean_certs protos clean_protos test lint gomodule gomodinit go

all: clean protos go

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
gomodule gomodinit ${GO_DATASERVER} ${GO_GEOMSERVER} ${GO_CLIENT}:
	$(MAKE) $@ -C ${GODIR}

clean_certs:
	-rm -f certs/*.pem
	-rm -f certs/*.key
	-rm -f certs/*.crt
	-rm -f certs/*.req

# Create a self-signed Root CA and use the Root CA to sign the server cert
certs: clean_certs
	mkdir -p certs
# Generation of self-signed(x509) Root Certificate (PEM-encodings .pem|.crt)
	openssl req -x509 -newkey rsa:4096 -sha256 -nodes \
	-keyout certs/root.ca.key.pem -out certs/root.ca.crt.pem -days 365 \
	-subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=root.ca.cert.com" \
	-addext "subjectAltName=DNS:localhost"
# Check the certificate
	echo "CA's self-signed certificate"
	openssl x509 -in certs/root.ca.crt.pem -noout -text
# Generation of server certificate signing request and key file (PEM-encodings .pem). It's not an x509 request
	openssl req -newkey rsa:4096 -sha256 -nodes \
	-keyout certs/server.grpc.key.pem -out certs/server.grpc.req.pem \
	-subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=*.geometry"
# Sign the server's CSR (.pem) using the Root CA and generate the cert
	openssl x509 -req -in certs/server.grpc.req.pem -days 60 -CA certs/root.ca.crt.pem \
	-CAkey certs/root.ca.key.pem -CAcreateserial -out certs/server.grpc.crt.pem \
	-extfile certs/server.grpc.ext
# Check the server's certificate
	echo "Server's signed certificate"
	openssl x509 -in certs/server.grpc.crt.pem -noout -text
