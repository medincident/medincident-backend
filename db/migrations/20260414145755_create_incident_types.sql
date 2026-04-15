-- migrate:up
CREATE TABLE domain.incident_types (
    id               UUID        PRIMARY KEY,
    organization_id  UUID        NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    category_id      UUID        NOT NULL
        REFERENCES domain.incident_categories(id) ON DELETE CASCADE,
    name             TEXT        NOT NULL,
    description      TEXT,
    is_active        BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX incident_types_organization_id_idx
    ON domain.incident_types (organization_id);

CREATE INDEX incident_types_category_id_idx
    ON domain.incident_types (category_id);

CREATE UNIQUE INDEX incident_types_org_active_name_uniq
    ON domain.incident_types (organization_id, name)
    WHERE is_active;

-- migrate:down
DROP TABLE IF EXISTS domain.incident_types;
