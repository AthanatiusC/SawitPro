/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/**
  id serial, primary key user identifier
  name varchar(60), bussiness requirements to limit name to 60 characters
  phone varchar(13), Indonesian phone number have 9-13 digits length e.g 628xxxxxxxxxx
  password varchar(72), store password in bcrypt+salt which have 72 character limit
  updated_at timestamp, to track last time data was updated
  created_at timestamp, to track when data was created
*/
CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	name VARCHAR(60) NOT NULL,
  phone VARCHAR(13) UNIQUE NOT NULL, 
  password VARCHAR(74) NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW(),
  created_at TIMESTAMP DEFAULT NOW()
);

/**
Create unique index for column id, id is frequently queried by endpoint
Create unique index for column users phone and users password, phone can be used as single column index since its the first entry
*/
CREATE UNIQUE INDEX index_user_id ON users(id);
CREATE UNIQUE INDEX index_user_phone_and_password ON users(phone,password);

-- I would like to make a audit trail but i think it is unecessary in this case

-- Seed users entry
INSERT INTO users(name,phone,password) VALUES ('user','6280000000000','$2a$06$bt380.sYY0HEAa1tz2eyfOOQDHarjgiABmv.ZJTXzKdXMU.hQFAyi');