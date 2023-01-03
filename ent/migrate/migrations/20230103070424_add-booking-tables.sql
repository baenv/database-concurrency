-- create "service_prodivers" table
CREATE TABLE "service_prodivers" ("id" uuid NOT NULL, "name" character varying NOT NULL, "email" character varying NOT NULL, "phone" character varying NOT NULL, "verdor_ref" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL, "name" character varying NOT NULL, "email" character varying NOT NULL, "phone" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- create "tickets" table
CREATE TABLE "tickets" ("id" uuid NOT NULL, "status" character varying NOT NULL, "metadata" jsonb NOT NULL, "versions" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "tickets_users_tickets" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE);
-- create "ticket_events" table
CREATE TABLE "ticket_events" ("id" uuid NOT NULL, "type" character varying NOT NULL, "metadada" jsonb NOT NULL, "versions" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "ticket_id" uuid NOT NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "ticket_events_tickets_ticket_events" FOREIGN KEY ("ticket_id") REFERENCES "tickets" ("id") ON DELETE CASCADE, CONSTRAINT "ticket_events_users_ticket_events" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE);
