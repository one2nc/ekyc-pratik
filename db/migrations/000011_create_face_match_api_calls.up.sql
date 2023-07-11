CREATE TABLE ekyc_schema.face_match_api_calls (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    customer_id UUID DEFAULT uuid_generate_v4() NOT NULL REFERENCES ekyc_schema.customers (id),
    score_id UUID NOT NULL REFERENCES ekyc_schema.face_match_score (id),
    api_call_charges FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);