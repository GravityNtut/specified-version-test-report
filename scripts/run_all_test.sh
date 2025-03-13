#!/bin/bash

# run gravity-cli-tests test
mkdir test_code
rm -rf test_code/*
cd test_code
# TODO: 改main branch
git clone https://github.com/BrobridgeOrg/gravity-cli-tests.git --branch GN-206_specified_version_test
cd gravity-cli-tests
earthly -P +specified-version-test

# TODO: 移除 改workflow觸發
sh ./create_test_report.sh