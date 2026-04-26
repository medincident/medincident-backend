-- migrate:up
CREATE TABLE projections.service_request_executor_history (
    id            UUID PRIMARY KEY,
    request_id    UUID NOT NULL,
    action        TEXT NOT NULL,
    employee_id   UUID NOT NULL,
    employee_name TEXT NOT NULL,
    actor_id      TEXT NOT NULL,
    actor_name    TEXT NOT NULL,
    changed_at    TIMESTAMPTZ NOT NULL
);

CREATE INDEX proj_sr_exec_hist_request_idx
    ON projections.service_request_executor_history (request_id, changed_at);

-- migrate:down
DROP TABLE IF EXISTS projections.service_request_executor_history;
