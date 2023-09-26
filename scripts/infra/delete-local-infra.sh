#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

cd ../../
cd deployments/kafka && docker-compose down && cd ../..
cd deployments/metrics && docker-compose down && cd ../..
cd deployments/postgres && docker-compose down && cd ../..

docker network rm f4allgo-network

cd -
