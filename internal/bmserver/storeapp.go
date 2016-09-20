package main

import (
	"net/http"
	"log"
	"github.com/xtracdev/goes"
	"fmt"
	"strconv"
	"github.com/xtracdev/goes/sample"
	_ "net/http/pprof"
	"github.com/xtracdev/goes/inmems"
)

var (
	eventStore goes.EventStore = inmemes.NewInMemoryEventStore()
	aggCount = 0
	eventCount = 0
)

func statsHandler(rw http.ResponseWriter, req *http.Request) {
	msg := fmt.Sprintf("Stored %v aggregates and %v events\n", aggCount, eventCount)
	rw.Write([]byte(msg))
}

func benchHandler(rw http.ResponseWriter, req *http.Request) {
	numAggregates, err := strconv.Atoi(req.FormValue("aggs"))
	if err != nil {
		http.Error(rw,err.Error(), http.StatusBadRequest)
		return
	}

	eventsPerAgg,err := strconv.Atoi(req.FormValue("eventsPerAgg"))
	if err != nil {
		http.Error(rw,err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("number of aggregates: %d events per agg: %d", numAggregates, eventsPerAgg)

	for i := 0; i < numAggregates; i++ {
		aggCount++
		user, _ := sample.NewUser("first", "last", "email")
		eventCount++ //NewUser generates a create event


		//We create eventsPerAgg - 1 because creating an aggregate means a create event
		//has been generated
		for j := 0; j < eventsPerAgg - 1; j++ {
			eventCount++
			user.UpdateFirstName("u1 new first")
		}

		err = user.Store(eventStore)
		if err != nil {
			http.Error(rw,err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	http.HandleFunc("/bench", benchHandler)
	http.HandleFunc("/stats", statsHandler)
	http.ListenAndServe(":8080", nil)
}
