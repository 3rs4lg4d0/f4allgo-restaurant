#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)"

cd "$BASE_DIR/deployments/kafka" && docker-compose stop
cd "$BASE_DIR/deployments/metrics" && docker-compose stop
cd "$BASE_DIR/deployments/postgres" && docker-compose stop

cd -
