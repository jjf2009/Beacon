CREATE TABLE endpoints (
id          UUID PRIMARY KEY,
project_id  UUID REFERENCES projects(id),
name        TEXT NOT NULL,
url         TEXT NOT NULL,
interval    INTEGER NOT NULL,  
created_at  TIMESTAMPTZ DEFAULT NOW()
);