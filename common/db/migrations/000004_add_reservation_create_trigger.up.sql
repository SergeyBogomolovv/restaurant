CREATE OR REPLACE FUNCTION check_reservation_active() 
RETURNS TRIGGER AS $$
BEGIN
	IF (SELECT TRUE FROM reservations
	    WHERE status = 'active'
		AND (NEW.start_time, NEW.end_time) OVERLAPS (start_time, end_time))
		THEN RAISE EXCEPTION 'table already reserved';
	END IF;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_table_free
BEFORE INSERT ON reservations
FOR EACH ROW
EXECUTE PROCEDURE check_reservation_active();