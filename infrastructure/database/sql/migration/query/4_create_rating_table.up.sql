CREATE TABLE IF NOT EXISTS rating (
      id         INTEGER PRIMARY KEY AUTOINCREMENT,
      song_id    INTEGER NOT NULL,
      author_id  INTEGER NOT NULL,
      rating     INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 10),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(song_id) REFERENCES song(id) ON DELETE CASCADE,
    FOREIGN KEY(author_id) REFERENCES user(id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_rating_song_id ON rating(song_id);

CREATE INDEX IF NOT EXISTS idx_rating_author_id ON rating(author_id);

CREATE TRIGGER IF NOT EXISTS update_rating_timestamp
    BEFORE UPDATE ON rating
    FOR EACH ROW BEGIN
    UPDATE rating SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

