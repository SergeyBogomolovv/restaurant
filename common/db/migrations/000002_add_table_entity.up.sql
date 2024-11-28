CREATE TABLE IF NOT EXISTS table_entity
(
  table_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  table_number INT NOT NULL,
  capacity INT NOT NULL,
  CONSTRAINT check_capacity CHECK (capacity > 0 AND capacity <= 8)
);