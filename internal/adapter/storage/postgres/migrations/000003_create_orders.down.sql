-- Drop the dependent table first
DROP TABLE IF EXISTS "order_items";

-- Then drop the table with foreign keys
DROP TABLE IF EXISTS "orders";

-- Finally, drop the custom enum type
DROP TYPE IF EXISTS "order_status";
