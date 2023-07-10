CREATE TABLE ekyc_schema.image_upload_api_calls (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    customer_id UUID NOT NULL REFERENCES ekyc_schema.customers (id),
    image_id UUID NOT NULL REFERENCES ekyc_schema.images (id),
    image_storage_charges FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);


