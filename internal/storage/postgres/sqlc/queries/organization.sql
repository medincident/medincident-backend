-- name: GetOrganization :one
SELECT
    id,
    name,
    description,
    legal_address_text,
    legal_address_lng,
    legal_address_lat,
    created_at,
    updated_at
FROM domain.organizations
WHERE id = $1;

-- name: UpsertOrganization :exec
INSERT INTO domain.organizations (
    id, name, description,
    legal_address_text, legal_address_lng, legal_address_lat,
    created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET
    name               = EXCLUDED.name,
    description        = EXCLUDED.description,
    legal_address_text = EXCLUDED.legal_address_text,
    legal_address_lng  = EXCLUDED.legal_address_lng,
    legal_address_lat  = EXCLUDED.legal_address_lat,
    updated_at         = EXCLUDED.updated_at;
