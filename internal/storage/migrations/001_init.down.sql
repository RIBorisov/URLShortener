BEGIN TRANSACTION ;

DROP INDEX IF EXISTS idx_short_is_not_deleted;

DROP INDEX IF EXISTS idx_long_is_not_deleted;

DROP TABLE urls;

COMMIT;