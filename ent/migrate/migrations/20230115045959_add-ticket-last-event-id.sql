-- modify "tickets" table
ALTER TABLE "tickets" ADD COLUMN "last_event_id" uuid NULL, ADD CONSTRAINT "tickets_ticket_events_last_event" FOREIGN KEY ("last_event_id") REFERENCES "ticket_events" ("id") ON DELETE SET NULL;
