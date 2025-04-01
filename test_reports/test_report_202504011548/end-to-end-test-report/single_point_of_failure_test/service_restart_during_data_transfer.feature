Feature: Gravity2 MSSQL to MySQL - Service restart during data transfer
    Background: Set up single point of failure test
        Given Load the initial configuration file
        Given Create all services
        Given Start the "source-mssql" service (timeout "60")
        Given Initialize the "source-mssql" table Accounts
        Given Start the "target-mysql" service (timeout "60")
        Given Initialize the "target-mysql" table Accounts
        Given Start the "nats-jetstream" service (timeout "60")
        Given Start the "gravity-dispatcher" service (timeout "60")
        Then Wait "10" seconds
        Given Create Data Product Accounts
        Given Set up atomic flow document
        Given Start the "atomic" service (timeout "60")
        Given Start the "gravity-adapter-mssql" service (timeout "60")
        
    Scenario Outline: Perform insertions, updates, or deletions of data, and restart services during data transfer. <ID> <Status>
        Then "source-mssql" table "Accounts" has "0" data (timeout "3")
        Given "source-mssql" table "Accounts" continuously inserting "3000" data (starting ID "1")
        Given docker compose "stop" service "<RestartService>" (in "foreground")
        Then container "<RestartService>" was "exited" (timeout "120")
        # Then Wait "20" seconds
        Given docker compose "start" service "<RestartService>" (in "foreground")
        When container "<RestartService>" ready (timeout "120")
        Then wait for "source-mssql" table "Accounts" insertion to complete (timeout "120")
        Then Check the nats stream DP has correct messages (timeout "10")
        And "target-mysql" has the same content as "source-mssql" in "Accounts" (timeout "90")
        
        Given "source-mssql" table "Accounts" continuously updating "3000" data - appending suffix 'updated' to each Name field (starting ID "1") and inserting "1000" data (starting ID "3001")
        Given docker compose "stop" service "<RestartService>" (in "foreground")
        Then container "<RestartService>" was "exited" (timeout "120")
        # Then Wait "20" seconds
        Given docker compose "start" service "<RestartService>" (in "foreground")
        When container "<RestartService>" ready (timeout "120")
        Then wait for "source-mssql" table "Accounts" update and insertion to complete (timeout "120")
        Then Check the nats stream DP has correct messages (timeout "10")
        And "target-mysql" has the same content as "source-mssql" in "Accounts" (timeout "90")

        Given "source-mssql" table "Accounts" continuous cleanup
        Given docker compose "stop" service "<RestartService>" (in "foreground")
        Then container "<RestartService>" was "exited" (timeout "120")
        # Then Wait "20" seconds
        Given docker compose "start" service "<RestartService>" (in "foreground")
        When container "<RestartService>" ready (timeout "120")
        Then wait for "source-mssql" table "Accounts" cleanup to complete (timeout "120")
        Then Check the nats stream DP has correct messages (timeout "10")
        And "target-mysql" table "Accounts" has "0" data (timeout "120")
        
        Examples:
            | ID   |   RestartService      | Status   |
            | M(1) | gravity-adapter-mssql | |
            | M(2) | gravity-dispatcher    | |
            | M(3) | nats-jetstream        | [waived] |
            | M(4) | atomic                | |
            | M(5) | target-mysql          | |
