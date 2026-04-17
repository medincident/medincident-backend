-- migrate:up

-- No FK to projections.organizations: the query service must tolerate
-- out-of-order event arrival (a ClinicCreated may be delivered before
-- the OrganizationCreated it references). Ordering guarantees come
-- from NATS per-subject, not cross-stream, so enforcing referential
-- integrity at the schema level would livelock the projector on
-- transient races. Consumers that need joined reads can LEFT JOIN on
-- organization_id and check for NULL.
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

-- migrate:down
DROP TABLE IF EXISTS projections.clinics;
