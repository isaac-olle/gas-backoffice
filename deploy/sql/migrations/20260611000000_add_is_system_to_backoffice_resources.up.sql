ALTER TABLE admins
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE tenants
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE protected_resources
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE scopes
    ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE admins              SET is_system = TRUE WHERE name = 'Gas-Seed-Admin';
UPDATE tenants             SET is_system = TRUE WHERE name = 'Gas-Tenant';
UPDATE protected_resources SET is_system = TRUE WHERE name = 'Gas-Backoffice';
UPDATE scopes              SET is_system = TRUE WHERE name IN ('tenants:create', 'tenants:read', 'tenants:read_all');
