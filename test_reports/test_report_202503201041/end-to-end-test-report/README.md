![workflow](https://github.com/BrobridgeOrg/End-to-End-test/actions/workflows/daily_test.yml/badge.svg)
![Test_time](https://byob.yarr.is/GravityNtut/BYOB_badge/test_time/E2E_test/E2E_test_services_versions.json)
![Dispatcher_version](https://byob.yarr.is/GravityNtut/BYOB_badge/dispatcher_version/E2E_test/E2E_test_services_versions.json)
![Atomic_version](https://byob.yarr.is/GravityNtut/BYOB_badge/atomic_version/E2E_test/E2E_test_services_versions.json)
![nats-server_version](https://byob.yarr.is/GravityNtut/BYOB_badge/nats_server_version/E2E_test/E2E_test_services_versions.json)
![adapter-mssql_version](https://byob.yarr.is/GravityNtut/BYOB_badge/adapter_mssql_version/E2E_test/E2E_test_services_versions.json)



# E2E_Test
This is end to end test repository for Gravity.

## View daily build test results
To see the test results please visit [this page](https://github.com/BrobridgeOrg/tests-report?tab=readme-ov-file#test-summary).  
To learn how to view the test report, you can refer to [this page](https://github.com/BrobridgeOrg/tests-report/blob/main/HOW_TO_USE.md).  

## Require services
- [Gravity-Dispatcher](https://github.com/BrobridgeOrg/gravity-dispatcher)
- [Gravity-Atomic](https://github.com/BrobridgeOrg/atomic)
- [Nats jetstream](https://github.com/BrobridgeOrg/gravity-nats-server)
- [Gravity-Adapter-mssql](https://github.com/BrobridgeOrg/gravity-adapter-mssql)
- sql-server 
- mysql

## File Structure
The repository is organized by test tasks, where each folder represents a specific test task:
- `.feature` files: Contain test scenarios and detailed test steps
- `config.json`: Environment settings for the test
- `assets/`: Directory containing required test files
- `docker-compose.yaml`: Service definitions for the test environment

For example, in `single_point_of_failure_test`
This test task includes:
- Two feature files:
    - `service_restart_during_data_transfer.feature`
    - `single_point_failure.feature`
- Supporting files and directories for test execution

Directory structure shown below:

```
.
├── regression_test
│   ├── assets
│   ├── tmp
│   ├── boolean_check.feature
│   ├── boolean_check_test.go
│   ├── config.json
│   └── docker-compose.yaml
├── single_point_of_failure_test
│   ├── assets
│   ├── tmp
│   ├── config.json
│   ├── docker-compose.yaml
│   ├── service_restart_during_data_transfer.feature
│   ├── single_point_failure.feature
│   ├── single_point_failure_test.go
│   └── UsefulCmd.txt
```

## Usage

Clone this repository:
```
git clone git@github.com:GravityNtut/E2E_Test.git
cd E2E_Test
```

### Execute through docker compose
1. Change to the specific test task directory:
    ```sh
    cd <test_task_directory>
    ```

2. Run the test (takes approximately 30 minutes):
    ```sh
    go test --timeout=60m
    ```

### Execute through kubernetes
- TODO
