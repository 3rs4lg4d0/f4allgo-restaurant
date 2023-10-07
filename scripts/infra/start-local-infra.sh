#!/bin/bash

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)"
NETWORK_NAME="f4allgo-network"

if ! docker network inspect $NETWORK_NAME >/dev/null 2>&1; then
    docker network create $NETWORK_NAME
fi

cd "$BASE_DIR/deployments/kafka" && docker-compose up -d

skip=false
cd "$BASE_DIR/deployments/metrics"
if docker-compose ps -a grafana | grep -q 'Exited'; then
    echo "Skipping resource creation in Grafana"
    skip=true
fi
docker-compose up -d && \
if $skip; then
    echo "Skipping resource creation in Grafana"
else
    docker run -it --rm --network $NETWORK_NAME --volume $BASE_DIR/assets/grafana:/root ubuntu bash -c ' \
        echo -e "Installing packages..." && apt update > /dev/null 2>&1 && apt install -y curl jq netcat > /dev/null 2>&1 && \
        while ! nc -zv prometheus 9090; do echo "Waiting for Prometheus to be ready..."; sleep 1; done && \
        while ! nc -zv grafana 3000; do echo "Waiting for Grafana to be ready..."; sleep 1; done && \
        response=$(curl -s -X POST http://grafana:3000/api/datasources -H "Content-Type: application/json" --data "{
            \"name\": \"Prometheus for f4allgo-restaurant\",
            \"type\": \"prometheus\",
            \"access\": \"proxy\",
            \"url\": \"http://prometheus:9090\"}") && \
        echo -e "Prometheus datasource created ✅" && \
        datasource_uid=$(echo "$response" | jq -r ".datasource.uid") && \
        dashboard_content=$(cat /root/f4allgo-restaurant-dashboard.json) && \
        payload=$(echo { \
            \"dashboard\": $dashboard_content, \
            \"overwrite\":true, \
            \"inputs\":[\{\"name\":\"DS_PROMETHEUS\",\"type\":\"datasource\",\"pluginId\":\"prometheus\",\"value\":\"$datasource_uid\"\}], \
            \"folderUid\":\"\"}) && \
        curl -s -o /dev/null -X POST http://grafana:3000/api/dashboards/import -H "Content-Type: application/json" --data "$payload" && \
        echo -e "Grafana dashboard imported ✅"
    '
fi

cd "$BASE_DIR/deployments/postgres" && docker-compose up -d && \
    docker-compose exec postgres bash -c \
        "while ! pg_isready; do echo 'Waiting for postgres to be ready...'; sleep 1; done"
cd -
