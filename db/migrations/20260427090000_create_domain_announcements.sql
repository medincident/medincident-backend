-- migrate:up

CREATE TABLE domain.announcements (
    id              UUID        PRIMARY KEY,
    organization_id UUID        NOT NULL REFERENCES domain.organizations ON DELETE CASCADE,
    clinic_id       UUID        REFERENCES domain.clinics,
    department_id   UUID        REFERENCES domain.departments,
    author_id       TEXT        NOT NULL,
    title           TEXT        NOT NULL,
    content         TEXT        NOT NULL,
    priority        announcement_priority NOT NULL DEFAULT 'normal',
    is_archived     BOOL        NOT NULL DEFAULT false,
    starts_at       TIMESTAMPTZ,
    ends_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT announcements_dept_requires_clinic CHECK (department_id IS NULL OR clinic_id IS NOT NULL)
);

CREATE INDEX announcements_organization_id_idx ON domain.announcements (organization_id);
CREATE INDEX announcements_clinic_id_idx ON domain.announcements (clinic_id) WHERE clinic_id IS NOT NULL;
CREATE INDEX announcements_department_id_idx ON domain.announcements (department_id) WHERE department_id IS NOT NULL;

-- migrate:down
DROP TABLE IF EXISTS domain.announcements;
