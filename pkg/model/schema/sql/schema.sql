-- Version: 1.01
-- Description: Create table users
CREATE TABLE users
(
    user_id       UUID,
    name          TEXT,
    phone         TEXT UNIQUE,
    role          TEXT,
    password_hash TEXT,
    date_created  TIMESTAMP,
    date_updated  TIMESTAMP,

    PRIMARY KEY (user_id)
);