-- Orchestra Database Initialization Script

-- Create the orchestra-table if it doesn't exist
CREATE TABLE IF NOT EXISTS "orchestra-table" (
    customer_id VARCHAR(255) NOT NULL PRIMARY KEY,
    first_name  VARCHAR(255) NOT NULL,
    last_name   VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on first_name and last_name for faster lookups
CREATE INDEX IF NOT EXISTS idx_customer_names
ON "orchestra-table"(first_name, last_name);

-- Insert some sample data for testing (optional)
INSERT INTO "orchestra-table" (customer_id, first_name, last_name)
VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'John', 'Doe'),
    ('550e8400-e29b-41d4-a716-446655440001', 'Jane', 'Smith')
ON CONFLICT (customer_id) DO NOTHING;
