# The "start-local-infra.sh" uses docker-compose to start a local infrastructure
# for development but also uses an ephemeral container to execute some post setup
# that requires some linux tools to be available (like netcat, curl or jq). Instead
# of using a generic image and install this tools for every run of the script, this
# image is provided (basically to speed up the process of starting a fresh local infra).
FROM ubuntu:latest

RUN apt update && apt install -y curl jq netcat
