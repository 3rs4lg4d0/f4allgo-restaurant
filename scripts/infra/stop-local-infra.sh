#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

cd ../../
cd deployments/kafka && docker-compose stop && cd ../..
cd deployments/metrics && docker-compose stop && cd ../..
cd deployments/postgres && docker-compose stop && cd ../..

cd -
