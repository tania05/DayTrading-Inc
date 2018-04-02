#!/bin/bash


set -o errexit

main() {
    init_user_and_db
}

init_user_and_db() {
    psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
        CREATE DATABASE quote;
EOSQL

    psql -v ON_ERROR_STOP=1 --username postgres --db quote <<-EOSQL
        CREATE TABLE users (
            id VARCHAR(64) PRIMARY KEY,
            money INTEGER CHECK (money >= 0)
        );

        CREATE TABLE stocks (
            user_id VARCHAR(64) REFERENCES users(id),
            stock_sym VARCHAR(3) NOT NULL,
            amount INTEGER CHECK (amount >= 0),
            UNIQUE(user_id, stock_sym)
        );

        CREATE TABLE triggers (
            execution_price INTEGER CHECK (execution_price >= 0),
            amount INTEGER CHECK (amount >= 0),
            stock_sym VARCHAR(3) NOT NULL,
            is_buy BOOLEAN
        );

        CREATE TABLE transactions (
            id SERIAL PRIMARY KEY,
            user_id VARCHAR(64) REFERENCES users(id),
            money_amount INTEGER,
            stock_sym VARCHAR(64),
            stock_amount INTEGER,
            is_buy BOOLEAN,
            created_at TIMESTAMP WITHOUT TIME ZONE
        );
EOSQL
}

main
