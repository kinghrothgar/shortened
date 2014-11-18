CREATE DATABASE shortened;
\c shortened

CREATE USER shortened WITH PASSWORD 'shortened';

GRANT ALL PRIVILEGES ON DATABASE shortened TO shortened;

CREATE TABLE urls (
  id bigserial primary key,
  url text not null
);

GRANT ALL PRIVILEGES ON TABLE urls TO shortened;
GRANT ALL PRIVILEGES ON SEQUENCE urls_id_seq TO SHORTENED;

CREATE TABLE url_usages (
  id serial primary key,
  url_id bigint not null,
  visits int not null default 0
);

GRANT ALL PRIVILEGES ON TABLE url_usages TO shortened;
GRANT ALL PRIVILEGES ON SEQUENCE url_usages_id_seq TO shortened;
