#!/bin/bash

# -----------------------------------------------------------------------------
# kafkactl.sh uses a very useful containerized version of 'kafkactl' to
# comunicate with the local containerized Kafka broker. To run the local
# containerized infra just execute "make start-infra".
# -----------------------------------------------------------------------------
docker run --platform linux/amd64 --rm --env BROKERS="domain-kafka:9092" \
    --network f4allgo-network -v ./kafkactl_config.yaml:/etc/kafkactl/config.yml deviceinsight/kafkactl "$@"
