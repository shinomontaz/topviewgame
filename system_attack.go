package main

import (
	"fmt"
	"topviewgame/state"

	"github.com/bytearena/ecs"
)

func ProcessAttacks(g *Game, attPos, defPos *Position) {
	var (
		attacker *ecs.QueryResult = nil
		defender *ecs.QueryResult = nil
	)

	for _, playerCombatant := range g.World.QueryPlayers() {
		pos := g.World.GetPosition(playerCombatant)
		if pos.IsEqual(attPos) {
			attacker = playerCombatant
		} else if pos.IsEqual(defPos) {
			defender = playerCombatant
		}
	}

	for _, cbt := range g.World.QueryMonsters() {
		if g.World.GetMonster(cbt).(*Monster).IsDead() {
			continue
		}
		pos := g.World.GetPosition(cbt)
		if pos.IsEqual(attPos) {
			attacker = cbt
		} else if pos.IsEqual(defPos) {
			defender = cbt
		}
	}

	if attacker == nil || defender == nil {
		return
	}

	defenderArmor := g.World.GetArmor(defender)
	defenderHealth := g.World.GetHealth(defender)
	defenderName := g.World.GetName(defender).Label
	attackerWeapon := g.World.GetMeleeWeapon(attacker)
	attackerName := g.World.GetName(attacker).Label
	defenderMessage := g.World.GetUserMessage(defender)
	attackerMessage := g.World.GetUserMessage(attacker)

	if g.World.GetHealth(attacker).Current <= 0 {
		return
	}

	if defenderName == "Player" {
		g.World.GetPlayer(defender).(*Player).SetState(state.STAND)
	} else {
		g.World.GetMonster(defender).(*Monster).SetState(state.STAND)
	}

	toHitRoll := GetDiceRoll(10)

	if toHitRoll+attackerWeapon.ToHitBonus > defenderArmor.Dodge {
		damageRoll := attackerWeapon.MinDamage + GetDiceRoll(attackerWeapon.MaxDamage-attackerWeapon.MinDamage)
		damageDone := damageRoll - defenderArmor.Defence

		if damageDone < 0 {
			damageDone = 0
		}

		defenderHealth.Current -= damageDone

		attackerMessage.AttackMessage = fmt.Sprintf("%s hit %s for %d damage\n", attackerName, defenderName, damageDone)
		if defenderHealth.Current <= 0 {
			if defenderName == "Player" {
				defenderMessage.GameStateMessage = "Game Over!\n"
				g.Turn = GameOver

				return
			}
			defenderMessage.DeadMessage = fmt.Sprintf("%s is dead\n", defenderName)

			l := g.Map.CurrentLevel
			t := l.Tiles[l.GetIndexFromXY(defPos.X, defPos.Y)]
			t.Blocked = false
			g.gm.updateMonsterPosition(defender.Entity, defPos, nil)
			g.World.GetMonster(defender).(*Monster).SetState(state.DEATH)
		}
	} else {
		attackerMessage.AttackMessage = fmt.Sprintf("%s missed %s", attackerName, defenderName)
	}
}
