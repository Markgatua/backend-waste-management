CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name varchar(255) NOT NULL,
  guard_name varchar(255) NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  description text DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL
);

CREATE TABLE permissions (
  id SERIAL PRIMARY KEY,
  name varchar(255) NOT NULL,
  guard_name varchar(255) NOT NULL UNIQUE,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  module varchar(255) NOT NULL,
  submodule VARCHAR(255) NULL
);

CREATE TABLE role_has_permissions (
  permission_id INTEGER,
  role_id INTEGER,
  FOREIGN Key (permission_id) REFERENCES permissions(id),
  FOREIGN Key (role_id) REFERENCES roles(id) on delete CASCADE
);


-- Create "users" table
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NULL,
    last_name VARCHAR(255) NULL,
    provider VARCHAR(255),
    role_id INTEGER,
    FOREIGN Key (role_id) REFERENCES roles(id),
    email VARCHAR(255) DEFAULT NULL UNIQUE,
    password TEXT DEFAULT NULL,
    avatar_url TEXT NULL,
    user_type SMALLINT,
    is_active BOOLEAN DEFAULT TRUE,
    calling_code VARCHAR(6) NULL,
    phone VARCHAR(15) NULL DEFAULT NULL,
    phone_confirmed_at TIMESTAMP NULL DEFAULT NULL,
    confirmed_at TIMESTAMP NULL DEFAULT NULL,
    confirmation_token VARCHAR(255) DEFAULT NULL,
    confirmation_sent_at TIMESTAMP NULL DEFAULT NULL,
    recovery_token VARCHAR(255) DEFAULT NULL,
    recovery_sent_at TIMESTAMP NULL DEFAULT NULL,
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create "phone_verification" table
CREATE TABLE phone_verification_token (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    calling_code TEXT NOT NULL,
    phone TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create "email_verification" table
CREATE TABLE email_verification_token (
    id SERIAL PRIMARY KEY,
  token TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
