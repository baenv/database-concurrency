-- create "transactions" table
CREATE TABLE "transactions" ("id" uuid NOT NULL, "hash" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "transactions_hash_key" to table: "transactions"
CREATE UNIQUE INDEX "transactions_hash_key" ON "transactions" ("hash");
