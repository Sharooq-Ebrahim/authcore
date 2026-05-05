CREATE TABLE  users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    password_hash TEXT,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW()
);