package main

type TurnState int

const (
	BeforePlayerAction TurnState = iota
	PlayerTurn
	EnemyTurn
	GameOver
)

func GetNextState(state TurnState) TurnState {
	switch state {
	case BeforePlayerAction:
		return PlayerTurn
	case PlayerTurn:
		return EnemyTurn
	case EnemyTurn:
		return BeforePlayerAction
	case GameOver:
		return GameOver
	default:
		return PlayerTurn
	}
}
