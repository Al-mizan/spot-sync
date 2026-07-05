# 🚗 SpotSync - Smart Parking & EV Charging Reservation

<div align="center">
  <img src="https://img.shields.io/badge/go-v1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/echo-v4-00ADD8?style=for-the-badge&logo=go" alt="Echo Framework" />
  <img src="https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/gorm-v1-00ADD8?style=for-the-badge&logo=go" alt="GORM" />
  <img src="https://img.shields.io/badge/JWT-Auth-black?style=for-the-badge&logo=JSON%20web%20tokens" alt="JWT Authentication" />
</div>

<br />

> A centralized platform for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots.

## ✨ Features

- **Robust Authentication:** Secure JWT-based authentication using `bcrypt` for password hashing.
- **Role-Based Access Control (RBAC):** Distinct `driver` and `admin` roles with specific permissions.
- **Concurrency-Safe Reservations:** Implements GORM Database Transactions and Row-Level Locking (`FOR UPDATE`) to prevent race conditions and overbooking (the "EV Spot Bottleneck" problem).
- **Domain-Driven Design (DDD):** Clean, modular architecture separating layers (Delivery/HTTP, Use Case/Service, Domain, Repository).
- **Dynamic Capacity Calculation:** Real-time tracking of available parking spots based on active reservations.

## 🛠️ Technology Stack

- **Language:** Go (Golang) v1.22+
- **Web Framework:** Echo (`github.com/labstack/echo/v4`)
- **ORM:** GORM (`gorm.io/gorm`)
- **Database:** PostgreSQL (NeonDB / Supabase)
- **Validation:** Go Playground Validator (`github.com/go-playground/validator/v10`)
- **Authentication:** JWT (`github.com/golang-jwt/jwt/v5`)
- **Live Reloading:** Air (`github.com/air-verse/air`)

## 📋 Prerequisites

Before you begin, ensure you have met the following requirements:
- Go 1.22 or higher installed.
- PostgreSQL running locally or a remote connection string.
- `make` installed (for running Makefile commands).
- [Air](https://github.com/air-verse/air) installed (optional, for hot-reloading during development).

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-username/spotsync.git
cd spotsync
```

### 2. Set up environment variables

Copy the example environment file and configure it with your local/remote database credentials and JWT secret.

```bash
cp .env.example .env
```

Ensure your `.env` looks like this:
```env
DSN="host=localhost user=postgres password=postgres dbname=spotsync port=5432 sslmode=disable TimeZone=Asia/Dhaka"
PORT=8080
JWT_SECRET=your_jwt_secret_key
```

### 3. Run the application

You can run the application directly using the provided Makefile:

```bash
# To run the app locally
make run

# To build the executable
make build

# To run tests
make test

# To format and vet the code
make fmt
make vet
```

Alternatively, if you have `air` installed, you can simply run `air` in the root directory for hot-reloading.

## 🌐 API Overview

### Authentication
- `POST /api/v1/auth/register` - Register a new user (`driver` or `admin`)
- `POST /api/v1/auth/login` - Login and receive a JWT

### Parking Zones
- `POST /api/v1/zones` - Create a new parking zone (Admin only)
- `GET /api/v1/zones` - List all parking zones with available spots (Public)
- `GET /api/v1/zones/:id` - Get details of a single parking zone (Public)

### Reservations
- `POST /api/v1/reservations` - Reserve a parking spot (Auth required)
- `GET /api/v1/reservations/my-reservations` - View your own reservations (Auth required)
- `DELETE /api/v1/reservations/:id` - Cancel your reservation (Auth required)
- `GET /api/v1/reservations` - View all reservations (Admin only)

## 📁 Project Structure

```text
.
├── cmd/                # Main application entry point (main.go)
├── internal/           # Private application and library code
│   ├── apperror/       # Custom application error definitions
│   ├── auth/           # Authentication utilities (JWT, password hashing)
│   ├── config/         # Environment variables and configuration loading
│   ├── ctxkeys/        # Context keys for Echo middleware
│   ├── domain/         # Core business logic and entities
│   ├── httpresponse/   # Standardized HTTP response wrappers
│   ├── middlewares/    # Echo middlewares (auth, logging, etc.)
│   └── server/         # Server initialization and routing
├── Makefile            # Build and development scripts
├── .env.example        # Example environment variables
└── go.mod              # Go module dependencies
```

## 🔒 Security & Concurrency

SpotSync places a strong emphasis on data integrity and security:
- **Transaction Locks:** When a reservation is created, the system locks the specific parking zone row to verify capacity dynamically before inserting the reservation. This completely eliminates race conditions where multiple users try to book the last available spot simultaneously.
- **Data Protection:** Passwords are never logged or returned in API responses.

## 📄 License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
