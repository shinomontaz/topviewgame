package state

const (
	STAND int = iota
	IDLE
	WALK
)

type owner interface {
	SetState(newId int)
}
