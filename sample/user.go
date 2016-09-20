package sample

import (
	"github.com/xtracdev/goes"
	"log"
)

//To be event sourced...
// The Aggregate type must be embedded
// A constructor for a brand new aggregate is available
// Constructing an aggregate produces an event
// Mutations occur via commands, which emit events handled by event handler, with events routed to handlers
// Events are recorded in event history
// An apply method routes an event to the event handler, and records the event
// When applying event history, only the route method is used - side effects occur in the command handlers

//User defines a simple domain object that will have methods that support event sourcing.
type User struct {
	*goes.Aggregate
	FirstName string
	LastName  string
	Email     string
}

//NewUser instantiates an instance of User, and initializes the embedded aggregate structure.
func NewUser(first, last, email string) (*User, error) {
	//Do validation... return an error if there's a problem
	var user = new(User)
	user.Aggregate = goes.NewAggregate()

	user.Version = 1
	user.Apply(
		goes.Event{
			Source:  user.ID,
			Version: user.Version,
			Payload: UserCreated{
				AggregateId: user.ID,
				FirstName:   first,
				LastName:    last,
				Email:       email,
			},
		})

	return user, nil
}

//NewUserFromHistory instantiates a User and applies its event history to derive the current
//state of hte aggregate.
func NewUserFromHistory(events []goes.Event) *User {
	user := new(User)
	user.Aggregate = goes.NewAggregate()

	for _, e := range events {
		log.Println("apply event", e)
		user.Version++
		user.Route(e)
	}

	return user
}

//UserCreated is the event generated when a user struct is first instantiated.
type UserCreated struct {
	AggregateId string
	FirstName   string
	LastName    string
	Email       string
}

//UserFirstNameUpdated is an event generated when the first name is updated.
type UserFirstNameUpdated struct {
	OldFirst string
	NewFirst string
}

//UserLastNameUpdated is an event generated when the last name is updated.
type UserLastNameUpdated struct {
	OldLast string
	NewLast string
}

//UpdateFirstName is a command handler that handles updating the user first name,
//generating a UserFirstNameUpdated event.
func (u *User) UpdateFirstName(first string) {
	u.Version++
	u.Apply(
		goes.Event{
			Source:  u.ID,
			Version: u.Version,
			Payload: UserFirstNameUpdated{
				OldFirst: u.FirstName,
				NewFirst: first,
			},
		})
}

func (u *User) handleUserCreated(event UserCreated) {
	u.Aggregate.ID = event.AggregateId
	u.FirstName = event.FirstName
	u.LastName = event.LastName
	u.Email = event.Email
}

func (u *User) handleUserFirstNameUpdate(event UserFirstNameUpdated) {
	u.FirstName = event.NewFirst
}

//Route is the standard method for routing events to event handlers.
func (u *User) Route(event goes.Event) {
	event.Version = u.Version
	switch event.Payload.(type) {
	case UserCreated:
		u.handleUserCreated(event.Payload.(UserCreated))
	case UserFirstNameUpdated:
		u.handleUserFirstNameUpdate(event.Payload.(UserFirstNameUpdated))
	default:
		panic("WARN: unknown event routed to User aggregate")
	}
}

//Apply is the standard event sourcing method that routes an event then records
//the event in the event history
func (u *User) Apply(event goes.Event) {
	u.Route(event)
	u.Events = append(u.Events, event)
}

//Store uses the event store passed to it to persistently recorded
//the current set of unperistent events.
func (u *User) Store(eventStore goes.EventStore) error {
	err := eventStore.StoreEvents(u.Aggregate)
	if err != nil {
		return err
	}

	u.Events = make([]goes.Event, 0)

	return nil
}
