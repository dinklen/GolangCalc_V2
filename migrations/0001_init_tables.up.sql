CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE expressions (
    id UUID PRIMARY KEY,
    user_id INT REFERENCES users(id),
    expression TEXT NOT NULL,
    result NUMERIC,
    status VARCHAR(20) NOT NULL CHECK (status IN ('calculating', 'calculated', 'error')),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE sub_expressions (
    id UUID PRIMARY KEY,
    expression_id UUID REFERENCES expressions(id),
    sub_expression TEXT NOT NULL,
    result NUMERIC,
    status VARCHAR(20) NOT NULL CHECK (status IN ('calculating', 'calculated', 'error'))
);
