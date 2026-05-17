ALTER TABLE recommendations
    ALTER COLUMN payload TYPE TEXT USING payload::text;
