#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB_NAME" <<-EOSQL
  CREATE SCHEMA IF NOT EXISTS image_service;
  CREATE TYPE enum_service AS ENUM('conversion', 'compression');
  ALTER TYPE enum_service SET SCHEMA image_service;
  CREATE TYPE enum_status AS ENUM ('queued', 'processing', 'done', 'processing failed');
  ALTER TYPE enum_status SET SCHEMA image_service;
  CREATE TABLE IF NOT EXISTS image_service.user_account (
      id uuid DEFAULT gen_random_uuid(),
      username character varying(50) NOT NULL,
      password character varying(60) NOT NULL,
      CONSTRAINT user_account_id PRIMARY KEY (id),
      CONSTRAINT user_account_username UNIQUE (username)
    );
  CREATE TABLE IF NOT EXISTS image_service.image(
      id uuid DEFAULT gen_random_uuid(),
      uploaded_name character varying(150) NOT NULL,
      uploaded_location character varying(150) NOT NULL,
      resulted_name character varying(150),
      resulted_location character varying(150),
      CONSTRAINT user_image_id PRIMARY KEY (id)
    );
  CREATE TABLE IF NOT EXISTS image_service.request (
      id uuid DEFAULT gen_random_uuid(),
      user_account_id uuid DEFAULT gen_random_uuid(),
      image_id uuid DEFAULT gen_random_uuid(),
      service_name image_service.enum_service,
      status image_service.enum_status,
      time_started TIMESTAMP,
      time_completed TIMESTAMP,
      CONSTRAINT fk_user_image_user_account_id FOREIGN KEY (user_account_id) REFERENCES image_service.user_account(id),
      CONSTRAINT fk_request_image_id FOREIGN KEY (image_id) REFERENCES image_service.image(id),
      CONSTRAINT request_id PRIMARY KEY (id)
    );
  CREATE ROLE $DB_USER WITH LOGIN ENCRYPTED PASSWORD '$DB_PASSWORD';
  GRANT USAGE ON SCHEMA image_service TO $DB_USER;
  GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA image_service TO $DB_USER;

EOSQL
