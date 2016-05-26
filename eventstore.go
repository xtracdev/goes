package goes

type EventPublishedCallback func(event Event)

type SubscriptionID string

type EventStore interface {
	StoreEvents(*Aggregate) error
	RetrieveEvents(aggID string) ([]Event, error)
}

type EventPublisher interface {
	SubscribeEvents(callback EventPublishedCallback) SubscriptionID
	Unsubscribe(sub SubscriptionID)
}

type EventSourced interface {
	Store(EventStore)
	Apply(Event)
	Route(Event)
}