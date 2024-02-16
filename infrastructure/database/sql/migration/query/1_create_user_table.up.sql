PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS user (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        name TEXT,
        auth_provider TEXT,
        provider_id TEXT,
        avatar TEXT,
        discriminator TEXT,
        authorization_code TEXT,
        access_token TEXT,
        refresh_token TEXT,
        token_type TEXT,
        expires_at TIMESTAMP,
        token_released_at TIMESTAMP,
        is_admin BOOLEAN DEFAULT 0,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_username ON user(username);
CREATE INDEX IF NOT EXISTS idx_user_email ON user(email);
CREATE INDEX IF NOT EXISTS idx_user_auth_provider ON user(auth_provider);
CREATE INDEX IF NOT EXISTS idx_user_provider_id ON user(provider_id);

CREATE TRIGGER IF NOT EXISTS update_user_timestamp
    BEFORE UPDATE ON user
    FOR EACH ROW
BEGIN
    UPDATE user SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
