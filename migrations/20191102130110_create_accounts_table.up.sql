CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS accounts (
    id      uuid    PRIMARY KEY DEFAULT uuid_generate_v4(),
    balance int     not null CHECK (balance >= 0)
);