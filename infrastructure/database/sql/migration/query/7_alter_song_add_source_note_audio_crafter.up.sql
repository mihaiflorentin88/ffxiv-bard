ALTER TABLE song ADD source TEXT DEFAULT 'N/A';
ALTER TABLE song ADD note TEXT DEFAULT 'N/A';
ALTER TABLE song ADD audio_crafter TEXT DEFAULT 'N/A' NOT NULL;
ALTER TABLE song add download_count INTEGER DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_song_source ON song(source);
CREATE INDEX IF NOT EXISTS idx_song_note ON song(note);
CREATE INDEX IF NOT EXISTS idx_song_audio_crafter ON song(audio_crafter);
