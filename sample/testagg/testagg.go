package testagg

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/xtracdev/goes"
)

var ErrUnknownType = errors.New("Unknown event type")

//The constants are used as unmarshalling hints when reconstructing the
//aggregate from its event history
const (
	TestAggCreatedTypeCode   = "TACRE"
	TestAggFooUpdateTypeCode = "TAFU"
)

//The aggregate for our example. In addition to the aggregate type, we need to capture
//the commands associated with the aggregate (implemented as exported methods that route
//events) and event types that are used to apply the state mutations.
type TestAgg struct {
	*goes.Aggregate
	Foo string
	Bar string
	Baz string
}

//Factory method for instantiating the aggregate, which is also the command method for 'create'
func NewTestAgg(foo, bar, baz string) (*TestAgg, error) {
	//Do validation... return an error if there's a problem
	var testAgg = new(TestAgg)
	testAgg.Aggregate = goes.NewAggregate()
	testAgg.Version = 1

	testAggCreated := TestAggCreated{
		AggregateId: testAgg.AggregateID,
		Foo:         foo,
		Bar:         bar,
		Baz:         baz,
	}

	testAgg.Apply(
		goes.Event{
			Source:  testAgg.AggregateID,
			Version: testAgg.Version,
			Payload: testAggCreated,
		})

	return testAgg, nil
}

//Factory method for recreating the aggregate state from the event history
func NewTestAggFromHistory(events []goes.Event) *TestAgg {
	if len(events) == 0 {
		return nil
	}

	testAgg := new(TestAgg)
	testAgg.Aggregate = goes.NewAggregate()

	unmarshalledEvents, err := unmarshallEvents(events)
	if err != nil {
		return nil
	}

	for _, e := range unmarshalledEvents {
		log.Debug("apply event", e)
		testAgg.Version += 1
		testAgg.Route(e)
	}

	return testAgg
}

//Command method for updating the foo attribute
func (ta *TestAgg) UpdateFoo(newfoo string) {
	ta.Version += 1
	ta.Apply(
		goes.Event{
			Source:  ta.AggregateID,
			Version: ta.Version,
			Payload: TestAggFooUpdated{
				AggregateId: ta.AggregateID,
				NewFoo:      newfoo,
			},
		})
}

//The required apply method, called only from commands to route and record events
func (ta *TestAgg) Apply(event goes.Event) {
	ta.Route(event)
	ta.Events = append(ta.Events, event)
}

//The required route method to route events to their handlers. Note the handlers may
//only change state - no other side effects are allowed.
func (ta *TestAgg) Route(event goes.Event) {
	event.Version = ta.Version
	switch event.Payload.(type) {
	case TestAggCreated:
		ta.handleTestAggCreated(event.Payload.(TestAggCreated))
	case TestAggFooUpdated:
		ta.handleTestFooUpdate(event.Payload.(TestAggFooUpdated))
	default:
		panic("WARN: unknown event routed to User aggregate")
	}
}

func (ta *TestAgg) handleTestAggCreated(event TestAggCreated) {
	ta.AggregateID = event.AggregateId
	ta.Foo = event.Foo
	ta.Bar = event.Bar
	ta.Baz = event.Baz
}

func (ta *TestAgg) handleTestFooUpdate(event TestAggFooUpdated) {
	ta.Foo = event.NewFoo
}

//Required implementation of the Store method.
func (ta *TestAgg) Store(eventStore goes.EventStore) error {

	marshalled, err := marshallEvents(ta.Events)
	if err != nil {
		return nil
	}

	log.Debug("Storing ", len(ta.Events), " events.")

	aggregateToStore := &goes.Aggregate{
		AggregateID: ta.AggregateID,
		Version:     ta.Version,
		Events:      marshalled,
	}

	err = eventStore.StoreEvents(aggregateToStore)
	if err != nil {
		return err
	}

	ta.Events = make([]goes.Event, 0)

	return nil
}

func marshallCreate(create TestAggCreated) ([]byte, error) {
	return proto.Marshal(&create)
}

func marshallFooUpdated(event TestAggFooUpdated) ([]byte, error) {
	return proto.Marshal(&event)
}

func unmarshallCreated(bytes []byte) (TestAggCreated, error) {
	var payload TestAggCreated
	err := proto.Unmarshal(bytes, &payload)
	return payload, err
}

func unmarshallFooUpdated(bytes []byte) (TestAggFooUpdated, error) {
	var payload TestAggFooUpdated
	err := proto.Unmarshal(bytes, &payload)
	return payload, err
}

func unmarshallEvents(events []goes.Event) ([]goes.Event, error) {
	var unmarshalled []goes.Event

	for _, e := range events {

		var err error
		var newEvent goes.Event
		newEvent.Source = e.Source
		newEvent.Version = e.Version
		newEvent.TypeCode = e.TypeCode

		switch e.TypeCode {
		case TestAggCreatedTypeCode:
			newEvent.Payload, err = unmarshallCreated(e.Payload.([]byte))
			if err != nil {
				return nil, err
			}
		case TestAggFooUpdateTypeCode:
			newEvent.Payload, err = unmarshallFooUpdated(e.Payload.([]byte))
			if err != nil {
				return nil, err
			}
		}

		unmarshalled = append(unmarshalled, newEvent)
	}

	return unmarshalled, nil
}

func marshallEvents(events []goes.Event) ([]goes.Event, error) {

	var updatedEvents []goes.Event

	for _, e := range events {

		var err error
		var newEvent goes.Event
		newEvent.Source = e.Source
		newEvent.Version = e.Version

		switch e.Payload.(type) {
		case TestAggCreated:
			newEvent.TypeCode = TestAggCreatedTypeCode
			newEvent.Payload, err = marshallCreate(e.Payload.(TestAggCreated))
			if err != nil {
				return nil, err
			}
		case TestAggFooUpdated:
			newEvent.TypeCode = TestAggFooUpdateTypeCode
			newEvent.Payload, err = marshallFooUpdated(e.Payload.(TestAggFooUpdated))
			if err != nil {
				return nil, err
			}
		default:
			return nil, ErrUnknownType
		}

		updatedEvents = append(updatedEvents, newEvent)
	}

	return updatedEvents, nil
}
