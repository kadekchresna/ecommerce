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

## How to Run the Service locally

1. Makesure **Docker** and **minikube** is already setup on your local

2. Deploy the **payroll service** and PostgreSQL

```
    make up
```

3. Access the service through `:8081` with this [postman collection](https://github.com/kadekchresna/payroll/blob/master/payroll.postman_collection.json)
4. To connect to the deployed PostgreSQL database, you can use this credential

```
    HOST: localhost
    PORT: 5432
    USER: postgres
    PASS: secret
    DB: payroll_db
```

## How to Run Unit Coverage

1. Run unit test coverage on main functionality

```
make unit-test
```

## User Credential

1. Login as Employee

```
    "username": "user1"
    "password": "secret123"
```

2. Login as Admin

```
    "username": "admin"
    "password": "1234567890"
```

## Note

1. Each endpoint must have `Authorization` header in order to access the enpoint

```
Authorization:Bearer xxxxx.xxxxxxxxxxxxx.xxxxxxx
```

2. In order the employee to able to submit overtime, they have to clocked out first by submit attendance second time at the same specified date
