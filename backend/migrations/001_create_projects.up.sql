CREATE TABLE projects (
id          UUID PRIMARY KEY,
name        TEXT NOT NULL,
created_at  TIMESTAMPTZ DEFAULT NOW()
);