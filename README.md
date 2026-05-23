<div align="center">

# SamaSalaire Backend

**Payroll, leave, and HR management REST API — production-ready in Go**

[![CI](https://img.shields.io/github/actions/workflow/status/utachicodes/SamaSalaireBackend/ci.yml?branch=main&style=flat-square&label=CI)](https://github.com/utachicodes/SamaSalaireBackend/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-1.12-00ADD8?style=flat-square&logo=go&logoColor=white)](https://gin-gonic.com)
[![MongoDB](https://img.shields.io/badge/MongoDB-7.x-47A248?style=flat-square&logo=mongodb&logoColor=white)](https://mongodb.com)
[![Docker](https://img.shields.io/badge/Docker-multi--stage-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docker.com)
[![License](https://img.shields.io/badge/license-MIT-black?style=flat-square)](LICENSE)

</div>

---

SamaSalaire is a REST API for managing employee payroll, leave, and HR operations. Built with Go and Gin, it ships with JWT authentication, role-based access control, an immutable audit trail on every mutating request, and a layered architecture that keeps handlers, services, and data access cleanly separated.

---

## Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Error Responses](#error-responses)
- [Roles and Permissions](#roles-and-permissions)
- [Docker](#docker)
- [Testing](#testing)
- [Seeding](#seeding)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)

---

## Features

- REST API designed around resource-oriented routes
- JWT authentication with configurable token expiry
- Role-based access control — `admin`, `hr`, `manager`, `employee`
- Employee management with department and reporting hierarchy
- Salary component builder — per-employee earnings and deductions
- Payroll period workflow: create → run calculations → finalize
- Payslip generation and retrieval
- Leave types, leave balances, request submission, and approval workflow
- Payroll and leave summary reports for HR and admin
- Immutable audit log on all write operations (actor, timestamp, payload)
- CORS configured for local and production frontends
- Automatic MongoDB index creation on startup
- Multi-stage Docker build producing a minimal Alpine image

---

## Tech Stack

| Layer            | Technology                         |
|------------------|------------------------------------|
| Language         | Go 1.26                            |
| HTTP Framework   | Gin 1.12                           |
| Database         | MongoDB 7.x (mongo-driver v2)      |
| Authentication   | JWT — golang-jwt/jwt v5            |
| Configuration    | godotenv                           |
| Containerization| Docker — Alpine 3.21 runtime       |
| Logging         | Go standard `log` package          |

---

## Project Structure

```
.
├── cmd/
│   ├── server/          # Application entry point
│   └── seed/            # Database seeder (development)
├── internal/
│   ├── config/          # Environment variable loading
│   ├── database/        # MongoDB connection, indexes, collection names
│   ├── handlers/        # HTTP handlers, one file per domain
│   ├── middleware/      # Auth, RBAC, CORS, audit logging
│   ├── models/          # BSON/JSON data models
│   ├── router/          # Route registration
│   └── services/        # Business logic (payroll engine)
└── Dockerfile           # Multi-stage build → Alpine runtime image
```

---

## Getting Started

**Prerequisites:** Go 1.26+, MongoDB 7+

```bash
git clone https://github.com/utachicodes/SamaSalaireBackend.git
cd SamaSalaireBackend

cp .env.example .env
# edit .env with your values — at minimum, set a strong JWT_SECRET

go mod download
go run ./cmd/server
```

Or, with the Makefile:

```bash
make run
```

The API starts at `http://localhost:8080`. The port is configurable via the `PORT` environment variable.

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## Architecture

```
HTTP request
    │
    ▼
┌───────────────┐    ┌──────────────┐    ┌──────────────┐
│  middleware   │ →  │   handlers   │ →  │   services   │
│ auth/RBAC/CORS│    │  (cmd/server)│    │ (payroll, …) │
└───────────────┘    └──────────────┘    └──────────────┘
                                                │
                                                ▼
                                         ┌──────────────┐
                                         │   MongoDB    │
                                         └──────────────┘
```

Handlers are kept thin — request parsing, validation, and response writing only. All business rules live in `internal/services` or in domain methods on the models.

---

## Configuration

| Variable           | Default                       | Description                         |
|--------------------|-------------------------------|-------------------------------------|
| `PORT`             | `8080`                        | Port the HTTP server listens on     |
| `MONGODB_URI`      | `mongodb://localhost:27017`   | MongoDB connection string           |
| `DB_NAME`          | `samasalaire`                 | MongoDB database name               |
| `JWT_SECRET`       | `change-me` *(insecure)*      | Secret used to sign JWTs            |
| `JWT_EXPIRY_HOURS` | `24`                          | Token lifetime in hours             |

> Always override `JWT_SECRET` in production with a high-entropy value (at least 32 random bytes). Rotating it invalidates every existing token.

---

## API Reference

All authenticated routes are prefixed with `/api`. The health endpoint is unprefixed. Protected routes require the header:

```
Authorization: Bearer <token>
```

### Auth

| Method | Endpoint          | Access | Description                       |
|--------|-------------------|--------|-----------------------------------|
| POST   | `/auth/login`     | Public | Exchange credentials for a JWT    |
| POST   | `/auth/logout`    | Auth   | Invalidate the current session    |

<details>
<summary>Example: login</summary>

```http
POST /api/auth/login
Content-Type: application/json

{ "email": "admin@example.com", "password": "admin123" }
```

```json
{ "token": "<jwt>", "user": { "id": "...", "role": "admin", "email": "admin@example.com" } }
```

</details>

### Employees

| Method | Endpoint           | Access                  | Description          |
|--------|--------------------|-------------------------|----------------------|
| GET    | `/employees`       | admin, hr, manager      | List employees       |
| POST   | `/employees`       | admin, hr               | Create an employee   |
| GET    | `/employees/:id`   | Auth                    | Fetch employee by ID |
| PUT    | `/employees/:id`   | admin, hr               | Update an employee   |
| DELETE | `/employees/:id`   | admin                   | Delete an employee   |

### Salary Components

| Method | Endpoint                            | Access     | Description                     |
|--------|-------------------------------------|------------|---------------------------------|
| GET    | `/salary-components/:employeeId`    | admin, hr  | List components for an employee |
| POST   | `/salary-components`                | hr         | Add a salary component          |
| PUT    | `/salary-components/:id`            | hr         | Update a salary component       |
| DELETE | `/salary-components/:id`            | hr         | Remove a salary component       |

### Payroll

| Method | Endpoint                          | Access     | Description                  |
|--------|-----------------------------------|------------|------------------------------|
| GET    | `/payroll-periods`                | admin, hr  | List payroll periods         |
| POST   | `/payroll-periods`                | admin, hr  | Create a payroll period      |
| POST   | `/payroll-periods/:id/run`        | hr         | Run payroll calculations     |
| POST   | `/payroll-periods/:id/finalize`   | hr         | Finalize a period            |
| GET    | `/payslips`                       | Auth       | List payslips                |
| GET    | `/payslips/:id`                   | Auth       | Get payslip by ID            |

### Leave

| Method | Endpoint                          | Access                   | Description               |
|--------|-----------------------------------|--------------------------|---------------------------|
| GET    | `/leave-types`                    | Auth                     | List leave types          |
| POST   | `/leave-types`                    | admin                    | Create leave type         |
| PUT    | `/leave-types/:id`                | admin                    | Update leave type         |
| GET    | `/leave-balances/:employeeId`     | Auth                     | Get leave balance         |
| GET    | `/leave-requests`                 | Auth                     | List leave requests       |
| POST   | `/leave-requests`                 | Auth                     | Submit a leave request    |
| PUT    | `/leave-requests/:id/decide`      | admin, hr, manager       | Approve or reject         |

### Reports

| Method | Endpoint                     | Access     | Description             |
|--------|------------------------------|------------|-------------------------|
| GET    | `/reports/payroll-summary`   | admin, hr  | Payroll summary report  |
| GET    | `/reports/leave-summary`     | admin, hr  | Leave summary report    |

### Users

| Method | Endpoint        | Access | Description      |
|--------|-----------------|--------|------------------|
| GET    | `/users`        | admin  | List all users   |
| POST   | `/users`        | admin  | Create a user    |
| PUT    | `/users/:id`    | admin  | Update a user    |

### Health

```
GET /health  →  {"status":"ok"}
```

Use this endpoint for liveness probes — it does not touch the database.

---

## Error Responses

All error responses use a consistent JSON shape:

```json
{ "error": "human-readable description" }
```

Common HTTP statuses:

- `400` — invalid request body or query parameters
- `401` — missing or invalid `Authorization` header
- `403` — authenticated but lacks the required role
- `404` — resource not found
- `409` — conflict (e.g. duplicate email)
- `500` — internal server error

---

## Roles and Permissions

| Role       | Capabilities                                                                |
|------------|-----------------------------------------------------------------------------|
| `admin`    | Full access. User management, leave type configuration, all HR operations.  |
| `hr`       | Manage employees, salary components, payroll lifecycle, and leave types.    |
| `manager`  | View the employee directory and approve or reject leave requests.           |
| `employee` | View own payslips and leave balance; submit leave requests.                 |

---

## Docker

```bash
# Build (multi-stage; takes ~30s on a warm cache)
docker build -t samasalaire-backend .

# Run
docker run --rm -p 8080:8080 \
  -e MONGODB_URI=mongodb://host.docker.internal:27017 \
  -e JWT_SECRET=your-secret-here \
  samasalaire-backend
```

The multi-stage build keeps the final image small — only the compiled binary and CA certificates are included in the runtime layer. The image runs as a non-root user by default.

---

## Testing

```bash
make test
```

Runs `go test ./... -race -count=1`. Disable the race detector on platforms where it is unavailable by overriding the `test` target.

---

## Seeding

Populate the database with demo data for development:

```bash
go run ./cmd/seed
```

Or:

```bash
make seed
```

This creates a set of employees across all four roles along with salary components and leave types. Re-running the seeder is idempotent: existing accounts are skipped rather than duplicated.

---

## Contributing

Pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow and code style.

---

## Security

Found a vulnerability? Please do **not** open a public issue. See [SECURITY.md](SECURITY.md) for the disclosure process.

---

## License

Released under the [MIT License](LICENSE).
