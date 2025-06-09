CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL, -- Assuming this links to customer IDs
    address TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_owner
        FOREIGN KEY(owner_id)
        REFERENCES customers(id)
        ON DELETE CASCADE -- Or ON DELETE RESTRICT / SET NULL depending on desired behavior
);

CREATE INDEX IF NOT EXISTS idx_shops_owner_id ON shops(owner_id);
