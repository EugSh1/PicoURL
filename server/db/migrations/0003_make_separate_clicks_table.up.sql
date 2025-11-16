ALTER TABLE links DROP COLUMN clicks;

CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    link_id VARCHAR(6) NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clicks_link_id ON clicks(link_id);