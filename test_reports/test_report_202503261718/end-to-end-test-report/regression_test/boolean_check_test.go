package boolcheck

import (
	"context"
	utils "e2e_test/testutils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/klauspost/compress/s2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/cucumber/godog"
	"github.com/nats-io/nats.go"
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

func TestFeature(t *testing.T) {
	utils.InitLog()

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type JSONData struct {
	Event   string `json:"event"`
	Payload string `json:"payload"`
}

type Products struct {
	ID       int       `json:"ID"`
	Name     string    `json:"Name"`
	Price    float64   `json:"Price"`
	Stock    int       `json:"Stock"`
	Obsolete bool      `json:"Obsolete"`
	ModTime  time.Time `json:"ModTime"`
}

var product Products

var Cmd *exec.Cmd

func InsertARecord(loc, tableName string) error {
	db, err := utils.GetDBInstance(loc)
	if err != nil {
		return err
	}
	log.Infof("Inserting a record to '%s' - '%s'",
		loc, tableName)

	query := "INSERT INTO Products (Name, Price, Stock, ModTime, Obsolete) VALUES ('pd1', 100, 100, GETDATE(), 0);"
	result := db.Exec(query)
	if result.Error != nil {
		log.Errorf("Failed to insert record: %v", result.Error)
	}

	return nil
}

func CheckDBData(loc string) error {
	db, err := utils.GetDBInstance(loc)
	if err != nil {
		return err
	}

	timeout := 10
	timeoutState := true
	for i := 0; i < timeout; i++ {
		time.Sleep(1 * time.Second)
		var count int64
		db.Table("Products").Count(&count)
		if count != 0 {
			timeoutState = false
			break
		}
	}

	if timeoutState {
		return fmt.Errorf("get %s data timeout", loc)
	}

	err = db.Table("Products").First(&product).Error
	if err != nil {
		return fmt.Errorf("failed to query Products table: %v", err)
	}

	if product.Obsolete {
		return fmt.Errorf("data in %s Obsolete is true", loc)
	}

	return nil
}

func CheckAtomicResult() error {
	data, err := os.ReadFile("./assets/atomic/check_result.json")
	if err != nil {
		return fmt.Errorf("check_result.json file not exist: \n%v", err)
	}
	var productContent Products
	err = json.Unmarshal(data, &productContent)
	if err != nil {
		return fmt.Errorf("failed to unmarshal check_result.json: \n%v", err)
	}
	if productContent.Obsolete == true {
		return fmt.Errorf("data in atomic Obsolete is true")
	}
	err = os.Remove("./assets/atomic/check_result.json")
	if err != nil {
		return fmt.Errorf("failed to remove check_result.json: \n%v", err)
	}
	return nil
}

func CheckNatsStreamResult() error {
	nc, _ := nats.Connect(fmt.Sprintf("nats://%s:%d", utils.ConnectionConfig.Nats.Host, utils.ConnectionConfig.Nats.Port))
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan *nats.Msg, 1)
	if _, err := js.ChanSubscribe("$GVT.default.EVENT.*", ch); err != nil {
		return fmt.Errorf("jetstream subscribe failed: %v", err)
	}

	var m *nats.Msg
	select {
	case m = <-ch:

	case <-time.After(30 * time.Second):
		return errors.New("subscribe out of time")
	}

	data := m.Data

	if m.Header.Get("Content-Encoding") == "s2" {
		decompressedMessage, err := s2.Decode(nil, m.Data)
		if err != nil {
			return fmt.Errorf("GVT_default s2 decompress failed: %v", err)
		}

		data = decompressedMessage
	}

	var jsonData JSONData
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("GVT_default json unmarshal failed: %v", err)
	}

	payload, err := utils.Base64ToString(jsonData.Payload)
	if err != nil {
		return fmt.Errorf("GVT_default base64 decode failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &result); err != nil {
		return fmt.Errorf("payload json unmarshal failed: %v", err)
	}

	obsolete, ok := result["Obsolete"].(bool)
	if !ok {
		return fmt.Errorf("Obsolete field not found or is not a boolean")
	}

	if obsolete {
		return fmt.Errorf("Obsolete is true")
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		if utils.IsSkipped(sc.Name) {
			return ctx, godog.ErrSkip
		}
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if err := utils.CloseAllServices(); err != nil {
			log.Errorf("failed to close all services: %v", err)
		}
		return ctx, nil
	})

	ctx.Given(`^Create all services$`, utils.CreateServices)
	ctx.Given(`^Load the initial configuration file$`, utils.LoadConnectionConfig)
	ctx.Given(`^Start the "([^"]*)" service \(timeout "(\d+)"\)$`, utils.DockerComposeServiceStart)
	ctx.Given(`^Initialize the "([^"]*)" table Products$`, utils.DBServerInit)
	ctx.Given(`^Create data product Products$`, utils.CreateDataProduct)
	ctx.Given(`^Set up atomic flow document$`, utils.InitAtomicService)

	ctx.Given(`^"([^"]*)" table "([^"]*)" inserted a record which has false boolean value$`, InsertARecord)
	ctx.Then(`^Check the "([^"]*)" table Products has a record with false value$`, CheckDBData)
	ctx.Then(`^Check the nats stream default domain has a record with false value$`, CheckNatsStreamResult)
	ctx.Then(`^Check the atomic has a record with false value$`, CheckAtomicResult)
	ctx.Then(`^Check the "([^"]*)" table Products has a record with false value$`, CheckDBData)

	ctx.Then(`^Wait "([^"]*)" seconds$`, utils.WaitSeconds)
}
