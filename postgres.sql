------------
-- Schema --
------------
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    title TEXT,
    import_time DOUBLE PRECISION,
    create_time DOUBLE PRECISION,
    update_time DOUBLE PRECISION,
    current_node TEXT,
    plugin_ids TEXT
);

CREATE TABLE IF NOT EXISTS nodes (
    id TEXT PRIMARY KEY,
    conversation_id TEXT REFERENCES conversations(id),
    parent TEXT,
    import_time DOUBLE PRECISION,
    children JSONB
);

CREATE TABLE messages (
  id TEXT PRIMARY KEY,
  node_id TEXT,
  author_role TEXT,
  author_name TEXT,
  import_time DOUBLE PRECISION,
  create_time DOUBLE PRECISION,
  update_time DOUBLE PRECISION,
  content_type TEXT,
  parts TEXT,
  end_turn BOOLEAN,
  weight DOUBLE PRECISION,
  recipient TEXT,
  FOREIGN KEY (node_id) REFERENCES nodes (id)
);
----------------------
-- Full-text search --
----------------------
ALTER TABLE messages ADD COLUMN parts_tsv tsvector;

CREATE INDEX parts_tsv_idx ON messages USING gin(parts_tsv);

CREATE OR REPLACE FUNCTION messages_parts_tsv_trigger() RETURNS trigger AS $$
BEGIN
  NEW.parts_tsv := to_tsvector('english', coalesce(NEW.parts::text, ''));
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER messages_parts_tsv_update
BEFORE INSERT OR UPDATE OF parts ON messages
FOR EACH ROW
EXECUTE PROCEDURE messages_parts_tsv_trigger();
---------------
-- Functions --
---------------
CREATE OR REPLACE FUNCTION search_conversations(query TEXT)
RETURNS TABLE (
  title TEXT,
  parts TEXT,
  rank REAL
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    c.title,
    m.parts,
    ts_rank(to_tsvector('english', c.title || ' ' || m.parts), to_tsquery('english', query)) AS rank
  FROM messages m
  JOIN nodes n ON m.node_id = n.id
  JOIN conversations c ON n.conversation_id = c.id
  WHERE to_tsvector('english', c.title || ' ' || m.parts) @@ to_tsquery('english', query)
  ORDER BY rank DESC;
END;
$$ LANGUAGE plpgsql;

