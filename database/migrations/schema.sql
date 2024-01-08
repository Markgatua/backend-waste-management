CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name varchar(255) NOT NULL,
  guard_name varchar(255) NOT NULL,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  description text DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL
);

CREATE TABLE uploads(
    id SERIAL PRIMARY KEY,
    item_id INTEGER, 
    type VARCHAR(100),
    path TEXT,
    related_table VARCHAR(150),
    meta JSON NULL
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

CREATE TABLE countries (
  id SERIAL PRIMARY KEY,
  name varchar(255) NOT NULL,
  currency_code varchar(255) NULL,
  capital varchar(255) NULL,
  citizenship  varchar(255) NOT NULL,
  country_code varchar(3) NOT NULL UNIQUE,
  currency varchar(255) NULL,
  currency_sub_unit varchar(255) NULL,
  currency_symbol varchar(3)  NULL,
  currency_decimals SMALLINT  NULL,
  full_name varchar(255)  NULL,
  iso_3166_2 varchar(2) NOT NULL DEFAULT '',
  iso_3166_3 varchar(3) NOT NULL DEFAULT '',
  region_code varchar(3) NOT NULL DEFAULT '',
  sub_region_code varchar(3) NOT NULL DEFAULT '',
  eea SMALLINT  DEFAULT 0,
  calling_code varchar(3) NULL,
  flag varchar(6) DEFAULT NULL
);

CREATE TABLE role_has_permissions (
  permission_id INTEGER,
  role_id INTEGER,
  FOREIGN Key (permission_id) REFERENCES permissions(id),
  FOREIGN Key (role_id) REFERENCES roles(id) on delete CASCADE
);

CREATE TABLE organizations(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  country_id INTEGER NOT NULL,
  FOREIGN Key (country_id) REFERENCES countries(id) on delete set null
);

CREATE TABLE companies (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  company_type INTEGER NOT NULL,  -- 1 FOR GREEN CORPORATES/CHAMPIONS || 2 FOR AGGREGATOR COMPANIES
  organization_id INTEGER NOT NULL,
  region VARCHAR(255) NULL,
  location VARCHAR(255),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN Key (organization_id) REFERENCES organizations(id)
);

-- Create "users" table
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NULL,
    last_name VARCHAR(255) NULL,
    provider VARCHAR(255),
    role_id INTEGER,
    FOREIGN Key (role_id) REFERENCES roles(id),
    user_company_id INTEGER NULL,
    FOREIGN Key (user_company_id) REFERENCES companies(id),
    email VARCHAR(255) DEFAULT NULL UNIQUE,
    password TEXT DEFAULT NULL,
    avatar_url TEXT NULL,
    user_type SMALLINT, -- 1 for TTNM ADMINS || 2 AGG GLOBAL ADMINS || 3 AGG ADMINS || 4 AGG USERS || 5 AGG COLLECTORS || 6 EXTERNAL COLLECTORS || 7 GREEN CHAMPIONS
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

CREATE TABLE waste_groups (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  category VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE waste_collections (
  id SERIAL PRIMARY KEY,
  date TIMESTAMP NOT NULL,
  champion_id INTEGER,
  FOREIGN Key (champion_id) REFERENCES users(id),
  collector_id INTEGER,
  FOREIGN Key (collector_id) REFERENCES users(id),
  waste JSON,
  is_collected BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE waste_for_sale (
  id SERIAL PRIMARY KEY,
  seller INTEGER,
  FOREIGN Key (seller) REFERENCES users(id),
  waste JSON
);

CREATE TABLE waste_buyers (
  id SERIAL PRIMARY KEY,
  buyer_id INTEGER,
  FOREIGN Key (buyer_id) REFERENCES users(id),
  rates JSON
);

CREATE TABLE payment_methods (
  id SERIAL PRIMARY KEY,
  payment_method VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE waste_transactions (
  id SERIAL PRIMARY KEY,
  date TIMESTAMP NOT NULL DEFAULT NOW(),
  buyer_id INTEGER,
  FOREIGN KEY (buyer_id) REFERENCES users(id),
  seller_id INTEGER,
  FOREIGN KEY (seller_id) REFERENCES users(id),
  waste_products JSON,
  total_amount VARCHAR NOT NULL,
  payment_method_id INTEGER,
  FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id),
  merchant_request_id VARCHAR(255) NULL,
  checkout_request_id VARCHAR(255) NULL,
  mpesa_result_code VARCHAR(255) NULL,
  mpesa_result_desc VARCHAR(255) NULL,
  mpesa_receipt_code VARCHAR(255) NULL,
  time_paid VARCHAR(255) NULL,
  is_paid BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- INVENTORY IS FOR COLLECTORS AND AGGREGATOR COMPANIES

--CREATE TABLE waste_inventory ()