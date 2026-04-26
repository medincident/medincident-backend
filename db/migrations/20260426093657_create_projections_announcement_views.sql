-- migrate:up

CREATE TABLE projections.announcement_views (
    announcement_id UUID   PRIMARY KEY REFERENCES domain.announcements ON DELETE CASCADE,
    view_count      BIGINT NOT NULL DEFAULT 0
);

-- migrate:down

DROP TABLE projections.announcement_views;
