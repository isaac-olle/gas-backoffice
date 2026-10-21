CREATE TABLE backoffice.tenant_profiles (
    user_id  TEXT NOT NULL PRIMARY KEY REFERENCES backoffice.users(id),
    admin_id TEXT NOT NULL
);

INSERT INTO backoffice.tenant_profiles (user_id, admin_id)
SELECT id, admin_id FROM backoffice.users WHERE role = 'tenant';

ALTER TABLE backoffice.users
    DROP COLUMN admin_id,
    DROP COLUMN author_id,
    DROP COLUMN author_role;
