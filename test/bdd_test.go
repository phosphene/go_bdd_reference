package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	appHost        string
	appPort        string
	dbConn         *sql.DB
	responseBody   string
	responseStatus int
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 1. Initialize Compose
	composeStack, err := compose.NewDockerCompose("testdata/docker-compose.yml")
	if err != nil {
		fmt.Printf("could not create compose stack: %v\n", err)
		os.Exit(1)
	}

	// 2. Configure Wait Strategy
	composeStack.WaitForService("app", wait.ForHTTP("/health").WithPort("8080/tcp"))

	// 3. Start Stack
	if err := composeStack.Up(ctx, compose.Wait(true)); err != nil {
		fmt.Printf("could not start compose stack: %v\n", err)
		os.Exit(1)
	}

	// 4. Discover Dynamic Port
	container, err := composeStack.ServiceContainer(ctx, "app")
	if err != nil {
		fmt.Printf("could not get app container: %v\n", err)
		_ = composeStack.Down(ctx, compose.RemoveOrphans(true))
		os.Exit(1)
	}

	host, err := container.Host(ctx)
	if err != nil {
		fmt.Printf("could not get app host: %v\n", err)
		_ = composeStack.Down(ctx, compose.RemoveOrphans(true))
		os.Exit(1)
	}

	mappedPort, err := container.MappedPort(ctx, "8080/tcp")
	if err != nil {
		fmt.Printf("could not get app mapped port: %v\n", err)
		_ = composeStack.Down(ctx, compose.RemoveOrphans(true))
		os.Exit(1)
	}

	appHost = host
	appPort = mappedPort.Port()

	// 5. Discover Postgres and Connect
	dbContainer, err := composeStack.ServiceContainer(ctx, "db")
	if err != nil {
		fmt.Printf("could not get db container: %v\n", err)
		_ = composeStack.Down(ctx, compose.RemoveOrphans(true))
		os.Exit(1)
	}

	dbHost, err := dbContainer.Host(ctx)
	if err != nil {
		fmt.Printf("could not get db host: %v\n", err)
		_ = composeStack.Down(ctx, compose.RemoveOrphans(true))
		os.Exit(1)
	}

	dbMappedPort, err := dbContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		fmt.Printf("could not get db mapped port: %v\n", err)
		_ = composeStack.Down(ctx, compose.RemoveOrphans(true))
		os.Exit(1)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=user password=password dbname=userdb sslmode=disable",
		dbHost, dbMappedPort.Port())
	dbConn, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("could not connect to db from test runner: %v\n", err)
		_ = composeStack.Down(ctx, compose.RemoveOrphans(true))
		os.Exit(1)
	}
	defer dbConn.Close()

	fmt.Printf("App running at %s:%s, DB at %s:%s\n", appHost, appPort, dbHost, dbMappedPort.Port())

	// 6. Run Tests
	exitCode := m.Run()

	// 6. Cleanup
	if err := composeStack.Down(ctx, compose.RemoveOrphans(true)); err != nil {
		fmt.Printf("could not tear down compose stack: %v\n", err)
	}

	os.Exit(exitCode)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func theApplicationIsRunning() error {
	if appHost == "" || appPort == "" {
		return fmt.Errorf("application is not running")
	}
	return nil
}

func iRequestTheHealthStatus() error {
	url := fmt.Sprintf("http://%s:%s/health", appHost, appPort)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseStatus = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	responseBody = string(body)
	return nil
}

func iRequestAnUnknownRoute() error {
	url := fmt.Sprintf("http://%s:%s/unknown", appHost, appPort)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseStatus = resp.StatusCode
	return nil
}

func theResponseStatusShouldBe(expected int) error {
	if responseStatus != expected {
		return fmt.Errorf("expected status %d, but got %d", expected, responseStatus)
	}
	return nil
}

func theResponseShouldBe(expected string) error {
	if responseBody != expected {
		return fmt.Errorf("expected response %q, but got %q", expected, responseBody)
	}
	return nil
}

func theRegistrationServiceIsAvailable() error {
	if appHost == "" || appPort == "" || dbConn == nil {
		return fmt.Errorf("registration service dependencies not found")
	}
	return dbConn.Ping()
}

func iRegisterAUser(email, name string) error {
	url := fmt.Sprintf("http://%s:%s/users", appHost, appPort)
	payload := map[string]string{
		"email": email,
		"name":  name,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseStatus = resp.StatusCode
	respBody, _ := io.ReadAll(resp.Body)
	responseBody = string(respBody)
	return nil
}

func theUserShouldExistInTheDatabase(email string) error {
	var count int
	query := "SELECT count(*) FROM users WHERE email = $1"
	err := dbConn.QueryRow(query, email).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("user with email %q not found in database", email)
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^the application is running$`, theApplicationIsRunning)
	sc.Step(`^the registration service is available$`, theRegistrationServiceIsAvailable)
	sc.Step(`^I request the health status$`, iRequestTheHealthStatus)
	sc.Step(`^I request an unknown route$`, iRequestAnUnknownRoute)
	sc.Step(`^I register a user with email "([^"]*)" and name "([^"]*)"$`, iRegisterAUser)
	sc.Step(`^the response status should be (\d+)$`, theResponseStatusShouldBe)
	sc.Step(`^the response should be "([^"]*)"$`, theResponseShouldBe)
	sc.Step(`^the user "([^"]*)" should exist in the database$`, theUserShouldExistInTheDatabase)
}
