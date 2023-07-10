CREATE TABLE ekyc_schema.daily_reports_table (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    customer_id UUID DEFAULT uuid_generate_v4() NOT NULL REFERENCES ekyc_schema.customers (id),
    date_of_report TIMESTAMP NOT NULL,
    daily_base_charges FLOAT NOT NULL,
    no_of_face_match INTEGER NOT NULL,
    total_cost_of_face_match FLOAT NOT NULL,
    number_of_ocr INTEGER NOT NULL,
    total_cost_of_ocr FLOAT NOT NULL,
    total_api_call_charges FLOAT NOT NULL,
    total_image_storage_size_mb FLOAT NOT NULL,
    total_image_storage_cost FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);