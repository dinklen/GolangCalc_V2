-- Users datatable
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Expressions datatable
CREATE TABLE expressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    expression TEXT NOT NULL,
    result NUMERIC DEFAULT 0.0,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'calculated', 'failed')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Subexpressions datatable
CREATE TABLE sub_expressions (
    id TEXT PRIMARY KEY,
    expression_id UUID REFERENCES expressions(id),
    sub_expression TEXT NOT NULL,
    result NUMERIC DEFAULT 0.0,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'calculated', 'failed'))
);
