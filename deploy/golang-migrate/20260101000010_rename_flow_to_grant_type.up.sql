ALTER TABLE backoffice.clients
    RENAME COLUMN allowed_flows TO allowed_grant_types;

ALTER TABLE protocol.protocol_sessions
    RENAME COLUMN flow TO grant_type;
