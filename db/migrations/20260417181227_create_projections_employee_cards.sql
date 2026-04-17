-- migrate:up
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
    event_sequence           bigint NOT NULL DEFAULT 0,
    updated_at               timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX employee_cards_zitadel_idx ON projections.employee_cards (zitadel_user_id);
CREATE INDEX employee_cards_org_idx     ON projections.employee_cards (organization_id);
CREATE INDEX employee_cards_clinic_idx  ON projections.employee_cards (clinic_id);
CREATE INDEX employee_cards_dept_idx    ON projections.employee_cards (department_id);
CREATE INDEX employee_cards_email_idx   ON projections.employee_cards (email);
CREATE INDEX employee_cards_display_name_prefix_idx
    ON projections.employee_cards (display_name text_pattern_ops);

-- migrate:down
DROP TABLE projections.employee_cards;
