BEGIN TRANSACTION ;
SET ROLE TO shortenerodmen;

SELECT current_user;

DROP TABLE urls;
COMMIT;