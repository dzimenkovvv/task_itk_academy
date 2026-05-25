CREATE EXTENSION IF NOT EXISTS "pgcrypto";
 
CREATE TABLE wallet (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance   DOUBLE PRECISION NOT NULL DEFAULT 0
        CONSTRAINT balance_non_negative CHECK (balance >= 0)
);
 