-- migrate:up

-- Read-model copy of domain.system_admins.
CREATE TABLE projections.system_admins (
    zitadel_user_id TEXT        PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS projections.system_admins;
