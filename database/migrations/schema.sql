CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    guard_name varchar(255) NOT NULL,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    description text DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE uploads(
    id SERIAL PRIMARY KEY,
    item_id INTEGER,
    type VARCHAR(100),
    path TEXT,
    related_table VARCHAR(150),
    meta JSON NULL,
    UNIQUE(item_id,related_table)
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    action varchar(255) NOT NULL UNIQUE,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    module varchar(255) NOT NULL,
    submodule VARCHAR(255) NULL
);

CREATE TABLE role_has_permissions (
      permission_id INTEGER NOT NULL,
      role_id INTEGER NOT NULL,
      FOREIGN Key (permission_id) REFERENCES permissions(id) on delete CASCADE,
      FOREIGN Key (role_id) REFERENCES roles(id) on delete CASCADE,
      UNIQUE (permission_id, role_id)
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

CREATE TABLE counties(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE sub_counties(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  county_id INTEGER NOT NULL,
  FOREIGN Key (county_id) REFERENCES counties(id) on delete set null
);

CREATE TABLE organizations(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  country_id INTEGER NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  organization_type INTEGER NOT NULL, --1 aggrigators, 2 -green champion
  FOREIGN Key (country_id) REFERENCES countries(id)
);

ALTER TABLE organizations ADD CONSTRAINT check_organizations_type CHECK (organization_type IN (1,2)); -- make sure organization type is either 1 or 2


CREATE TABLE main_organization(
  id SERIAL PRIMARY KEY,
  organization_id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  tag_line TEXT NOT NULL,
  about_us TEXT NOT NULL,
  logo_path TEXT NOT NULL,
  app_appstore_link TEXT NOT NULL,
  app_google_playstore_link TEXT NOT NULL,
  website_url TEXT NOT NULL,
  city TEXT NOT NULL,
  state TEXT NOT NULL,
  zip TEXT NOT NULL,
  country TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE companies (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  country_id INTEGER NOT NULL,
  company_type INTEGER NOT NULL,  -- 1 FOR GREEN CORPORATES/CHAMPIONS || 2 FOR AGGREGATOR COMPANIES
  organization_id INTEGER NULL,
  region VARCHAR NULL, --Specific region this company is in
  location VARCHAR(255) NULL, -- the location of this company ie citadel muthithi road
  administrative_level_1_location VARCHAR(255) NULL, -- in kenya, this will be county, in uganda it will be a different value , ie Nairobi county
  lat float NULL,
  lng float NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  contact_person1_first_name VARCHAR NULL,
  contact_person1_last_name VARCHAR NULL,
  contact_person1_phone VARCHAR NULL,
  contact_person1_email VARCHAR NULL,
  contact_person2_email VARCHAR NULL,
  contact_person2_first_name VARCHAR NULL,
  contact_person2_last_name VARCHAR NULL,
  contact_person2_phone VARCHAR NULL,
  
  FOREIGN Key (organization_id) REFERENCES organizations(id),
  FOREIGN Key (country_id) REFERENCES countries(id)
);

ALTER TABLE companies ADD CONSTRAINT check_company_type CHECK (company_type IN (1,2)); -- make sure company type is either 1 or 2


-- Create "users" table
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NULL,
    last_name VARCHAR(255) NULL,
    provider VARCHAR(255),
    role_id INTEGER,
    FOREIGN Key (role_id) REFERENCES roles(id),
    user_company_id INTEGER NULL,
    user_organization_id INTEGER NULL,
    FOREIGN Key (user_company_id) REFERENCES companies(id),
    FOREIGN Key (user_organization_id) REFERENCES organizations(id),
    is_main_organization_user BOOLEAN DEFAULT false not null,
    is_organization_super_admin BOOLEAN DEFAULT false not null,
    is_company_super_admin BOOLEAN DEFAULT false not null,

    email VARCHAR(255) DEFAULT NULL UNIQUE,
    password TEXT DEFAULT NULL,
    avatar_url TEXT NULL,
    user_type SMALLINT, -- 1 for TTNM ADMINS || 2 AGG GLOBAL ADMINS || 3 AGG ADMINS || 4 AGG USERS || 5 AGG COLLECTORS || 6 EXTERNAL COLLECTORS || 7 GREEN CHAMPIONS 8|| Green chamption global admin, || 9 Global aggregator admins, ||10 green chamption super admin
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


-- Create "Champion_assigned_aggregator" table
CREATE TABLE champion_aggregator_assignments (
  id SERIAL PRIMARY KEY,
  champion_id INTEGER NOT NULL,
  FOREIGN Key (champion_id) REFERENCES companies(id),
  collector_id INTEGER NOT NULL,
  FOREIGN Key (collector_id) REFERENCES companies(id),
  --pickup_day VARCHAR NULL,
  --pickup_time VARCHAR NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

--ALTER TABLE champion_aggregator_assignments ADD CONSTRAINT check_pickup_day CHECK (pickup_day IN ('Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday')); 


CREATE TABLE pickup_time_stamps (
  id SERIAL PRIMARY KEY,
  stamp VARCHAR(255) NOT NULL,
  time_range VARCHAR(255) NOT NULL
);
ALTER TABLE pickup_time_stamps ADD CONSTRAINT check_stamp CHECK (stamp IN ('Morning','Afternoon','Evening'));

-- Create "Champion pickup times"
CREATE TABLE champion_pickup_times(
  id SERIAL PRIMARY KEY,
  champion_aggregator_assignment_id INTEGER NOT NULL,
  pickup_time_stamp_id INTEGER NOT NULL,
  ExactPickupTime VARCHAR(255) NULL,
  FOREIGN Key (champion_aggregator_assignment_id) REFERENCES champion_aggregator_assignments(id) on delete cascade,
  FOREIGN Key (pickup_time_stamp_id) REFERENCES pickup_time_stamps(id),
  pickup_day VARCHAR NOT NULL
);

ALTER TABLE champion_pickup_times ADD CONSTRAINT champion_pickup_times_check_pickup_day CHECK (pickup_day IN ('Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday')); 


CREATE TABLE waste_types (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  is_active BOOLEAN not null DEFAULT true,
  parent_id INTEGER NULL,
  FOREIGN Key (parent_id) REFERENCES waste_types(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE collection_requests (
  id SERIAL PRIMARY KEY,
  producer_id INTEGER NOT NULL,
  FOREIGN Key (producer_id) REFERENCES companies(id),
  collector_id INTEGER NOT NULL,
  FOREIGN Key (collector_id) REFERENCES companies(id),
  request_date DATE NOT NULL,
  pickup_time_stamp_id INTEGER NOT NULL,
  FOREIGN Key (pickup_time_stamp_id) REFERENCES pickup_time_stamps(id),
  location VARCHAR(255) NULL, -- the location of this company ie citadel muthithi road
  administrative_level_1_location VARCHAR(255) NULL, -- in kenya, this will be county, in uganda it will be a different value , ie Nairobi county
  lat float NULL,
  lng float NULL,
  pickup_date TIMESTAMP NULL,
  status INTEGER NOT NULL, --1 pending, 2 confirmed, 3 on the way, 4 cancelled, 5 completed  
  first_contact_person VARCHAR(255) NOT NULL,
  second_contact_person VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE collection_requests ADD CONSTRAINT collection_requests_status CHECK (status IN (1,2,3,4,5)); 

CREATE TABLE waste_items (
  id SERIAL PRIMARY KEY,
  collection_request_id INTEGER NOT NULL,
  FOREIGN Key (collection_request_id) REFERENCES collection_requests(id),
  waste_type_id INTEGER NOT NULL,
  FOREIGN Key (waste_type_id) REFERENCES waste_types(id),
  weight DECIMAL NOT NULL
);

CREATE UNIQUE INDEX waste_types_unique_name_idx on waste_types (LOWER(name));  


CREATE TABLE notifications (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  FOREIGN Key (user_id) REFERENCES users(id),
  subject VARCHAR(255) NOT NULL,
  message VARCHAR(255) NOT NULL,
  status BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- CREATE TABLE "aggregator_waste_types"
CREATE TABLE aggregator_waste_types(
  id SERIAL PRIMARY KEY,
  aggregator_id INTEGER NOT NULL,
  waste_id INTEGER NOT NULL,
  alert_level FLOAT NULL DEFAULT 0,
  FOREIGN Key (aggregator_id) REFERENCES companies(id),
  FOREIGN Key (waste_id) REFERENCES waste_types(id),
  UNIQUE(aggregator_id,waste_id)
);



-- Create "buyers" table, these are the ones that buy waste from aggregators
CREATE TABLE buyers(
  id SERIAL PRIMARY KEY,
  company_id int not null,
  company VARCHAR NULL,
  country_id INTEGER NULL,
  FOREIGN Key (country_id) REFERENCES countries(id),
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  is_active BOOLEAN not null, 
  region VARCHAR(255) NULL,
  calling_code VARCHAR(6) NULL,
  location VARCHAR(255) NULL, -- the location of this company ie citadel muthithi road
  administrative_level_1_location VARCHAR(255) NULL, -- in kenya, this will be county, in uganda it will be a different value , ie Nairobi county
  lat float NULL,
  lng float NULL,
  phone VARCHAR(15) NULL DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN Key (company_id) REFERENCES companies(id)
);

-- Create "buyers" table, these are the ones that buy waste from aggregators
CREATE TABLE suppliers(
  id SERIAL PRIMARY KEY,
  company_id int not null,
  company VARCHAR NULL,
  first_name VARCHAR(255) NOT NULL,
  country_id INTEGER NULL,
  FOREIGN Key (country_id) REFERENCES countries(id),
  last_name VARCHAR(255) NOT NULL,
  is_active BOOLEAN not null, 
  region VARCHAR(255) NULL,
  calling_code VARCHAR(6) NULL,
  location VARCHAR(255) NULL, -- the location of this company ie citadel muthithi road
  administrative_level_1_location VARCHAR(255) NULL, -- in kenya, this will be county, in uganda it will be a different value , ie Nairobi county
  lat float NULL,
  lng float NULL,
  phone VARCHAR(15) NULL DEFAULT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN Key (company_id) REFERENCES companies(id)
);



CREATE TABLE inventory(
      id SERIAL PRIMARY KEY,
  company_id int not null,
  FOREIGN Key (company_id) REFERENCES companies(id),
  waste_type_id INTEGER NULL,
  FOREIGN Key (waste_type_id) REFERENCES waste_types(id),
  total_weight FLOAT not null --in kgs
);

CREATE TABLE inventory_adjustments(
    id SERIAL PRIMARY KEY,
    adjusted_by INTEGER not null,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    company_id int not null,
    adjustment_amount DECIMAL not NULL,
    is_positive_adjustment BOOLEAN not null,
    FOREIGN Key (company_id) REFERENCES companies(id),
    FOREIGN Key (adjusted_by) REFERENCES users(id)
);

CREATE TABLE sales(
  id SERIAL PRIMARY KEY,
  ref VARCHAR not null,
  company_id int not null,
  buyer_id int not null,
  FOREIGN Key (buyer_id) REFERENCES buyers(id),
  FOREIGN Key (company_id) REFERENCES companies(id),
  total_weight FLOAT NULL, --in kgs
  total_amount FLOAT NULL, --ksh
  date TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
  dump json NULL
);

CREATE TABLE purchases(
  id SERIAL PRIMARY KEY,
  ref VARCHAR not null,
  company_id int not null,
  supplier_id int not null,
  FOREIGN Key (supplier_id) REFERENCES suppliers(id),
  FOREIGN Key (company_id) REFERENCES companies(id),
  total_weight FLOAT NULL, --in kgs
  total_amount FLOAT NULL, --ksh
  date TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
  dump json NULL
);

CREATE TABLE purchase_items(
  id SERIAL PRIMARY KEY,
  company_id INTEGER NOT NULL,
  purchase_id INTEGER NOT NULL,
  FOREIGN Key (company_id) REFERENCES companies(id),
  waste_type_id INTEGER NOT NULL,
  FOREIGN Key (purchase_id) REFERENCES purchases(id) on delete cascade,
  FOREIGN Key (waste_type_id) REFERENCES waste_types(id),
  weight FLOAT NULL,
  cost_per_kg FLOAT NULL,
  total_amount FLOAT NOT NULL
);

CREATE TABLE sale_items(
  id SERIAL PRIMARY KEY,
  company_id INTEGER NOT NULL,
  sale_id INTEGER NOT NULL,
  FOREIGN Key (company_id) REFERENCES companies(id),
  waste_type_id INTEGER NOT NULL,
  FOREIGN Key (sale_id) REFERENCES sales(id) on delete cascade,
  FOREIGN Key (waste_type_id) REFERENCES waste_types(id),
  weight FLOAT NULL,
  cost_per_kg FLOAT NULL,
  total_amount FLOAT NOT NULL
);

CREATE TABLE purchase_transactions(
  ref VARCHAR NOT NULL,
  id SERIAL PRIMARY KEY,
  purchase_id INTEGER not NULL,
  FOREIGN Key (purchase_id) REFERENCES purchases(id) on delete cascade,
  company_id int not null,
  FOREIGN Key (company_id) REFERENCES companies(id),
  payment_method VARCHAR NOT NULL,
  checkout_request_id VARCHAR NULL,
  merchant_request_id VARCHAR NULL,
  card_mask VARCHAR NULL,
  msisdn_idnum VARCHAR NULL,
  transaction_date TIMESTAMP NULL,
  receipt_no VARCHAR NULL,
  amount DECIMAL not NULL,
  mpesa_result_code VARCHAR NULL,
  mpesa_result_desc VARCHAR NULL,
  ipay_status VARCHAR NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE sale_transactions(
  ref VARCHAR NOT NULL,
  id SERIAL PRIMARY KEY,
  sale_id INTEGER not NULL,
  FOREIGN Key (sale_id) REFERENCES sales(id) on delete cascade,
  company_id int not null,
  FOREIGN Key (company_id) REFERENCES companies(id),
  payment_method VARCHAR NOT NULL,
  checkout_request_id VARCHAR NULL,
  merchant_request_id VARCHAR NULL,
  card_mask VARCHAR NULL,
  msisdn_idnum VARCHAR NULL,
  transaction_date TIMESTAMP NULL,
  receipt_no VARCHAR NULL,
  amount DECIMAL not NULL,
  mpesa_result_code VARCHAR NULL,
  mpesa_result_desc VARCHAR NULL,
  ipay_status VARCHAR NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- INVENTORY IS FOR COLLECTORS AND AGGREGATOR COMPANIES

--CREATE TABLE waste_inventory ()
