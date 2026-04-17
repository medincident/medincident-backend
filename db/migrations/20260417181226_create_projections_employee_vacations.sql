-- migrate:up
CREATE TABLE projections.employee_vacations (
    id             uuid        PRIMARY KEY,
    employee_id    uuid        NOT NULL,
    state          text        NOT NULL CHECK (state IN ('scheduled','active','ended','cancelled')),
    starts_at      timestamptz NOT NULL,
    ends_at        timestamptz NULL,
    event_sequence bigint      NOT NULL DEFAULT 0,
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
DROP TABLE projections.employee_vacations;
