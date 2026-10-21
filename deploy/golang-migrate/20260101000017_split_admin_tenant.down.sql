ALTER TABLE backoffice.users
    ADD COLUMN admin_id    TEXT,
    ADD COLUMN author_id   TEXT NOT NULL DEFAULT '',
    ADD COLUMN author_role TEXT NOT NULL DEFAULT '';

UPDATE backoffice.users u
SET admin_id = tp.admin_id
FROM backoffice.tenant_profiles tp
WHERE tp.user_id = u.id;

-- An admin is its own admin (see decision 2026-10-12): backfill admin_id for
-- admin rows, which have no tenant_profiles row.
UPDATE backoffice.users SET admin_id = id WHERE role = 'admin' AND admin_id IS NULL;

DROP TABLE backoffice.tenant_profiles;
