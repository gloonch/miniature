-- products table schema
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0), -- Assuming 2 decimal places for price
    sku VARCHAR(100), -- Stock Keeping Unit
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_shop
        FOREIGN KEY(shop_id)
        REFERENCES shops(id)
        ON DELETE CASCADE, -- If a shop is deleted, its products are also deleted.

    CONSTRAINT uq_shop_sku UNIQUE (shop_id, sku) -- SKU should be unique within a single shop
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_products_shop_id ON products(shop_id);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku); -- Index SKU for faster lookups if needed globally
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name); -- Index name for searching
