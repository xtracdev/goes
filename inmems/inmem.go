package inmemes

import (
	"errors"
	"sync"

	"github.com/xtracdev/goes"
)

type subscriberStorage struct {
	subscriberID goes.SubscriptionID
	callback     goes.EventPublishedCallback
}

type eventStorage struct {
	events         []goes.Event
	currentVersion int
}

//InMemoryEventStore implements the
type InMemoryEventStore struct {
	sync.RWMutex
	storage     map[string]eventStorage
	subscribers []subscriberStorage
}

//NewInMemoryEventStore is a factory method for creating InMemoryEventStore
//instances.
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		storage: make(map[string]eventStorage),
	}
}

func (im *InMemoryEventStore) publishEvent(event goes.Event) {
	for _, sub := range im.subscribers {
		sub.callback(event)
	}
}

//StoreEvents stores the events for the given aggregate in the event
//store.
func (im *InMemoryEventStore) StoreEvents(agg *goes.Aggregate) error {
	im.Lock()
	defer im.Unlock()

	//Do we have events for this aggregate?
	aggStorage, ok := im.storage[agg.AggregateID]
	if !ok {
		aggStorage = eventStorage{}
	}

	//Has someone update the aggregate before the current caller?
	if !(aggStorage.currentVersion < agg.Version) {
		return errors.New("Concurrency exception")
	}

	//Set the new version, and append the events
	aggStorage.currentVersion = agg.Version
	for _, e := range agg.Events {
		aggStorage.events = append(aggStorage.events, e)
		im.publishEvent(e)
	}

	im.storage[agg.AggregateID] = aggStorage

	return nil
}

//RetrieveEvents retrieves the events in the event store assocaited with the given
//aggregate id.
func (im *InMemoryEventStore) RetrieveEvents(aggregateID string) ([]goes.Event, error) {
	im.RLock()
	defer im.RUnlock()

	eventStorage, ok := im.storage[aggregateID]
	if !ok {
		return nil, errors.New("No events stored for aggregate")
	}

	return eventStorage.events, nil
}

//SubscribeEvents registers the provided callback as an event subscriber.
func (im *InMemoryEventStore) SubscribeEvents(callback goes.EventPublishedCallback) goes.SubscriptionID {
	im.Lock()
	defer im.Unlock()
	subscriptionID := goes.SubscriptionID(goes.GenerateID())
	im.subscribers = append(im.subscribers, subscriberStorage{subscriberID: subscriptionID, callback: callback})
	return subscriptionID

}

//Unsubscribe removes the event subscription associated with the provided
//subscription id.
func (im *InMemoryEventStore) Unsubscribe(subscriptionID goes.SubscriptionID) {
	im.Lock()
	remainingSubs := make([]subscriberStorage, 0, len(im.subscribers)-1)
	for _, sub := range im.subscribers {
		if sub.subscriberID != subscriptionID {
			remainingSubs = append(remainingSubs, sub)
		}
	}
	im.subscribers = remainingSubs
	im.Unlock()
}
