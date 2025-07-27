# 💸 Money Transfer API

A lightweight, containerized REST API for transferring money between users, featuring logging, monitoring, graceful shutdown, and queue emulation.

---

## 🚀 Features

- ✅ **Money transfers between users**
- ✅ **User balance retrieval**
- 🔁 **Queue emulation for background processing**
- 📦 **PostgreSQL with auto migrations**
- 📈 **Prometheus + Grafana monitoring**
- 📘 **Swagger UI for API documentation**
- 🪵 **Structured logging with custom logger**
- 🛑 **Graceful shutdown**
- 🧪 **Unit tests for repository, service, and handler layers**

---

## 📁 Project Structure

- `api/handler` – HTTP controllers
- `internal/domain` – DTOs, models, contracts, and business logic
- `internal/repository` – PostgreSQL repositories
- `internal/queue` – Kafka-like job queue simulation
- `pkg/logger` – centralized logger
- `pkg/metrics` – Prometheus middleware
- `cmd/main.go` – entrypoint with graceful shutdown and routing

---

## 🔧 Endpoints

| Method | Path                        | Description                    |
|--------|-----------------------------|--------------------------------|
| GET    | `/transfers/{userId}`       | Get all transactions by user  |
| POST   | `/transfers`                | Create a new money transfer   |
| GET    | `/balance/{userId}`         | Get balance for a specific user |
| GET    | `/swagger/index.html`       | Swagger UI                    |
| GET    | `/metrics`                  | Prometheus metrics            |

---

## 🧪 Technologies

- Go 1.23
- PostgreSQL 14
- Prometheus + Grafana
- Gorilla Mux
- Docker & Docker Compose
- Swagger (via swaggo/http-swagger)

---

## 🐳 Docker Compose Setup

```
docker-compose up --build
```
#### This will start:

🟦 API on http://localhost:8080

🟩 Swagger UI on http://localhost:8080/swagger/index.html

🟥 Prometheus on http://localhost:9090

🟨 Grafana on http://localhost:3000 (login: admin / admin)

#### When Grafana starts:

- It will ask you to change the password.
- Then go to "Add Data Source" → Prometheus → Set URL to http://money_transfer_prometheus:9090
- Small dashboard is located in 'Dashboards' section 

---

## 🛠️ Environment Configuration

Create a **.env** file in root:

```
POSTGRES_HOST=money_transfer_db
POSTGRES_PORT=5432
POSTGRES_DB=money_transfer
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
SERVER_PORT=8080
```

---

## 🧹 Graceful Shutdown

The server catches OS signals (SIGINT, SIGTERM) and shuts down:

- Closes DB connection
- Waits for ongoing requests
- Logs shutdown event

---

## ✅ Unit Tests

Test coverage for:

- transfer_repository, user_repository
- transfer_service, user_service
- transfer_controller, user_controller

Once you're ready to test:

```
go test ./...
```

---

## 🛣️ Future Improvements

- Kafka integration (real message broker)
- Rate limiting & retries
- Authentication and authorization

---

## 📌 Notes

Swagger docs are generated with swag:

```
swag init
```

The database is initialized using migrations in:

```
./migrations/init.sql
```

---

#### 🧑‍💻 Author
Made with ❤️ by Dora as part of a pet project