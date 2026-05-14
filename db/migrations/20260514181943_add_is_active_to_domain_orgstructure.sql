-- migrate:up
ALTER TABLE domain.organizations ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE domain.clinics       ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE domain.departments   ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE domain.employees     ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- migrate:down
ALTER TABLE domain.organizations DROP COLUMN is_active;
ALTER TABLE domain.clinics       DROP COLUMN is_active;
ALTER TABLE domain.departments   DROP COLUMN is_active;
ALTER TABLE domain.employees     DROP COLUMN is_active;
