package event

const (
	EventNone EventType = iota
	EventPass
	EventClick
	EventKey
)

type EventType int

type Event struct {
	Type EventType
	Pos  [2]int
}
