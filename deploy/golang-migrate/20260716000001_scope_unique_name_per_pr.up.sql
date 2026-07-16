CREATE UNIQUE INDEX IF NOT EXISTS scopes_name_protected_resource_id_unique
    ON backoffice.scopes (protected_resource_id, name)
    WHERE deleted_at IS NULL;
