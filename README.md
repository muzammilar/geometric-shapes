# geometric-shapes
An example of gRPC Servers and multi-lingual clients. A basic golang gRPC example (with server certificates).

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
docker-compose run --rm basebuilder
# To access the container itself, use any of the following command(s)
docker-compose run --rm basebuilder bash
docker-compose run --rm basebuilder sh
docker-compose run --rm basebuilder /bin/sh

## Step 3: (Can be skipped) (one time)
# Build the certs using basebuilder
docker-compose run --rm  basebuilder make certs


## Step 4:
# Run the containers.
# Run the grpc servers and clients
docker-compose up --build --detach

# Shutdown everything (and remove volumes)
docker-compose rm --force --stop -v
```

## Examples:
Some good gRPC examples are available here:
* https://github.com/grpc/grpc-web

## Notes:
Due to the go module being called `github.com/muzammilar/geometric-shapes` and not `github.com/muzammilar/geometric-shapes/go`, calling this go module from external repos, might lead to inconsistency.
One way to solve this problem is to use multiple repositories. Another way to solve this is to move the go code to `github.com/muzammilar/geometric-shapes/go/geometric-shapes` and make that a new module (however, using `-` in package name can cause other side effects).
