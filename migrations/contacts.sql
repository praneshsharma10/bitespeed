CREATE TABLE IF NOT EXISTS contacts (
    id              SERIAL PRIMARY KEY,
    phone_number    TEXT,
    email           TEXT,
    linked_id       INT REFERENCES contacts(id),
    link_precedence TEXT NOT NULL DEFAULT 'primary',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
