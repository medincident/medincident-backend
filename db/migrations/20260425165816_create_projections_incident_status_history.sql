-- migrate:up
CREATE TABLE projections.incident_status_history (
    id                  UUID PRIMARY KEY,
    incident_id         UUID NOT NULL,
    old_status          domain.incident_status,
    new_status          domain.incident_status NOT NULL,
    actor_employee_id   UUID,
    actor_display_name  TEXT,
    changed_at          TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_status_hist_incident_idx
    ON projections.incident_status_history (incident_id, changed_at);

-- migrate:down
DROP TABLE IF EXISTS projections.incident_status_history;
