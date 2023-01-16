-- modify "ticket_events" table
ALTER TABLE "ticket_events" ALTER COLUMN "metadada" DROP NOT NULL, ALTER COLUMN "versions" DROP NOT NULL;
-- modify "tickets" table
ALTER TABLE "tickets" ALTER COLUMN "metadata" DROP NOT NULL, ALTER COLUMN "versions" DROP NOT NULL;
