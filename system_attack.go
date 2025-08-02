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

	if attacker.Components[healthC].(*Health).Current <= 0 {
		return
	}

	toHitRoll := GetDiceRoll(10)

	if toHitRoll+attackerWeapon.ToHitBonus > defenderArmor.Dodge {
		damageRoll := attackerWeapon.MinDamage + GetDiceRoll(attackerWeapon.MaxDamage-attackerWeapon.MinDamage)
		damageDone := damageRoll - defenderArmor.Defence

		if damageDone < 0 {
			damageDone = 0
		}

		defenderHealth.Current -= damageDone

		fmt.Println(attackerName, "hit", defenderName, "for", damageDone, "damage")

		if defenderHealth.Current <= 0 {
			if defenderName == "Player" {
				fmt.Println("You died!")
				g.Turn = GameOver

				return
			}

			fmt.Println(defenderName, "is dead")

			//			g.World.DisposeEntity(defender.Entity)

			l := g.Map.CurrentLevel
			t := l.Tiles[l.GetIndexFromXY(defPos.X, defPos.Y)]
			t.Blocked = false
			defender.Components[monsterC].(*Monster).SetState(state.DEATH)
		}
	} else {
		fmt.Println(attackerName, "missed", defenderName)
	}
}
