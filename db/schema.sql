CREATE DATABASE shred_db;
\c shred_db

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS shred_user (
  user_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  first_name VARCHAR(45) NOT NULL,
  last_name VARCHAR(45) NOT NULL,
  email VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by UUID NOT NULL,
  PRIMARY KEY (user_uuid),
  FOREIGN KEY (created_by) REFERENCES shred_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS muscle_type (
  muscle_code VARCHAR(45) NOT NULL,
  muscle_name VARCHAR(45) NOT NULL,
  muscle_description VARCHAR(2500) NULL, -- description of the muscle as markdown
  muscle_group VARCHAR(45) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by UUID NOT NULL,
  PRIMARY KEY (muscle_code),
  FOREIGN KEY (created_by) REFERENCES shred_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS category_type (
  category_code VARCHAR(45) NOT NULL,
  category_name VARCHAR(45) NOT NULL,
  category_description VARCHAR(2500) NULL, -- description of the category as markdown
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by UUID NOT NULL,
  PRIMARY KEY (category_code),
  FOREIGN KEY (created_by) REFERENCES shred_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS apparatus_type (
  apparatus_code VARCHAR(45) NOT NULL,
  apparatus_name VARCHAR(45) NOT NULL,
  apparatus_description VARCHAR(2500) NULL, -- description of the apparatus as markdown
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by UUID NOT NULL,
  PRIMARY KEY (apparatus_code),
  FOREIGN KEY (created_by) REFERENCES shred_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS license (
  license_short_name VARCHAR(45) NOT NULL,
  license_full_name VARCHAR(45) NOT NULL,
  url VARCHAR(250) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by UUID NOT NULL,
  PRIMARY KEY (license_short_name),
  FOREIGN KEY (created_by) REFERENCES shred_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS exercise (
  exercise_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
  exercise_name VARCHAR(100) NOT NULL,
  exercise_description VARCHAR(2500) NOT NULL,
  instructions VARCHAR(2500),
  cues VARCHAR(2500),
  video_url VARCHAR(256),
  category_code VARCHAR(45) NOT NULL,
  license_short_name VARCHAR(45) NULL,
  license_author VARCHAR(100) NULL,
  created_by UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (exercise_uuid),
  FOREIGN KEY (created_by) REFERENCES shred_user(user_uuid),
  FOREIGN KEY (category_code) REFERENCES category_type(category_code),
  FOREIGN KEY (license_short_name) REFERENCES license(license_short_name)
);

CREATE TABLE IF NOT EXISTS exercise_apparatus (
  exercise_uuid UUID NOT NULL,
  apparatus_code VARCHAR(45) NOT NULL,
  PRIMARY KEY (exercise_uuid, apparatus_code),
  FOREIGN KEY (exercise_uuid) REFERENCES exercise(exercise_uuid),
  FOREIGN KEY (apparatus_code) REFERENCES apparatus_type(apparatus_code)
);

CREATE TABLE IF NOT EXISTS exercise_muscle (
  exercise_uuid UUID NOT NULL,
  muscle_code VARCHAR(45) NOT NULL,
  PRIMARY KEY (exercise_uuid, muscle_code),
  FOREIGN KEY (exercise_uuid) REFERENCES exercise(exercise_uuid),
  FOREIGN KEY (muscle_code) REFERENCES muscle_type(muscle_code)
);
