# eCommerce

This is a collective of backend services for managing product, shops, users, warehouses, and orders for ecommerce users. It is designed by implementing **Separation of Concerns** and **Domain-Driven Design (DDD)** principles to promote scalability, testability, and maintainability.

## Features

- Employee attendance tracking
- Overtime and reimbursement management
- Payslip and payroll generation
- JWT-based authentication
- User role seperation
- GORM-powered database access layer
- Unit test support with mocking

## Tech Stack

| Layer            | Tech/Library                                                                                  |
| ---------------- | --------------------------------------------------------------------------------------------- |
| Language         | Go (Golang)                                                                                   |
| Web Framework    | [Echo v4](https://echo.labstack.com)                                                          |
| ORM              | [GORM](https://gorm.io/)                                                                      |
| Testing          | [Testify](https://github.com/stretchr/testify) + [Mockery](https://github.com/vektra/mockery) |
| RDBMS            | PostgreSQL                                                                                    |
| NoSQL            | Redis                                                                                         |
| Message Queue    | Kafka                                                                                         |
| Auth             | JWT (JSON Web Tokens)                                                                         |
| Logging          | [Zap](https://github.com/uber-go/zap)                                                         |
| Containerization | Docker + k8s + Helm Chart                                                                     |

## Folder Structure

```bash
.
├── Makefile
├── README.md
├── infra
│   └── charts
├── order-service
│   ├── Dockerfile
│   ├── cmd                 # Main service entry point
│   ├── config              # Configuration service
│   ├── env.example
│   ├── go.mod
│   ├── go.sum
│   ├── helper/
│   └── infrastructure/      # driver and dependency
│       ├── db/
│       ├── db/
│       ├── kafka/
│       ├── lock/
│       └── messaging/
│   ├── internal/
│   └── v1/                   # core logic versioning
│       ├── delivery/
│       ├── dto/
│       ├── model/
│       ├── repository/
│       └── usecase/
│   └── v2/
│       ├── ...
│   ├── main.go
│   └── migration
├── product-service
│   └── ...
├── shop-service
│   └── ...
├── user-service
│   └── ...
└── warehouse-service
│   └── ...
```

This project dividing code into layers:

- Model (Domain) – pure business objects and logic
- DTO (Application) – for external data exchange
- Repository – actual DB queries using GORM
- Usecase – orchestrates business logic
- Delivery – handles routing and request parsing

## Prerequisites

- [Docker](https://www.docker.com/)
- [minikube](https://minikube.sigs.k8s.io/)
