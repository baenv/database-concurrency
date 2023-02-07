package transducer

// States
const (
	// an invalid state for empty states
	Invalid State = "Invalid"

	Idle       State = "Idle"
	Reserved   State = "Reserved"
	Booked     State = "Booked"
	Cancelled  State = "Cancelled"
	CheckedIn  State = "CheckedIn"
	CheckedOut State = "CheckedOut"
)

// Inputs
const (
	Reserve  Input = "Reserve"
	Book     Input = "Book"
	Cancel   Input = "Cancel"
	CheckIn  Input = "CheckIn"
	CheckOut Input = "CheckOut"
)

// Effects
const (
	UpdateBookingStatus Effect = iota
	EmailUser
	EmailClient
	CallClient
	SMSUser
	CreateBookingEvent
)

// Names
const (
	BookingTransducerName = "booking"
)
