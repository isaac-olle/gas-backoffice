-- Disposable credentials for this Compose project only. The normal migrations
-- grant this role the same schema/table privileges used by the application.
CREATE ROLE "gas-backoffice-user" LOGIN PASSWORD 'integration-app';
