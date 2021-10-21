#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB_NAME" <<-EOSQL
  CREATE SCHEMA IF NOT EXISTS image_service;
  CREATE TYPE enum_service AS ENUM('conversion', 'compression');
  ALTER TYPE enum_service SET SCHEMA image_service;
  CREATE TYPE enum_status AS ENUM ('queued', 'processing', 'done');
  ALTER TYPE enum_status SET SCHEMA image_service;
  CREATE TABLE IF NOT EXISTS image_service.user_account (
      id uuid DEFAULT gen_random_uuid(),
      username character varying(50) NOT NULL,
      password character varying(60) NOT NULL,
      CONSTRAINT user_account_id PRIMARY KEY (id),
      CONSTRAINT user_account_username UNIQUE (username)
    );
  CREATE TABLE IF NOT EXISTS image_service.uploaded_image(
      id uuid DEFAULT gen_random_uuid(),
      uploaded_name character varying(100) NOT NULL,
      uploaded_location character varying(150) NOT NULL,
      CONSTRAINT uploaded_image_id PRIMARY KEY (id)
    );
  CREATE TABLE IF NOT EXISTS image_service.resulted_image(
      id uuid DEFAULT gen_random_uuid(),
      resulted_name character varying(100) NOT NULL,
      resulted_location character varying(150) NOT NULL,
      service image_service.enum_service,
      CONSTRAINT resulted_image_id PRIMARY KEY (id)
    );
  CREATE TABLE IF NOT EXISTS image_service.user_image(
      id uuid DEFAULT gen_random_uuid(),
      user_account_id uuid DEFAULT gen_random_uuid(),
      uploaded_image_id uuid DEFAULT gen_random_uuid(),
      resulting_image_id uuid DEFAULT gen_random_uuid(),
      status image_service.enum_status,
      CONSTRAINT fk_user_image_user_account_id FOREIGN KEY (user_account_id) REFERENCES image_service.user_account(id),
      CONSTRAINT fk_user_image_uploaded_image_id FOREIGN KEY (uploaded_image_id) REFERENCES image_service.uploaded_image(id),
      CONSTRAINT fk_user_image_resulting_image_id FOREIGN KEY (resulting_image_id) REFERENCES image_service.resulted_image(id),
      CONSTRAINT user_image_id PRIMARY KEY (id)
    );
  CREATE TABLE IF NOT EXISTS image_service.request (
      id uuid DEFAULT gen_random_uuid(),
      user_image_id uuid DEFAULT gen_random_uuid(),
      time_start date NOT NULL,
      end_of_time date NOT NULL,
      CONSTRAINT fk_request_user_image_id FOREIGN KEY (user_image_id) REFERENCES image_service.user_image(id),
      CONSTRAINT request_id PRIMARY KEY (id)
    );
  CREATE ROLE $DB_USER WITH LOGIN ENCRYPTED PASSWORD '$DB_PASSWORD';
  GRANT USAGE ON SCHEMA image_service TO $DB_USER;
  GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA image_service TO $DB_USER;
EOSQL