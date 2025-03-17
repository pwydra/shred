CREATE DATABASE shred_db;
\c shred_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS exercise (
  uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  exercise_name VARCHAR(45) NOT NULL,
  description VARCHAR(45) NOT NULL,
  instructions VARCHAR(250),
  category VARCHAR(45),
  cues VARCHAR(250),
  primary_muscles VARCHAR(45),
  secondary_muscles VARCHAR(45),
  front_image VARCHAR(256),
  back_image VARCHAR(256),
  video_url VARCHAR(256),
  apparatus VARCHAR(256),
  license VARCHAR(256),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_uuid UUID NOT NULL,
  PRIMARY KEY (uuid)
);

