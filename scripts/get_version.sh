#!/bin/bash

# $1 = version number

if [ "$#" -ne 1 ]; then
    echo "use $0 <config.json path>"
    exit 1
fi

CONFIG_FILE=configs/$1.json

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Can't find $CONFIG_FILE"
    exit 1
fi

GRAVITY_SDK=$(jq -r '.package["gravity-sdk"]' "$CONFIG_FILE")
NATS_JETSTREAM=$(jq -r '.gravity["nats-jetstream"]' "$CONFIG_FILE")
GRAVITY_DISPATCHER=$(jq -r '.gravity["gravity-dispatcher"]' "$CONFIG_FILE")
ATOMIC=$(jq -r '.gravity["atomic"]' "$CONFIG_FILE")
GRAVITY_ADAPTER_MSSQL=$(jq -r '.gravity["gravity_adapter_mssql"]' "$CONFIG_FILE")

echo "gravity_sdk_version=$GRAVITY_SDK"
echo "nats_jetstream_version=$NATS_JETSTREAM"
echo "gravity_dispatcher_version=$GRAVITY_DISPATCHER"
echo "atomic_version=$ATOMIC"
echo "gravity_adapter_mssql_version=$GRAVITY_ADAPTER_MSSQL"