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
- [x] Seed default roles (admin, preparation, accueil) and admin user on first startup (empty DB only)
- [x] Prevent deletion or deactivation of the last active admin account

### Route Permission Matrix

| Route Group | Method | Endpoint | Admin | Accueil | Preparation |
| --- | --- | --- | --- | --- | --- |
| **Users** | POST | `/users/login` | Public | Public | Public |
| | POST | `/users/` | x | | |
| | GET | `/users/`, `/users/:id` | x | | |
| | DELETE | `/users/:id` | x | | |
| | PATCH | `/users/:id/reset-password` | x | | |
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
- [x] Add `PATCH /users/:id/reset-password` endpoint for admin-initiated password reset (generates random temp password)

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
- [x] Write tests for **controllers/users.go**: Login (success/wrong pw/no email/deactivated/invalid), CreateUser (success/dup email/weak pw/bad role), GetUsers, GetUser, DeleteUser (success/last admin blocked/non-admin allowed/soft delete preserves record/email reusable/cannot login after delete), ToggleUserStatus (success/last admin blocked), ChangePassword (own/other/admin), ResetPassword (success/not found/invalid ID/login with temp password)
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
- [x] Write tests for **utils/pwdvalidator.go**: valid password, too short, missing upper/lower/number/special, all special chars, GenerateTempPassword (length/validation/min floor/randomness)
- [x] Ensure all tests pass — 148 tests passing across 12 test files. Run with `CGO_ENABLED=1 go test ./... -v`

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

- [x] Verify ERD diagrams in `references/` are up-to-date with the current models (especially if new fields are added)
- [x] Create or update a functional flow diagram showing view navigation and user interactions per role
- [x] Update `README.md` with: project description, setup instructions, environment variables, API overview, tech stack, how to run tests
- [x] Document the API endpoints (Swagger already covers this — ensure it's regenerated after any changes with `swag init`)

---

# 2. Frontend

## 2.1 User Management

- [x] Add activate/deactivate toggle in the frontend (Users page)
- [x] Show logged-in user's role in topbar
- [x] Hide sidebar navigation based on user role (admin/accueil/preparation)
- [x] Add password reset button in the Users page (shows temp password in modal)
- [x] Handle last-admin guard errors gracefully (toast + table reload)
- [x] Update permissions table with Reset Password action

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

| Category                  | Done / Total | Status    |
| ------------------------- | ------------ | --------- |
| **Backend**               |              |           |
| RBAC                      | 8/8          | Done      |
| User Management           | 3/4          | Optional remains (`PUT /users/:id`) |
| Unit Tests                | 14/14        | Done (148 tests) |
| Security Hardening        | 2/4          | HTTPS redirect pending, per-IP rate limit optional |
| Product Images            | 3/3          | Done      |
| AutoMigrate Fix           | 1/1          | Done      |
| Deployment                | 0/5          | Not started |
| Code Quality              | 6/6          | Done      |
| Documentation/Deliverables| 4/4          | Done      |
| **Frontend**              |              |           |
| User Management           | 6/7          | Optional remains (user update form) |
| Product Images            | 2/2          | Done      |
| Order Sorting             | 3/3          | Done      |
| GDPR                      | 5/5          | Done      |
| Security                  | 2/2          | Done      |
| Frontend Polish           | 4/4          | Done      |


