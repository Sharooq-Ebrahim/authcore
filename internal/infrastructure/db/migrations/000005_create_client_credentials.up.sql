CREATE TABLE client_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL,
    key TEXT UNIQUE NOT NULL,
    name TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);