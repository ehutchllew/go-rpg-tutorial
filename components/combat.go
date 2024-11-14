package components

type Combat interface {
	Attack() bool
	Attacking() bool
	AttackPower() int
	Damage(amount int)
	Health() int
	Update()
}

type BasicCombat struct {
	attacking   bool
	attackPower int
	health      int
}

func NewBasicCombat(atkPwr, health int) *BasicCombat {
	return &BasicCombat{
		attacking:   false,
		attackPower: atkPwr,
		health:      health,
	}
}

func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return true
}

func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

func (b *BasicCombat) Health() int {
	return b.health
}

func (b *BasicCombat) Update() {}

var _ Combat = (*BasicCombat)(nil)

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
}

func NewEnemyCombat(atkCdwn, atkPwr, health int) *EnemyCombat {
	return &EnemyCombat{
		BasicCombat:     NewBasicCombat(atkPwr, health),
		attackCooldown:  atkCdwn,
		timeSinceAttack: 0,
	}
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
		return true
	}
	return false
}

func (e *EnemyCombat) Update() {
	e.timeSinceAttack += 1
}

var _ Combat = (*EnemyCombat)(nil)
