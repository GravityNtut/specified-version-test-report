#!/bin/bash

# $1 nats_jetstream_version
# $2 gravity_dispatcher_version
# $3 gravity_sdk_version

# run gravity-cli-tests test
mkdir -p test_code
rm -rf test_code/*
cd test_code

# TODO: 改main branch
git clone https://github.com/BrobridgeOrg/gravity-cli-tests.git --branch GN-206_specified_version_test
cd gravity-cli-tests
# TODO: 改shell參數透過action傳入
earthly -P +specified-version-test --nats_jetstream_version=$1 --gravity_dispatcher_version=$2 --gravity_sdk_version=$3
