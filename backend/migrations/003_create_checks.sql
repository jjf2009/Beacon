CREATE TABLE checks (
id            UUID PRIMARY KEY,
endpoint_id   UUID REFERENCES endpoints(id),
status        TEXT NOT NULL ,
status_code   INTEGER,
response_time INTEGER,     
error         TEXT,
checked_at    TIMESTAMPTZ DEFAULT NOW()
);