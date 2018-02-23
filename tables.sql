CREATE DATABASE urlShorter;

use urlShorter;

CREATE TABLE urls (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  url VARCHAR(1024) NOT NULL
);

CREATE TABLE access_log (
  id                  INTEGER PRIMARY KEY AUTO_INCREMENT,
  url_id              INTEGER,
  time                TIMESTAMP,
  remote_ip           VARCHAR(64),
  forwarded_ip        VARCHAR(64),
  UA                  VARCHAR(1024),
  referer             VARCHAR(512)
);

CREATE TABLE insert_log (
  id                  INTEGER PRIMARY KEY AUTO_INCREMENT,
  url_id              INTEGER,
  time                TIMESTAMP,
  remote_ip           VARCHAR(64),
  forwarded_ip        VARCHAR(64),
  UA                  VARCHAR(1024),
  referer             VARCHAR(512)
);
