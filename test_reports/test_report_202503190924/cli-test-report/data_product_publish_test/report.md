# Test Report

**Total tests:** 63
**Failures:** 22
**Errors:** 0
**Time:** 95.33 seconds

## Data Product publish


### [Publish for data product of the event (Success scenario)](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_publish_test/data_product_publish_test.feature#L9)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| Publish for data product of the event (Success scenario) M(1)  | ✅Passed | 2.15 |  |
| Publish for data product of the event (Success scenario) M(2)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(3)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(4)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(5)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(6)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(7)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(8)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(9)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(10)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(11)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(12)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(13)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(14)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(15)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(16)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(17)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(18)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(19)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(20)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(21)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(22)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(23)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(24)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(25)  | ✅Passed | 3.21 |  |
| Publish for data product of the event (Success scenario) M(26)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(27)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(28)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(29)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(30)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(31)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(32)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(33)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(34)  | ✅Passed | 2.16 |  |
| Publish for data product of the event (Success scenario) M(35)  | ✅Passed | 2.16 |  |

### [Similar to the successful scenario, but after publishing, update the ruleset to enabled first, then update the data product to enabled](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_publish_test/data_product_publish_test.feature#L57)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| Similar to the successful scenario, but after publishing, update the ruleset to enabled first, then update the data product to enabled E1(1)  | ✅Passed | 0.16 |  |

### [publish for data product of the event (failure scenario for the publish command)](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_publish_test/data_product_publish_test.feature#L69)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| publish for data product of the event (failure scenario for the publish command) E2(1)  | ✅Passed | 0.64 |  |
| publish for data product of the event (failure scenario for the publish command) E2(2)  | ❌Failed | 0.14 | Step CLI returns create failed: publish should be failed |
| publish for data product of the event (failure scenario for the publish command) E2(3)  | ❌Failed | 0.14 | Step CLI returns create failed: publish should be failed |
| publish for data product of the event (failure scenario for the publish command) E2(4)  | ❌Failed | 0.14 | Step CLI returns create failed: publish should be failed |

### [The command executes successfully but does not publish to the specified DP](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_publish_test/data_product_publish_test.feature#L83)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| The command executes successfully but does not publish to the specified DP E3(1)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(2)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(3)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(4)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(5)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(6)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(7)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(8)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(9)  | ✅Passed | 5.14 |  |
| The command executes successfully but does not publish to the specified DP E3(10)  | ⏭️Skipped | 0.01 |  |
| The command executes successfully but does not publish to the specified DP E3(11)  | ❌Failed | 0.22 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(12)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(13)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(14)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(15)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(16)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(17)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(18)  | ✅Passed | 5.14 |  |
| The command executes successfully but does not publish to the specified DP E3(19)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |
| The command executes successfully but does not publish to the specified DP E3(20)  | ❌Failed | 0.23 | Step Query GVT_default_DP_"'drink'" has no "'drinkEvent'": expected not publish in GVT_default_DP，but now in GVT_default_DP |

### [publish for data product of the event (The same event is published to multiple data products)](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_publish_test/data_product_publish_test.feature#L114)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| publish for data product of the event (The same event is published to multiple data products) E4(1)  | ✅Passed | 0.24 |  |

### [publish for data product of the event (Continuously publish two events with the same PK value, but the number and content of other fields are different.)](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_publish_test/data_product_publish_test.feature#L127)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| publish for data product of the event (Continuously publish two events with the same PK value, but the number and content of other fields are different.) E5(1)  | ✅Passed | 0.26 |  |

### [verify that payloads are not combined during rapid consecutive event publishing](https://github.com/BrobridgeOrg/gravity-cli-tests/tree/main/data_product_publish_test/data_product_publish_test.feature#L138)

| Test Case | Status | Time (s) | Failure Reason |
|-----------|--------|----------|----------------|
| verify that payloads are not combined during rapid consecutive event publishing E6(1)  | ❌Failed | 0.47 | Step Query GVT_default_DP_drink has 3 events in sequence: {"id":1}, {"id":2}, and {"id":3}: payload2 is not matched, expected: {"id":2}, actual: {"name":"test"} |

