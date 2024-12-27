CREATE TYPE order_status AS ENUM ('Pending', 'Processing', 'Shipped', 'Delivered', 'Cancelled');

CREATE TABLE orders (
    "id" UUID PRIMARY KEY default uuid_generate_v4(),
    "user_id" UUID NOT NULL,
    "status" order_status DEFAULT 'Pending',
    "total_amount" NUMERIC(10,2) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE order_items (
    "id" UUID PRIMARY KEY default uuid_generate_v4(),
    "order_id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "product_name" varchar NOT NULL,
    "quantity" INTEGER NOT NULL CHECK (quantity > 0),
    "unit_price" numeric(10,2) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE INDEX "orders_status_idx" ON "orders" ("status");
CREATE INDEX "order_items_order_id_idx" ON "order_items" ("order_id");