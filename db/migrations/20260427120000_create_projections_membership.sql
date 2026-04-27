-- migrate:up

CREATE TABLE projections.employees (
    id              uuid        PRIMARY KEY,
    zitadel_user_id text        NOT NULL,
    organization_id uuid        NOT NULL,
    clinic_id       uuid        NULL,
    department_id   uuid        NOT NULL,
    position        text        NULL,
    hired_at        timestamptz NOT NULL,
    terminated_at   timestamptz NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX employees_zitadel_user_id_idx
    ON projections.employees (zitadel_user_id);

CREATE INDEX employees_org_dept_idx
    ON projections.employees (organization_id, department_id);

CREATE INDEX employees_active_org_idx
    ON projections.employees (organization_id)
    WHERE terminated_at IS NULL;

CREATE INDEX employees_pending_backfill_idx
    ON projections.employees (department_id)
    WHERE clinic_id IS NULL;

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

CREATE TABLE projections.employee_cards (
    employee_id              uuid PRIMARY KEY,
    zitadel_user_id          text NOT NULL,
    first_name               text NULL,
    last_name                text NULL,
    display_name             text NULL,
    email                    text NULL,
    organization_id          uuid NOT NULL,
    organization_name        text NULL,
    clinic_id                uuid NULL,
    clinic_name              text NULL,
    department_id            uuid NOT NULL,
    department_name          text NULL,
    position                 text NULL,
    terminated_at            timestamptz NULL,
    current_vacation_id      uuid NULL,
    current_vacation_ends_at timestamptz NULL,
    next_vacation_id         uuid NULL,
    next_vacation_starts_at  timestamptz NULL,
    updated_at               timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX employee_cards_zitadel_idx ON projections.employee_cards (zitadel_user_id);
CREATE INDEX employee_cards_org_idx     ON projections.employee_cards (organization_id);
CREATE INDEX employee_cards_clinic_idx  ON projections.employee_cards (clinic_id);
CREATE INDEX employee_cards_dept_idx    ON projections.employee_cards (department_id);
CREATE INDEX employee_cards_email_idx   ON projections.employee_cards (email);
CREATE INDEX employee_cards_display_name_prefix_idx
    ON projections.employee_cards (display_name text_pattern_ops);

CREATE TABLE projections.organization_counters (
    organization_id   uuid PRIMARY KEY,
    employees_total   int NOT NULL DEFAULT 0,
    clinics_total     int NOT NULL DEFAULT 0,
    departments_total int NOT NULL DEFAULT 0,
    updated_at        timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE projections.clinic_counters (
    clinic_id         uuid PRIMARY KEY,
    organization_id   uuid NOT NULL,
    employees_total   int NOT NULL DEFAULT 0,
    departments_total int NOT NULL DEFAULT 0,
    updated_at        timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX clinic_counters_org_idx
    ON projections.clinic_counters (organization_id);

CREATE TABLE projections.department_counters (
    department_id   uuid PRIMARY KEY,
    clinic_id       uuid NULL,
    organization_id uuid NOT NULL,
    employees_total int NOT NULL DEFAULT 0,
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX department_counters_clinic_idx
    ON projections.department_counters (clinic_id);

CREATE INDEX department_counters_org_idx
    ON projections.department_counters (organization_id);

-- migrate:down
DROP TABLE IF EXISTS projections.department_counters;
DROP TABLE IF EXISTS projections.clinic_counters;
DROP TABLE IF EXISTS projections.organization_counters;
DROP TABLE IF EXISTS projections.employee_cards;
DROP TABLE IF EXISTS projections.employee_vacations;
DROP TABLE IF EXISTS projections.employees;
