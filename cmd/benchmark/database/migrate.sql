-- Drop tables if they exist
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;

-- Create tables fresh
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    email TEXT
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    amount INT,
    created_at TIMESTAMP
);

-- Optionally insert a single test row
INSERT INTO users(name, email) VALUES ('example', 'example@mail.com');
INSERT INTO orders(user_id, amount, created_at) VALUES (1, 100, NOW());

