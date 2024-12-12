# Fashionista Backend

Fashionista Backend is a RESTful API server built using Go. It supports user authentication, store and inventory management, and integration with Google Cloud Storage (GCS).

---

## TODO
### API Enhancements
- [ ] **Pagination**:
    - Add pagination support for `ListStores` and `ListInventories` endpoints.
- [ ] **Search and Filtering**:
    - Implement search and filter capabilities for inventory items based on criteria like `name` or `store_id`.

### Security
- [ ] **Rate Limiting**:
    - Protect API endpoints from abuse by adding rate-limiting middleware.
- [ ] **Token Management**:
    - Enhance JWT authentication by introducing refresh tokens.
- [ ] **Secure HTTP Headers**:
    - Configure secure HTTP headers to protect against common web vulnerabilities.

### Testing
- [ ] **Unit Tests**:
    - Write unit tests for core services and controllers.
- [ ] **Integration Tests**:
    - Add integration tests for major API workflows.
- [ ] **CI/CD**:
    - Set up a CI/CD pipeline for automated testing and deployment.

### Code Refactoring
- [ ] **Constants Management**:
    - Extract reusable constants into a dedicated `constants.go` file.
- [ ] **Error Handling**:
    - Standardize error handling and logging across the project.
- [ ] **Code Cleanup**:
    - Remove duplicate code and improve modularity in services and controllers.

### Documentation
- [ ] **API Documentation**:
    - Add OpenAPI/Swagger documentation for all endpoints.
- [ ] **Developer Guide**:
    - Create a detailed guide for contributors to understand the codebase and set up their environment.

### External Integrations
- [ ] **TryOn API**:
    - Extend the functionality to support multi-item try-on.
- [ ] **Image Processing**:
    - Add pre-upload image transformations like resizing and cropping.

### Deployment
- [ ] **Kubernetes Deployment**:
    - Create Kubernetes manifests for production-ready deployment.
- [ ] **Infrastructure as Code**:
    - Use Terraform to provision infrastructure.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
    - [1. Clone the Repository](#1-clone-the-repository)
    - [2. API Endpoints](#5-api-endpoints)
- [Testing](#testing)
- [License](#license)

---

## Features

- **User Authentication**: Registration and login with JWT.
- **Store Management**: Create and manage stores and their items.
- **Image Upload**: Upload and manage images with Google Cloud Storage.
- **Inventory Management**: Add, list, and delete inventory items.
- **External API Integration**: Integration with external APIs like TryOn.

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
│   ├── gcs.go           # GCS helper functions
│   ├── response.go      # Standardized API responses
│   ├── token.go         # JWT utilities
│   ├── valid.go         # Validation functions
├── .env                 # Environment variables
├── go.mod               # Go module file
├── LICENSE              # License information
├── smart-exchange-441906-p0-c0277d140202.json
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
go run app/main.go 

```
---

## Testing

You can test the API endpoints using tools like **Postman** or **cURL**. Make sure the backend is running locally or deployed to an accessible server.

---

## License

This project is licensed under the [MIT License](LICENSE).

---