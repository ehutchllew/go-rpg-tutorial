package entities

import "github.com/ev-the-dev/rpg-tutorial/components"

type Enemy struct {
	*Sprite
	CombatComp    *components.EnemyCombat
	FollowsPlayer bool
}
