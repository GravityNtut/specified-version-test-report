#!/bin/bash

# create config file in test reports directory

# $1 nats_jetstream_version
# $2 gravity_dispatcher_version
# $3 atomic_version
# $4 gravity_adapter_mssql_version
# $5 gravity_sdk_version
# $6 test start time

save_path=test_reports/test_report_$6/config.json

echo '{
    "package": {
        "gravity-sdk": "'$5'"
    },
    "gravity": {
        "nats-jetstream": "'$1'",
        "gravity-dispatcher": "'$2'",
        "atomic": "'$3'",
        "gravity_adapter_mssql": "'$4'"
    }
}' > $save_path

echo "Config file created at $save_path"