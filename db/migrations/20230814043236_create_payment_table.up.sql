CREATE TABLE "payments" (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    booking_date DATE NOT NULL,
    qty INT NOT NULL,
    bill_no VARCHAR (255) NOT NULL,
    bill_total FLOAT NOT NULL,
    status_id int NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);