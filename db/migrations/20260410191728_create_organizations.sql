-- migrate:up
CREATE TABLE domain.organizations (
    id                 UUID PRIMARY KEY,
    name               TEXT NOT NULL,
    description        TEXT,
    legal_address_text TEXT,
    legal_address_lng  DOUBLE PRECISION,
    legal_address_lat  DOUBLE PRECISION,
    created_at         TIMESTAMPTZ NOT NULL,
    updated_at         TIMESTAMPTZ NOT NULL,

    -- Address text, if present, must be non-empty. Mirrors
    -- geo.NewAddress which trims and rejects empty input.
    CONSTRAINT organizations_legal_address_text_non_empty
        CHECK (legal_address_text IS NULL OR btrim(legal_address_text) <> ''),

    -- Point coordinates are both present or both absent. Never one without the other.
    CONSTRAINT organizations_legal_address_point_both_or_neither
        CHECK (
            (legal_address_lng IS NULL AND legal_address_lat IS NULL)
            OR (legal_address_lng IS NOT NULL AND legal_address_lat IS NOT NULL)
        ),

    -- A point without a text makes no sense: the domain Address VO always
    -- has text, and Point is optional INSIDE Address, not independent of it.
    CONSTRAINT organizations_legal_address_point_requires_text
        CHECK (
            legal_address_lng IS NULL
            OR legal_address_text IS NOT NULL
        )
);

CREATE INDEX organizations_legal_address_text_idx
    ON domain.organizations (legal_address_text)
    WHERE legal_address_text IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS domain.organizations;
