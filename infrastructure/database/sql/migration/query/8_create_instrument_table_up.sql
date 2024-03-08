CREATE TABLE IF NOT EXISTS instrument (
             id INTEGER PRIMARY KEY AUTOINCREMENT,
             name       TEXT NOT NULL,
             created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
             updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_instrument_name ON instrument(name);

CREATE TABLE IF NOT EXISTS song_instrument (
                                          song_id    INTEGER NOT NULL,
                                          instrument_id   INTEGER NOT NULL,
                                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                          PRIMARY KEY (song_id, instrument_id),
    FOREIGN KEY(song_id) REFERENCES song(id) ON DELETE CASCADE,
    FOREIGN KEY(instrument_id) REFERENCES instrument(id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_song_instrument_song_id ON song_instrument(song_id);
CREATE INDEX IF NOT EXISTS idx_song_instrument_instrument_id ON song_instrument(instrument_id);

CREATE TRIGGER IF NOT EXISTS update_instrument_timestamp
    BEFORE UPDATE ON instrument
                                 FOR EACH ROW
BEGIN
UPDATE instrument SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_song_instrument_timestamp
    BEFORE UPDATE ON song_instrument
                                 FOR EACH ROW
BEGIN
UPDATE song_instrument SET updated_at = CURRENT_TIMESTAMP WHERE song_id = NEW.song_id AND instrument_id = NEW.instrument_id;
END;
