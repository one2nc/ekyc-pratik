CREATE TABLE ekyc_schema.plans (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    plan_name ekyc_schema.plan_type NOT NULL,
    is_active BOOLEAN NOT NULL,
    image_upload_cost INTEGER NOT NULL,
    face_match_cost INTEGER NOT NULL,
    ocr_cost INTEGER NOT NULL,
    daily_base_cost INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);