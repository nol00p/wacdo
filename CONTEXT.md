# WacDo — Project Context File

> Give this file to Claude at the start of any new session so it can resume work immediately.
> Last updated: 2026-03-25

---

## What is this project?

**WacDo** is a restaurant back-office / order management system built as a Go training project (formation). It manages users, products, menus, customers, and orders for a fast-food ordering kiosk. The project is evaluated by a professional jury against a formal brief (see `brief` section below).

---

## Tech Stack

| Layer      | Tech                                      |
| ---------- | ----------------------------------------- |
| Language   | Go 1.25                                   |
| Framework  | Gin (HTTP router)                         |
| ORM        | GORM                                      |
| Database   | PostgreSQL                                |
| Auth       | JWT (HS256, 2h expiry, `golang-jwt/v5`)   |
| Passwords  | bcrypt                                    |
| Docs       | Swagger via `swag`                        |
| Frontend   | Vanilla JS SPA (no build step)            |
| Deployment | Render (planned), no Dockerfile yet       |

---

## Project Structure

```
wacdo/
├── main.go                  # Entry point: middleware, routes, DB migration, seed defaults
├── config/
│   ├── db.go                # PostgreSQL connection (GORM)
│   ├── cors.go              # CORS middleware
│   ├── secure.go            # Security headers (X-Frame, CSP, XSS filter)
│   └── rate.go              # Rate limiter (global, not per-IP)
├── middlewares/
│   └── auth.go              # JWT Bearer authentication middleware
├── models/
│   ├── users.go             # Users + UserInput
│   ├── roles.go             # Roles (permissions stored as text)
│   ├── products.go          # Category, Products, ProductOptions, OptionValues
│   ├── menu.go              # Menu, MenuProduct
│   ├── customer.go          # Customer
│   └── orders.go            # Order, OrderItem, OrderItemOption
├── controllers/
│   ├── users.go             # Login + CRUD + password reset (no update endpoint)
│   ├── roles.go             # CRUD
│   ├── products.go          # CRUD + availability toggle + stock update
│   ├── product_categories.go# CRUD
│   ├── product_options.go   # CRUD + get by product
│   ├── product_option_values.go # CRUD (batch create) + get by option
│   ├── menu.go              # CRUD + availability + add/remove/list products
│   ├── customers.go         # CRUD
│   └── orders.go            # Create (transactional, server-side pricing),
│                            #   GetOrders (status filter), GetOrder,
│                            #   UpdateStatus (state machine), Cancel,
│                            #   GetOrdersByCustomer
├── routes/
│   ├── users.go             # /users (login public, rest protected)
│   ├── roles.go             # /roles (all protected)
│   ├── products.go          # /products, /categories, /options (all protected)
│   ├── menu.go              # /menus (all protected)
│   ├── customers.go         # /customers (all protected)
│   └── orders.go            # /orders + /customers/:id/orders (all protected)
├── utils/
│   └── pwdvalidator.go      # Password strength validation
├── frontend/
│   ├── index.html           # Single-page app shell
│   ├── css/style.css        # Dark theme, responsive
│   └── js/
│       ├── app.js           # Router, API helper, toast, modal
│       └── pages/           # login, dashboard, products, menus,
│                            #   orders (kanban + table + new order form),
│                            #   customers, users
├── docs/                    # Swagger auto-generated (swagger.json/yaml)
├── references/
│   ├── erd.md               # ERD documentation with table examples
│   ├── user_stories.md      # 50 user stories across 6 epics
│   ├── wacdo_ERD.png        # ERD diagram image
│   └── *.drawio             # Draw.io source files
├── .env                     # DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME,
│                            #   JWT_SECRET, CORS_ORIGINS
├── todo.md                  # Full TODO list (73 items) for brief coverage
└── README.md                # Basic readme (outdated)
```

---

## Data Models (11 tables)

```
Users         → has one Role (FK: RolesID), has IsActive bool, soft delete (DeletedAt)
Roles         → name, description, permissions (text)
Category      → name, description, display_order, image_url
Products      → belongs to Category, price, stock, is_available, prep_time, image_url
ProductOptions → belongs to Product, name, is_unique (single/multiple), is_required
OptionValues  → belongs to ProductOption, value, option_price
Menu          → name, description, price, is_available
MenuProduct   → join: Menu ↔ Product, quantity, is_optional, display_order
Customer      → name, phone, email
Order         → belongs to Customer (optional), created_by User,
                order_type (counter/phone), status, notes, scheduled_time, total_price
OrderItem     → belongs to Order, has Product OR Menu, quantity, unit_price, item_total
OrderItemOption → belongs to OrderItem, has OptionValue, price_applied
```

**Order status workflow:** `pending → preparing → prepared → delivered` (cancel only from pending)

---

## API Endpoints (~40 routes)

All protected routes require `Authorization: Bearer <jwt>` header.

| Group      | Endpoints                                                                 |
| ---------- | ------------------------------------------------------------------------- |
| Users      | `POST /users/login` (public), `POST /users/`, `GET /users/`, `GET /users/:id`, `DELETE /users/:id`, `PATCH /users/:id/reset-password` (admin) |
| Roles      | `GET/POST /roles/`, `GET/DELETE /roles/:id`                               |
| Categories | `GET/POST /categories/`, `GET/PUT/DELETE /categories/:id`                 |
| Products   | `GET/POST /products/`, `GET/PUT/DELETE /products/:id`, `GET /products/category/:id`, `PATCH /products/:id/availability`, `PATCH /products/:id/stock` |
| Options    | `GET/POST /options/`, `GET/PUT/DELETE /options/:id`, `GET /options/product/:id` |
| Opt Values | `POST/GET /options/:id/values/`, `GET/PUT/DELETE /options/values/:id`     |
| Menus      | `GET/POST /menus/`, `GET/PUT/DELETE /menus/:id`, `PATCH /menus/:id/availability`, `POST/GET /menus/:id/products/`, `DELETE /menus/products/:id` |
| Customers  | `GET/POST /customers/`, `GET/PUT/DELETE /customers/:id`                   |
| Orders     | `POST/GET /orders/`, `GET /orders/:id`, `PATCH /orders/:id/status`, `PATCH /orders/:id/cancel`, `GET /customers/:id/orders` |

---

## What's Working

- Full CRUD on all entities
- JWT login with bcrypt password hashing + admin password reset + last-admin protection
- Server-side order price calculation inside a DB transaction
- Order status state machine with valid transition enforcement
- Swagger docs at `/swagger/index.html`
- Frontend SPA: login, dashboard (stats), products (categories/products/options tabs), menus (expandable product lists), orders (kanban + table + new order modal with live pricing), customers (search + order history), users (users + roles tabs + password reset + last-admin error handling)
- Security: CORS, CSP headers, XSS filter, frame deny, rate limiting

---

## What's NOT Working / Missing (see todo.md for full list)

### Remaining
1. **No user update endpoint** (optional) — `PUT /users/:id` for editing username/email/role.
2. **HTTPS redirect** — not yet enabled for production.
3. **Rate limiter is global** (optional) — single bucket shared across all IPs.
4. **No Dockerfile or deployment config** — Render deployment pending.
5. **Frontend not served by Go** — standalone, hardcoded to `localhost:8000`.

### Resolved (since CONTEXT.md creation)
- RBAC, unit tests (148), product images, order sorting, GDPR, code quality, typos, ERD/README, MenuProduct AutoMigrate, role seeding, soft delete, password reset, last-admin guard.

---

## Key Design Decisions Already Made

- **Server-side price computation** for orders (no client trust) — done in a DB transaction.
- **JWT claims contain UserID and RoleName** — used by RBAC middleware for authorization.
- **Permissions stored as text in Roles table** — the ERD doc describes a proper permission/glue-table system but it was simplified to a text field. For RBAC, we can either check `RolesID` directly (simpler) or build the full permission system.
- **Order items can be either a Product or a Menu** (exactly one, validated).
- **Soft delete on Users** — `gorm.DeletedAt` field enables soft delete. Deleted users are hidden from queries but preserved for order audit trails. Email reuse is enabled via partial unique index. `IsActive` is for deactivation (reversible), soft delete is for removal (irreversible from UI). Last active admin cannot be deleted or deactivated.
- **Frontend uses global `window.*` functions** for event handlers in dynamically rendered HTML.

---

## Environment Variables (.env)

```
DB_HOST=...
DB_PORT=...
DB_USER=...
DB_PASS=...
DB_NAME=...
JWT_SECRET=...
CORS_ORIGINS=...
```

Also supports `DATABASE_URL` (single connection string for Render) and `PORT` (defaults to 8000).

---

## Git State

Check `git status` and `git log --oneline -10` for current state — this section is no longer maintained inline.

---

## How to Resume Work

1. Read this file to get full context.
2. Check `todo.md` for the current task list and what's been completed.
3. Check `git status` and `git log --oneline -5` for the latest state.
4. Ask me what to work on next, or continue from where we left off.

---

## Brief Summary (French original)

The project is a **back-office for a fast-food ordering kiosk** (borne de commande). Key requirements:
- 3 user roles with enforced permissions (admin, préparation, accueil)
- Product & menu management (CRUD + images + availability)
- Order creation (counter/phone), preparation tracking (sorted by delivery time), delivery
- Security (auth, sessions, data protection, injection prevention)
- GDPR compliance (data notice, right to consult/modify/delete)
- Unit tests
- Deployed and functional
- Clean code, documented, with ERD and functional diagrams
- Evaluated by professional jury on ~30 criteria
