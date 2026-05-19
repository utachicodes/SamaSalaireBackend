<div align="center">

# SamaSalaire Backend

**Modern payroll management API built for West African businesses**

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com)
[![MongoDB](https://img.shields.io/badge/MongoDB-Database-47A248?style=for-the-badge&logo=mongodb&logoColor=white)](https://mongodb.com)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)](LICENSE)

</div>

---

## Overview

SamaSalaire is a production-ready REST API for managing employee payroll, leave, and HR operations. Built with Go and Gin, it features JWT authentication, role-based access control, full audit logging, and a clean layered architecture.

---

## Features

- **JWT Authentication** — secure login/logout with configurable token expiry
- **Role-Based Access Control** — four roles: `admin`, `hr`, `manager`, `employee`
- **Employee Management** — full CRUD with department and hierarchy support
- **Payroll Engine** — create periods, run payroll calculations, generate payslips, finalize periods
- **Salary Components** — configurable allowances and deductions per employee
- **Leave Management** — leave types, balances, requests, and approval workflows
- **Reports** — payroll summaries and leave summaries for HR and admin
- **Audit Logging** — every mutating action is logged with actor, timestamp, and payload
- **Docker Ready** — multi-stage build with minimal Alpine runtime image

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26 |
| HTTP Framework | Gin |
| Database | MongoDB (via mongo-driver v2) |
| Auth | JWT (golang-jwt/jwt v5) |
| Config | godotenv |
| Containerization | Docker (multi-stage) |

---

## Project Structure

```
samasalaire-backend/
├── cmd/
│   ├── server/         # Application entrypoint
│   └── seed/           # Database seeder (dev)
└── internal/
    ├── config/         # Environment configuration
    ├── database/       # MongoDB connection, indexes, collections
    ├── handlers/       # HTTP request handlers
    ├── middleware/     # Auth, RBAC, CORS, audit logging
    ├── models/         # Data models
    ├── router/         # Route definitions
    └── services/       # Business logic
```

---

## Getting Started

### Prerequisites

- Go 1.26+
- MongoDB 6+
- Docker (optional)

### Run Locally

```bash
# Clone
git clone https://github.com/utachicodes/SamaSalaireBackend.git
cd SamaSalaireBackend

# Configure
cp .env.example .env
# Edit .env with your values

# Install dependencies
go mod download

# Start
go run ./cmd/server/main.go
```

### Run with Docker

```bash
docker build -t samasalaire-backend .
docker run -p 8080:8080 --env-file .env samasalaire-backend
```

### Seed the Database

```bash
go run ./cmd/seed/main.go
```

This creates demo employees across `admin`, `hr`, `manager`, and `employee` roles.

---

## Configuration

Copy `.env.example` to `.env` and set the following:

| Variable | Description | Default |
|---|---|---|
| `PORT` | HTTP server port | `8080` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DB_NAME` | Database name | `samasalaire` |
| `JWT_SECRET` | Secret key for signing JWTs | — |
| `JWT_EXPIRY_HOURS` | Token lifetime in hours | `24` |

---

## API Reference

All routes are prefixed with `/api`. Protected routes require `Authorization: Bearer <token>`.

### Auth

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/api/auth/login` | Public | Obtain JWT token |
| `POST` | `/api/auth/logout` | Auth | Invalidate session |

### Employees

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/api/employees` | HR / Admin / Manager | List all employees |
| `POST` | `/api/employees` | HR / Admin | Create employee |
| `GET` | `/api/employees/:id` | Auth | Get single employee |
| `PUT` | `/api/employees/:id` | HR / Admin | Update employee |
| `DELETE` | `/api/employees/:id` | Admin | Delete employee |

### Salary Components

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/api/salary-components/:employeeId` | HR / Admin | List components |
| `POST` | `/api/salary-components` | HR | Add component |
| `PUT` | `/api/salary-components/:id` | HR | Update component |
| `DELETE` | `/api/salary-components/:id` | HR | Remove component |

### Payroll

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/api/payroll-periods` | HR / Admin | List periods |
| `POST` | `/api/payroll-periods` | HR / Admin | Create period |
| `POST` | `/api/payroll-periods/:id/run` | HR | Run payroll calculations |
| `POST` | `/api/payroll-periods/:id/finalize` | HR | Finalize period |
| `GET` | `/api/payslips` | Auth | List payslips |
| `GET` | `/api/payslips/:id` | Auth | Get payslip |

### Leave

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/api/leave-types` | Auth | List leave types |
| `POST` | `/api/leave-types` | Admin | Create leave type |
| `PUT` | `/api/leave-types/:id` | Admin | Update leave type |
| `GET` | `/api/leave-balances/:employeeId` | Auth | Get leave balance |
| `GET` | `/api/leave-requests` | Auth | List leave requests |
| `POST` | `/api/leave-requests` | Auth | Submit leave request |
| `PUT` | `/api/leave-requests/:id/decide` | HR / Admin / Manager | Approve or reject |

### Reports

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/api/reports/payroll-summary` | HR / Admin | Payroll summary report |
| `GET` | `/api/reports/leave-summary` | HR / Admin | Leave summary report |

### Users

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/api/users` | Admin | List users |
| `POST` | `/api/users` | Admin | Create user |
| `PUT` | `/api/users/:id` | Admin | Update user |

### Health

```
GET /health  →  { "status": "ok" }
```

---

## Roles & Permissions

| Role | Capabilities |
|---|---|
| `admin` | Full access to all resources including user management and destructive operations |
| `hr` | Manage employees, salary components, payroll, leave types and requests |
| `manager` | View employees, approve/reject leave requests |
| `employee` | View own data, submit leave requests |

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes
4. Open a pull request

---

## License

MIT © [utachicodes](https://github.com/utachicodes)
