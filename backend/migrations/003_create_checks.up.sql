CREATE TABLE checks (
id            UUID PRIMARY KEY,
endpoint_id   UUID NOT NULL REFERENCES endpoints(id),
status        TEXT NOT NULL ,
status_code   INTEGER,
response_time INTEGER,     
error         TEXT,
checked_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_checks_endpoint_checked_at
ON checks (endpoint_id, checked_at DESC);