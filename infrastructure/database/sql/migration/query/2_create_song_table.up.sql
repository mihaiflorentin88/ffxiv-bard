CREATE TABLE IF NOT EXISTS song (
        id              INTEGER PRIMARY KEY AUTOINCREMENT,
        title           TEXT NOT NULL,
        artist          TEXT NOT NULL,
        ensemble_size   INTEGER NOT NULL,
        filename       TEXT NOT NULL UNIQUE,
        uploader_id     INTEGER NOT NULL,
        status          INTEGER NOT NULL,
        status_message  TEXT,
        checksum        TEXT NOT NULL UNIQUE,
        lock_expire_ts  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
