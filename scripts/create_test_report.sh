#!/bin/bash

# $1 nats_jetstream_version
# $2 gravity_dispatcher_version
# $3 atomic_version
# $4 gravity_adapter_mssql_version
# $5 gravity_sdk_version
# $6 test start time

report_path=test_reports/test_report_$6
# cli
mkdir -p $report_path/cli-test-report
rm -rf $report_path/cli-test-report/*
cp -r test_code/gravity-cli-tests/test_report/* $report_path/cli-test-report
# e2e
mkdir -p $report_path/end-to-end-test-report
rm -rf $report_path/end-to-end-test-report/*
cp -r test_code/e2e-tests/test_report/* $report_path/end-to-end-test-report

# create test report
output_file=$report_path/README.md
rm -f $output_file
touch $output_file
# TODO: 改網址到BB
base_url="https://github.com/GravityNtut/specified-version-test-report/blob/main" 

cat <<EOT > $output_file
# Test Summary
To learn how to view the test report, you can refer to [this page](https://github.com/GravityNtut/specified-version-test-report/blob/main/HOW_TO_USE.md).

## Service Version
| Service | Version |
|---------|---------|
| nats-jetstream | $1 |
| gravity-dispatcher | $2 |
| atomic | $3 |
| gravity-adapter-mssql | $4 |
| gravity-sdk | $5 |
EOT

for folder in $(find ./$report_path/ -mindepth 1 -maxdepth 1 -type d -name "*-test-report"); do
  folder_name=$(basename "$folder")
  echo "## ${folder_name}" >> $output_file
  echo "" >> $output_file
  echo "| Folder Name | Total Tests | Passed | Failed | Link |" >> $output_file
  echo "|-------------|-------------|--------|--------|------|" >> $output_file

  find "$folder" -type f -name "report.xml" | sort -t '/' -k1,2 | while read -r xml; do
    relative_path="${xml%/report.xml}" 
    relative_folder_path="${relative_path#./$folder_name/}"  
    total_tests=$(xmllint --xpath "string(/testsuites/@tests)" "$xml")
    failures=$(xmllint --xpath "string(/testsuites/@failures)" "$xml")
    passed=$((total_tests - failures))

    report_link="${base_url}/${relative_folder_path}/report.md"

    test_folder_name=$(basename "$relative_folder_path")

    echo "| ${test_folder_name} | ${total_tests} | ${passed} | ${failures} | [View Report](${report_link}) |" >> "$output_file"
  done

  echo "" >> "$output_file"
done

echo "Test summary saved to $output_file"
