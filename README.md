# Go Event Sourcing

This project defines the base types and interfaces for working with event sourcing. The approach for this and related projects is loosely defined - packages using event sourcing are responsible for proving certain methods and observing conventions, as opposed to having a general framework that provides the hooks for packages to tap into for event sourcing.

## Go Event Source Sample

This sample package shows how to implement an event sourced aggregate using the Go Event Store project (goes).

To be event sourced...

* The Aggregate type must be embedded
* A method for creating a new aggregate must be provided
* A method for creating an aggregate in memory from its event history must be provided.
* Constructing a new aggregate produces an event
* Mutations occur via commands, which emit events handled by event handler, with events routed to handlers
* Events are recorded in event history
* An apply method routes an event to the event handler, and records the event
* When applying event history, only the route method is used - side effects occur in the command handlers

The User object is used in internal feature tests run using gucumber. The testagg package
has a test aggregate that uses protobufs to marshall and unmarshall events, and is used
in the feature tests in the pgevent store project.

### Dependencies

<pre>
go get github.com/lib/pq
go get -u github.com/golang/protobuf/protoc-gen-go
</pre>

### Generate protobuf code

(In testagg)

<pre>
protoc --go_out=. *.proto
</pre>

### Integration Tests

<pre>
gucumber
</pre>

### Detailed Walkthrough

Consider a simple aggregate with three properties, an operation to
create an instance of the aggregate, and an operation to update
a specific aggregate on the property.

<pre>
model TestAggs

class TestAgg
attributes
    Foo: String
    Bar: String
    Baz: String
operations
    NewTestAgg():TestAgg
    UpdateFoo(newFooVal:String)
end
</pre>

In Go we'd define a type with the attributes, a factory method,
and a method receiver on that type.

<pre>
type TestAgg struct {
	Foo string
	Bar string
	Baz string
}

func NewTestAgg(foo, bar, baz string) (*TestAgg, error) {
    ...stuff...
}

func (ta *TestAgg) UpdateFoo(newfoo string) {
    ...stuff...
}
</pre>

To event source this type, we need to embed the aggregate type:

<pre>
type TestAgg struct {
	*goes.Aggregate
	Foo string
	Bar string
	Baz string
}
</pre>

Creating a new instance of the aggregate falls naturally into
how we modeled the aggregate, but since we want to use event
sourcing on the aggregate, we have to think about how to
construct it from the event history.

We'll need a method for recreating state from event history:

<pre>
func NewTestAggFromHistory(events []goes.Event) *TestAgg {
    ...stuff...
}
</pre>


An important thing to remember is when we reconstruct the state
of an aggregate from event history, we do not want to have all
the side effects that have occured of the history of the aggregate
instance to occur again.

To enable this, we think of the methods on an aggregate as being
commands, the execution of which generates events which mutate
aggregate state. The events can be used to modify state when
a command has been executed, or when loading events from
history.

So in addition to defining the methods on the aggregate, we
need to define the events associated with the methods that
will be used in event sourcing to mutate state. In this example,
we modeled the events using protobuf definitions and code
generated the Go representation. The events are `TestAggCreated`
and `TestUpdateFoo`

In the implemenation, we use two methods for use events to
mutate state. `Apply` is used directly from commands, and `Route`
is used when loading from history. The `Apply` command merely
records the event in the in-memory event history then calls
`Route` with the event for state mutation.

<pre>
func (ta *TestAgg) Apply(event goes.Event) {
	ta.Route(event)
	ta.Events = append(ta.Events, event)
}
</pre>

The `Route` method routes events to a method to handle the event
based on the event type. For example:

<pre>
func (ta *TestAgg) Route(event goes.Event) {
	fmt.Printf("Test aggregate: %vEvent: %v\n", ta, event)
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
</pre>

In the above we see the `TestAggCreate` event is routed to
`handleTestAggCreated`, and `TestAggFooUpdate` is routed to
`handleTestFooUpdated`.

Each event requires an event handler method to perform the state
mutations. The handlers may only change aggregate state -- no
other side effects are allowed or they will happen everytime
state is loaded.

Getting back to the split between side effects and state mutation,
behavior is split between two methods - the command method, and the
event handler. Command methods do their side effects, then
`Apply` and event which is handed to the event handler.

Some methods might not have side effects. To support recreation from
event history, then simply construct the event and have the event
handler change state. Factory methods are like this, for exampe:

<pre>
func NewTestAgg(foo, bar, baz string) (*TestAgg, error) {
	//Do validation... return an error if there's a problem
	var testAgg = new(TestAgg)
	testAgg.Aggregate = goes.NewAggregate()
	testAgg.Version = 1

	testAggCreated := TestAggCreated{
		AggregateId: testAgg.ID,
		Foo:         foo,
		Bar:         bar,
		Baz:         baz,
	}

	testAgg.Apply(
		goes.Event{
			Source:  testAgg.ID,
			Version: testAgg.Version,
			Payload: testAggCreated,
		})

	return testAgg, nil
}

func (ta *TestAgg) handleTestAggCreated(event TestAggCreated) {
	ta.ID = event.AggregateId
	ta.Foo = event.Foo
	ta.Bar = event.Bar
	ta.Baz = event.Baz
}
</pre>


Events are flushed from memory to the persistent event store
when Store is called on the aggregate. The Store method must
be supplied by the aggregate - it is called with an EventStore
interface, which is uses to store events via the StoreEvents
method, afterwhich the in memory events list in the aggregate
is cleared.

Note that in the implementation, we also provide an unmarshalling
hint for the persistent storage of the event: when reading back
the bytes that represent the event from the underlying event
storage, we need something that indicates the type to use to
unmarshall the event.

### Summary

To apply event sourcing using this minimal toolkit, the
aggregate responsibilities includes:

* Embedding a pointer to the goes.Aggregate type in the aggregate
type definition.
* Supplying both a factory method to instantiate a new aggregate
instance, and a method to load the events from history.
* Defining Apply and Route methods. The Apply method stores
the event to be routed in memory then calls Route.
* Defining the 'command' methods and the events they generate,
concretely defining types to represent the events.
* Providing an event handler method for each event, and
ensuring the Route method will route event types to the
appropriate handler.
* Providing a Store method that takes an EventStore interface as
an object, storing the events via the StoreEvents method on the
supplied EventStore, then clearing the list of events on the
in memory aggregate.

## Inmems - in memory event store

Example implementation of the Go Event Source event store and event publisher interfaces, with the implementation being in-memory.
