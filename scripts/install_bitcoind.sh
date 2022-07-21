#!/usr/bin/env bash

set -ev

BROCOIND_VERSION=${BROCOIN_VERSION:-22.0}

docker pull lightninglabs/brocoin-core:$BROCOIND_VERSION
CONTAINER_ID=$(docker create lightninglabs/brocoin-core:$BROCOIND_VERSION)
sudo docker cp $CONTAINER_ID:/opt/brocoin-$BROCOIND_VERSION/bin/brocoind /usr/local/bin/brocoind
docker rm $CONTAINER_ID
