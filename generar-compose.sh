#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <yaml_file> <amount_clients>"
    exit 1
fi

YAML_FILE="$1"
AMOUNT_CLIENTS="$2"

cat > docker-compose-dev.yaml << EOF
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
EOF

for i in $(seq 1 "$AMOUNT_CLIENTS"); do
    cat >> docker-compose-dev.yaml << EOF

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
    networks:
      - testing_net
    depends_on:
      - server
EOF
done

cat >> docker-compose-dev.yaml << EOF

networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
EOF
