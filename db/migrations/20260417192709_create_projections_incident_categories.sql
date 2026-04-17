-- migrate:up

-- Read-model copy of domain.incident_categories. No cross-schema FKs.
-- The unique index on (organization_id, name) WHERE is_active is
-- preserved because it is load-bearing for read-side queries that
-- already rely on it.
CREATE TABLE projections.incident_categories (
    id                  UUID        PRIMARY KEY,
    organization_id     UUID        NOT NULL,
    parent_category_id  UUID,
    name                TEXT        NOT NULL,
    description         TEXT,
    is_active           BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX incident_categories_organization_id_idx
    ON projections.incident_categories (organization_id);

CREATE INDEX incident_categories_parent_category_id_idx
    ON projections.incident_categories (parent_category_id);

-- Matches the domain-side partial unique index so read queries can
-- rely on the same invariant (one active category name per org).
CREATE UNIQUE INDEX incident_categories_org_active_name_uniq
    ON projections.incident_categories (organization_id, name)
    WHERE is_active;

-- migrate:down
DROP TABLE IF EXISTS projections.incident_categories;
