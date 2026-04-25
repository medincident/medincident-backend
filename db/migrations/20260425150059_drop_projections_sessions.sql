-- migrate:up
DROP TABLE IF EXISTS projections.sessions;

-- migrate:down
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
