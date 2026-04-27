-- migrate:up

CREATE TABLE domain.organizations (
    id                UUID PRIMARY KEY,
    name              TEXT        NOT NULL,
    description       TEXT,
    legal_address     domain.address NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

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

CREATE TABLE domain.departments (
    id          UUID PRIMARY KEY,
    clinic_id   UUID        NOT NULL
        REFERENCES domain.clinics(id) ON DELETE RESTRICT,
    name        TEXT        NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX departments_clinic_id_idx
    ON domain.departments (clinic_id);

-- migrate:down
DROP TABLE IF EXISTS domain.departments;
DROP TABLE IF EXISTS domain.clinics;
DROP TABLE IF EXISTS domain.organizations;
