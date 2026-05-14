-- migrate:up
ALTER TABLE projections.organizations ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE projections.clinics       ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE projections.departments   ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE projections.employees     ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- migrate:down
ALTER TABLE projections.organizations DROP COLUMN is_active;
ALTER TABLE projections.clinics       DROP COLUMN is_active;
ALTER TABLE projections.departments   DROP COLUMN is_active;
ALTER TABLE projections.employees     DROP COLUMN is_active;
