CREATE TABLE IF NOT EXISTS recommendations (
    id          VARCHAR(36) PRIMARY KEY,
    sender_id   VARCHAR(36) NOT NULL,
    payload     TEXT        NOT NULL,
    score       NUMERIC(5, 4) NOT NULL DEFAULT 0,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recommendations_sender_id ON recommendations(sender_id);
CREATE INDEX IF NOT EXISTS idx_recommendations_created_at ON recommendations(created_at DESC);
