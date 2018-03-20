#!/bin/bash

main() {
    psql -v ON_ERROR_STOP=1 --username postgres --db quote <<-EOSQL
        CREATE OR REPLACE FUNCTION add_funds(p_user_id VARCHAR(64), p_funds INTEGER) RETURNS refcursor AS \$\$
            DECLARE
                ref refcursor;
            BEGIN
                UPDATE users SET money = money + p_funds WHERE id = p_user_id RETURNING money INTO ref;

                IF NOT FOUND THEN
                    INSERT INTO users(id, money) VALUES (p_user_id, p_funds) RETURNING money INTO ref;
                END IF;
                RETURN ref;
            END
        \$\$ LANGUAGE plpgsql;
EOSQL
}

main