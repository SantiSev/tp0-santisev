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
      - CLI_AGENCY_FILEPATH=/data/agency.csv
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-$i.csv:/data/agency.csv
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
