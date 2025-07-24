CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    balance NUMERIC DEFAULT 0
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    sender_id UUID REFERENCES users(id),
    receiver_id UUID REFERENCES users(id),
    amount NUMERIC NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO users (id, first_name, last_name, email, balance) VALUES
   ('7141b92f-a8c8-471e-83e5-7fc72da61cb9', 'Alice', 'Doe', 'alice@example.com', 1000),
   ('861d7697-b717-43e8-95a2-1a74f9a36ab1', 'Joe', 'Brook', 'joe@example.com', 10800),
   ('939cb506-0d70-4791-8c9e-d4284d87c749', 'Bob', 'Sink', 'bob@example.com', 500),
   ('9d02adbc-27ca-4695-9d92-10cb35db67f4', 'Carol', 'Smith', 'carol@example.com', 750),
   ('befeef21-1475-4a13-a0de-3943d2eb0910', 'Dave', 'Johnson', 'dave@example.com', 1200),
   ('39ede33f-2a57-44bb-a563-7a37dda46bdf', 'Eve', 'Miller', 'eve@example.com', 3000),
   ('3f2fdcde-cb17-488e-819e-99cafea3f984', 'Frank', 'White', 'frank@example.com', 640),
   ('595e4e71-ad88-4a65-85d2-be98718f36df', 'Grace', 'Taylor', 'grace@example.com', 9800);

-- 1. Alice send to Joe
-- 2. Bob send to carol
-- 3. Dave sends to Eve
-- 4. Frank sends to Grace
INSERT INTO transactions(id, sender_id, receiver_id, amount, status, created_at) VALUES
    (gen_random_uuid(), '7141b92f-a8c8-471e-83e5-7fc72da61cb9', '861d7697-b717-43e8-95a2-1a74f9a36ab1', 10, 'SUCCESS', current_timestamp),
    (gen_random_uuid(), '939cb506-0d70-4791-8c9e-d4284d87c749', '9d02adbc-27ca-4695-9d92-10cb35db67f4', 1, 'SUCCESS', current_timestamp),
    (gen_random_uuid(), 'befeef21-1475-4a13-a0de-3943d2eb0910', '39ede33f-2a57-44bb-a563-7a37dda46bdf', 5, 'SUCCESS', current_timestamp),
    (gen_random_uuid(), '3f2fdcde-cb17-488e-819e-99cafea3f984', '595e4e71-ad88-4a65-85d2-be98718f36df', 7, 'SUCCESS', current_timestamp);


CREATE EXTENSION IF NOT EXISTS "pgcrypto";
