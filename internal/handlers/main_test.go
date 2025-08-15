package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"scorer/internal/client"
	metric "scorer/internal/metrics"
	"scorer/internal/repositories"
	"scorer/internal/services"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MainTest struct {
	app                 *fiber.App
	db                  *gorm.DB
	productRepo         repositories.ProductRepository
	merchantRepo        repositories.MerchantRepository
	merchantProductRepo repositories.MerchantProductRepository
	client              client.QueServiceClient
	service             services.ServicesFuncs
	handler             Handler
}

var MTest MainTest

func TestMain(m *testing.M) {
	MTest.setupApp()
	code := m.Run()
	os.Exit(code)
}

func (m *MainTest) setupApp() {
	m.app = fiber.New()
	m.db = m.createContainer(context.Background())
	m.productRepo = repositories.NewProductRepository(m.db)
	m.merchantRepo = repositories.NewMerchantRepository(m.db)
	m.merchantProductRepo = repositories.NewMerchantProductRepository(m.db)
	m.client = client.NewQueServiceClient("http://127.0.0.1:3000", nil)

	//metrics
	reg := prometheus.NewRegistry()
	metric := metric.NewMetric(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	m.service = services.NewServicesFuncs(m.productRepo, m.merchantRepo, m.merchantProductRepo, m.client, metric)
	m.handler = NewHandler(m.service)
	m.app.Post("/score", m.handler.HandleScore)
	
	// expose /metrics on the same Fiber app/port
	m.app.Get("/metrics", adaptor.HTTPHandler(promHandler))
}

func (m MainTest) createContainer(ctx context.Context) *gorm.DB {
	dbPass := "DbPass"
	dbUser := "DbUser"
	dbName := "DbName"
	var env = map[string]string{
		"POSTGRES_PASSWORD": dbPass,
		"POSTGRES_USER":     dbUser,
		"POSTGRES_DB":       dbName,
	}
	postgresPort := nat.Port("5432/tcp")

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{postgresPort.Port()},
			Env:          env,
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort(postgresPort),
			),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		panic(err)
	}

	p, err := container.MappedPort(ctx, postgresPort)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", "0.0.0.0", dbUser, dbPass, dbName, p.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile("./db.sql")
	if err != nil {
		panic(err)
	}

	db.Exec(string(content))

	return db
}

func (m *MainTest) sendRequest(app *fiber.App, t *testing.T, body []byte, apiKey string) (int, string) {
	req := httptest.NewRequest(http.MethodPost, "/score", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return resp.StatusCode, buf.String()
}
