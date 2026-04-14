# WacDo — Restaurant Back-Office & Ordering System

WacDo is a back-office application for a fast-food ordering kiosk (borne de commande). It handles user management, product catalogs, menu composition, customer records, and order lifecycle — from creation through preparation to delivery.

## Tech Stack

| Layer      | Technology                        |
| ---------- | --------------------------------- |
| Language   | Go 1.25                          |
| Framework  | Gin                              |
| ORM        | GORM                             |
| Database   | PostgreSQL                       |
| Auth       | JWT (HS256, golang-jwt/v5)       |
| Passwords  | bcrypt                           |
| API Docs   | Swagger (swag)                   |
| Frontend   | Vanilla JS SPA                   |

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- [swag](https://github.com/swaggo/swag) CLI (for regenerating Swagger docs)

### Setup

1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd wacdo
   ```

2. Create a `.env` file at the project root:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_user
   DB_PASS=your_password
   DB_NAME=wacdo
   JWT_SECRET=your_secret_key
   CORS_ORIGINS=http://localhost:5500
   ```
   For single connection string (e.g. Render): set `DATABASE_URL` instead of individual DB_ vars.

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run the application:
   ```bash
   go run .
   ```
   The server starts on port `8000` by default (override with `PORT` env var). GORM auto-migrates all tables on startup. On first launch (empty database), default roles and an admin user are seeded automatically:
   - **Email:** `admin@wacdo.fr`
   - **Password:** `Admin@1234`
   - Change the admin password after first login.

5. Open the frontend by serving the `frontend/` directory (e.g. with VS Code Live Server) or any static file server.

### Running Tests

Tests use an in-memory SQLite database — no external DB required.

```bash
CGO_ENABLED=1 go test ./... -v
```

148 tests across 12 test files covering all controllers, middlewares, and utilities.

### Regenerating Swagger Docs

```bash
swag init
```

Swagger UI is available at `/swagger/index.html` when the server is running.

## API Overview

All endpoints except `POST /users/login` require a JWT Bearer token in the `Authorization` header.

| Group      | Key Endpoints                                                              |
| ---------- | -------------------------------------------------------------------------- |
| Users      | `POST /users/login`, `POST/GET /users/`, `GET/DELETE /users/:id`, `PATCH /users/:id/status`, `PATCH /users/:id/password`, `PATCH /users/:id/reset-password` |
| Roles      | `GET/POST /roles/`, `GET/DELETE /roles/:id`                                |
| Categories | `GET/POST /categories/`, `GET/PUT/DELETE /categories/:id`                  |
| Products   | `GET/POST /products/`, `GET/PUT/DELETE /products/:id`, `PATCH .../availability`, `PATCH .../stock` |
| Options    | `GET/POST /options/`, `GET/PUT/DELETE /options/:id`, `GET /options/product/:id` |
| Opt Values | `POST/GET /options/:id/values/`, `GET/PUT/DELETE /options/values/:id`      |
| Menus      | `GET/POST /menus/`, `GET/PUT/DELETE /menus/:id`, `PATCH .../availability`, menu products CRUD |
| Customers  | `GET/POST /customers/`, `GET/PUT/DELETE /customers/:id`                    |
| Orders     | `POST/GET /orders/`, `GET /orders/:id`, `PATCH .../status`, `PATCH .../cancel`, `GET /customers/:id/orders` |

Full details available in the Swagger documentation.

## Role-Based Access Control

Three roles with enforced permissions:

| Capability             | Admin | Accueil | Preparation |
| ---------------------- | :---: | :-----: | :---------: |
| User & role management | x     |         |             |
| Product/menu CRUD      | x     |         |             |
| View products/menus    | x     | x       | x           |
| Customer management    | x     | x       |             |
| Create orders          | x     | x       |             |
| View orders            | x     | x       | x           |
| Update order status    | x     | x       | x           |
| Cancel orders          | x     | x       |             |

## Order Lifecycle

```
pending → preparing → prepared → delivered
   ↓
cancelled (only from pending)
```

Orders use server-side price computation within a database transaction. Each order item can reference either a product or a menu.

## Project Structure

```
wacdo/
├── main.go              # Entry point, middleware, routes, DB migration, seed defaults
├── config/              # DB connection, CORS, security headers, rate limiter
├── middlewares/          # JWT auth + RBAC middleware
├── models/              # GORM models (12 tables)
├── controllers/         # Business logic for all entities
├── routes/              # Route definitions with role restrictions
├── utils/               # Password validator + temp password generator
├── frontend/            # Vanilla JS SPA (login, dashboard, CRUD pages)
├── docs/                # Auto-generated Swagger files
├── references/          # ERD, user stories, diagrams
└── testutils/           # Test setup with in-memory SQLite
```

## Security

- JWT authentication with 2-hour token expiry
- bcrypt password hashing
- Password strength validation (length, uppercase, lowercase, number, special char)
- Admin password reset (generates cryptographically random temp password)
- Last-admin protection (cannot delete or deactivate the only active admin)
- Soft delete on users (preserves order audit trails, frees email for reuse)
- Automatic role and admin seeding on first install
- CORS configuration
- Security headers (X-Frame-Options, CSP, XSS filter)
- Rate limiting
- GORM parameterized queries (SQL injection prevention)
- Frontend input sanitization (XSS prevention)

## Documentation

- **API**: Swagger UI at `/swagger/index.html`
- **ERD**: `references/erd.md` (schema + mermaid diagram)
- **User Stories**: `references/user_stories.md` (63 stories across 9 epics)
- **Functional Flows**: `references/functional_flows.md`
