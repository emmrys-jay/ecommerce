CREATE TYPE "product_status_enum" AS ENUM ('active', 'inactive', 'out_of_stock');

CREATE TABLE "products" (
   "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   "name" varchar NOT NULL,
   "description" text,
   "price" NUMERIC(10,2) NOT NULL,
   "quantity" int NOT NULL DEFAULT 0,
   "status" product_status_enum DEFAULT 'active',
   "created_at" timestamptz NOT NULL DEFAULT (now()),
   "updated_at" timestamptz NOT NULL DEFAULT (now()),
   "deleted_at" timestamptz
);

-- Indexes
CREATE INDEX "products_name_idx" ON "products" ("name");
CREATE INDEX "products_status_idx" ON "products" ("status");
CREATE INDEX "products_price_idx" ON "products" ("price");
CREATE INDEX "products_deleted_at_idx" ON "products" ("deleted_at") WHERE deleted_at IS NOT NULL;
