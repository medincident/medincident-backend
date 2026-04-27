-- migrate:up

-- Department responsibles.
CREATE TABLE domain.department_responsibles (
    department_id      UUID NOT NULL
        REFERENCES domain.departments(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (department_id, employee_id)
);

CREATE UNIQUE INDEX department_responsibles_holder_uniq
    ON domain.department_responsibles (employee_id);

CREATE INDEX department_responsibles_deputy_employee_id_idx
    ON domain.department_responsibles (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- Clinic heads.
CREATE TABLE domain.clinic_heads (
    clinic_id          UUID NOT NULL
        REFERENCES domain.clinics(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (clinic_id, employee_id)
);

-- Load-bearing invariant: an employee can hold at most one ClinicHead
-- row across the whole table. The service layer relies on this to
-- detect 'already assigned' via the unique-violation SQLSTATE. Do not
-- drop this index without also rewriting the duplicate-detection logic
-- in services/membership/clinic_head_assign.go.
CREATE UNIQUE INDEX clinic_heads_holder_uniq
    ON domain.clinic_heads (employee_id);

CREATE INDEX clinic_heads_deputy_employee_id_idx
    ON domain.clinic_heads (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- Organization admins.
CREATE TABLE domain.org_admins (
    organization_id    UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE UNIQUE INDEX org_admins_holder_uniq
    ON domain.org_admins (employee_id);

CREATE INDEX org_admins_deputy_employee_id_idx
    ON domain.org_admins (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- Organization heads.
CREATE TABLE domain.org_heads (
    organization_id    UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE UNIQUE INDEX org_heads_holder_uniq
    ON domain.org_heads (employee_id);

CREATE INDEX org_heads_deputy_employee_id_idx
    ON domain.org_heads (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- Organization dispatchers.
CREATE TABLE domain.org_dispatchers (
    organization_id    UUID NOT NULL
        REFERENCES domain.organizations(id) ON DELETE RESTRICT,
    employee_id        UUID NOT NULL
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    deputy_employee_id UUID
        REFERENCES domain.employees(id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE UNIQUE INDEX org_dispatchers_holder_uniq
    ON domain.org_dispatchers (employee_id);

CREATE INDEX org_dispatchers_deputy_employee_id_idx
    ON domain.org_dispatchers (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

-- System admins.
CREATE TABLE domain.system_admins (
    zitadel_user_id TEXT        PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS domain.system_admins;
DROP TABLE IF EXISTS domain.org_dispatchers;
DROP TABLE IF EXISTS domain.org_heads;
DROP TABLE IF EXISTS domain.org_admins;
DROP TABLE IF EXISTS domain.clinic_heads;
DROP TABLE IF EXISTS domain.department_responsibles;
