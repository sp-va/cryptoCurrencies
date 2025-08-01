CREATE TABLE IF NOT EXISTS currency_prices (
    id SERIAL PRIMARY KEY,
    coin TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    price DECIMAL(20, 10) NOT NULL
);

CREATE TABLE IF NOT EXISTS track_currencies (
    id SERIAL PRIMARY KEY,
    coin TEXT NOT NULL UNIQUE
);