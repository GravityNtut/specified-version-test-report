package e2e

import (
	"bytes"
	"context"
	utils "e2e_test/testutils"
	"fmt"
	"os/exec"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/cucumber/godog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"gorm.io/gorm"
)

var opts = godog.Options{
	Format:        "pretty",
	Paths:         []string{"./"},
	StopOnFailure: false,
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
	pflag.Parse()
}

func TestFeatures(t *testing.T) {
	utils.InitLog()

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

var streamRecordCount int

func reset() {
	streamRecordCount = 0
}

// Table Account schema
type Account struct {
	ID    int
	Name  string `gorm:"size:50"`
	Phone string `gorm:"size:16"`
}

var Cmd *exec.Cmd

type SharedState struct {
	Insertion insertionState
	Update    updateState
	Delete    deleteState
	WG        sync.WaitGroup
}

type insertionState struct {
	Done         chan bool
	DataCount    int
	CurrentTotal int
}

type updateState struct {
	Done         chan bool
	DataCount    int
	CurrentTotal int
}

type deleteState struct {
	Done chan bool
}

var state = &SharedState{
	Insertion: insertionState{
		Done:         make(chan bool),
		DataCount:    0,
		CurrentTotal: 0,
	},
	Update: updateState{
		Done:         make(chan bool),
		DataCount:    0,
		CurrentTotal: 0,
	},
	Delete: deleteState{
		Done: make(chan bool),
	},
}

func InitInsertionState(total int) {
	state.Insertion.Done = make(chan bool)
	state.Insertion.DataCount = total
	state.Insertion.CurrentTotal = 0
}

func InitUpdateState(total int) {
	state.Update.Done = make(chan bool)
	state.Update.DataCount = total
	state.Update.CurrentTotal = 0
}

func InitDeleteState() {
	state.Delete.Done = make(chan bool)

}

func GetContainerStateByName(ctName string) (*types.ContainerState, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containerInfo, err := cli.ContainerInspect(context.Background(), ctName)
	if err != nil {
		if client.IsErrNotFound(err) {
			return nil, fmt.Errorf("container name %s is not found", ctName)
		}
		return nil, err
	}

	return containerInfo.State, nil
}

func VerifyRowCountTimeoutSeconds(loc, tableName string, expectedRowCount, timeoutSec int) error {
	db, err := utils.GetDBInstance(loc)
	if err != nil {
		return err
	}

	var currRowCount int64
	var retry int
	for retry = 0; retry < timeoutSec; retry++ {
		db.Table(tableName).Count(&currRowCount)
		log.Infof("Waiting for '%s' table '%s' to has %d records.. (%d sec), current total: %d",
			loc, tableName, expectedRowCount, retry, currRowCount)
		if currRowCount == int64(expectedRowCount) {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("expected %d records, but got %d", expectedRowCount, currRowCount)
}

func InsertDummyDataFromID(loc, tableName string, total int, beginID int) error {
	db, err := utils.GetDBInstance(loc)
	if err != nil {
		return err
	}
	log.Infof("Inserting total %d records to '%s' - '%s', begin ID '%d'",
		total, loc, tableName, beginID)
	// Insert dummy data from beginID
	start := time.Now()

	for i := beginID; i < total+beginID; i++ {
		account := Account{
			ID:    i,
			Name:  fmt.Sprintf("Name %d", i),
			Phone: fmt.Sprintf("Phone %d", i),
		}
		query := fmt.Sprintf("INSERT INTO %s (id, name, phone) VALUES (%d, '%s', '%s')",
			tableName, account.ID, account.Name, account.Phone)
		result := db.Exec(query)
		if result.Error != nil {
			log.Printf("Failed to insert '%d th' record: %v", i, result.Error)
		}
	}

	elapsed := time.Since(start)
	log.Infof("Inserted total %d records to '%s', ID '%d ~ %d' (elapsed: %s)",
		total, loc, beginID, beginID+total-1, elapsed)
	return nil
}

func InsertDummyDataFromIDGoroutine(loc, tableName string, total int, beginID int) error {
	InitInsertionState(total)
	state.WG.Add(1)
	go func() {
		defer state.WG.Done()
		db, err := utils.GetDBInstance(loc)
		if err != nil {
			log.Error(err)
			state.Insertion.Done <- false
			return
		}
		log.Infof("Inserting total %d records to '%s' - '%s', begin ID '%d'",
			total, loc, tableName, beginID)
		// Insert dummy data from beginID
		start := time.Now()

		for i := beginID; i < total+beginID; i++ {
			account := Account{
				ID:    i,
				Name:  fmt.Sprintf("Name %d", i),
				Phone: fmt.Sprintf("Phone %d", i),
			}
			query := fmt.Sprintf("INSERT INTO %s (id, name, phone) VALUES (%d, '%s', '%s')",
				tableName, account.ID, account.Name, account.Phone)
			result := db.Exec(query)
			if result.Error != nil {
				log.Printf("Failed to insert '%d th' record: %v", i, result.Error)
			}
			state.Insertion.CurrentTotal++
		}

		elapsed := time.Since(start)
		log.Infof("Inserted total %d records to '%s', ID '%d ~ %d' (elapsed: %s)",
			total, loc, beginID, beginID+total-1, elapsed)
		state.Insertion.Done <- true
		CountStreamRecords(total)
	}()
	return nil
}

func CountStreamRecords(num int) {
	streamRecordCount += num
}

func CompareRecords(sourceDB, targetDB *gorm.DB) (int, error) {
	var (
		limit    = 10000
		offset   = 0
		moreData = true

		lastMatchID = 0
	)

	for moreData {
		var (
			records1 []Account
			records2 []Account
		)

		err := sourceDB.Table("Accounts").Order("id ASC").Limit(limit).Offset(offset).Find(&records1).Error
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve records from source model: %v", err)
		}

		err = targetDB.Table("Accounts").Order("id ASC").Limit(limit).Offset(offset).Find(&records2).Error
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve records from target model: %v", err)
		}

		for i := range records1 {
			if records1[i].Name != records2[i].Name {
				return lastMatchID, fmt.Errorf("ID: %d source has '%s', target has '%s'", records1[i].ID, records1[i].Name, records2[i].Name)
			}

			if records1[i].Phone != records2[i].Phone {
				return lastMatchID, fmt.Errorf("ID: %d source has '%s', target has '%s'", records1[i].ID, records1[i].Phone, records2[i].Phone)
			}
			lastMatchID = records1[i].ID
		}

		offset += limit

		if len(records1) < limit {
			moreData = false
		}
	}

	return lastMatchID, nil
}

func VerifyFromToRowCountAndContentTimeoutSeconds(locTo, locFrom, tableName string, timeoutSec int) error {
	// compare source/target table row count
	sourceDB, err := utils.GetDBInstance(locFrom)
	if err != nil {
		return err
	}
	targetDB, err := utils.GetDBInstance(locTo)
	if err != nil {
		return err
	}
	srcRowCount, err := utils.GetCount(sourceDB, tableName)
	if err != nil {
		return err
	}

	// Check number of records
	var retry int
	var targetRowCount int64
	for retry = 0; retry < timeoutSec; retry++ {
		targetRowCount, err = utils.GetCount(targetDB, tableName)
		if err != nil {
			return err
		}
		log.Infof("Waiting for '%s' table '%s' to has %d records.. (%d sec), current total: %d",
			locTo, tableName, srcRowCount, retry, targetRowCount)
		if targetRowCount == srcRowCount {
			break
		} else if targetRowCount > srcRowCount {
			return fmt.Errorf("number of records in table '%s' is %d, expected %d after %d second",
				tableName, targetRowCount, srcRowCount, timeoutSec)
		}
		time.Sleep(1 * time.Second)
	}

	if retry == timeoutSec {
		return fmt.Errorf("number of records in table '%s' is %d, expected %d after %d second",
			tableName, targetRowCount, srcRowCount, timeoutSec)
	}

	for retry = 0; retry < timeoutSec; retry++ {
		lastMatchID, err := CompareRecords(sourceDB, targetDB)
		if err == nil {
			return nil
		}
		log.Infof("Waiting for '%s' table '%s' to has %d same content.. (%d sec), last match ID %d",
			locTo, tableName, srcRowCount, retry, lastMatchID)
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("content of table '%s' is not the same after %d second", tableName, timeoutSec)
}

func DockerComposeServiceIn(action, serviceName, executionMode string) error {
	if action != "start" && action != "stop" && action != "restart" {
		return fmt.Errorf("invalid docker-compose action '%s'", action)
	}

	if executionMode != "foreground" && executionMode != "background" {
		return fmt.Errorf("invalid docker compose execution mode '%s'", executionMode)
	}

	cmd := exec.Command("docker", "compose", "-f", utils.ConnectionConfig.DockerComposeFilePath, action, serviceName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	switch executionMode {
	case "foreground":
		// Execute the command in foreground and wait for its completion
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	case "background":
		// Execute the command in the background without waiting for its completion
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
	}

	log.Infof("docker compose create stdout: %s", stdout.String())
	log.Debugf("docker compose create stderr: %s", stderr.String())
	return nil
}

func ContainerStateWasTimeoutSeconds(ctName string, expectedState string, timeoutSec int) error {
	var (
		err   error
		i     int
		state *types.ContainerState
	)

	for i = 0; i < timeoutSec; i++ {
		state, err = GetContainerStateByName(ctName)
		if err == nil {
			log.Debugf("Container '%s' state: %s", ctName, state.Status)
			if expectedState == "exited" && state.Running == false {
				return nil
			}
			if expectedState == "running" && state.Running == true {
				return nil
			}
		} else {
			log.Errorf("failed to get container '%s' state: %s", ctName, err.Error())
		}
		if i%10 == 0 {
			log.Infof("Waiting for container '%s' state to be '%s'.. (%d sec)", ctName, expectedState, i)
		}
		time.Sleep(1 * time.Second)
	}

	if state != nil {
		log.Errorf("container '%s' state expected '%s', but got '%s' after %d sec", ctName, expectedState, state.Status, i)
	} else {
		log.Errorf("container '%s' state expected '%s', but not found after %d sec", ctName, expectedState, i)
	}
	return err
}

func UpdateRowDummyDataFromID(loc, tableName string, total, beginID int) error {
	db, err := utils.GetDBInstance(loc)
	if err != nil {
		return err
	}

	log.Infof("Updating %d records to '%s' - '%s'", total, loc, tableName)
	// Update dummy data
	opFailed := 0
	start := time.Now()

	for i := beginID; i < total+beginID; i++ {
		// Update the Name field
		err = db.Table("Accounts").Where("ID = ?", i).Update("Name", gorm.Expr("CONCAT(Name, ?)", " updated")).Error
		if err != nil {
			log.Errorf("failed to insert '%d th' record: %v", i, err)
			opFailed++
		}
	}

	elapsed := time.Since(start)
	log.Infof("Updated total %d records to '%s' table '%s', ID '%d ~ %d', failed count %d (elapsed: %s)",
		total, loc, tableName, beginID, beginID+total-1, opFailed, elapsed)
	return nil
}

func WaitForOperationDone(loc, tableName, op string, timeoutSec int) error {
	var doneChan chan bool
	var currentTotal *int

	switch op {
	case "insert":
		doneChan = state.Insertion.Done
		currentTotal = &state.Insertion.CurrentTotal
	case "update":
		doneChan = state.Update.Done
		currentTotal = &state.Update.CurrentTotal
	case "delete":
		doneChan = state.Delete.Done
		currentTotal = nil
	default:
		return fmt.Errorf("invalid operation '%s'", op)
	}

	for retry := 0; retry < timeoutSec; retry++ {
		select {
		case done := <-doneChan:
			if !done {
				return fmt.Errorf("'%s' table '%s' %s failed", loc, tableName, op)
			}
			log.Infof("'%s' table '%s' %s done.. (%d sec)", loc, tableName, op, retry)
			return nil
		default:
			log.Infof("Waiting for '%s' table '%s' %s done.. (%d sec), current total: %d", loc, tableName, op, retry, *currentTotal)
			time.Sleep(1 * time.Second)
		}
	}
	return fmt.Errorf("'%s' table '%s' update done timeout", loc, tableName)
}

func WaitForInsertionDone(loc, tableName string, timeoutSec int) error {
	if err := WaitForOperationDone(loc, tableName, "insert", timeoutSec); err != nil {
		return err
	}
	return nil
}

func WaitForDeleteDone(loc, tableName string, timeoutSec int) error {
	if err := WaitForOperationDone(loc, tableName, "delete", timeoutSec); err != nil {
		return err
	}
	return nil
}

func WaitForUpdateAndInsertDone(loc, tableName string, timeoutSec int) error {
	if err := WaitForOperationDone(loc, tableName, "update", timeoutSec); err != nil {
		return err
	}
	if err := WaitForOperationDone(loc, tableName, "insert", timeoutSec); err != nil {
		return err
	}
	return nil
}

func UpdateRowAndInsertDummyDataFromIDGoroutine(loc, tableName string, updateTotal, updatebeginID, insertionTotal, insertionBeginID int) error {
	InitUpdateState(updateTotal)
	InitInsertionState(insertionTotal)
	state.WG.Add(1)
	go func() {
		defer state.WG.Done()
		db, err := utils.GetDBInstance(loc)
		if err != nil {
			log.Error(err)
			state.Update.Done <- false
			return
		}

		log.Infof("Updating %d records to '%s' - '%s'", updateTotal, loc, tableName)
		// Update dummy data
		opFailed := 0
		start := time.Now()

		for i := updatebeginID; i < updateTotal+updatebeginID; i++ {
			// Update the Name field
			err = db.Table("Accounts").Where("ID = ?", i).Update("Name", gorm.Expr("CONCAT(Name, ?)", " updated")).Error
			if err != nil {
				log.Errorf("failed to insert '%d th' record: %v", i, err)
				opFailed++
			}
			state.Update.CurrentTotal++
		}

		elapsed := time.Since(start)
		log.Infof("Updated total %d records to '%s' table '%s', ID '%d ~ %d', failed count %d (elapsed: %s)",
			updateTotal, loc, tableName, updatebeginID, updatebeginID+updateTotal-1, opFailed, elapsed)
		state.Update.Done <- true
		CountStreamRecords(updateTotal)

		start = time.Now()
		for i := insertionBeginID; i < insertionTotal+insertionBeginID; i++ {
			account := Account{
				ID:    i,
				Name:  fmt.Sprintf("Name %d", i),
				Phone: fmt.Sprintf("Phone %d", i),
			}
			query := fmt.Sprintf("INSERT INTO %s (id, name, phone) VALUES (%d, '%s', '%s')",
				tableName, account.ID, account.Name, account.Phone)
			result := db.Exec(query)
			if result.Error != nil {
				log.Printf("Failed to insert '%d th' record: %v", i, result.Error)
			}
			state.Insertion.CurrentTotal++
		}
		elapsed = time.Since(start)
		log.Infof("Inserted total %d records to '%s', ID '%d ~ %d' (elapsed: %s)",
			insertionTotal, loc, insertionBeginID, insertionBeginID+insertionTotal-1, elapsed)

		state.Insertion.Done <- true
		CountStreamRecords(insertionTotal)
	}()
	return nil
}

func CheckNatsStreamMessages(timeoutSec int) error {
	nc, _ := nats.Connect(fmt.Sprintf("nats://%s:%d", utils.ConnectionConfig.Nats.Host, utils.ConnectionConfig.Nats.Port))
	defer nc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	js, _ := jetstream.New(nc)
	stream, _ := js.Stream(ctx, "GVT_default_DP_accounts")

	log.Infof("check nats DP stream")
	var retry int
	var msgs uint64
	for retry = 0; retry < timeoutSec; retry++ {
		info, err := stream.Info(ctx)
		if err != nil {
			return err
		}
		msgs = info.State.Msgs
		log.Infof("Waiting for DP stream to has %d records.. (%d sec), current total: %d",
			streamRecordCount, retry, msgs)
		if msgs == uint64(streamRecordCount) {
			break
		} else if msgs > uint64(streamRecordCount) {
			return fmt.Errorf("number of records in DP stream is %d, expected %d", msgs, streamRecordCount)
		}
		time.Sleep(1 * time.Second)
	}

	if retry == timeoutSec {
		return fmt.Errorf("number of records in DP stream is %d, expected %d after %d second", msgs, streamRecordCount, timeoutSec)
	}

	return nil
}

func CleanUpTable(loc, tableName string) error {
	db, err := utils.GetDBInstance(loc)
	if err != nil {
		return err
	}
	// Clean up
	result := db.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
	if result.Error != nil {
		return fmt.Errorf("failed to exec clean up '%s' table '%s': %v", loc, tableName, result.Error)
	}

	var rowCount int64
	db.Table(tableName).Count(&rowCount)
	if rowCount != 0 {
		return fmt.Errorf("failed to clean up '%s' table '%s'", loc, tableName)
	}
	return nil
}

func CleanUpTableGoroutine(loc, tableName string) error {
	InitDeleteState()
	state.WG.Add(1)
	go func() {
		defer state.WG.Done()
		db, err := utils.GetDBInstance(loc)
		if err != nil {
			log.Error(err)
			state.Delete.Done <- false
			return
		}
		var deleteRowCount int64
		db.Table(tableName).Count(&deleteRowCount)
		// Clean up
		result := db.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
		if result.Error != nil {
			log.Errorf("failed to exec clean up '%s' table '%s': %v", loc, tableName, result.Error)
			state.Delete.Done <- false
			return
		}

		var rowCount int64
		db.Table(tableName).Count(&rowCount)
		if rowCount != 0 {
			log.Errorf("failed to clean up '%s' table '%s'", loc, tableName)
			state.Delete.Done <- false
			return
		}
		state.Delete.Done <- true
		CountStreamRecords(int(deleteRowCount))
	}()
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		reset()
		if utils.IsSkipped(sc.Name) {
			return ctx, godog.ErrSkip
		}
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		state.WG.Wait()
		if err := utils.CloseAllServices(); err != nil {
			log.Errorf("failed to close all services: %v", err)
		}
		return ctx, nil
	})

	ctx.Given(`^Create all services$`, utils.CreateServices)
	ctx.Given(`^Close all services$`, utils.CloseAllServices)
	ctx.Given(`^Load the initial configuration file$`, utils.LoadConnectionConfig)
	ctx.Given(`^Start the "([^"]*)" service \(timeout "(\d+)"\)$`, utils.DockerComposeServiceStart)
	ctx.Given(`^Initialize the "([^"]*)" table Accounts$`, utils.DBServerInit)
	ctx.Given(`^Create Data Product Accounts$`, utils.CreateDataProduct)
	ctx.Given(`^Set up atomic flow document$`, utils.InitAtomicService)

	ctx.Then(`^"([^"]*)" table "([^"]*)" has "(\d+)" data \(timeout "([^"]*)"\)$`, VerifyRowCountTimeoutSeconds)
	ctx.Given(`^"([^"]*)" table "([^"]*)" inserted "([^"]*)" data \(starting ID "(\d+)"\)$`, InsertDummyDataFromID)
	ctx.Given(`^"([^"]*)" table "([^"]*)" continuously inserting "([^"]*)" data \(starting ID "(\d+)"\)$`, InsertDummyDataFromIDGoroutine)
	ctx.Then(`^"([^"]*)" has the same content as "([^"]*)" in "([^"]*)" \(timeout "([^"]*)"\)$`, VerifyFromToRowCountAndContentTimeoutSeconds)
	ctx.Given(`^docker compose "([^"]*)" service "([^"]*)" \(in "([^"]*)"\)$`, DockerComposeServiceIn)
	ctx.Then(`^container "([^"]*)" was "([^"]*)" \(timeout "(\d+)"\)$`, ContainerStateWasTimeoutSeconds)
	ctx.When(`^container "([^"]*)" ready \(timeout "(\d+)"\)$`, utils.ContainerAndProcessReadyTimeoutSeconds)
	ctx.Then(`wait for "([^"]*)" table "([^"]*)" insertion to complete \(timeout "([^"]*)"\)$`, WaitForInsertionDone)
	ctx.Then(`^Check the nats stream DP has correct messages \(timeout "([^"]*)"\)$`, CheckNatsStreamMessages)
	ctx.Then(`^Wait "([^"]*)" seconds$`, utils.WaitSeconds)

	ctx.Given(`^"([^"]*)" table "([^"]*)" updated "([^"]*)" data - appending suffix 'updated' to each Name field \(starting ID "(\d+)"\)$`, UpdateRowDummyDataFromID)
	ctx.Given(`^"([^"]*)" table "([^"]*)" continuously updating "([^"]*)" data - appending suffix 'updated' to each Name field \(starting ID "(\d+)"\) and inserting "([^"]*)" data \(starting ID "(\d+)"\)$`, UpdateRowAndInsertDummyDataFromIDGoroutine)
	ctx.Given(`^"([^"]*)" table "([^"]*)" cleared$`, CleanUpTable)
	ctx.Then(`^wait for "([^"]*)" table "([^"]*)" update and insertion to complete \(timeout "([^"]*)"\)$`, WaitForUpdateAndInsertDone)
	ctx.Given(`^"([^"]*)" table "([^"]*)" continuous cleanup$`, CleanUpTableGoroutine)
	ctx.Then(`^wait for "([^"]*)" table "([^"]*)" cleanup to complete \(timeout "([^"]*)"\)$`, WaitForDeleteDone)
}
