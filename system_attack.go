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

	for _, playerCombatant := range g.World.Query(g.WorldTags["players"]) {
		pos := playerCombatant.Components[positionC].(*Position)
		if pos.IsEqual(attPos) {
			attacker = playerCombatant
		} else if pos.IsEqual(defPos) {
			defender = playerCombatant
		}
	}

	for _, cbt := range g.World.Query(g.WorldTags["monsters"]) {
		if cbt.Components[monsterC].(*Monster).IsDead() {
			continue
		}
		pos := cbt.Components[positionC].(*Position)
		if pos.IsEqual(attPos) {
			attacker = cbt
		} else if pos.IsEqual(defPos) {
			defender = cbt
		}
	}

	if attacker == nil || defender == nil {
		return
	}

	defenderArmor := defender.Components[armorC].(*Armor)
	defenderHealth := defender.Components[healthC].(*Health)
	defenderName := defender.Components[nameC].(*Name).Label
	attackerWeapon := attacker.Components[meleeWeaponC].(*MeleeWeapon)
	attackerName := attacker.Components[nameC].(*Name).Label
	defenderMessage := defender.Components[userMessage].(*UserMessage)
	attackerMessage := attacker.Components[userMessage].(*UserMessage)

	if attacker.Components[healthC].(*Health).Current <= 0 {
		return
	}

	if defenderName == "Player" {
		defender.Components[playerC].(*Player).SetState(state.STAND)
	} else {
		defender.Components[monsterC].(*Monster).SetState(state.STAND)
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
			defender.Components[monsterC].(*Monster).SetState(state.DEATH)
		}
	} else {
		attackerMessage.AttackMessage = fmt.Sprintf("%s missed %s", attackerName, defenderName)
	}
}
