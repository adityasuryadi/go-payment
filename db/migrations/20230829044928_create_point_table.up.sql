DROP TABLE IF EXISTS "payments";CREATE TABLE "payments" (
    id UUID NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    point FLOAT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);