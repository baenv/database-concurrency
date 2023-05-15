package payload

type ConsumerRegisterRequest struct {
	ConsumerName string `json:"consumer_name"`
	HealthURL    string `json:"health_url"`
}

type TicketAcknowledgeRequest struct {
	TicketID  string `json:"ticket_id"`
	MessageID string `json:"message_id"`
}
