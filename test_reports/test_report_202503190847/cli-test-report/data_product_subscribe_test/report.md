# Test Report

**Total tests:** 24
**Failures:** 6
**Errors:** 0
**Time:** 17.51 seconds

## Data Product subscribe


### [Successful scenario. Use the `product sub` command to receive all data published to the specified data product](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_subscribe_test/data_product_subscribe_test.feature#L9)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(1)  | ❌Failed | 1.19 | Step The CLI returns all events data within the "'-1'" and after "'1'": <br>expect: create <br>actual: INSERT |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(2)  | ✅Passed | 1.19 |  |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(3)  | ❌Failed | 1.19 | Step The CLI returns all events data within the "'200'" and after "'1'": <br>expect: create <br>actual: INSERT |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(4)  | ✅Passed | 1.18 |  |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(5)  | ✅Passed | 1.18 |  |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(6)  | ❌Failed | 1.18 | Step The CLI returns all events data within the "'[ignore]'" and after "'1'": <br>expect: create <br>actual: INSERT |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(7)  | ✅Passed | 1.18 |  |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(8)  | ✅Passed | 1.18 |  |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(9)  | ❌Failed | 1.18 | Step The CLI returns all events data within the "'131,200'" and after "'1'": <br>expect: create <br>actual: INSERT |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(10)  | ✅Passed | 1.18 |  |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(11)  | ❌Failed | 1.18 | Step The CLI returns all events data within the "'-1'" and after "'[ignore]'": <br>expect: create <br>actual: INSERT |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(12)  | ✅Passed | 1.18 |  |
| Successful scenario. Use the `product sub` command to receive all data published to the specified data product M(13)  | ❌Failed | 1.19 | Step The CLI returns all events data within the "'-1'" and after "'5'": <br>expect: create <br>actual: INSERT |

### [Failure scenario. Use the `product sub` command to receive all data published to the specified data product](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_subscribe_test/data_product_subscribe_test.feature#L34)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(1)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(2)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(3)  | ✅Passed | 0.20 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(4)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(5)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(6)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(7)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(8)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(9)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(10)  | ✅Passed | 0.19 |  |
| Failure scenario. Use the `product sub` command to receive all data published to the specified data product E1(11)  | ✅Passed | 0.19 |  |

