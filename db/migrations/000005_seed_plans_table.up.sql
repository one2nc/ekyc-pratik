INSERT INTO ekyc_schema.plans (plan_name,is_active, image_upload_cost, face_match_cost,ocr_cost,daily_base_cost,created_at,updated_at)
VALUES
    ('basic', true,0.1,0.1 , 0.1, 10,current_timestamp,current_timestamp),
    ('advanced',true,  0.05,0.05 , 0.05, 15,current_timestamp,current_timestamp),
    ('enterprise',true, 0.01,0.01 , 0.01, 20,current_timestamp,current_timestamp);