1. create a keyspace

CREATE KEYSPACE hanzo WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};

2. create filemeta table

 USE hanzo;

 CREATE TABLE filemeta (
    dirhash bigint,
    directory varchar,
    name varchar,
    meta blob,
    PRIMARY KEY ((dirhash, directory), name)
 ) WITH CLUSTERING ORDER BY (name ASC);
