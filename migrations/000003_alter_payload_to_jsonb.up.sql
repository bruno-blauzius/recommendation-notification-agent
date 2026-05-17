ALTER TABLE recommendations
    ALTER COLUMN payload TYPE JSONB USING payload::jsonb;
