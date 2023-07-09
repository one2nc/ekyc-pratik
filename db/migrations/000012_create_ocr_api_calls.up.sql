CREATE TABLE ekyc_schema.ocr_api_calls (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY NOT NULL,
    customer_id UUID NOT NULL REFERENCES ekyc_schema.customers (id),
    image_id UUID NOT NULL REFERENCES ekyc_schema.images (id),
    ocr_id UUID NOT NULL REFERENCES ekyc_schema.ocr_data (id),
    api_call_charges FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);