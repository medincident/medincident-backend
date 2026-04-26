-- migrate:up
CREATE TABLE domain.service_request_executors (
    id              UUID PRIMARY KEY,
    request_id      UUID NOT NULL
        REFERENCES domain.service_requests(id) ON DELETE CASCADE,
    employee_id     UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    assigned_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    assigned_by_id  TEXT NOT NULL
);

CREATE UNIQUE INDEX service_request_executors_request_employee_idx
    ON domain.service_request_executors (request_id, employee_id);

CREATE INDEX service_request_executors_request_id_idx
    ON domain.service_request_executors (request_id);

-- migrate:down
DROP TABLE IF EXISTS domain.service_request_executors;
