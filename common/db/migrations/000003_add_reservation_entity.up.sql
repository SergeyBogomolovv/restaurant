CREATE TYPE reservation_status AS
ENUM ('active', 'closed', 'cancelled');

CREATE TABLE IF NOT EXISTS reservations
(
	reservation_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	customer_id UUID NOT NULL REFERENCES customers(customer_id),
	start_time TIMESTAMP WITH TIME ZONE NOT NULL,
	end_time TIMESTAMP WITH TIME ZONE NOT NULL,
	status reservation_status DEFAULT 'active',
	table_id UUID NOT NULL REFERENCES tables(table_id),
	CONSTRAINT check_end_time CHECK (end_time > start_time)
);