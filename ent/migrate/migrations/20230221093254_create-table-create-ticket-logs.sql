-- create "create_ticket_logs" table
CREATE TABLE "create_ticket_logs" ("ticket_id" uuid NOT NULL, "unique_id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("ticket_id"));
-- create index "create_ticket_logs_unique_id_key" to table: "create_ticket_logs"
CREATE UNIQUE INDEX "create_ticket_logs_unique_id_key" ON "create_ticket_logs" ("unique_id");
