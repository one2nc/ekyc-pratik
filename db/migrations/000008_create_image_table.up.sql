CREATE TABLE ekyc_schema.images (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    customer_id UUID NOT NULL REFERENCES ekyc_schema.customers (id),
    file_path VARCHAR NOT NULL,
    file_extension VARCHAR NOT NULL,
    file_size_mb FLOAT NOT NULL,
    image_type ekyc_schema.image_type NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
