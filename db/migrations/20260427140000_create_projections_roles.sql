-- migrate:up

-- Read-model copies of domain role tables. No cross-schema FKs.

CREATE TABLE projections.clinic_heads (
    clinic_id          UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (clinic_id, employee_id)
);

CREATE INDEX clinic_heads_employee_id_idx
    ON projections.clinic_heads (employee_id);
CREATE INDEX clinic_heads_deputy_employee_id_idx
    ON projections.clinic_heads (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

CREATE TABLE projections.department_responsibles (
    department_id      UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (department_id, employee_id)
);

CREATE INDEX department_responsibles_employee_id_idx
    ON projections.department_responsibles (employee_id);
CREATE INDEX department_responsibles_deputy_employee_id_idx
    ON projections.department_responsibles (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

CREATE TABLE projections.org_admins (
    organization_id    UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE INDEX org_admins_employee_id_idx
    ON projections.org_admins (employee_id);
CREATE INDEX org_admins_deputy_employee_id_idx
    ON projections.org_admins (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

CREATE TABLE projections.org_dispatchers (
    organization_id    UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE INDEX org_dispatchers_employee_id_idx
    ON projections.org_dispatchers (employee_id);
CREATE INDEX org_dispatchers_deputy_employee_id_idx
    ON projections.org_dispatchers (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

CREATE TABLE projections.org_heads (
    organization_id    UUID NOT NULL,
    employee_id        UUID NOT NULL,
    deputy_employee_id UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (organization_id, employee_id)
);

CREATE INDEX org_heads_employee_id_idx
    ON projections.org_heads (employee_id);
CREATE INDEX org_heads_deputy_employee_id_idx
    ON projections.org_heads (deputy_employee_id)
    WHERE deputy_employee_id IS NOT NULL;

CREATE TABLE projections.system_admins (
    zitadel_user_id TEXT        PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS projections.system_admins;
DROP TABLE IF EXISTS projections.org_heads;
DROP TABLE IF EXISTS projections.org_dispatchers;
DROP TABLE IF EXISTS projections.org_admins;
DROP TABLE IF EXISTS projections.department_responsibles;
DROP TABLE IF EXISTS projections.clinic_heads;
