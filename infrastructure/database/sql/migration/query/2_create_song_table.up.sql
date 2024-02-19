CREATE TABLE IF NOT EXISTS song (
        id              INTEGER PRIMARY KEY AUTOINCREMENT,
        title           TEXT NOT NULL,
        artist          TEXT NOT NULL,
        ensemble_size   INTEGER NOT NULL,
        file_code       TEXT NOT NULL UNIQUE,
        uploader_id     INTEGER NOT NULL,
        status          INTEGER NOT NULL,
        status_message  TEXT,
        checksum        TEXT NOT NULL UNIQUE,
        lock_expire_ts  TIMESTAMP,
        created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(uploader_id) REFERENCES user(id)
    );

CREATE INDEX IF NOT EXISTS idx_song_title ON song(title);
CREATE INDEX IF NOT EXISTS idx_song_artist ON song(artist);
CREATE INDEX IF NOT EXISTS idx_song_uploader_id ON song(uploader_id);
CREATE INDEX IF NOT EXISTS idx_song_status ON song(status);

CREATE TRIGGER IF NOT EXISTS update_song_timestamp
    BEFORE UPDATE ON song
    FOR EACH ROW
BEGIN
    UPDATE song SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER make_artist_lowercase_after_insert
    AFTER INSERT ON song
    FOR EACH ROW
    WHEN NEW.artist <> LOWER(NEW.artist)
BEGIN
    UPDATE song SET artist = LOWER(NEW.artist) WHERE id = NEW.id;
END;

CREATE TRIGGER make_artist_lowercase_after_update
    AFTER UPDATE ON song
    FOR EACH ROW
    WHEN NEW.artist <> LOWER(NEW.artist)
BEGIN
    UPDATE song SET artist = LOWER(NEW.artist) WHERE id = NEW.id;
END;
