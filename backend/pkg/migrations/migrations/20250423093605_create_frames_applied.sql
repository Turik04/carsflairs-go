CREATE TABLE frames (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    price INT NOT NULL,
    image TEXT,
    material TEXT
);