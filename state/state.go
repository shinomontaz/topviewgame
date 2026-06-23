package state

const (
	STAND int = iota
	IDLE
	WALK
	ATTACK
	HURT
	DEATH
)

type owner interface {
	SetState(newId int)
	GetDirection() (int, int)
}
