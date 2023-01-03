package transducer

import (
	"testing"
)

func TestRehearseBookingTransducer(t *testing.T) {
	bookingConfig := CreateConfig().SetState(Idle)
	bookingTransducer := NewBookingTransducer(bookingConfig)

	transitionQuery, aggregateQuery := bookingTransducer.ToSQL(Idle)
	t.Logf("Transition Query: %v\n", transitionQuery)
	t.Logf("Aggregate Query: %v\n", aggregateQuery)

	digraph := bookingTransducer.ToDiGraph()
	t.Logf("Digraph: %v\n", digraph)
}
