CREATE TABLE ekyc_schema.ocr_data (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    customer_id UUID NOT NULL REFERENCES ekyc_schema.customers (id),
    image_id UUID NOT NULL REFERENCES ekyc_schema.images (id),
    ocr_data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);