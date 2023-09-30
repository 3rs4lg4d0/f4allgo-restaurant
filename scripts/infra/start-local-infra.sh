#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)"
NETWORK_NAME="f4allgo-network"

if ! docker network inspect $NETWORK_NAME >/dev/null 2>&1; then
    docker network create $NETWORK_NAME
fi

cd "$BASE_DIR/deployments/kafka" && docker-compose up -d
cd "$BASE_DIR/deployments/metrics" && docker-compose up -d
cd "$BASE_DIR/deployments/postgres" && docker-compose up -d && \
    docker-compose exec postgres bash -c \
        "while ! pg_isready; do echo 'Waiting for postgres to be ready...'; sleep 1; done" && cd ../..
cd -
