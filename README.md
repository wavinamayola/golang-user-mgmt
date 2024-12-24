
# User Management API

This is a simple Go application that provides a RESTful API for user management, including basic authentication, CRUD operations, and basic input validation.

## Features
-  **Create a user**: Add a new user to the system.
-  **Retrieve user details**: Fetch user information by ID.
-  **Update user details**: Modify user information.
-  **Delete user**: Remove a user from the system.
-  **Basic Authentication**: All endpoints require basic authentication (`admin:admin`).
-  **Input Validation**: Ensure proper fields (e.g., required fields, valid email) are provided when creating or updating a user.

## Prerequisites

Make sure you have the following installed on your system:
- [Go](https://golang.org/dl/) version 1.16 or higher
- [Git](https://git-scm.com/)

## Getting Started

### Clone the repository

```bash
git clone https://github.com/your-repository-url.git
cd your-repository-name
```

### Install dependencies

Run the following to download the required Go dependencies:
```bash
go mod tidy
```

### Update config file

Set the necessary config inside `env/config.yaml` (a `sample.config.yaml` file is provided for the config) based on your setup


### Run the migrations

Once the config is set, you can now run the migrations.

```bash
go run ./internal/migrations/migrate.go up
```


### Run the application

Start the server by running the following command:

```bash
go run cmd/user-mgmt/*
```

## Authentication

All endpoints require HTTP Basic Authentication. The default username and password are:

-   **Username**: `admin`
-   **Password**: `admin`

You can pass the credentials using `curl` or any HTTP client as shown below:
```bash
curl -u admin:admin http://localhost:8080/home
```
You should see something like `Welcome to the home page. This is just sample`

## Running Tests

To run the tests for the project, use:
```bash
go test ./...
```
It will run the tests inside service and storage level.

## Endpoints

Once the service is running, you can access http://localhost:8080/swagger/ locally to see the available endpoints.