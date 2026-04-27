-- migrate:up

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE domain.employee_vacations (
    id           UUID        PRIMARY KEY,
    employee_id  UUID        NOT NULL
        REFERENCES domain.employees(id) ON DELETE CASCADE,
    starts_at    TIMESTAMPTZ NOT NULL,
    ends_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT employee_vacations_no_overlap
        EXCLUDE USING GIST (
            employee_id WITH =,
            tstzrange(starts_at, coalesce(ends_at, 'infinity'::timestamptz), '[)') WITH &&
        )
);

CREATE INDEX employee_vacations_employee_id_idx
    ON domain.employee_vacations (employee_id);

-- migrate:down
DROP TABLE IF EXISTS domain.employee_vacations;
