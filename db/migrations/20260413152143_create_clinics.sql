-- migrate:up
CREATE TABLE domain.clinics (
    id                UUID PRIMARY KEY,
    organization_id   UUID        NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    name              TEXT        NOT NULL,
    description       TEXT,
    physical_address  domain.address NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX clinics_organization_id_idx
    ON domain.clinics (organization_id);

-- migrate:down
DROP TABLE IF EXISTS domain.clinics;
