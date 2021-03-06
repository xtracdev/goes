package aggregates

import (
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/goes/sample"
)

var u1, u2 *sample.User

func init() {
	Given(`^an event sourced aggregate$`, func() {
	})

	When(`^I create new instances$`, func() {
		u1, _ = sample.NewUser("first", "last", "one@foo.com")
		u2, _ = sample.NewUser("first", "last", "two@foo.com")

		assert.NotNil(T, u1)
		assert.NotNil(T, u2)
	})

	Then(`^the instances have unique IDs$`, func() {
		assert.NotEqual(T, u1.AggregateID, u2.AggregateID, "IDs were not unique")
	})

	And(`^there's an uncommitted event$`, func() {
		assert.Equal(T, 1, len(u1.Events))
		assert.Equal(T, 1, len(u2.Events))
	})

	And(`^the uncommited event's source ID is the aggregate ID$`, func() {
		assert.Equal(T, u1.AggregateID, u1.Events[0].Source)
		assert.Equal(T, u2.AggregateID, u2.Events[0].Source)
	})
}
