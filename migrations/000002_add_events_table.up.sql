-- Create events table
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMP,
    location VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add event_id foreign key to guests table
ALTER TABLE guests ADD COLUMN event_id INTEGER REFERENCES events(id);

-- Insert a default event for existing guests
INSERT INTO events (title, description, event_date, location) 
VALUES ('Default Event', 'Default event for existing guests', NOW(), 'TBD');

-- Update existing guests to use the default event
UPDATE guests SET event_id = 1 WHERE event_id IS NULL;

-- Make event_id required for future inserts
ALTER TABLE guests ALTER COLUMN event_id SET NOT NULL;
