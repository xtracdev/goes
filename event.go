package goes

//Event defines a structure that carries information related to an event, including the 
//source aggregate of the event, the aggregate version the event is associated with, the payload, 
//and the typecode indicating the type of the event.
//
//Note that events for an aggregate can be ordered by version; version is incremented for each event
//associated with an aggregate. Event storage will also typically include a timestamp column for
//the absolute ordering of events in terms of their storage date.
type Event struct {
	Source   string
	Version  int
	Payload  interface{}
	TypeCode string
}
