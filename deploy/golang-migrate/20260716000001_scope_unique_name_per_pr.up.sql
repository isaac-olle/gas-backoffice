CREATE UNIQUE INDEX scopes_name_protected_resource_id_unique
    ON scopes (protected_resource_id, name)
    WHERE deleted_at IS NULL;
