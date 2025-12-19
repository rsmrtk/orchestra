-- Orchestra Database Schema
-- ============================================
-- Table: orchestra-table (customers)
-- ============================================

CREATE TABLE IF NOT EXISTS "orchestra-table" (
    customer_id VARCHAR(255) NOT NULL PRIMARY KEY,
    first_name  VARCHAR(255) NOT NULL,
    last_name   VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_customer_name ON "orchestra-table" (first_name, last_name);
