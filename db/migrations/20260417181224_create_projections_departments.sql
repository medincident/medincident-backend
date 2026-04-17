-- migrate:up

-- See projections.clinics migration for the rationale behind omitting
-- cross-aggregate foreign keys.
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
