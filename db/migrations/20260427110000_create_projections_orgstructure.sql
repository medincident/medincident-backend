-- migrate:up

CREATE TABLE projections.organizations (
    id                       UUID PRIMARY KEY,
    name                     TEXT        NOT NULL,
    description              TEXT,
    legal_address_text       TEXT        NOT NULL,
    legal_address_longitude  DOUBLE PRECISION,
    legal_address_latitude   DOUBLE PRECISION,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX organizations_name_idx
    ON projections.organizations USING btree (name);

-- No FK to projections.organizations: the query service must tolerate
-- out-of-order event arrival. Ordering guarantees come from NATS
-- per-subject, not cross-stream.
CREATE TABLE projections.clinics (
    id                          UUID PRIMARY KEY,
    organization_id             UUID        NOT NULL,
    name                        TEXT        NOT NULL,
    description                 TEXT,
    physical_address_text       TEXT        NOT NULL,
    physical_address_longitude  DOUBLE PRECISION,
    physical_address_latitude   DOUBLE PRECISION,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX clinics_organization_id_idx
    ON projections.clinics (organization_id);
CREATE INDEX clinics_name_idx
    ON projections.clinics USING btree (name);

CREATE TABLE projections.departments (
    id          UUID PRIMARY KEY,
    clinic_id   UUID        NOT NULL,
    name        TEXT        NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX departments_clinic_id_idx
    ON projections.departments (clinic_id);
CREATE INDEX departments_name_idx
    ON projections.departments USING btree (name);

-- migrate:down
DROP TABLE IF EXISTS projections.departments;
DROP TABLE IF EXISTS projections.clinics;
DROP TABLE IF EXISTS projections.organizations;
