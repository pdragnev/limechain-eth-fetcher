CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_hash        VARCHAR(66) NOT NULL UNIQUE,
    transaction_status      INTEGER NOT NULL,
    block_hash              VARCHAR(66) NOT NULL,
    block_number            BIGINT NOT NULL,
    from_address            VARCHAR(42) NOT NULL,
    to_address              VARCHAR(42),
    contract_address        VARCHAR(42),
    logs_count              INTEGER,
    input                   TEXT,
    value                   VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(64) NOT NULL
);

-- Insert the test users
INSERT INTO users (username, password) VALUES
('alice', '$2a$10$TAQe.aeP6bWSxXrvOPGwou5tmbitJ0akraYt3cbP1Isen9EhqQ07u'),
('bob', '$2a$10$W4sAcXpPOdUMPjS/GPHH5O0go3lEVSI4/ldpJJSi.1ILvnzK7.sM.'),
('carol', '$2a$10$NV.RQDMW6h/gSQEuXP7RkeRzKtpHa4jksuDP4ALdxZjcfPWUbJrvG'),
('dave', '$2a$10$FdBSyC.2iLXS7KwDfeMXDOpC1UPgJVIYKSeNEEeMC05pRC3tXtbo2')
ON CONFLICT (username) DO NOTHING;