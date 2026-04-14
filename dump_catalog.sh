#!/bin/bash
# WacDo — Dump product catalog from local DB
# Usage: bash dump_catalog.sh
# Output: catalog_dump.sql

set -e

source .env

{
echo "ALTER TABLE menu_products DISABLE TRIGGER ALL;"
echo "ALTER TABLE option_values DISABLE TRIGGER ALL;"
echo "ALTER TABLE product_options DISABLE TRIGGER ALL;"
echo "ALTER TABLE menus DISABLE TRIGGER ALL;"
echo "ALTER TABLE products DISABLE TRIGGER ALL;"
echo "ALTER TABLE categories DISABLE TRIGGER ALL;"
echo "TRUNCATE menu_products, option_values, product_options, menus, products, categories CASCADE;"
PGPASSWORD=$DB_PASS pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
  --data-only \
  --inserts \
  --column-inserts \
  --no-owner \
  --no-privileges \
  -t categories \
  -t products \
  -t product_options \
  -t option_values \
  -t menus \
  -t menu_products
echo "ALTER TABLE menu_products ENABLE TRIGGER ALL;"
echo "ALTER TABLE option_values ENABLE TRIGGER ALL;"
echo "ALTER TABLE product_options ENABLE TRIGGER ALL;"
echo "ALTER TABLE menus ENABLE TRIGGER ALL;"
echo "ALTER TABLE products ENABLE TRIGGER ALL;"
echo "ALTER TABLE categories ENABLE TRIGGER ALL;"
} > catalog_dump.sql

echo "Catalog dumped to catalog_dump.sql"
echo "Import on remote with: PGPASSWORD=<pass> psql -h localhost -U wacdo -d wacdo < catalog_dump.sql"
