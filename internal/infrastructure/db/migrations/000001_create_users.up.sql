CREATE TABLE  users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW()
);