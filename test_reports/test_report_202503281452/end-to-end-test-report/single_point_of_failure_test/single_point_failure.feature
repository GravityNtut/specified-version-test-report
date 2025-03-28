Feature: Gravity2 MSSQL to MySQL - Component restart while No data changes
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
        
    Scenario Outline: After synchronizing changes, restart the component, wait for it ready, then add, update, or delete data <ID> <Status>
        Then "source-mssql" table "Accounts" has "0" data (timeout "3")
        Given "source-mssql" table "Accounts" inserted "1000" data (starting ID "1")
        Then "target-mysql" has the same content as "source-mssql" in "Accounts" (timeout "90")
        Given docker compose "stop" service "<RestartService>" (in "foreground")
        Then container "<RestartService>" was "exited" (timeout "120")
        Given docker compose "start" service "<RestartService>" (in "foreground")
        When container "<RestartService>" ready (timeout "120")
        Given "source-mssql" table "Accounts" updated "1000" data - appending suffix 'updated' to each Name field (starting ID "1")
        Given "source-mssql" table "Accounts" inserted "1000" data (starting ID "1001")
        Then "target-mysql" has the same content as "source-mssql" in "Accounts" (timeout "90")
        Given "source-mssql" table "Accounts" cleared
        Then "target-mysql" table "Accounts" has "0" data (timeout "120")

        Examples:
            | ID   |   RestartService      | Status |
            | M(1) | gravity-adapter-mssql | |
            | M(2) | gravity-dispatcher    | |
            | M(3) | nats-jetstream        | |
            | M(4) | atomic                | |
            | M(5) | source-mssql          | |
            | M(6) | target-mysql          | |

