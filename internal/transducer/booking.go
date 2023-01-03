package transducer

func NewBookingTransducer(setupConfig *Config, childTransducers ...Transducer) *Transducer {
	transitionTable := TransitionTable{
		{Idle, Book}: func() *Outputs {
			return CreateOutputs().SetState(Reserved).
				AddEffect(UpdateBookingStatus).
				AddEffect(EmailUser)
		},

		{Reserved, Book}: func() *Outputs {
			return CreateOutputs().SetState(Booked).
				AddEffect(UpdateBookingStatus).
				AddEffect(EmailUser).
				AddEffect(CallClient).
				AddEffect(SMSUser)
		},

		{Reserved, Cancel}: func() *Outputs {
			return CreateOutputs().SetState(Cancelled).
				AddEffect(UpdateBookingStatus).
				AddEffect(EmailUser)
		},

		{Booked, CheckIn}: func() *Outputs {
			return CreateOutputs().SetState(CheckedIn).
				AddEffect(UpdateBookingStatus).
				AddEffect(EmailClient)
		},

		{Booked, CheckOut}: func() *Outputs {
			return CreateOutputs().SetState(CheckedOut).
				AddEffect(UpdateBookingStatus).
				AddEffect(EmailClient)
		},
	}

	return &Transducer{
		Name:            BookingTransducerName,
		TransitionTable: transitionTable,
	}
}

func NewBookingMachine(BookStateStr string) (*Config, *Transducer) {
	BookState := State(BookStateStr)
	BookConfig := CreateConfig().SetState(BookState)
	BookTransducer := NewBookingTransducer(BookConfig)

	return BookConfig, BookTransducer
}
