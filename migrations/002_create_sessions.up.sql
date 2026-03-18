-- Sessions table for SCS (alexedwards/scs) with pgxstore backend.
-- Session data is stored server-side; only an opaque token is sent to the client.
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
