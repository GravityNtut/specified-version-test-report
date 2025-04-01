# Test Report

**Total tests:** 11
**Failures:** 1
**Errors:** 0
**Time:** 2176.44 seconds

## Gravity2 MSSQL to MySQL - Component restart while No data changes


### [After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data](https://github.com/BrobridgeOrg/End-to-End-test/tree/main/single_point_of_failure_test/single_point_failure.feature#L17)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(1)  | ✅Passed | 175.55 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(2)  | ✅Passed | 172.87 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(3)  | ✅Passed | 173.94 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(4)  | ✅Passed | 163.72 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(5)  | ✅Passed | 171.31 |  |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(6)  | ✅Passed | 169.84 |  |

## Gravity2 MSSQL to MySQL - Service restart during data transfer


### [Perform insertions, updates, or deletions of data, and restart services during data transfer.](https://github.com/BrobridgeOrg/End-to-End-test/tree/main/single_point_of_failure_test/service_restart_during_data_transfer.feature#L17)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(1)  | ✅Passed | 252.55 |  |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(2)  | ✅Passed | 254.62 |  |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(3)  | ⚠️Waived | 130.59 | Step "target-mysql" has the same content as "source-mssql" in "Accounts" (timeout "90"): number of records in table 'Accounts' is 3031, expected 3000 after 90 second |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(4)  | ✅Passed | 253.57 |  |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(5)  | ✅Passed | 257.88 |  |

