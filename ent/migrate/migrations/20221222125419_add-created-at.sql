-- modify "transactions" table
ALTER TABLE "transactions" ADD COLUMN "created_at" timestamptz NOT NULL;
