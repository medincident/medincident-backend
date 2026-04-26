-- migrate:up
CREATE TABLE domain.request_types (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    name            TEXT NOT NULL,
    description     TEXT,
    is_active       BOOL NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX request_types_org_name_active_idx
    ON domain.request_types (organization_id, name)
    WHERE is_active = true;

CREATE INDEX request_types_organization_id_idx
    ON domain.request_types (organization_id);

-- migrate:down
DROP TABLE IF EXISTS domain.request_types;
