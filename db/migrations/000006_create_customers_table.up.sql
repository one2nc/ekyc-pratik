-- +migrate Up
CREATE TABLE ekyc_schema.customers (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    plan_id UUID NOT NULL REFERENCES ekyc_schema.plans (id),
    access_key VARCHAR NOT NULL,
    secret_key VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);


