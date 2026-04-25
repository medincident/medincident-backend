-- migrate:up
-- Reverse-lookup of incidents from a published buffer entry. Without
-- this partial index any "find the incident materialised from this
-- buffer row" query falls back to a full sequential scan.
CREATE INDEX incidents_source_buffer_id_idx
    ON domain.incidents (source_buffer_id)
    WHERE source_buffer_id IS NOT NULL;

CREATE INDEX proj_incidents_source_buffer_id_idx
    ON projections.incidents (source_buffer_id)
    WHERE source_buffer_id IS NOT NULL;

-- migrate:down
DROP INDEX IF EXISTS projections.proj_incidents_source_buffer_id_idx;
DROP INDEX IF EXISTS domain.incidents_source_buffer_id_idx;
