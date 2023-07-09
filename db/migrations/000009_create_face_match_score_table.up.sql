CREATE TABLE ekyc_schema.face_match_score (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    customer_id UUID NOT NULL REFERENCES ekyc_schema.customers (id),
    image_id_1 UUID NOT NULL REFERENCES ekyc_schema.images (id),
    image_id_2 UUID NOT NULL REFERENCES ekyc_schema.images (id),
    score INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);