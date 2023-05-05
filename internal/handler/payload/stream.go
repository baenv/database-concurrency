package payload

type ConsumerRegisterRequest struct {
	ConsumerName string `json:"consumer_name"`
}

type TicketAcknowledgeRequest struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
}
