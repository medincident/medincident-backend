-- migrate:up

-- Switch classifier FKs from ON DELETE CASCADE to ON DELETE RESTRICT
-- so direct row deletions can no longer silently wipe out subtrees
-- (and their outbox events) without going through the explicit
-- IncidentCategoryService.Delete cleanup. The service now iterates the
-- subtree children-first, emitting one Deleted event per row and
-- issuing an explicit DELETE for each row under the same transaction.

ALTER TABLE domain.incident_categories
    DROP CONSTRAINT incident_categories_parent_category_id_fkey;

ALTER TABLE domain.incident_categories
    ADD CONSTRAINT incident_categories_parent_category_id_fkey
        FOREIGN KEY (parent_category_id)
        REFERENCES domain.incident_categories(id)
        ON DELETE RESTRICT;

ALTER TABLE domain.incident_types
    DROP CONSTRAINT incident_types_category_id_fkey;

ALTER TABLE domain.incident_types
    ADD CONSTRAINT incident_types_category_id_fkey
        FOREIGN KEY (category_id)
        REFERENCES domain.incident_categories(id)
        ON DELETE RESTRICT;

-- migrate:down

ALTER TABLE domain.incident_categories
    DROP CONSTRAINT incident_categories_parent_category_id_fkey;

ALTER TABLE domain.incident_categories
    ADD CONSTRAINT incident_categories_parent_category_id_fkey
        FOREIGN KEY (parent_category_id)
        REFERENCES domain.incident_categories(id)
        ON DELETE CASCADE;

ALTER TABLE domain.incident_types
    DROP CONSTRAINT incident_types_category_id_fkey;

ALTER TABLE domain.incident_types
    ADD CONSTRAINT incident_types_category_id_fkey
        FOREIGN KEY (category_id)
        REFERENCES domain.incident_categories(id)
        ON DELETE CASCADE;
