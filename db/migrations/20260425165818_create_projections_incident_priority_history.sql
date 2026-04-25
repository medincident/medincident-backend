-- migrate:up
CREATE TABLE projections.incident_priority_history (
    id                  UUID PRIMARY KEY,
    incident_id         UUID NOT NULL,
    old_priority        domain.incident_priority NOT NULL,
    new_priority        domain.incident_priority NOT NULL,
    actor_employee_id   UUID,
    actor_display_name  TEXT,
    changed_at          TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_priority_hist_incident_idx
    ON projections.incident_priority_history (incident_id, changed_at);

-- migrate:down
DROP TABLE IF EXISTS projections.incident_priority_history;
