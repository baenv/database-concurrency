package transducer

import (
	"database-concurrency/internal/transducer/graph"
	"fmt"
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

func TestBookingShortestPaths(t *testing.T) {
	characterConfig := CreateConfig().SetState(Idle)
	characterTransducer := NewBookingTransducer(characterConfig)

	paths, edges := characterTransducer.GetShortestPaths()
	currentState := graph.Vertex(Idle)

	for _, path := range paths {
		for j, state := range path {
			if path[0] != currentState {
				stateLog := fmt.Sprintf("\n'%v' state", path[0])
				currentState = path[0]
				fmt.Println(stateLog)
			}

			if j < len(path)-1 {
				nextState := path[j+1]

				input := edges[state][nextState]
				transitionLog := fmt.Sprintf("\t\t%-30v + %-30v -> %-30v", state, input, nextState)
				fmt.Println(transitionLog)
			}
		}
	}
}