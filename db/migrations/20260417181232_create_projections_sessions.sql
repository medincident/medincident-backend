-- migrate:up

-- Zitadel session ids are opaque strings, so id is TEXT. user_agent is
-- stored as JSONB because it is a nested message (IP, header map,
-- optional fingerprint) that readers currently only need to round-trip.
-- Keeping it opaque avoids a schema migration every time UserAgent
-- evolves on the proto side.
CREATE TABLE projections.sessions (
    id                    TEXT PRIMARY KEY,
    user_id               TEXT,
    user_resource_owner   TEXT,
    preferred_language    TEXT,
    checked_at            TIMESTAMPTZ,
    user_agent            JSONB       NOT NULL DEFAULT '{}'::jsonb,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX sessions_user_id_idx
    ON projections.sessions (user_id)
    WHERE user_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS projections.sessions;
