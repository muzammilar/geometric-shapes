

# Variables
COMMIT := $(shell /usr/bin/git describe --always)
PWD := $(shell pwd)
DEFAULT_VERSION := 0.0.1 # the default application version
BENCH_CPUS := 1          # number of cpus for benchmark testing
BENCH_ITERATIONS := 1000 # number of iterations for benchmark testing
 # For now, we skip the "go" (for simplicty). If an external program imports it, they'll provide full path
GOMODULENAME := "github.com/muzammilar/geometric-shapes"

SERVER_VERSION := ${DEFAULT_VERSION}
CLIENT_VERSION := ${DEFAULT_VERSION}

.PHONY: all clean test lint gomodule clean_gomodule gomodinit gomodtidy goget

all: clean gomodtidy test

clean:
	echo "Nothing to clean in ${PWD}"

test:
# check for raceconditions
	go test -race ./...
# run the benchmark tests
	go test -cpu ${BENCH_CPUS} -benchmem -benchtime ${BENCH_ITERATIONS}x -bench=. ./...

lint:
	golint ./...

# clean the go module
clean_gomodule:
	-rm -f go.mod go.sum

# initialize module
gomodinit:
	go mod init ${GOMODULENAME}

# add module requirements and sums
gomodtidy:
	go mod tidy

gomodule: clean_gomodule gomodinit gomodtidy

goget:
	go get -d -v ./...
