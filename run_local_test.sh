#!/bin/bash

make docker-compose-down
make docker-compose-up

rm -f logs.txt

make docker-compose-logs >> logs.txt