# E-commerce API
This project is a RESTful API for an e-commerce application, built using Golang. It supports user authentication, product management, and order management with role-based access control.

## Assessment Requirement
- User Management: Register and login with JWT authentication.
- Product Management: CRUD operations for products (restricted to admin users).
- Order Management: Place orders, view user orders, cancel orders, and update order status (admin-only).
- Role-based Access: Admin and user roles with specific permissions.
- Validation & Error Handling: Complete input validation and appropriate HTTP status codes.
- Swagger Documentation: Each endpoint is documented for easy reference.

## Running
First, clone the repo and install the dependencies:

```bash
git clone https://github.com/funmi4194/ecommerce.git
cd instashop
go mod tidy
create .env file
go run main.go
```

## Render Workflow

- `RENDER_DEPLOY_HOOK` refers to the hook to trigger a render deployment for the service
# ecommerce
