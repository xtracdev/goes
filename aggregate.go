package goes

import (
	"github.com/nu7hatch/gouuid"
)

//Aggregate represents data every persistent domain object or aggregate object
//must track for event sourcing. 
type Aggregate struct {
	ID      string
	Events  []Event
	Version int
}

//NewAggregate returns a pointer to an Aggregate initialized with a 
//uique ID
func NewAggregate() *Aggregate {
	return &Aggregate{
		ID: GenerateID(),
	}
}

//GenerateID generates a unique ID using UUID v4.
func GenerateID() string {
	u, _ := uuid.NewV4()
	return u.String()
}
