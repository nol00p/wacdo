# WacDo — Functional Flows

This document describes the navigation structure and user interactions per role.

---

## Authentication Flow

```
┌─────────────┐     POST /users/login      ┌──────────────┐
│  Login Page  │ ──────────────────────────→ │  API Server  │
│  (email +    │                             │  validates   │
│   password)  │ ←────────────────────────── │  credentials │
└─────────────┘     JWT token (2h expiry)    └──────────────┘
       │
       │ Token stored in localStorage
       │
       ▼
┌─────────────────────────────────────────┐
│  Role-based redirect:                    │
│  • Admin       → Full dashboard          │
│  • Accueil     → Order-focused dashboard │
│  • Preparation → Preparation dashboard   │
└─────────────────────────────────────────┘
```

On 401 response → session expired → redirect to login.

---

## Navigation Per Role

### Admin — Full Access

```
Sidebar:
├── Dashboard ──→ Stats overview (total orders, products, customers, users)
├── Products  ──→ Tabs: Categories | Products | Options
│                  • CRUD on categories, products, product options, option values
│                  • Toggle product availability
│                  • Update stock
│                  • Product image URL management
├── Menus     ──→ CRUD on menus
│                  • Add/remove products to menus
│                  • Toggle menu availability
├── Orders    ──→ Kanban board (pending → preparing → prepared → delivered)
│                  • Table view with status filter
│                  • Create new order (modal)
│                  • Update order status
│                  • Cancel pending orders
├── Customers ──→ Customer list with search
│                  • CRUD on customers
│                  • View order history per customer
│                  • GDPR data management
└── Users     ──→ Tabs: Users | Roles
                   • Create/delete users
                   • Activate/deactivate users
                   • Manage roles
```

### Accueil — Order Entry & Customer Facing

```
Sidebar:
├── Dashboard ──→ Order creation prominent, delivered orders tracking
├── Products  ──→ View only (categories, products, options)
├── Menus     ──→ View only
├── Orders    ──→ Kanban board + table
│                  • Create new orders (counter/phone)
│                  • Set scheduled delivery time
│                  • Cancel pending orders
│                  • Mark orders as delivered
└── Customers ──→ CRUD on customers
                   • View order history
```

### Preparation — Kitchen View

```
Sidebar:
├── Dashboard ──→ Orders sorted by scheduled time
│                  Only pending and preparing orders visible
├── Products  ──→ View only
├── Menus     ──→ View only
└── Orders    ──→ Kanban board (filtered: pending + preparing + prepared)
                   • Mark orders as preparing
                   • Mark orders as prepared
```

---

## Order Creation Flow

```
┌──────────────────┐
│  New Order Modal  │
│                   │
│  1. Select type   │──→ "counter" or "phone"
│  2. Select        │──→ Optional: link existing customer
│     customer      │    or create new one
│  3. Add items     │──→ Choose product or menu
│     ├── Product   │    Select options (size, toppings...)
│     └── Menu      │    Fixed price, no options
│  4. Set quantity  │
│  5. Set scheduled │──→ Optional delivery time
│     time          │
│  6. Add notes     │──→ Special instructions for kitchen
│  7. Review total  │──→ Live price calculation (client-side preview)
│                   │
│  [Submit Order]   │
└────────┬──────────┘
         │
         ▼  POST /orders/
┌──────────────────┐
│  Server computes  │──→ Fetches current prices from DB
│  total in a DB    │    Validates availability & stock
│  transaction      │    Creates order + items + options
│                   │    Returns order with total_price
└──────────────────┘
```

---

## Order Status Flow

```
          ┌───────────┐
          │  pending   │ ← Order created
          └─────┬─────┘
                │
        ┌───────┴───────┐
        │               │
        ▼               ▼
┌──────────────┐  ┌───────────┐
│  preparing   │  │ cancelled │  ← Only from pending
│  (kitchen    │  └───────────┘
│   started)   │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   prepared   │ ← Ready for pickup
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  delivered   │ ← Customer received order
└──────────────┘
```

**Who can do what:**
- **Accueil:** Create orders, cancel pending, mark as delivered
- **Preparation:** Mark as preparing, mark as prepared
- **Admin:** All status transitions

---

## Customer Data Flow (GDPR)

```
┌─────────────────────────────────────────────┐
│  Customer Management (Accueil / Admin)       │
│                                              │
│  Create  ──→ Consent acknowledgement shown   │
│  Read    ──→ View customer details + orders  │
│  Update  ──→ Modify name, phone, email       │
│  Delete  ──→ Hard delete (right to erasure)  │
│                                              │
│  GDPR notice accessible from frontend        │
│  Explains: what data, why, who has access    │
└─────────────────────────────────────────────┘
```

---

## Product Management Flow (Admin Only)

```
Categories ──→ Products ──→ Product Options ──→ Option Values
                  │
                  ▼
               Menus (via MenuProducts join)

1. Create categories to organize the catalog
2. Create products within categories (name, price, stock, image)
3. Define options for products (e.g., "Size" — single/required)
4. Add values to options (e.g., "Large" — +€3.00)
5. Create menus and assign products to them
6. Toggle availability on products and menus as needed
```
