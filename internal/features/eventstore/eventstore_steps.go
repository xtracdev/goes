package eventstore

import (
	"github.com/xtracdev/goes"
	"github.com/xtracdev/goes/sample"
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/goes/inmems"
)

func init() {
	var user, user2 *sample.User
	var eventStore goes.EventStore

	Given(`^an aggregate$`, func() {
		var err error
		user, err = sample.NewUser("first", "last", "email")
		assert.Nil(T, err)
	})

	And(`^an event store$`, func() {
		eventStore = inmemes.NewInMemoryEventStore()
	})

	When(`^the aggregate has uncommitted events$`, func() {

	})

	And(`^the events are stored$`, func() {
		err := user.Store(eventStore)
		assert.Nil(T, err)
		assert.Equal(T, 0, len(user.Events))
	})

	Then(`^the events for that aggregate can be retrieved$`, func() {
		eventSets, err := eventStore.RetrieveEvents(user.ID)
		assert.Nil(T, err)
		assert.Equal(T, 1, len(eventSets), "Expected one event set to be retrieved")
	})

	And(`^the aggregate state can be recreated using the events$`, func() {
		events, err := eventStore.RetrieveEvents(user.ID)
		assert.Nil(T, err)
		retUser := sample.NewUserFromHistory(events)
		assert.NotNil(T, retUser)
		assert.Equal(T, user.FirstName, retUser.FirstName)
		assert.Equal(T, user.LastName, retUser.LastName)
		assert.Equal(T, user.Email, retUser.Email)
	})

	Given(`^two aggregates$`, func() {
		user2, _ = sample.NewUser("first2", "last2", "email2")
		user2.UpdateFirstName("new first")
		assert.Equal(T, "new first", user2.FirstName)
	})

	And(`^an event store containing all events for both aggregates$`, func() {
		err := user2.Store(eventStore)
		assert.Nil(T, err)
		assert.Equal(T, 0, len(user2.Events))
	})

	When(`^the events for an aggregate are retrieved$`, func() {
		events, err := eventStore.RetrieveEvents(user2.ID)
		assert.Nil(T, err)
		assert.Equal(T, 2, len(events))

	})

	Then(`^only the events associated with the specific aggregate are retrieved$`, func() {
		events, err := eventStore.RetrieveEvents(user2.ID)
		assert.Nil(T, err)
		retUser := sample.NewUserFromHistory(events)
		assert.NotNil(T, retUser)
		assert.Equal(T, user2.FirstName, retUser.FirstName)
		assert.Equal(T, user2.LastName, retUser.LastName)
		assert.Equal(T, user2.Email, retUser.Email)
	})

	When(`^I add an event$`, func() {
		user.UpdateFirstName("new first")
	})

	Then(`^the aggregate version is incremented$`, func() {
		assert.Equal(T, 2, user.Version, "Expected user to be at version 2")
	})

	And(`^the aggregate version is correct when built from event history$`, func() {
		user.Store(eventStore)
		events, err := eventStore.RetrieveEvents(user.ID)
		assert.Nil(T, err)
		retUser := sample.NewUserFromHistory(events)
		assert.Equal(T, 2, retUser.Version, "Expected retrieved user to be at version 2")

	})

	var u1, u2 *sample.User
	When(`^it is modified by two concurrent threads of control$`, func() {
		u1 = user
		u2 = user
		u1.UpdateFirstName("u1 new first")
		u2.UpdateFirstName("u2 new first")

	})

	Then(`^the second thread that stored the aggregate gets a concurrency error$`, func() {
		err := u1.Store(eventStore)
		assert.Nil(T, err)
		err = u2.Store(eventStore)
		assert.NotNil(T, err)
	})

	And(`^all the events in the event history have the aggregate id as their source$`, func() {
		events, err := eventStore.RetrieveEvents(user.ID)
		assert.Nil(T, err)
		for _, e := range events {
			assert.Equal(T, user.ID, e.Source)
		}
	})

}
