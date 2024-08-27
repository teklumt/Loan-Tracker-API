# 📑 Loan Tracker API - Documentation

The Loan Tracker API is a backend service developed in Golang using the Gin framework. It enables users to apply for loans and provides admin functionalities for managing users and loans. The API follows clean architecture principles and uses MongoDB as the database. Below is the documentation of the API endpoints, including user and admin functionalities, as well as the environmental setup required for development and deployment.

## 🌐 Overview

The Loan Tracker API allows users to register, log in, apply for loans, and manage their accounts. Admins can manage users and loans, ensuring efficient handling of data and system settings. The API is designed to be secure, efficient, and scalable.

## ⚙️ Environmental Setup

### 1. System Requirements

- **Golang**: Version 1.18+
- **MongoDB**: Version 4.4+
- **Docker**: (For containerized deployment)
- **Make**: (For build automation)

### 2. Environment Variables

You need to configure the following environment variables in a `.env` file at the root of your project. These are essential for connecting to the database, handling JWT authentication, and setting up server configuration.

```plaintext
# Server Configuration
SERVER_PORT=8080
SERVER_ENV=development

# Database Configuration
DB_URI=mongodb://localhost:27017
DB_NAME=loan_tracker

# JWT Configuration
JWT_SECRET=mysecretkey
JWT_EXPIRATION_HOURS=24

# Admin Credentials (For Admin Access)
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=adminpassword

# Email Service Configuration (Optional, if using email verification and password reset)
EMAIL_SERVICE=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=email@example.com
EMAIL_PASSWORD=your-email-password
```

### 3. Docker Setup

A `Dockerfile` and `docker-compose.yml` are provided to easily set up the Loan Tracker API in a containerized environment. To build and run the application using Docker, follow these steps:

1. **Build the Docker Image:**

   ```bash
   docker build -t loan-tracker-api .
   ```

2. **Run the Application with Docker Compose:**

   ```bash
   docker-compose up -d
   ```

   This will start both the API and MongoDB in Docker containers.

### 4. Running Locally

To run the API locally on your machine:

1. **Install Dependencies:**

   ```bash
   go mod tidy
   ```

2. **Run the Application:**

   ```bash
   go run main.go
   ```

3. The API should now be running at `http://localhost:8080`.

### 5. Testing

- **Unit Tests**: Run the following command to execute unit tests:

  ```bash
  go test ./...
  ```

## 🔗 Base URL

- **Production:** `https://api.loantracker.com`
- **Development:** `http://localhost:8080`

## 🛠 Endpoints

### 👤 1. User Management

#### 📝 User Registration

- **Endpoint:** `POST /users/register`
- **Description:** Register a new user with email, password, and profile details.
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }
  ```
- **Response:**
  - **Status Code:** `201 Created`
  - **Body:**
  ```json
  {
    "status_code": 201,
    "message": "Registration successful. Please verify your email."
  }
  ```

#### ✅ Email Verification

- **Endpoint:** `GET /users/verify-email`
- **Description:** Verify the user's email address using a token sent to their email.
- **Query Parameters:**
  - `token`: Verification token sent via email
  - `email`: User's email address
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
  ```json
  {
    "status_code": 200,
    "message": "Email verified successfully."
  }
  ```

#### 🔑 User Login

- **Endpoint:** `POST /users/login`
- **Description:** Authenticate user and provide access and refresh tokens.
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
  ```json
  {
    "status_code": 200,
    "access_token": "jwt-access-token",
    "refresh_token": "jwt-refresh-token"
  }
  ```

#### 🔄 Token Refresh

- **Endpoint:** `POST /users/token/refresh`
- **Description:** Refresh access token using the refresh token.
- **Request Body:**
  ```json
  {
    "refresh_token": "jwt-refresh-token"
  }
  ```
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
  ```json
  {
    "status_code": 200,
    "access_token": "new-jwt-access-token"
  }
  ```

#### 🛂 User Profile

- **Endpoint:** `GET /users/profile`
- **Description:** Retrieve the authenticated user profile.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
  ```json
  {
    "status_code": 200,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
  ```

#### 🔒 Password Reset Request

- **Endpoint:** `POST /users/password-reset`
- **Description:** Send password reset link to the user's email.
- **Request Body:**
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
  ```json
  {
    "status_code": 200,
    "message": "Password reset link sent."
  }
  ```

#### 🔑 Password Update After Reset

- **Endpoint:** `POST /users/password-update`
- **Description:** Update the user's password using the token received in the password reset email.
- **Request Body:**
  ```json
  {
    "token": "reset-token",
    "new_password": "newpassword123"
  }
  ```
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
  ```json
  {
    "status_code": 200,
    "message": "Password updated successfully."
  }
  ```

### 🛠 2. Admin Functionalities

#### 👥 View All Users

- **Endpoint:** `GET /admin/users`
- **Description:** Retrieve a list of all users.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
  ```json
  {
    "status_code": 200,
    "users": [
      {
        "email": "user1@example.com",
        "first_name": "John",
        "last_name": "Doe"
      },
      ...
    ]
  }
  ```

#### 🗑 Delete User Account

- **Endpoint:** `DELETE /admin/users/{id}`
- **Description:** Delete a specific user account.
- **Response:**
  - **Status Code:** `204 No Content`
  - **Body:**
  ```json
  {
    "status_code": 204,
    "message": "User account deleted successfully."
  }
  ```

## 📚 Documentation

- **API Documentation:** Available on Postman [here](https://documenter.getpostman.com/view/32898780/2sAXjGdEjE).

This API enables a secure and scalable loan management system with user and admin functionalities, built on Golang with the Gin framework and MongoDB.

---
