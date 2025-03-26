#!/bin/bash

# $1 nats_jetstream_version
# $2 gravity_dispatcher_version
# $3 atomic_version
# $4 gravity_adapter_mssql_version
# $5 gravity_sdk_version

cd test_code
root_path=$(pwd)

# run gravity-cli-tests test
cd gravity-cli-tests
earthly -P +specified-version-test --nats_jetstream_version=$1 --gravity_dispatcher_version=$2 --gravity_sdk_version=$5

# run e2e-test
cd $root_path
cd e2e-tests
earthly -P +specified-version-test --nats_jetstream_version=$1 --gravity_dispatcher_version=$2 --atomic_version=$3 --gravity_adapter_mssql_version=$4
