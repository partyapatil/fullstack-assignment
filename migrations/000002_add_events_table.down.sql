-- Make event_id nullable first
ALTER TABLE guests ALTER COLUMN event_id DROP NOT NULL;

-- Remove foreign key and column
ALTER TABLE guests DROP COLUMN event_id;

-- Drop events table
DROP TABLE events;
