CREATE DATABASE IF NOT EXISTS s3;
CREATE USER IF NOT EXISTS 's3'@'%' IDENTIFIED BY 'secret';
GRANT ALL PRIVILEGES ON s3.* TO 's3'@'%';
FLUSH PRIVILEGES;
USE s3;
CREATE TABLE IF NOT EXISTS `filemeta` (
    `dirhash`   BIGINT NOT NULL       COMMENT 'first 64 bits of MD5 hash value of directory field',
    `name`      VARCHAR(766) NOT NULL COMMENT 'directory or file name',
    `directory` TEXT NOT NULL         COMMENT 'full path to parent directory',
    `meta`      LONGBLOB,
    PRIMARY KEY (`dirhash`, `name`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;