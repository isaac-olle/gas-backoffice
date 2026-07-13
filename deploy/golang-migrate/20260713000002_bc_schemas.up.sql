CREATE SCHEMA IF NOT EXISTS backoffice;
CREATE SCHEMA IF NOT EXISTS protocol;
CREATE SCHEMA IF NOT EXISTS security;

ALTER TABLE admins SET SCHEMA backoffice;
ALTER TABLE users SET SCHEMA backoffice;
ALTER TABLE protected_resources SET SCHEMA backoffice;
ALTER TABLE scopes SET SCHEMA backoffice;
ALTER TABLE clients SET SCHEMA backoffice;

ALTER TABLE protocol_sessions SET SCHEMA protocol;
ALTER TABLE protocol_client_credentials_sessions SET SCHEMA protocol;

ALTER TABLE key_pairs SET SCHEMA security;
ALTER TABLE admin_assertion_keys SET SCHEMA security;
ALTER TABLE client_assertion_keys SET SCHEMA security;
ALTER TABLE client_secrets SET SCHEMA security;

GRANT ALL ON SCHEMA backoffice TO gas_backoffice_admin;
GRANT ALL ON SCHEMA protocol TO gas_backoffice_admin;
GRANT ALL ON SCHEMA security TO gas_backoffice_admin;
GRANT ALL ON ALL TABLES IN SCHEMA backoffice TO gas_backoffice_admin;
GRANT ALL ON ALL TABLES IN SCHEMA protocol TO gas_backoffice_admin;
GRANT ALL ON ALL TABLES IN SCHEMA security TO gas_backoffice_admin;

GRANT USAGE ON SCHEMA backoffice TO gas_backoffice_service;
GRANT USAGE ON SCHEMA protocol TO gas_backoffice_service;
GRANT USAGE ON SCHEMA security TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA backoffice TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA protocol TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA security TO gas_backoffice_service;
