CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL,
    coins    INT DEFAULT 100     NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions
(
    id         SERIAL PRIMARY KEY,
    from_user  VARCHAR(255) NOT NULL,
    to_user    VARCHAR(255) NOT NULL,
    amount     INT          NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inventory
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    type     VARCHAR(255) NOT NULL,
    quantity INT          NOT NULL,
    UNIQUE (username, type)
);

CREATE TABLE IF NOT EXISTS merch
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) UNIQUE NOT NULL,
    price INT                 NOT NULL
);

INSERT INTO merch (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500)
ON CONFLICT (name) DO NOTHING;