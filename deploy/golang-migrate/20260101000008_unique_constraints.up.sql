CREATE UNIQUE INDEX IF NOT EXISTS protected_resources_subscription_id_name_unique
    ON backoffice.protected_resources (subscription_id, name)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS clients_subscription_id_name_unique
    ON backoffice.clients (subscription_id, name)
    WHERE deleted_at IS NULL;
