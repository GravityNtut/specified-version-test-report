package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ConnectionConfig struct {
	SourceDB                 serverInfo  `json:"source-mssql"`
	TargetDB                 serverInfo  `json:"target-mysql"`
	Nats                     serverInfo  `json:"nats"`
	DockerComposeFilePath    string      `json:"dockerComposeFilePath"`
	DockerComposeServiceName serviceName `json:"dockerComposeServiceNames"`
}

type serviceName struct {
	SourceMSSQL   string
	TargetMySQL   string
	Dispatcher    string
	Atomic        string
	Adapter       string
	NatsJetstream string
}

type serverInfo struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func LoadConnectionConfig() error {
	str, err := os.ReadFile("./connection_config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, &ConnectionConfig)
	if err != nil {
		return err
	}
	return nil
}

func InitLog() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func DatabaseLifeCheck(dialector gorm.Dialector, timeout int) (*gorm.DB, error) {
	for i := 0; i < timeout; i += 5 {
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Infof("Attempting to connect to %s Database (%d sec)", dialector.Name(), i)
			time.Sleep(5 * time.Second)
			continue
		}

		sqlDB, err := db.DB()
		sqlDB.SetMaxIdleConns(10)
		if err != nil {
			log.Infof("Failed to get sql.DB from gorm.DB: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := sqlDB.Ping(); err != nil {
			log.Infof("Waiting for %s Database to become available (%d sec)", dialector.Name(), i)
			time.Sleep(5 * time.Second)
			continue
		}
		return db, nil
	}
	return nil, fmt.Errorf("timeout connecting to the %s Database", dialector.Name())
}

func CreateTestDB(dialector gorm.Dialector, createTestDBFilePath string) error {
	db, err := DatabaseLifeCheck(dialector, 60)
	if err != nil {
		return err
	}
	str, err := os.ReadFile(createTestDBFilePath)
	if err != nil {
		return fmt.Errorf("failed to read create_test_db.sql: %v", err)
	}
	db.Exec(string(str))
	return nil
}

func ConnectToDB(s *serverInfo) (*gorm.DB, error) {
	if s.Type == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			s.Username, s.Password, s.Host, s.Port, s.Database)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %v", err)
		}
		return db, nil
	} else if s.Type == "mssql" {
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
			s.Username, s.Password, s.Host, s.Port, s.Database)
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %v", err)
		}
		return db, nil
	}
	return nil, fmt.Errorf("invalid database type '%s'", s.Type)
}

func GetDBInstance(loc string) (*gorm.DB, error) {
	switch loc {
	case ConnectionConfig.DockerComposeServiceName.SourceMSSQL:
		return ConnectToDB(&ConnectionConfig.SourceDB)
	case ConnectionConfig.DockerComposeServiceName.TargetMySQL:
		return ConnectToDB(&ConnectionConfig.TargetDB)
	default:
		return nil, fmt.Errorf("invalid database location '%s'", loc)
	}

}

func CheckProcessRunningInContainer(containerName, processName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			processes, err := cli.ContainerTop(context.Background(), container.ID, []string{})
			if err != nil {
				return err
			}
			for _, processInfo := range processes.Processes {
				lastElement := len(processInfo) - 1
				cmdLine := processInfo[lastElement]
				if cmdLine == "/"+processName || cmdLine == processName {
					return nil
				}

				if len(cmdLine) > 0 {
					cmds := strings.Split(cmdLine, " ")
					if len(cmds) > 0 {
						for _, cmd := range cmds {
							if cmd == processName || cmd == "/"+processName {
								return nil
							}
						}
					}
				}
			}
		}
	}

	return fmt.Errorf("process %s is not running in container %s", processName, containerName)
}

func ContainerLifeCheck(ctName, psName string, timeout int) error {
	for i := 0; i < timeout; i++ {
		err := CheckProcessRunningInContainer(ctName, psName)
		if err == nil {
			log.Infof("container '%s' is ready.. %d", ctName, i)
			return nil
		}
		if i%10 == 0 {
			log.Infof("Waiting for container '%s' to be ready.. (%d sec)", ctName, i)
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("the service '%s' timed out within %d seconds", ctName, timeout)

}

func NatsLifeCheck(timeout int) (*nats.Conn, error) {
	for i := 0; i < timeout; i++ {
		nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", ConnectionConfig.Nats.Host, ConnectionConfig.Nats.Port))
		if err != nil {
			log.Infoln("Unable to connect to the NATS server. Retry after 1 second")
			time.Sleep(1 * time.Second)
			continue
		}
		return nc, nil
	}
	return nil, fmt.Errorf("timeout connecting to NATS server")
}

func DockerComposeServiceStart(serviceName string, timeout int) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %v", err)
	}
	if err := cli.ContainerStart(context.Background(), serviceName, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container '%s': %v", serviceName, err)
	}
	return ContainerAndProcessReadyTimeoutSeconds(serviceName, timeout)
}

func ContainerAndProcessReadyTimeoutSeconds(ctName string, timeoutSec int) error {
	switch ctName {
	case ConnectionConfig.DockerComposeServiceName.Atomic:
		return ContainerLifeCheck(ctName, "node-red", timeoutSec)
	case ConnectionConfig.DockerComposeServiceName.Adapter:
		fallthrough
	case ConnectionConfig.DockerComposeServiceName.Dispatcher:
		return ContainerLifeCheck(ctName, ctName, timeoutSec)
	case ConnectionConfig.DockerComposeServiceName.NatsJetstream:
		nc, err := NatsLifeCheck(timeoutSec)
		if err != nil {
			return err
		}
		defer nc.Close()
		return nil
	case ConnectionConfig.DockerComposeServiceName.TargetMySQL:
		s := ConnectionConfig.TargetDB
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
			s.Username, s.Password, s.Host, s.Port)
		db := mysql.Open(dsn)
		_, err := DatabaseLifeCheck(db, timeoutSec)
		if err != nil {
			return err
		}
		return nil
	case ConnectionConfig.DockerComposeServiceName.SourceMSSQL:
		s := ConnectionConfig.SourceDB
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?encrypt=disable",
			s.Username, s.Password, s.Host, s.Port)
		db := sqlserver.Open(dsn)
		_, err := DatabaseLifeCheck(db, timeoutSec)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("invalid container name '%s'", ctName)
	}
}

func InitProductsTable(s *serverInfo, createTableFilePath string) error {
	var err error
	sourceDB, err := ConnectToDB(s)
	if err != nil {
		return fmt.Errorf("failed to connect to '%s' database: %v", s.Type, err)
	}

	db, err := sourceDB.DB()
	if err != nil {
		return fmt.Errorf("failed to connect to '%s' database: %v", s.Type, err)
	}
	str, err := os.ReadFile(createTableFilePath)
	if err != nil {
		return fmt.Errorf("failed to read create_table.sql: %v", err)
	}
	if _, err := db.Exec(string(str)); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

func InitAtomicService() error {
	token, err := GetToken()
	if err != nil {
		return err
	}
	inputFileName := "assets/unprocessed_cred.json"
	byteValue, err := os.ReadFile(inputFileName)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	var data map[string]map[string]string
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return fmt.Errorf("failed to parse JSON file: %v", err)
	}

	for _, component := range data {
		if _, exist := component["accessToken"]; exist {
			component["accessToken"] = token
		}
	}

	outputFileName := "tmp/unencrypted_cred.json"
	modifiedFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("failed to create output JSON file: %v", err)
	}

	defer func() {
		if err := modifiedFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	modifiedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal modified JSON: %v", err)
	}

	if _, err := modifiedFile.Write(modifiedJSON); err != nil {
		return fmt.Errorf("failed to write to %s JSON file: %v", outputFileName, err)
	}

	cmd := exec.Command("sh", "./assets/flowEnc.sh", outputFileName,
		"./assets/atomic", ">", "./assets/atomic/flows_cred.json")
	credFile, err := os.Create("./assets/atomic/flows_cred.json")
	if err != nil {
		return fmt.Errorf("failed to create flows_cred.json: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stdout = credFile
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute flowEnc.sh: %s", stderr.String())
	}
	return nil
}

func DBServerInit(dbStr string) error {
	var (
		dialector            gorm.Dialector
		createTestDBFilePath string
		serverInfo           *serverInfo
		createTableFilePath  string
	)

	switch dbStr {
	case ConnectionConfig.DockerComposeServiceName.SourceMSSQL:
		info := &ConnectionConfig.SourceDB
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?encrypt=disable",
			info.Username, info.Password, info.Host, info.Port)

		dialector = sqlserver.Open(dsn)
		createTestDBFilePath = "./assets/mssql/create_test_db.sql"
		serverInfo = &ConnectionConfig.SourceDB
		createTableFilePath = "./assets/mssql/create_table.sql"
	case ConnectionConfig.DockerComposeServiceName.TargetMySQL:
		info := &ConnectionConfig.TargetDB
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
			info.Username, info.Password, info.Host, info.Port)

		dialector = mysql.Open(dsn)
		createTestDBFilePath = "./assets/mysql/create_test_db.sql"

		serverInfo = &ConnectionConfig.TargetDB
		createTableFilePath = "./assets/mysql/create_table.sql"
	default:
		return fmt.Errorf("invalid database type '%s'", dbStr)
	}

	if err := CreateTestDB(dialector, createTestDBFilePath); err != nil {
		return err
	}
	if err := InitProductsTable(serverInfo, createTableFilePath); err != nil {
		return err
	}
	return nil
}

func GetCount(db *gorm.DB, tableName string) (int64, error) {
	var count int64
	err := db.Table(tableName).Count(&count).Error
	return count, err
}

func ExecuteContainerCommand(containerID string, cmd []string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("error creating Docker client: %v", err)
	}

	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("error creating exec instance: %v", err)

	}

	resp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("error attaching to exec instance: %v", err)
	}
	defer resp.Close()

	scanner := bufio.NewScanner(resp.Reader)
	result := ""
	for scanner.Scan() {
		result += scanner.Text() + "\n"
	}
	return result, nil
}

func GetToken() (string, error) {
	cmdString := []string{"/gravity-cli", "token", "create", "-s", ConnectionConfig.DockerComposeServiceName.NatsJetstream + ":4222"}
	result, err := ExecuteContainerCommand(ConnectionConfig.DockerComposeServiceName.Dispatcher, cmdString)
	if err != nil {
		return "", err
	}
	regexp := regexp.MustCompile(`Token: (.*)`)
	parts := regexp.FindStringSubmatch(result)
	if parts == nil {
		return "", fmt.Errorf("failed to get token from result: %s", result)
	}
	return parts[1], nil
}

func Base64ToString(base64Str string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func CreateServices() error {
	cmd := exec.Command("docker", "compose", "-f", ConnectionConfig.DockerComposeFilePath, "create")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("docker compose create fail: %v", err)
	}

	log.Infof("docker compose create stdout: %s", stdout.String())
	log.Debugf("docker compose create stderr: %s", stderr.String())
	return nil
}

func CloseAllServices() error {
	cmd := exec.Command("docker", "compose", "-f", ConnectionConfig.DockerComposeFilePath, "down")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	log.Infof("docker compose create stdout: %s", stdout.String())
	log.Debugf("docker compose create stderr: %s", stderr.String())
	return nil
}

func WaitSeconds(seconds int) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return nil
}

func CreateDataProduct() error {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%d", ConnectionConfig.Nats.Host, ConnectionConfig.Nats.Port))
	if err != nil {
		return err
	}
	defer nc.Close()
	containerID := ConnectionConfig.DockerComposeServiceName.Dispatcher

	cmd := []string{"sh", "/assets/dispatcher/create_product.sh"}
	result, err := ExecuteContainerCommand(containerID, cmd)
	if err != nil {
		return err
	}
	log.Infoln(result)
	return nil
}

func IsSkipped(s string) bool {
	re := regexp.MustCompile(`\[skipped\]`)
	return re.MatchString(s)
}
