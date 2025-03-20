# Test Report

**Total tests:** 11
**Failures:** 11
**Errors:** 0
**Time:** 1266.38 seconds

## Gravity2 MSSQL to MySQL - Component restart while No data changes


### [After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data](https://github.com/BrobridgeOrg/End-to-End-test/tree/main/single_point_of_failure_test/single_point_failure.feature#L17)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(1)  | ❌Failed | 115.10 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(2)  | ❌Failed | 115.00 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(3)  | ❌Failed | 115.17 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(4)  | ❌Failed | 115.22 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(5)  | ❌Failed | 115.12 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data M(6)  | ❌Failed | 115.19 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |

## Gravity2 MSSQL to MySQL - Service restart during data transfer


### [Perform insertions, updates, or deletions of data, and restart services during data transfer.](https://github.com/BrobridgeOrg/End-to-End-test/tree/main/single_point_of_failure_test/service_restart_during_data_transfer.feature#L17)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(1)  | ❌Failed | 115.26 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(2)  | ❌Failed | 115.09 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(3)  | ⚠️Waived | 115.12 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(4)  | ❌Failed | 115.02 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |
| Perform insertions, updates, or deletions of data, and restart services during data transfer. M(5)  | ❌Failed | 115.07 | Step Start the "atomic" service (timeout "60"): the service 'atomic' timed out within 60 seconds |

