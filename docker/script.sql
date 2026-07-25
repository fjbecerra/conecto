CREATE TABLE connections (

    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL,

    provider TEXT NOT NULL,

    external_id TEXT,

    metadata JSONB,

    status TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),

    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE sync_jobs (
    id UUID PRIMARY KEY,

    connection_id UUID NOT NULL,

    pipeline_id TEXT NOT NULL,

    status TEXT NOT NULL,

    attempt INT NOT NULL,

    max_retries INT NOT NULL,

    next_retry_at TIMESTAMP,

    last_error TEXT,

    created_at TIMESTAMP NOT NULL,

    started_at TIMESTAMP,

    finished_at TIMESTAMP
);

