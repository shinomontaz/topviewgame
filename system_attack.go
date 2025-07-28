package main

import (
	"fmt"

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
			fmt.Println(defenderName, "is dead")
			if defenderName == "Player" {
				fmt.Println("You died!")
				g.Turn = GameOver
			}
			g.World.DisposeEntity(defender.Entity)
		}
	} else {
		fmt.Println(attackerName, "missed", defenderName)
	}
}
