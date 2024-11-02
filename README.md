# E-commerce API
This project is a RESTful API for an e-commerce application, built using Golang. It supports user authentication, product management, and order management with role-based access control.

## Features
- User Management: Register and login with JWT authentication.
- Product Management: CRUD operations for products (restricted to admin users).
- Order Management: Place orders, view user orders, cancel orders, and update order status (admin-only).
- Role-based Access: Admin and user roles with specific permissions.
- Validation & Error Handling: Complete input validation and appropriate HTTP status codes.
- Swagger Documentation: Each endpoint is documented for easy reference.

## Prerequisites
- Go 1.16+ installed on your machine.
- A running instance of the E-commerce server.


## Installation
To install and use the E-commerce server, first, clone the repository, install the dependencies, create a .env file in the project root and add the environmental variables in .env.example and then run the server

1. Clone this repository:
```bash
git clone https://github.com/Funmi4194/ecommerce.git
cd ecommerce
```
2. Install dependencies
```bash
go mod tidy
```
3. Create a .env file
```bash
create .env file
```
4. Run the server
```bash
go run main.go
```

## Render Workflow

- `RENDER_DEPLOY_HOOK` refers to the hook to trigger a render deployment for the service
# ecommerce
