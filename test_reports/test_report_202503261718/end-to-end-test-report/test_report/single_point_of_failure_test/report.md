# Test Report

**Total tests:** 11
**Failures:** 1
**Errors:** 0
**Time:** 2195.58 seconds

## Gravity2 MSSQL to MySQL - Component restart while No data changes


### [After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data](https://github.com/BrobridgeOrg/End-to-End-test/tree/main/single_point_of_failure_test/single_point_failure.feature#L17)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(1)  | ✅Passed | 175.45 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(2)  | ✅Passed | 173.74 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(3)  | ✅Passed | 172.99 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(4)  | ✅Passed | 166.76 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(5)  | ✅Passed | 172.64 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(6)  | ✅Passed | 175.22 |  |

## Gravity2 MSSQL to MySQL - Service restart during data transfer


### [Perform insertions, updates, or deletions of data, and restart services during data transfer.](https://github.com/BrobridgeOrg/End-to-End-test/tree/main/single_point_of_failure_test/service_restart_during_data_transfer.feature#L17)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(1)  | ✅Passed | 253.87 |  |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(2)  | ✅Passed | 257.68 |  |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(3)  | ⚠️Waived | 129.70 | Step "target-mysql" has the same content as "source-mssql" in "Accounts" (timeout "90"): number of records in table 'Accounts' is 3002, expected 3000 after 90 second |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(4)  | ✅Passed | 258.62 |  |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(5)  | ✅Passed | 258.91 |  |

