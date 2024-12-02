CREATE TABLE IF NOT EXISTS tables
(
  table_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  table_number INT NOT NULL UNIQUE,
  capacity INT NOT NULL,
  CONSTRAINT check_capacity CHECK (capacity > 0 AND capacity <= 8)
);