# Geometric Shapes RPC
An example of multi-lingual gRPC clients and servers to implement basic geometric computation on different shapes. The example illustrates usage both server certificates and metrics (using `stats.Handler` interface in golang). The performance impact of metrics calculation has not been evaluated.

```sh
## Step 1: (Can be skipped) (one time)
# Build the protobufs (and generate the *.pb.go files).
# The second command will run "make protos" in the container as well as recreate go module/update dependencies
docker-compose up --detach --build protobuilder
docker-compose run --rm protobuilder

## Step 2
# Build the base image.
# Since docker-compose currently doesn't support copying from local images (between multiple Dockerfiles),
# we use the host Operating system to copy the files as needed.
docker-compose up --detach --build basebuilder
docker-compose run --rm basebuilder make   # This needs to be called whenever code is modified
# Debug: To access the container itself, use any of the following command(s)
docker-compose run --rm basebuilder bash
docker-compose run --rm basebuilder sh
docker-compose run --rm basebuilder /bin/sh
docker-compose run --rm --entrypoint sh goclient

## Step 3: (One time only - Usually an offline stage)
# Build the certs using basebuilder
docker-compose run --rm  basebuilder make certs


## Step 4:
# Run the containers.
# Run the grpc servers and clients
# Note: This will start the protobuilder and basebuilder containers as well.
docker-compose up --build --detach

## Step 5:
# Connect to prometheus and query for metrics.
# Open your browser and go to `localhost:9090`

## Step 6:
# Shutdown everything (and remove networks and local images). Networks are removed in this.
docker-compose down --volumes
# Use `docker-compose down --rmi all --volumes` with above to images as well
# Remove everything (and remove volumes). Networks are not removed here.
docker-compose rm --force --stop -v
```

## Examples:
Some good gRPC examples are available here:
* https://github.com/grpc/grpc-web

## Notes:
Due to the go module being called `github.com/muzammilar/geomrpc` and not `github.com/muzammilar/geomrpc/go`, calling this go module from external repos, might lead to inconsistency.
One way to solve this problem is to use multiple repositories. Another way to solve this is to move the go code to `github.com/muzammilar/geomrpc/go/geomrpc` and make that a new module (however, using `-` in package name can cause other side effects).
