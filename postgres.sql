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

CREATE TABLE IF NOT EXISTS messages (
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
ALTER TABLE messages ADD COLUMN parts_tsv TSVECTOR
  GENERATED ALWAYS AS (to_tsvector('english', coalesce(parts::TEXT, ''))) STORED;

CREATE INDEX parts_tsv_idx ON messages USING gin(parts_tsv);
---------------
-- Functions --
---------------
CREATE OR REPLACE FUNCTION search_conversations(query TEXT)
RETURNS TABLE (
  conversation_id TEXT,
  headline_title TEXT,
  headline_parts TEXT,
  rank REAL
) AS $$
DECLARE
  max_fragments       INT       := 10;
  min_words           INT       := 10;
  max_words           INT       := 25;
  ts_config           REGCONFIG := 'english';
  ts_headline_options TEXT      := format('MaxFragments=%s, MinWords=%s, MaxWords=%s', max_fragments, min_words, max_words);
  search_query        TSQUERY   := websearch_to_tsquery(ts_config, query);
BEGIN
  RETURN QUERY
  SELECT
    c.id conversation_id,
    ts_headline(ts_config, c.title, search_query, ts_headline_options),
    ts_headline(ts_config, m.parts, search_query, ts_headline_options),
    ts_rank(m.parts_tsv, search_query) AS rank
  FROM messages m
  JOIN nodes n ON m.node_id = n.id
  JOIN conversations c ON n.conversation_id = c.id
  WHERE m.parts_tsv @@ search_query
  ORDER BY rank DESC;
END
$$ LANGUAGE plpgsql;
