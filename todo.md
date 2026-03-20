# WacDo - TODO for 100% Brief Coverage

## Legend

- [ ] Not started
- [x] Done

---

# 1. Backend

## 1.1 Role-Based Access Control (RBAC)

> Brief: "Mise en place des comptes utilisateurs (utilisateurs internes), avec la prise en compte des autorisations"

- [x] Create an authorization middleware that checks the user's role before granting access to endpoints
- [x] Restrict **Admin** role: full access to user management, product management, menu management
- [x] Restrict **Preparation** role: can only view orders and mark them as "prepared"
- [x] Restrict **Accueil** role: can create orders (counter/phone), deliver orders to customers, view orders
- [x] Apply role middleware to all protected routes with appropriate role requirements
- [x] Store the user's role (or role ID) in the JWT claims so the middleware can read it without extra DB queries
- [ ] Seed default roles (admin, preparation, accueil) on first startup or provide a migration script

### Route Permission Matrix

| Route Group | Method | Endpoint | Admin | Accueil | Preparation |
| --- | --- | --- | --- | --- | --- |
| **Users** | POST | `/users/login` | Public | Public | Public |
| | POST | `/users/` | x | | |
| | GET | `/users/`, `/users/:id` | x | | |
| | DELETE | `/users/:id` | x | | |
| **Roles** | GET/POST/DELETE | `/roles/...` | x | | |
| **Products** | GET | `/products/...` | x | x | x |
| | POST/PUT/DELETE/PATCH | `/products/...` | x | | |
| **Categories** | GET | `/categories/...` | x | x | x |
| | POST/PUT/DELETE | `/categories/...` | x | | |
| **Options** | GET | `/options/...` | x | x | x |
| | POST/PUT/DELETE | `/options/...` | x | | |
| **Menus** | GET | `/menus/...` | x | x | x |
| | POST/PUT/DELETE/PATCH | `/menus/...` | x | | |
| **Customers** | GET/POST/PUT/DELETE | `/customers/...` | x | x | |
| **Orders** | GET | `/orders/`, `/orders/:id` | x | x | x |
| | POST | `/orders/` | x | x | |
| | PATCH | `/orders/:id/cancel` | x | x | |
| | PATCH | `/orders/:id/status` | x | x | x |
| | GET | `/customers/:id/orders` | x | x | |

---

## 1.2 User Management

> Brief: "Gestion des Utilisateurs", "désactiver/réactiver"

- [x] Add `PATCH /users/:id/status` endpoint to toggle `IsActive` (deactivate/reactivate a user without deleting)
- [x] Block login for deactivated users (`IsActive == false`)

### Optional
- [ ] Add `PUT /users/:id` endpoint to update user information (username, email, role)

---

## 1.3 Product Images

> Brief: "y compris les noms, descriptions, prix et images"

- [x] Add `ImageURL string` field to the `Products` model
- [x] Run migration to add the column to the database (AutoMigrate handles it)
- [x] Update `CreateProduct` and `UpdateProduct` controllers to handle the new field (GORM handles it automatically via the model)

---

## 1.4 Unit Tests

> Brief: "Des tests unitaires sont réalisés et validés"

- [x] Set up test infrastructure (test database or mocks, test helpers) — `testutils/setup.go` with in-memory SQLite, seed helpers, JSON request builder
- [x] Write tests for **controllers/users.go**: Login (success/wrong pw/no email/deactivated/invalid), CreateUser (success/dup email/weak pw/bad role), GetUsers, GetUser, DeleteUser, ToggleUserStatus, ChangePassword (own/other/admin)
- [x] Write tests for **controllers/roles.go**: CreateRole (success/duplicate/invalid), GetRoles, GetRole (success/not found/invalid ID), DeleteRole (success/not found/still in use)
- [x] Write tests for **controllers/products.go**: CRUD + availability toggle + stock update + category filter + duplicates + not found
- [x] Write tests for **controllers/product_categories.go**: CRUD + duplicate + name conflict + still in use protection
- [x] Write tests for **controllers/product_options.go**: CRUD + invalid is_unique + get by product + not found
- [x] Write tests for **controllers/product_option_values.go**: Create batch + duplicate value + CRUD + get by option
- [x] Write tests for **controllers/menu.go**: CRUD + availability toggle + add/get/remove products + name conflict
- [x] Write tests for **controllers/customers.go**: CRUD + duplicate phone + phone conflict on update
- [x] Write tests for **controllers/orders.go**: CreateOrder (product/menu/invalid type/no items/both product+menu/unavailable), GetOrders + status filter, GetOrder, UpdateStatus (valid/invalid/full workflow), CancelOrder (pending/non-pending), GetOrdersByCustomer
- [x] Write tests for **middlewares/auth.go**: valid token, expired token, missing token, no Bearer prefix, invalid token, wrong signing method, context values
- [x] Write tests for **RBAC middleware**: role allowed, role denied, missing role, multiple roles, empty string role
- [x] Write tests for **utils/pwdvalidator.go**: valid password, too short, missing upper/lower/number/special, all special chars
- [x] Ensure all tests pass — 134 tests passing across 12 test files. Run with `CGO_ENABLED=1 go test ./... -v`

---

## 1.5 Security Hardening

> Brief: "Mise en oeuvre de mesures de sécurité robustes", "Protection des données", "empêchant toute injection"

- [ ] Add HTTPS redirect for production (`SSLRedirect: true` when not in dev mode)
- [x] Ensure `.env` is in `.gitignore` (already done) and that no secrets are committed
- [x] Add password change endpoint for users

### Optional
- [ ] Rate limiter should be **per-IP** instead of global (current implementation uses a single shared bucket)

---

## 1.6 Fix AutoMigrate

- [x] Add `&models.MenuProduct{}` to the `AutoMigrate()` call in `main.go`

---

## 1.7 Deployment

> Brief: "L'application fonctionnelle déployée sur le serveur"

- [ ] Create a `Dockerfile` (multi-stage build: Go compile + binary)
- [ ] Create a `render.yaml` or `Procfile` for Render deployment
- [ ] Set up environment variables on the hosting platform (DATABASE_URL, JWT_SECRET, CORS_ORIGINS)
- [ ] Verify the deployed application works end-to-end (login, CRUD, orders)
- [ ] Ensure Swagger UI is accessible on the deployed version

---

## 1.8 Code Quality & Documentation

> Brief: "Le code est indenté, les commentaires aident à la compréhension du code", "conventions de nommage"

- [x] Add comments to all controller functions explaining the business logic (not just Swagger annotations)
- [x] Add comments to models explaining each field's purpose and constraints
- [x] Add a package-level comment to each Go package (`config`, `controllers`, `middlewares`, `models`, `routes`, `utils`)
- [x] Review and fix typos in error messages (`"couln't not be reated"` → `"couldn't be created"`, `"Category couln't not be reated"` → `"Category couldn't be created"`)
- [x] Review and fix typos in middleware (`"Unauthorized Accesss"` → `"Unauthorized access"`, `tokeString` → `tokenString`)
- [x] Ensure consistent error message casing and formatting across all controllers — standardized to sentence case, replaced `"Can't get..."` with `"Failed to retrieve..."`, normalized `"Not found"/"not found"` patterns

---

## 1.9 Deliverables Documentation

> Brief: "schémas conceptuels et physiques du modèle de données", "schémas fonctionnels"

- [ ] Verify ERD diagrams in `references/` are up-to-date with the current models (especially if new fields are added)
- [ ] Create or update a functional flow diagram showing view navigation and user interactions per role
- [ ] Update `README.md` with: project description, setup instructions, environment variables, API overview, tech stack, how to run tests
- [ ] Document the API endpoints (Swagger already covers this — ensure it's regenerated after any changes with `swag init`)

---

# 2. Frontend

## 2.1 User Management

- [x] Add activate/deactivate toggle in the frontend (Users page)
- [x] Show logged-in user's role in topbar
- [x] Hide sidebar navigation based on user role (admin/accueil/preparation)

### Optional
- [ ] Add user update form in the frontend (Users page)

---

## 2.2 Product Images

- [x] Update the frontend product form to include an image URL field
- [x] Display product images in the frontend product list

---

## 2.3 Order Sorting & Filtering for Preparation View

> Brief: "liste des commandes à préparer (triées par heure de livraison croissante)"

- [x] Sort orders client-side by `scheduled_time ASC` (nulls last) in preparation view (pending/preparing filters)
- [x] Show scheduled delivery time in the kanban cards and order table
- [x] Add scheduled time field in the new order form

---

## 2.4 GDPR & Data Protection

> Brief: "L'application informe l'utilisateur du stockage, de l'utilisation et du cadre de partage de ses données personnelles" + "droit de consultation, modification et de suppression"

- [x] Add a GDPR / Privacy notice page or modal accessible from the frontend
- [x] The notice must explain: what data is stored, how it is used, and with whom it is shared
- [x] Ensure customers can be viewed, updated, and deleted (already have endpoints — verify they work)
- [x] Add a "Delete my data" or "Export my data" mechanism (or document how staff can do it via the back-office)
- [x] Add a consent acknowledgement step when creating a customer record (frontend)

---

## 2.5 Security

- [x] Add session expiration handling on the frontend (already redirects on 401, but add a "session expired" message)
- [x] Sanitize user-supplied text before rendering with `innerHTML` to prevent XSS (SQL injection is already handled by GORM)

---

## 2.6 Frontend Polish

- [x] Add a **preparation-specific dashboard** view: orders sorted by scheduled time, only "pending" and "preparing" visible
- [x] Add a **accueil-specific dashboard** view: order creation prominent, delivered orders tracking
- [x] Add confirmation dialogs for all destructive actions (some already exist)
- [x] Add loading states/spinners for all async operations

---

# Summary

| Category                  | Items | Priority  |
| ------------------------- | ----- | --------- |
| **Backend**               |       |           |
| RBAC                      | 7     | Critical  |
| User Management           | 3     | Critical  |
| Unit Tests                | 14    | Done      |
| Security Hardening        | 4     | High      |
| Product Images            | 3     | Done      |
| AutoMigrate Fix           | 1     | Done      |
| Deployment                | 5     | Medium    |
| Code Quality              | 6     | Done      |
| Documentation/Deliverables| 4     | Low       |
| **Frontend**              |       |           |
| User Management           | 4     | Done      |
| Product Images            | 2     | Done      |
| Order Sorting             | 3     | Done      |
| GDPR                      | 5     | Done      |
| Security                  | 2     | Done      |
| Frontend Polish           | 4     | Done      |

---

# Open Questions

### Users & Roles
- **User hard delete vs deactivation** (needs PO validation): Users with orders can never be hard-deleted due to FK constraint on `orders.created_by_id`. Three options:
  1. **Deactivation only** for users with order history, hard delete only for users without orders. Downside: deactivated user's email can never be reused.
  2. **SET NULL on delete** — make `created_by_id` nullable with `ON DELETE SET NULL`. Hard delete works, email is freed, order history preserved but loses "who created it".
  3. **Soft delete** (`gorm.DeletedAt`) — row stays but is hidden from queries. Downside: email still blocked by unique constraint.
- **Role seeding**: CLI flag (`go run . --seed`) vs separate command (`cmd/seed/main.go`) vs manual SQL?
