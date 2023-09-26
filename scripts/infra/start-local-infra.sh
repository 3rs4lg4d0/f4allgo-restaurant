#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

NETWORK_NAME="f4allgo-network"

if ! docker network inspect $NETWORK_NAME >/dev/null 2>&1; then
    docker network create $NETWORK_NAME
fi

cd ../../
cd deployments/kafka && docker-compose up -d && cd ../..
cd deployments/metrics && docker-compose up -d && cd ../..
cd deployments/postgres && docker-compose up -d && \
    docker-compose exec postgres bash -c \
        "while ! pg_isready; do echo 'Waiting for postgres to be ready...'; sleep 1; done" && cd ../..
cd -
