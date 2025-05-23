CREATE TABLE customers
(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone            VARCHAR(20) UNIQUE NOT NULL,
    name             TEXT,
    total_spent      NUMERIC          DEFAULT 0,
    cashback_balance NUMERIC          DEFAULT 0,
    role             TEXT             DEFAULT,
    created_at       TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);
