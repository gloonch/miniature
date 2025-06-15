CREATE TABLE customers
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone            VARCHAR(20) UNIQUE NOT NULL,
    name             TEXT,
    role             TEXT             DEFAULT 'OWNER',
    total_spent      NUMERIC          DEFAULT 0,
    cashback_balance NUMERIC          DEFAULT 0,
    is_active        BOOLEAN          DEFAULT TRUE,
    created_at       TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shops
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID REFERENCES customers (id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT,
    is_active   BOOLEAN          DEFAULT TRUE,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shop_users
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id    UUID REFERENCES shops (id) ON DELETE CASCADE,
    user_id    UUID REFERENCES customers (id) ON DELETE CASCADE,
    role       TEXT CHECK (role IN ('OWNER', 'SELLER')) NOT NULL,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (shop_id, user_id)
);


CREATE TABLE products
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT               NOT NULL,
    description TEXT,
    price       NUMERIC            NOT NULL,
    image_url   TEXT,
    stock       INTEGER          DEFAULT 0,
    category    TEXT,
    shop_id     UUID REFERENCES shops (id) ON DELETE CASCADE,
    is_active   BOOLEAN          DEFAULT TRUE,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders
(
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id           UUID REFERENCES shops (id) ON DELETE CASCADE,
    customer_id       UUID REFERENCES customers (id) ON DELETE CASCADE,
    status            TEXT CHECK (status IN
                                  ('PENDING', 'PAID', 'CONFIRMED', 'SHIPPED', 'DELIVERED', 'CANCELLED')) NOT NULL,
    total_amount      NUMERIC                                                                            NOT NULL,
    delivery_estimate TEXT,
    delivery_address  TEXT,
    payment_proof_url TEXT,
    cashback_applied  NUMERIC          DEFAULT 0,
    created_at        TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    confirmed_at      TIMESTAMP
);

CREATE TABLE order_items
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id       UUID REFERENCES orders (id) ON DELETE CASCADE,
    product_id     UUID REFERENCES products (id) ON DELETE CASCADE,
    quantity       INTEGER NOT NULL,
    price_at_order NUMERIC NOT NULL
);

CREATE TABLE cashback_logs
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id      UUID REFERENCES customers (id) ON DELETE CASCADE,
    related_order_id UUID REFERENCES orders (id),
    amount           NUMERIC                                 NOT NULL,
    type             TEXT CHECK (type IN ('EARNED', 'USED')) NOT NULL,
    created_at       TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_items
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id   TEXT NOT NULL,
    product_id   UUID REFERENCES products (id) ON DELETE CASCADE,
    quantity     INTEGER          DEFAULT 1,
    scanned_at   TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    is_finalized BOOLEAN          DEFAULT FALSE
);

CREATE TABLE feedback
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id     UUID REFERENCES customers (id) ON DELETE SET NULL,
    order_id        UUID REFERENCES orders (id) ON DELETE SET NULL,
    message         TEXT,
    sentiment_score NUMERIC,
    created_at      TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Materialized view or table for analytics (if needed)
CREATE TABLE product_stats
(
    product_id        UUID PRIMARY KEY REFERENCES products (id) ON DELETE CASCADE,
    total_sold        INTEGER DEFAULT 0,
    total_revenue     NUMERIC DEFAULT 0,
    last_purchased_at TIMESTAMP
);
