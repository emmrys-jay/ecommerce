-- Drop the indexes
DROP INDEX IF EXISTS "products_name_idx";
DROP INDEX IF EXISTS "products_status_idx";
DROP INDEX IF EXISTS "products_price_idx";
DROP INDEX IF EXISTS "products_deleted_at_idx";

-- Drop the table
DROP TABLE IF EXISTS "products";

-- Drop the custom enum type
DROP TYPE IF EXISTS "product_status_enum";
