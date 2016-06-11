package goes

//Event defines a structure that carries information related to an event, including the 
//source aggregate of the event, the aggregate version the event is associated with, the payload, 
//and the typecode indicating the type of the event.
type Event struct {
	Source   string
	Version  int
	Payload  interface{}
	TypeCode string
}
