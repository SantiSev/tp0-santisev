#!/bin/bash

SERVER_PORT=$(grep 'SERVER_PORT' server/config.ini | cut -d'=' -f2)

PING_MESSAGE="hello darkness my old friend..."

SERVER_RESPONSE=$(echo "$PING_MESSAGE" | docker run --rm -i --network container:server busybox nc server $SERVER_PORT)

echo "Received response: '$SERVER_RESPONSE'"

RESULT="fail"
if [ "$SERVER_RESPONSE" = "$PING_MESSAGE" ]; then
    RESULT="success"
fi

echo "test_echo_server | result: $RESULT"