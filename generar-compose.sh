#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <yaml_file> <amount_clients>"
    exit 1
fi

YAML_FILE="$1"
AMOUNT_CLIENTS="$2"

cat > "$YAML_FILE" << EOF
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
    volumes:
      - ./server/config.ini:/config.ini
EOF

for i in $(seq 1 "$AMOUNT_CLIENTS"); do
    cat >> "$YAML_FILE" << EOF

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
      - BET_NUMBER=67890
      - CLIENT_DOCUMENT=1
      - CLIENT_FIRST_NAME=Santi
      - CLIENT_LAST_NAME=Sev
      - CLIENT_BIRTHDATE=2000-08-10
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
    depends_on:
      - server
EOF
done

cat >> "$YAML_FILE" << EOF

networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
EOF
