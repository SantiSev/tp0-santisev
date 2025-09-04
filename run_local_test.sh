#!/bin/bash

make docker-compose-down
make docker-compose-up

make docker-compose-logs > logs.txt