package goes

type Event struct {
	Source  string
	Version int
	Payload interface{}
}