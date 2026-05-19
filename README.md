<div align="center">

# SamaSalaire Backend

**Production-ready payroll and HR management API**

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-1.12-00ADD8?style=flat-square&logo=go&logoColor=white)](https://gin-gonic.com)
[![MongoDB](https://img.shields.io/badge/MongoDB-7.x-47A248?style=flat-square&logo=mongodb&logoColor=white)](https://mongodb.com)
[![Docker](https://img.shields.io/badge/Docker-multi--stage-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docker.com)
[![License](https://img.shields.io/badge/license-MIT-black?style=flat-square)](LICENSE)

</div>

---

SamaSalaire is a REST API for managing employee payroll, leave, and HR operations. Built with Go and Gin, it features JWT authentication, role-based access control, a full audit trail on every mutating action, and a clean layered architecture.

---

## Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Roles and Permissions](#roles-and-permissions)
- [Docker](#docker)
- [Seeding](#seeding)

---

## Features

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
- Multi-stage Docker build producing a minimal Alpine image

---

## Tech Stack

| Layer           | Technology                         |
|-----------------|------------------------------------|
| Language        | Go 1.26                            |
| HTTP Framework  | Gin 1.12                           |
| Database        | MongoDB 7.x (mongo-driver v2)      |
| Authentication  | JWT — golang-jwt/jwt v5            |
| Configuration   | godotenv                           |
| Containerization| Docker — Alpine 3.21 runtime       |

---

## Project Structure

```
.
├── cmd/
│   ├── server/          # Application entry point
│   └── seed/            # Database seeder (development)
└── internal/
    ├── config/          # Environment variable loading
    ├── database/        # MongoDB connection, indexes, collection names
    ├── handlers/        # HTTP handlers, one file per domain
    ├── middleware/      # Auth, RBAC, CORS, audit logging
    ├── models/          # BSON/JSON data models
    ├── router/          # Route registration
    └── services/        # Business logic (payroll engine)
```

---

## Getting Started

**Prerequisites:** Go 1.22+, MongoDB 6+

```bash
git clone https://github.com/utachicodes/SamaSalaireBackend.git
cd SamaSalaireBackend

cp .env.example .env
# edit .env with your values

go mod download
go run ./cmd/server/main.go
```

The API starts at `http://localhost:8080`.

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

---

## Configuration

| Variable           | Default                       | Description                         |
|--------------------|-------------------------------|-------------------------------------|
| `PORT`             | `8080`                        | Port the HTTP server listens on     |
| `MONGODB_URI`      | `mongodb://localhost:27017`   | MongoDB connection string           |
| `DB_NAME`          | `samasalaire`                 | MongoDB database name               |
| `JWT_SECRET`       | —                             | Secret used to sign JWTs            |
| `JWT_EXPIRY_HOURS` | `24`                          | Token lifetime in hours             |

---

## API Reference

All routes are prefixed with `/api`. Protected routes require the header:

```
Authorization: Bearer <token>
```

### Auth

| Method | Endpoint          | Access | Description         |
|--------|-------------------|--------|---------------------|
| POST   | `/auth/login`     | Public | Obtain a JWT        |
| POST   | `/auth/logout`    | Auth   | Invalidate session  |

### Employees

| Method | Endpoint           | Access                  | Description         |
|--------|--------------------|-------------------------|---------------------|
| GET    | `/employees`       | admin, hr, manager      | List all employees  |
| POST   | `/employees`       | admin, hr               | Create employee     |
| GET    | `/employees/:id`   | Auth                    | Get employee by ID  |
| PUT    | `/employees/:id`   | admin, hr               | Update employee     |
| DELETE | `/employees/:id`   | admin                   | Delete employee     |

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

---

## Roles and Permissions

| Role       | Capabilities                                                                |
|------------|-----------------------------------------------------------------------------|
| `admin`    | Full access. User management, leave type configuration, all HR operations.  |
| `hr`       | Manage employees, salary components, payroll lifecycle, and leave.          |
| `manager`  | View employee directory, approve or reject leave requests.                  |
| `employee` | View own payslips and leave balance, submit leave requests.                 |

---

## Docker

```bash
# Build
docker build -t samasalaire-backend .

# Run
docker run -p 8080:8080 \
  -e MONGODB_URI=mongodb://host.docker.internal:27017 \
  -e JWT_SECRET=your-secret-here \
  samasalaire-backend
```

The multi-stage build keeps the final image small — only the compiled binary and CA certificates are included in the runtime layer.

---

## Seeding

Populate the database with demo data for development:

```bash
go run ./cmd/seed/main.go
```

This creates a set of employees across all four roles along with salary components and leave types.

---

## License

MIT — [utachicodes](https://github.com/utachicodes)
