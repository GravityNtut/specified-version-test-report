#!/bin/bash

# run gravity-cli-tests test
mkdir -p test_code
rm -rf test_code/*
cd test_code

# TODO: 改main branch
git clone https://github.com/BrobridgeOrg/gravity-cli-tests.git --branch GN-206_specified_version_test
cd gravity-cli-tests
# TODO: 改shell參數透過action傳入
earthly -P +specified-version-test --gravity_sdk_version=v2.0.7 --nats_jetstream_version=v1.3.21-20250117 --gravity_dispatcher_version=v0.0.31-20250220
