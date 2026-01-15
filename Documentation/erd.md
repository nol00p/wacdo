# Project data management

## User Management
in order to avoid having to hardcode most of the permisison and increase maintenability we will split the roles from the actual permissions and have a "glue" table 

**Roles**
```
roles:
┌────┬──────────────┐
│ id │ name         │
├────┼──────────────┤
│ 1  │ admin        │
│ 2  │ preparation  │
│ 3  │ accueil      │
└────┴──────────────┘
```
**Permissions**
```
permissions:
┌────┬────────────────┬──────────┬────────┐
│ id │ name           │ resource │ action │
├────┼────────────────┼──────────┼────────┤
│ 1  │ manage_users   │ users    │ create │
│ 2  │ manage_products│ products │ create │
│ 3  │ view_orders    │ orders   │ read   │
│ 4  │ create_orders  │ orders   │ create │
│ 5  │ prepare_orders │ orders   │ prepare│
│ 6  │ deliver_orders │ orders   │ deliver│
└────┴────────────────┴──────────┴────────┘
```
**Glue table**
```
role_permissions:
┌────┬─────────┬────────────────┬─────────────────────────────┐
│ id │ role_id │ permission_id  │ Meaning                     │
├────┼─────────┼────────────────┼─────────────────────────────┤
│ 1  │ 1       │ 1              │ admin can manage_users      │
│ 2  │ 1       │ 2              │ admin can manage_products   │
│ 3  │ 1       │ 3              │ admin can view_orders       │
│ 4  │ 1       │ 4              │ admin can create_orders     │
│ 5  │ 1       │ 5              │ admin can prepare_orders    │
│ 6  │ 1       │ 6              │ admin can deliver_orders    │
│ 7  │ 2       │ 3              │ preparation can view_orders │
│ 8  │ 2       │ 5              │ preparation can prepare_orders│
│ 9  │ 3       │ 3              │ accueil can view_orders     │
│ 10 │ 3       │ 4              │ accueil can create_orders   │
│ 11 │ 3       │ 6              │ accueil can deliver_orders  │
└────┴─────────┴────────────────┴─────────────────────────────┘
```
##  Product & Menu Management
#### Example
 - Option 1: Size (REQUIRED, SINGLE choice)
 ```
Product_options:
- id: 1
- product_id: 100 (Margherita Pizza)
- name: "Size"
- type: "single"
- is_required: true

option_values:
- id: 1,  option_id: 1, value: "Small (8 inch)",   price_modifier: -2.00  → €8.00
- id: 2,  option_id: 1, value: "Medium (12 inch)", price_modifier:  0.00  → €10.00
- id: 3,  option_id: 1, value: "Large (16 inch)",  price_modifier: +3.00  → €13.00
- id: 4,  option_id: 1, value: "XL (20 inch)",     price_modifier: +6.00  → €16.00
 ```

 - Option 2: Extra Toppings (OPTIONAL, MULTIPLE choices)
 ```
product_options:
- id: 2
- product_id: 100 (Margherita Pizza)
- name: "Extra Toppings"
- type: "multiple"
- is_required: false

option_values:
- id: 5,  option_id: 2, value: "Extra Cheese",    price_modifier: +1.50
- id: 6,  option_id: 2, value: "Pepperoni",       price_modifier: +2.00
- id: 7,  option_id: 2, value: "Mushrooms",       price_modifier: +1.00
- id: 8,  option_id: 2, value: "Olives",          price_modifier: +1.00
- id: 9,  option_id: 2, value: "Onions",          price_modifier: +0.50
- id: 10, option_id: 2, value: "Jalapeños",       price_modifier: +1.50
 ```
 
 - Option 3: Crust Type (OPTIONAL, SINGLE choice)
```
product_options:
- id: 3
- product_id: 100 (Margherita Pizza)
- name: "Crust Type"
- type: "single"
- is_required: false

option_values:
- id: 11, option_id: 3, value: "Regular",         price_modifier:  0.00
- id: 12, option_id: 3, value: "Thin Crust",      price_modifier:  0.00
- id: 13, option_id: 3, value: "Thick Crust",     price_modifier: +1.00
- id: 14, option_id: 3, value: "Stuffed Crust",   price_modifier: +3.00
- id: 15, option_id: 3, value: "Gluten-Free",     price_modifier: +2.50
```


## Customer & Order Management

One order can contain serveral items. Each item can serveral option. 
```
  1 ORDER                                                    
    └─→ N ORDER_ITEMS (different products)                  
           └─→ N ORDER_ITEM_OPTIONS (selected choices)      

  Example:                                                   
  Order #1001                                                
    ├─→ Item 1: Pizza (Large, Cheese, Pepperoni)            
    │     ├─→ Option: Size = Large (+€3.00)                 
    │     ├─→ Option: Topping = Cheese (+€1.50)             
    │     └─→ Option: Topping = Pepperoni (+€2.00)          
    │                                                         
    └─→ Item 2: Pizza (Large, Cheese)                        
          ├─→ Option: Size = Large (+€3.00)                  
          └─→ Option: Topping = Cheese (+€1.50)        
```      

## Security & Session Management
- for security all action made by users are logged in the DB
- The user session are all time stamped. 

## ERD 

![ERD](/Documentation/wacdo_ERD.png)



