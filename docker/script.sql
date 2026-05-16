CREATE TABLE oauth_tokens (
    tenant_id      TEXT NOT NULL,
    provider       TEXT NOT NULL,

    ciphertext     BYTEA NOT NULL,
    nonce          BYTEA NOT NULL,

    key_version    TEXT NOT NULL,

    expires_at     TIMESTAMPTZ,

    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (tenant_id, provider)
);