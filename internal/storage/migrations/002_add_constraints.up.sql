BEGIN TRANSACTION;
CREATE UNIQUE INDEX IF NOT EXISTS idx_long_url ON urls (long);

COMMIT;