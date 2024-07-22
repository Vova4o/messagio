CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	message TEXT NOT NULL,
	data TEXT NOT NULL,
	"time_created" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	"time_modified" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	processed BOOLEAN DEFAULT FALSE
	);
	
	
	CREATE OR REPLACE FUNCTION update_modified_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW."time_modified" = now();
		RETURN NEW;
	END;
	$$ language 'plpgsql';
	
	CREATE TRIGGER update_messages_modtime
	BEFORE UPDATE ON messages
	FOR EACH ROW
	EXECUTE FUNCTION update_modified_column();

