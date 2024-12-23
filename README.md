# Fashionista Backend

Fashionista Backend is a RESTful API server built using Go. It supports user authentication, store and inventory management, and integration with Google Cloud Storage (GCS).

---
## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [1. Clone the Repository](#1-clone-the-repository)
  - [2. Configure Environment Variables](#2-configure-environment-variables)
  - [3. Run the Application with Docker](#3-run-the-application-with-docker)
  - [4. Access the API](#4-access-the-api)
- [Development Process](#development-process)
  - [Commit Standards](#commit-standards)
  - [Code Review](#code-review)
  - [Continuous Integration and Deployment](#continuous-integration-and-deployment)
- [Testing](#testing)
- [License](#license)
- [Notes](#notes)


## Features

1. **User Authentication**:
   1. Support for user registration and login using JWT.
   2. Use FireBase to verify via OTP
2. **Store Management**:
   1. Create and manage stores and their inventory.
3. **Inventory Management**:
   1. Add, list, and delete items from the inventory.
4. **Try-one**:
   1. Integration with Google Cloud Storage for image upload and management.
   2. Seamless integration with other , such as the TryOn API.
6. **CI/CD**:
   1. GoLint
   2. CD
---

## Project Structure

```plaintext
Fashionista/
├── app/
│   ├── main.go          # Entry point of the application
│   ├── routes.go        # Route definitions for public and protected APIs
├── config/
│   ├── config.go        # Centralized configuration
│   ├── db_config.go     # Database configuration
│   ├── gsc_config.go    # Google Cloud Storage (GCS) config
│   ├── jwt_config.go    # JWT-related configuration
│   ├── load_env.go      # Environment variable loader
│   ├── model_config.go  # Model-related config
│   ├── server_config.go # Server setup configuration
├── controllers/         # API controllers
├── docker/
│   ├── postgresql/
│   │   ├── init_db.sql  # Initial SQL setup for PostgreSQL
│   ├── docker-compose.yml
│   ├── init.sh          # Script for initializing the application with Docker
├── middlewares/
│   ├── auth.go          # Middleware for authentication
├── models/
│   ├── dto.go           # Data Transfer Objects
│   ├── service_account.go # GCS service account model
│   ├── setup.go         # Database and model initialization
├── run/
│   ├── start.sh         # Script to start the project
├── services/
│   ├── auth_service/
│   │   ├── auth.go      # User authentication service
│   ├── external/
│   │   ├── external.go  # External service integration
│   │   ├── gcp.go       # Google Cloud Platform utility functions
│   │   ├── page_service.go
│   │   ├── try_on_api.go
│   ├── user_service/
│       ├── user.go      # User-related business logic
├── templates/
│   ├── add_item.html    # Frontend template for adding items
│   ├── create_store.html
│   ├── home.html
├── utils/
│   ├── response.go      # Standardized API responses
│   ├── token.go         # JWT utilities
│   ├── valid.go         # Validation functions
├── .env                 # Environment variables
├── go.mod               # Go module file
├── LICENSE              # License information
├── file config Firebase and GCS
└── README.md            # Project documentation
```
## Prerequisites

- **Go**: Installed on your system ([Installation Guide](https://golang.org/doc/install)).
- **Docker**: Required for database setup.
- **Google Cloud**: A configured service account with storage permissions.

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Fashionista
File .env và smart-exchange connect to Kante
bash start_server.sh

```
---

## Testing

You can test the API endpoints using tools like **Postman** or **cURL**. Make sure the backend is running locally or deployed to an accessible server.

---
## Development Process

### Commit Standards
- Use **Conventional Commits** to ensure clear and consistent commit messages:
    - `feat`: Adding new features (e.g., `feat: add user authentication endpoint`).
    - `fix`: Fixing bugs (e.g., `fix: resolve token expiration issue`).
    - `docs`: Updating documentation (e.g., `docs: update README with API examples`).
    - `refactor`: Refactoring code without adding features or fixing bugs (e.g., `refactor: optimize database queries`).
    - `test`: Adding or updating tests (e.g., `test: add unit tests for user service`).
- Avoid vague commit messages like "update code" or "fix bug."

### Code Review


### Continuous Integration and Deployment
- **Pre-Merge Checks**:
    - PRs must pass all CI checks, including linting, testing, and build validation.
- **Development Environment Deployment**:
    - Merging into the `dev` branch triggers automatic deployment to the development environment via GitHub Actions.
- **Production Environment Deployment**:
    - Changes merged into the `main` branch are automatically deployed to production after passing all tests.
- Use **GitHub Actions** with workflows for:
    - Running linting checks (e.g., `golangci-lint`).
    - Running unit and integration tests.
    - Building and deploying Docker containers.

---

This process ensures high-quality code, smooth collaboration, and reliable deployments.
## License

This project is licensed under the [MIT License](LICENSE).

---