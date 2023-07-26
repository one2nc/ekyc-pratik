CREATE TABLE ekyc_schema.cron_registry (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL,
    unique_identifier VARCHAR UNIQUE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
