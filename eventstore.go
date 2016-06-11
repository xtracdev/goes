package goes

//EventPublishedCallback definte the type of a callback function invoked
//on behalf of a subscriber when an event is published.
type EventPublishedCallback func(event Event)

//SubscriptionID represents the ID assocaited with a subscriber
type SubscriptionID string

//EventStore defines the methods offered by an EventStore
type EventStore interface {
	StoreEvents(*Aggregate) error
	RetrieveEvents(aggID string) ([]Event, error)
}

//EventPublisher defines the methods an EventPublisher must implement
type EventPublisher interface {
	SubscribeEvents(callback EventPublishedCallback) SubscriptionID
	Unsubscribe(sub SubscriptionID)
}

//EventSourced specifies the methods an event sourced domain object must
//implement.
type EventSourced interface {
	Store(EventStore)
	Apply(Event)
	Route(Event)
}
