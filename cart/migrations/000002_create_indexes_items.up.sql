CREATE INDEX user_id_hash_index
ON items
USING hash (user_id);

