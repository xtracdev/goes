package benchmark

import (
	"testing"
	"github.com/xtracdev/goes"
	"github.com/xtracdev/goes/sample"
	"fmt"
	"github.com/xtracdev/goes/inmems"
)

var eventStore goes.EventStore = inmemes.NewInMemoryEventStore()

func BenchmarkStoreAgg(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		user, err := sample.NewUser("first", "last", "email")
		if err != nil {
			b.Fatal(err.Error())
		}

		for j := 0; j < 10; j++ {
			user.UpdateFirstName("u1 new first")
		}

		user.Store(eventStore)
	}

	fmt.Println("Exit after",b.N, "iterations")
}
