package goes

import (
	"github.com/nu7hatch/gouuid"
)

type Aggregate struct {
	ID      string
	Events  []Event
	Version int
}

func NewAggregate() *Aggregate {
	return &Aggregate{
		ID: GenerateID(),
	}
}

func GenerateID() string {
	u, _ := uuid.NewV4()
	return u.String()
}
