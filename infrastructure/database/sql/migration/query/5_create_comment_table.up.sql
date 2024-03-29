CREATE TABLE IF NOT EXISTS comment (
       id         INTEGER PRIMARY KEY AUTOINCREMENT,
       author_id  INTEGER NOT NULL,
       song_id    INTEGER NOT NULL,
       content    TEXT NOT NULL,
       status     BOOLEAN DEFAULT 0,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY(author_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY(song_id) REFERENCES song(id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_comment_author_id ON comment(author_id);

CREATE INDEX IF NOT EXISTS idx_comment_song_id ON comment(song_id);

CREATE TRIGGER IF NOT EXISTS update_comment_timestamp
    BEFORE UPDATE ON comment
    FOR EACH ROW BEGIN
    UPDATE comment SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
