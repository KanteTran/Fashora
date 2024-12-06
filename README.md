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

## 2.API Endpoints

The following are the available APIs provided by the Fashionista Backend, grouped by their functionalities.

---

### **Authentication APIs**

1. **Register a New User**
    - **Endpoint**: `POST /auth/register`
    - **Description**: Register a new user with phone number and password.
    - **Request Body**:
      ```json
      {
        "phone_number": "123456789",
        "password": "password123",
        "user_name": "John Doe",
        "birthday": "1990-01-01",
        "address": "123 Main Street",
        "device_id": "device123",
        "gender": 1
      }
      ```
    - **Response**:
      ```json
      {
        "success": true,
        "status": 201,
        "message": "User created successfully",
        "data": {
          "token": "JWT_TOKEN",
          "user": {
            "id": "1",
            "phone_number": "123456789"
          }
        }
      }
      ```

2. **Login**
    - **Endpoint**: `POST /auth/login`
    - **Description**: Log in with phone number and password to retrieve a JWT token.
    - **Request Body**:
      ```json
      {
        "phone_number": "123456789",
        "password": "password123"
      }
      ```
    - **Response**:
      ```json
      {
        "success": true,
        "status": 200,
        "message": "Login successful",
        "data": {
          "token": "JWT_TOKEN",
          "user": {
            "id": "1",
            "phone_number": "123456789"
          }
        }
      }
      ```

3. **Check Phone Number Existence**
    - **Endpoint**: `POST /auth/check_phone`
    - **Description**: Check if a phone number is already registered.
    - **Request Body**:
      ```json
      {
        "phone_number": "123456789"
      }
      ```
    - **Response**:
      ```json
      {
        "success": true,
        "status": 200,
        "message": "Phone number exists",
        "data": {
          "exists": true
        }
      }
      ```

---

### **Image APIs**

1. **Upload Image**
    - **Endpoint**: `POST /image/push`
    - **Description**: Upload an image to Google Cloud Storage (GCS).
    - **Request**:
        - Content-Type: `multipart/form-data`
        - Form Field: `image`
    - **Response**:
      ```json
      {
        "message": "Image uploaded successfully",
        "url": "https://storage.googleapis.com/<bucket_name>/images/<filename>"
      }
      ```

2. **Get Signed Image URL**
    - **Endpoint**: `GET /image/get`
    - **Description**: Generate a signed URL for an uploaded image.
    - **Query Parameter**:
        - `filename`: The image's GCS object name.
    - **Response**:
      ```json
      {
        "message": "Image URL generated successfully",
        "url": "https://storage.googleapis.com/<bucket_name>/signed-url"
      }
      ```

---

### **Store APIs**

1. **Create a Store**
    - **Endpoint**: `POST /stores/create-store`
    - **Description**: Create a new store with basic details.
    - **Request Body**:
      ```json
      {
        "phone": "123456789",
        "store_name": "Fashionista Store",
        "address": "123 Main Street",
        "description": "A cool store for fashionistas"
      }
      ```
    - **Response**:
      ```json
      {
        "success": true,
        "status": 201,
        "message": "Store created successfully",
        "data": {
          "id": "store_id"
        }
      }
      ```

2. **List All Stores**
    - **Endpoint**: `GET /stores/list-all-store`
    - **Description**: Retrieve a list of all stores.
    - **Response**:
      ```json
      {
        "success": true,
        "status": 200,
        "message": "Stores fetched successfully",
        "data": [
          {
            "id": "store_id",
            "store_name": "Fashionista Store",
            "address": "123 Main Street",
            "description": "A cool store for fashionistas",
            "url_image": "https://storage.googleapis.com/<bucket_name>/store_image.jpg"
          }
        ]
      }
      ```

---

### **Inventory APIs**

1. **Add Inventory Item**
    - **Endpoint**: `POST /inventory/add-item`
    - **Description**: Add a new inventory item to a store.
    - **Request Body**:
      ```json
      {
        "store_id": "store_id",
        "name": "Fashionable Dress",
        "url": "https://fashionista.com/products/123",
        "image_url": "https://storage.googleapis.com/<bucket_name>/dress_image.jpg",
        "user_id": "user_id"
      }
      ```
    - **Response**:
      ```json
      {
        "success": true,
        "status": 201,
        "message": "Item added successfully",
        "data": {
          "item_id": "inventory_item_id"
        }
      }
      ```

2. **List All Inventory Items**
    - **Endpoint**: `POST /inventory/all-items`
    - **Description**: List all inventory items for a specific user.
    - **Request Body**:
      ```json
      {
        "user_id": "user_id"
      }
      ```
    - **Response**:
      ```json
      {
        "success": true,
        "status": 200,
        "message": "Items fetched successfully",
        "data": [
          {
            "id": "inventory_item_id",
            "store_id": "store_id",
            "name": "Fashionable Dress",
            "url": "https://fashionista.com/products/123",
            "image_url": "https://storage.googleapis.com/<bucket_name>/dress_image.jpg",
            "user_id": "user_id"
          }
        ]
      }
      ```

3. **Delete Inventory Item**
    - **Endpoint**: `DELETE /inventory/del-item`
    - **Description**: Delete an inventory item by its ID.
    - **Request Body**:
      ```json
      {
        "id": "inventory_item_id"
      }
      ```
    - **Response**:
      ```json
      {
        "success": true,
        "status": 200,
        "message": "Item deleted successfully",
        "data": null
      }
      ```

---

## Testing

You can test the API endpoints using tools like **Postman** or **cURL**. Make sure the backend is running locally or deployed to an accessible server.

---

## License

This project is licensed under the [MIT License](LICENSE).

---