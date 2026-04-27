-- migrate:up

-- Zitadel user ids are opaque strings (snowflake-like), so id is TEXT
-- rather than UUID. gender is stored as smallint preserving the proto
-- enum value (0 = unspecified, 1 = female, 2 = male, 3 = diverse).
CREATE TABLE projections.users (
    id                  TEXT PRIMARY KEY,
    user_name           TEXT        NOT NULL,
    first_name          TEXT        NOT NULL,
    last_name           TEXT        NOT NULL,
    display_name        TEXT        NOT NULL,
    nick_name           TEXT,
    email               TEXT        NOT NULL,
    email_verified      BOOLEAN     NOT NULL DEFAULT FALSE,
    preferred_language  TEXT        NOT NULL DEFAULT '',
    gender              SMALLINT    NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX users_email_idx
    ON projections.users USING btree (email);
CREATE INDEX users_user_name_idx
    ON projections.users USING btree (user_name);

-- migrate:down
DROP TABLE IF EXISTS projections.users;
