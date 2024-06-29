BEGIN TRANSACTION;

DROP INDEX IF EXISTS idx_short_is_not_deleted;

ALTER TABLE urls
DROP COLUMN is_deleted,
DROP COLUMN user_id;

COMMIT;