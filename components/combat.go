package components

type Combat interface {
	AttackPower() int
	Damage(amount int)
	Health() int
}

type BasicCombat struct {
	attackPower int
	health      int
}

func NewBasicCombat(atkPwr, health int) *BasicCombat {
	return &BasicCombat{
		attackPower: atkPwr,
		health:      health,
	}
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

var _ Combat = (*BasicCombat)(nil)
