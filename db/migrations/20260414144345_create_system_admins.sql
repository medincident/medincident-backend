-- migrate:up
CREATE TABLE domain.system_admins (
    zitadel_user_id TEXT        PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS domain.system_admins;
