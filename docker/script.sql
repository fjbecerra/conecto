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

