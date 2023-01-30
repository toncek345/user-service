CREATE TABLE users (
  id UUID primary key,
  first_name text,
  last_name text,
  email text,
  country text,
  password text,
  created_at timestamp,
  updated_at timestamp
  );
