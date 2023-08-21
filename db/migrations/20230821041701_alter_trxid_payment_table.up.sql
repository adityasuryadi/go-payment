ALTER TABLE "payments" 
ADD trx_id VARCHAR(255) NULL,
ADD payment_channel_uid INT NULL,
ADD payment_channel VARCHAR(255) NULL;