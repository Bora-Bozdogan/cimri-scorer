# Cimri Scorer Microservice

This service receives product update events from merchants, calculates a score based on business rules, enqueues the scored requests into a queue service, and exposes Prometheus metrics for monitoring. It is part of the Cimri product update pipeline.

## High-Level Flow

- **HTTP server (Fiber):** starts a Fiber app and exposes REST endpoints.  
- **Config:** loads YAML (and env) configuration via Viper.  
- **Persistence:** connects to Postgres with GORM and wires repositories for Product, Merchant, and MerchantProduct entities.  
- **Services layer:** computes the score, enqueues the request into the Queue service, and updates Prometheus metrics.  

## HTTP Endpoints

### POST /score
Parses and validates the request, computes a score using the scoring algorithm, enqueues the request with the score to the queue service, and updates metrics.

### GET /metrics
Exposes Prometheus metrics via promhttp from the same Fiber app.

## Architecture

Merchant Update  
  -> [Fiber /score]  
      -> Body parse & Validate  
      -> services.CalculateScore(req)  
          -> Look up Product, Merchant, MerchantProduct in Postgres  
          -> Determine store tier & update type  
          -> Compute score (weighted sum of factors)  
      -> services.EnqueueRequest(req, score)  
          -> HTTP POST to QueueService /enqueue  
      -> Update Prometheus metrics  

Main wiring: HTTP app, DB connection (DSN from config), metrics registry, repositories, Redis client, Queue client, and services are created in `main.go`.  

Repositories: `ProductRepository`, `MerchantRepository`, and `MerchantProductRepository` read entities via GORM.  
Queue client: `QueServiceClient.EnqueueRequest` sends scored requests to the Queue service.  
Metrics: Gauges incremented for requests and valid requests.  

## Service Logic

- **CalculateScore(req):** assigns a numeric score using:  
  - Store tier (Amazon=10, Ebay=9, Trendyol=8, etc.)  
  - Update type (new item, price drop, stock change, description/image change, etc.)  
  - Product popularity and urgency  
  - (future) Time decay  

- **EnqueueRequest(req, score):** forwards the request to the Queue microservice.  
- **IncrementRequestCount, IncrementValidRequestCount:** update Prometheus gauges.  

## API

### POST /score

**Headers**  
`Content-Type: application/json`  

**Request body (example used in tests)**  
```json
{
  "ApiKey": "amazon-key",
  "ProductName": "iPhone 16",
  "ProductDescription": "Latest Apple flagship phone",
  "ProductImage": "https://example.com/iphone16.jpg",
  "StoreName": "Amazon",
  "Price": 1500,
  "Stock": 50,
  "PopularityScore": 5,
  "UrgencyScore": 5
}
```

**Responses**  
- `200 OK` — "600 sent to que" (score calculated and enqueued).  
- `400 Bad Request` — parse/validation errors.  

### GET /metrics

Prometheus metrics endpoint exposed via the same Fiber app.  

**Metrics:**  
- `scorer_requests_made` — total /score requests.  
- `scorer_valid_requests_made` — valid requests forwarded to the queue.  

## Development

### Requirements
- Go toolchain  
- Postgres (for local development)  
- Redis (for Queue client testing)  
- For tests: Docker running (Testcontainers spins up a `postgres:16-alpine` container)  

### Run
Set up `internal/config/config.yaml` (or environment variables).  
Start the app:  
```bash
go run ./...
```  
`main.go` builds the DSN from the config and opens the DB connection; the app listens on `server_params.listen_port`.  

### Test
The test suite exercises the `/score` handler with multiple scenarios (price drops, stock changes, new products, image changes).  
```bash
go test ./...
```  

## Notable Files
- `cmd/main.go` — app wiring.  
- `internal/handlers/handlers.go` — HTTP handler for /score.  
- `internal/services/calculate_score.go` — core scoring algorithm.  
- `internal/repositories/*.go` — repositories via GORM.  
- `internal/client/worker_service_client.go` — queue client (POST /enqueue).  
- `internal/metrics/metrics.go` — Prometheus metrics.  
- `internal/config/load_config.go` — Viper config loader.  
- `internal/redis_client/redis_client.go` — Redis-backed queue client.  
- `*_test.go` — handler tests with sample scoring scenarios.  

## Notes / Future Work
- Add time decay scoring factor (based on last queued event).  

## Tech
Built with **Fiber, GORM, Viper, Prometheus client, Redis, and Testcontainers**.
