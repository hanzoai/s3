
1. create "hanzo" database

export PGHOME=/Library/PostgreSQL/10
$PGHOME/bin/createdb --username=postgres --password hanzo

2. create "filemeta" table
$PGHOME/bin/psql --username=postgres --password hanzo

CREATE TABLE IF NOT EXISTS filemeta (
  dirhash     BIGINT,
  name        VARCHAR(65535),
  directory   VARCHAR(65535),
  meta        bytea,
  PRIMARY KEY (dirhash, name)
);

