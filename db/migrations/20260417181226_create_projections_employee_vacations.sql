-- migrate:up
CREATE TABLE projections.employee_vacations (
    id             uuid        PRIMARY KEY,
    employee_id    uuid        NOT NULL,
    state          text        NOT NULL,
    starts_at      timestamptz NOT NULL,
    ends_at        timestamptz NULL,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX employee_vacations_employee_state_idx
    ON projections.employee_vacations (employee_id, state);

CREATE INDEX employee_vacations_active_idx
    ON projections.employee_vacations (employee_id)
    WHERE state = 'active';

CREATE INDEX employee_vacations_scheduled_idx
    ON projections.employee_vacations (employee_id)
    WHERE state = 'scheduled';

-- migrate:down
DROP TABLE IF EXISTS projections.employee_vacations;
