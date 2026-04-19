-- migrate:up
CREATE INDEX employees_zitadel_user_id_idx
    ON domain.employees (zitadel_user_id);

-- migrate:down
DROP INDEX IF EXISTS domain.employees_zitadel_user_id_idx;
