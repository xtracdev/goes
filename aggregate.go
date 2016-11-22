package goes

import (
	"github.com/xtracdev/goes/uuid"
)

//Aggregate represents data every persistent domain object or aggregate object
//must track for event sourcing.
type Aggregate struct {
	AggregateID string
	Events      []Event
	Version     int
}

//NewAggregate returns a pointer to an Aggregate initialized with a
//uique ID
func NewAggregate() (*Aggregate, error) {
	aggId, err := GenerateID()
	if err != nil {
		return nil, err
	}
	return &Aggregate{
		AggregateID: aggId,
	}, nil
}

//GenerateID generates a unique ID using UUID v4.
func GenerateID() (string, error) {
	u, err := uuid.GenerateUuidV4()
	if err != nil {
		return "", err
	}
	return u, nil
}
