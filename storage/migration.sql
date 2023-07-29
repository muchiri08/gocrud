CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20),
    email VARCHAR(20),
    password_hash VARCHAR(150)
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20)
    price NUMERIC(6, 2)
);