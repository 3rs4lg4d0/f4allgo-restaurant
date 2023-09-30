#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)"

cd "$BASE_DIR/deployments/kafka" && docker-compose down
cd "$BASE_DIR/deployments/metrics" && docker-compose down
cd "$BASE_DIR/deployments/postgres" && docker-compose down

docker network rm f4allgo-network

cd -
