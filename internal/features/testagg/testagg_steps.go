package testagg

import (
	. "github.com/gucumber/gucumber"
	"github.com/xtraclabs/goes/sample/testagg"

	"github.com/xtraclabs/goes"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/goes/inmems"
)

func init() {
	var ta *testagg.TestAgg
	var eventStore goes.EventStore

	Given(`^a TestAgg aggregate$`, func() {
		var err error
		ta,err = testagg.NewTestAgg("f","b","b")
		assert.Nil(T,err)
	})

	And(`^an event store for storing TestAgg events$`, func() {
		eventStore = inmemes.NewInMemoryEventStore()
	})

	When(`^the TestAgg aggregate has uncommitted events$`, func() {
		assert.Equal(T, 1, len(ta.Events))
	})

	And(`^the TestAgg events are stored$`, func() {
		err := ta.Store(eventStore)
		assert.Nil(T, err)
		assert.Equal(T, 0, len(ta.Events))
	})

	Then(`^the events for the TestAgg aggregate can be retrieved$`, func() {
		eventSets, err := eventStore.RetrieveEvents(ta.ID)
		assert.Nil(T, err)
		assert.Equal(T, 1, len(eventSets), "Expected one event set to be retrieved")
	})

	And(`^the TestAgg aggregate state can be recreated using the events$`, func() {
		events, err := eventStore.RetrieveEvents(ta.ID)
		assert.Nil(T, err)
		retTestAgg := testagg.NewTestAggFromHistory(events)
		assert.Equal(T, "f",retTestAgg.Foo)
		assert.Equal(T, "b",retTestAgg.Bar)
		assert.Equal(T, "b",retTestAgg.Baz)
	})

	When(`^I update foo$`, func() {
		ta.UpdateFoo("new foo for you")
	})

	Then(`^the TestAgg aggregate version is incremented$`, func() {
		assert.Equal(T, 2, ta.Version, "TestAgg version not expected value of 2")
	})

	And(`^the TestAgg aggregate version is correct when built from event history$`, func() {
		err := ta.Store(eventStore)
		assert.Nil(T, err)
		assert.Equal(T, 0, len(ta.Events))
		events, err := eventStore.RetrieveEvents(ta.ID)
		assert.Nil(T, err)
		retTestAgg := testagg.NewTestAggFromHistory(events)
		assert.Equal(T, 2, retTestAgg.Version, "Rebuilt TestAgg version not expected value of 2")
	})

	And(`^all the events in the Test event history have the aggregate id as their source$`, func() {
		events, err := eventStore.RetrieveEvents(ta.ID)
		assert.Nil(T, err)
		for _, e := range events {
			assert.Equal(T, ta.ID, e.Source)
		}
	})

}

