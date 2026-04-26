-- migrate:up
CREATE TABLE projections.request_types (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT,
    is_active       BOOL NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_request_types_org_idx
    ON projections.request_types (organization_id);

-- migrate:down
DROP TABLE IF EXISTS projections.request_types;
