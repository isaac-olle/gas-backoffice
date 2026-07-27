GRANT USAGE ON SCHEMA backoffice TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA backoffice TO gas_backoffice_service;

GRANT USAGE ON SCHEMA protocol TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA protocol TO gas_backoffice_service;

GRANT USAGE ON SCHEMA security TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA security TO gas_backoffice_service;

GRANT USAGE ON SCHEMA events TO gas_backoffice_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA events TO gas_backoffice_service;
