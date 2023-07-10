CREATE TABLE ekyc_schema.plans (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    plan_name ekyc_schema.plan_type NOT NULL,
    is_active BOOLEAN NOT NULL,
    image_upload_cost FLOAT NOT NULL,
    face_match_cost FLOAT NOT NULL,
    ocr_cost FLOAT NOT NULL,
    daily_base_cost FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);