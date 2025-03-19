#!/bin/bash

# $1 nats_jetstream_version
# $2 gravity_dispatcher_version
# $3 gravity_sdk_version

root_path=$(pwd)
echo "root_path: $root_path"
ls
cd test_code

# run gravity-cli-tests test
cd gravity-cli-tests
earthly -P +specified-version-test --nats_jetstream_version=$1 --gravity_dispatcher_version=$2 --gravity_sdk_version=$3

cd $root_path