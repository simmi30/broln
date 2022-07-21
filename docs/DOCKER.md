# Docker Instructions

There are two flavors of Dockerfiles available:
 - `Dockerfile`: Used for production builds. Checks out the source code from
   GitHub during build. The build argument `--build-arg checkout=v0.x.x-beta`
   can be used to specify what git tag or commit to check out before building.
 - `dev.Dockerfile` Used for development or testing builds. Uses the local code
   when building and allows local changes to be tested more easily.

## Development/testing

To build a standalone development image from the local source directory, use the
following command:

```shell
⛰  docker build --tag=myrepository/broln-dev -f dev.Dockerfile .
```

There is also a `docker-compose` setup available for development or testing that
spins up a `brond` backend alongside `broln`. Check out the documentation at
[docker/README.md](../docker/README.md) to learn more about how to use that
setup to create a small local Lightning Network.

## Production (manual build)

To use Docker in a production environment, you can run `broln` by creating a
Docker container, adding the appropriate command-line options as parameters.

You first need to build the `broln` docker image:

```shell
⛰  docker build --tag=myrepository/broln --build-arg checkout=v0.11.1-beta .
```

It is recommended that you checkout the latest released tag.

You can continue by creating and running the container:

```shell
⛰  docker run myrepository/broln [command-line options]
```

## Production (official images)

Starting with `broln v0.12.0-beta`, there are official, automatically built docker
images of `broln` available in the
[`lightninglabs/broln` repository on Docker Hub](https://hub.docker.com/r/lightninglabs/broln).

You can just pull those images by specifying a release tag:

```shell
⛰  docker pull lightninglabs/broln:v0.12.0-beta
⛰  docker run lightninglabs/broln [command-line options]
```

### Verifying docker images

To verify the `broln` and `brolncli` binaries inside the docker images against the
signed, [reproducible release binaries](release.md), there is a verification
script in the image that can be called (before starting the container for
example):

```shell
⛰  docker run --rm --entrypoint="" lightninglabs/broln:v0.12.1-beta /verify-install.sh v0.12.1-beta
⛰  OK=$?
⛰  if [ "$OK" -ne "0" ]; then echo "Verification failed!"; exit 1; done
⛰  docker run lightninglabs/broln [command-line options]
```

## Volumes

A Docker volume will be created with your `.broln` directory automatically, and will
persist through container restarts.

You can also optionally manually specify a local folder to be used as a volume:

```shell
⛰  docker create --name=mybrolncontainer -v /media/broln-docker/:/root/.broln myrepository/broln [command-line options]
```

## Example

Here is an example testnet `broln` that uses Neutrino:

```shell
⛰  docker run --name broln-testnet myrepository/broln --brocoin.active --brocoin.testnet --brocoin.node=neutrino --neutrino.connect=faucet.lightning.community
```

Create a wallet (and write down the seed):

```shell
⛰  docker exec -it broln-testnet brolncli create
```

Confirm `broln` has begun to synchronize:

```shell
⛰  docker logs broln-testnet
[snipped]
2018-05-01 02:28:01.201 [INF] RPCS: RPC server listening on 127.0.0.1:10019
2018-05-01 02:28:01.201 [INF] LTND: Waiting for chain backend to finish sync, start_height=2546
2018-05-01 02:28:01.201 [INF] RPCS: gRPC proxy started at 127.0.0.1:8080
2018-05-01 02:28:08.999 [INF] LNWL: Caught up to height 10000
2018-05-01 02:28:09.872 [INF] BTCN: Processed 10547 blocks in the last 10.23s (height 10547, 2012-05-28 05:02:32 +0000 UTC)
```

This is a simple example, it is possible to use any command-line options necessary
to expose RPC ports, use `brond` or `brocoind`, or add additional chains.

## broln Development and Testing

To test the Docker production image locally, run the following from the project root:

```shell
⛰  docker build . -t myrepository/broln:master
```

To choose a specific [branch](https://github.com/brsuite/broln/branches) or [tag](https://hub.docker.com/r/lightninglabs/broln/tags?page=1&ordering=last_updated) instead, use the "checkout" build-arg. For example, to build the latest tagged commit:

```shell
⛰  docker build . --build-arg checkout=v0.13.0-beta -t myrepository/broln:v0.13.0-beta
```

To build the image using the most current tag:

```shell
⛰  docker build . --build-arg checkout=$(git describe --tags `git rev-list --tags --max-count=1`) -t myrepository/broln:latest-tag
```

Once the image has been built and tagged locally, start the container:

```shell
⛰  docker run --name=broln-testnet -it myrepository/broln:latest-tag --brocoin.active --brocoin.testnet --brocoin.node=neutrino --neutrino.connect=faucet.lightning.community
```
