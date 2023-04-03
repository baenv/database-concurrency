INSERT INTO "users" ("id", "name", "email", "phone", "created_at", "updated_at") VALUES ('91272a62-c537-42ed-948c-bb2a91af2051', 'John Doe', 'john@gmail.com', '+1 234 567 890', '2019-01-01 00:00:00', '2019-01-01 00:00:00');

INSERT INTO "tickets" ("id", "status", "metadata", "versions", "user_id", "created_at", "updated_at") VALUES ('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 'Idle', '{"foo": "bar"}', '0', '91272a62-c537-42ed-948c-bb2a91af2051', '2019-01-01 00:00:00', '2019-01-01 00:00:00');

