CREATE TABLE IF NOT EXISTS urls(
    ID SERIAL PRIMARY KEY,
    shorted_url TEXT,
    full_url TEXT,
    user_id INTEGER,
    deleted bool
);