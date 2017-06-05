package eventpub

import (
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/goes"
	"github.com/xtracdev/goes/inmems"
	"github.com/xtracdev/goes/sample"
)

func init() {

	var user *sample.User
	var events, republishedEvents []goes.Event
	var inMemEventStore = inmemes.NewInMemoryEventStore()
	var eventStore goes.EventStore = inMemEventStore
	var eventPublisher goes.EventPublisher = inMemEventStore
	var eventRepublisher goes.EventRepublisher = inMemEventStore
	var subId goes.SubscriptionID

	var callback = func(event goes.Event) {
		events = append(events, event)
	}

	var republishedCallback = func(event goes.Event) {
		republishedEvents = append(republishedEvents, event)
	}

	When(`^I create and modify an instance of the aggregate$`, func() {
		user, _ = sample.NewUser("first", "last", "email")
		user.UpdateFirstName("updated")
		subId = eventPublisher.SubscribeEvents(callback)
		user.Store(eventStore)
	})

	Then(`^all the events are published$`, func() {
		assert.Equal(T, 2, len(events))
	})

	Then(`^no events are published$`, func() {
		eventHistory, _ := eventStore.RetrieveEvents(user.AggregateID)
		sample.NewUserFromHistory(eventHistory)
		assert.Equal(T, 2, len(events))
	})

	Given(`^an event store with a registered subscriber$`, func() {
	})

	When(`^the subscriber unsubscribes$`, func() {
		eventPublisher.Unsubscribe(subId)
	})

	Then(`^previously subscribed callback is not invoked when events are published$`, func() {
		user, _ = sample.NewUser("first", "last", "email")
		user.Store(eventStore)
		assert.Equal(T, 2, len(events))
	})

	Given(`^a populated event store$`, func() {
		eventPublisher.SubscribeEvents(republishedCallback)
	})

	When(`^republish all events is called$`, func() {
		eventRepublisher.RepublishAllEvents()
	})

	Then(`^the events store events are republished$`, func() {
		assert.Equal(T, 3, len(republishedEvents))
	})

}
