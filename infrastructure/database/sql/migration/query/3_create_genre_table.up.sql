CREATE TABLE IF NOT EXISTS genre (
     id INTEGER PRIMARY KEY AUTOINCREMENT,
     name       TEXT NOT NULL,
     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_genre_name ON genre(name);

CREATE TABLE IF NOT EXISTS song_genre (
          song_id    INTEGER NOT NULL,
          genre_id   INTEGER NOT NULL,
          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
          PRIMARY KEY (song_id, genre_id),
          FOREIGN KEY (song_id) REFERENCES song(id),
          FOREIGN KEY (genre_id) REFERENCES genre(id)
);

CREATE INDEX IF NOT EXISTS idx_song_genre_song_id ON song_genre(song_id);
CREATE INDEX IF NOT EXISTS idx_song_genre_genre_id ON song_genre(genre_id);

CREATE TRIGGER IF NOT EXISTS update_genre_timestamp
    BEFORE UPDATE ON genre
    FOR EACH ROW
BEGIN
    UPDATE genre SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_song_genre_timestamp
    BEFORE UPDATE ON song_genre
    FOR EACH ROW
BEGIN
    UPDATE song_genre SET updated_at = CURRENT_TIMESTAMP WHERE song_id = NEW.song_id AND genre_id = NEW.genre_id;
END;
