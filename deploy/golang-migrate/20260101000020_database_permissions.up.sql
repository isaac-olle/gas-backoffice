-- gas_admin is managed by the parent infrastructure project and needs migration privileges.
GRANT ALL PRIVILEGES ON DATABASE gas_backoffice TO gas_admin;
GRANT ALL PRIVILEGES ON SCHEMA public TO gas_admin;
GRANT USAGE, CREATE ON SCHEMA backoffice, protocol, security, events TO gas_admin;

-- The application connects as gas-backoffice-user.
GRANT CONNECT ON DATABASE gas_backoffice TO "gas-backoffice-user";

GRANT USAGE ON SCHEMA backoffice, protocol, security, events, public TO "gas-backoffice-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA backoffice TO "gas-backoffice-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA protocol TO "gas-backoffice-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA security TO "gas-backoffice-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA events TO "gas-backoffice-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO "gas-backoffice-user";
GRANT USAGE ON ALL SEQUENCES IN SCHEMA backoffice, protocol, security, events, public TO "gas-backoffice-user";

ALTER DEFAULT PRIVILEGES FOR ROLE doadmin IN SCHEMA backoffice
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "gas-backoffice-user";
ALTER DEFAULT PRIVILEGES FOR ROLE doadmin IN SCHEMA protocol
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "gas-backoffice-user";
ALTER DEFAULT PRIVILEGES FOR ROLE doadmin IN SCHEMA security
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "gas-backoffice-user";
ALTER DEFAULT PRIVILEGES FOR ROLE doadmin IN SCHEMA events
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "gas-backoffice-user";
ALTER DEFAULT PRIVILEGES FOR ROLE doadmin IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO "gas-backoffice-user";

