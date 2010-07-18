package main

type Creature interface {
	SetMaxHp(int) bool
	GetMaxHp() int
	
        DealDamage(int, Element) int
        HealDamage(int) int
	
	GetHp() int
	SetHp(int) bool
	
	GetStrength() int
	GetAgility() int
	GetVitality() int
	GetEnergy() int

	SetStrength(int) bool
	SetAgility(int) bool
	SetVitality(int) bool
	SetEnergy(int) bool
	
        GetAccuracy() float
        GetDefense() float
        
        SetAccuracy( float ) bool
        SetDefense( float ) bool

	GetName() string
	SetName(string) bool

	AddPower(Power) bool
	GetPowers() [4]Power
	SetPower(int, Power) bool

	GetExp() int
	SetExp(int) bool

	GetLevel() int
	SetLevel(int) bool

	GetExpForLevel(int) (int, bool)
	GetLevelForExp(int) (int, bool)

	GetElementList() ElementList

	BoxLevel(int) int
        BoxHp(int) int
        BoxStat(int) int
        
        BoxAccuracy(float) float
        BoxDefense(float) float
        SyncAccess(func())
}

const (
	MaxStat  = 999
	MaxHp    = 999
	MaxLevel = 99
	
	MinAccuracy, MaxAccuracy = 0.5, 1.5
	MinDefense, MaxDefense = 0, 0.75
)


type testcreature struct {
	hp, maxhp  int
	s, a, v, e int
	lvl, exp   int
	powers     [4]Power
	name       string
	exp_req    *[MaxLevel + 1]int
	elems      ElementList
        keyman     KeyMan
}

func TestCreature(name string, s, a, v, e int, lvl int, elems ...Element) Creature {
	c := new(testcreature)
	c.s, c.a, c.v, c.e = s, a, v, e
	c.powers = [4]Power{nil, nil, nil, nil}
	
        c.exp_req = new([MaxLevel + 1]int)
        c.name = name
        c.elems = TestElementList(elems)
        c.keyman = TestKeyMan()
        
	hp := c.BoxMaxHp( 5 + int(1.5*float(c.v+c.lvl)+0.2*float(c.s+c.a+c.e)) )
	
	c.SetMaxHp(hp)
        c.SetLevel(lvl)
	return c
}

func (c *testcreature) GetMaxHp() int {
	return c.maxhp
}

func (c *testcreature) SetMaxHp(maxhp int) bool {
	if maxhp >= 1 && maxhp <= MaxHp {
		delta := maxhp - c.maxhp
		c.maxhp += delta
		if ok := c.SetHp(c.GetHp() + delta); !ok {
			c.SetHp(1)
		}
		return true
	}
	return false
}

func (c *testcreature) DealDamage(dam int, elemType Element) int {
        dam_real := float(dam)                               
        for e := range c.elems.Iter() {                         
                dam_real *= elemType.GetDamageMod(e)            
        }                                                          
        hp := c.BoxHp( c.GetHp() - int(dam_real) )
        dam = hp - c.GetHp()
        c.SetHp(hp)
        return dam
}

func (c *testcreature) HealDamage(dam int) int {
        hp := c.BoxHp( c.GetHp() + dam )
        dam = hp - c.GetHp()
        c.SetHp(hp)
        return dam
}                            


func (c *testcreature) GetHp() int {
	return c.hp
}

func (c *testcreature) SetHp(hp int) bool {
	if hp >= 0 && hp <= c.maxhp {
		c.hp = hp
		return true
	}
	return false
}


func (c *testcreature) GetStrength() int {
	return c.s
}
func (c *testcreature) GetAgility() int {
	return c.a
}
func (c *testcreature) GetVitality() int {
	return c.v
}
func (c *testcreature) GetEnergy() int {
	return c.e
}

func (c *testcreature) SetStrength(s int) bool {
	if s >= 0 && s <= MaxStat {
		c.s = s
	}
	return false
}
func (c *testcreature) SetAgility(a int) bool {
	if a >= 0 && a <= MaxStat {
		c.a = a
	}
	return false
}
func (c *testcreature) SetVitality(v int) bool {
	if v >= 0 && v <= MaxStat {
		c.v = v
	}
	return false
}
func (c *testcreature) SetEnergy(e int) bool {
	if e >= 0 && e <= MaxStat {
		c.e = e
	}
	return false
}


func (c *testcreature) GetAccuracy() float {
        return 1.0
}
func (c *testcreature) GetDefense() float {
        return 0.5
}


func (c *testcreature) SetAccuracy( acc float ) bool {
        return false
}
func (c *testcreature) SetDefense( def float ) bool {
        return false
}


func (c *testcreature) GetName() string {
	return c.name
}

func (c *testcreature) setName(name string) {
	c.name = name
}
func (c *testcreature) SetName(name string) bool {
	if len(name) >= 3 && len(name) <= 10 {
		c.setName(name)
		return true
	}
	return false
}


func (c *testcreature) AddPower(newp Power) bool {
	for id, p := range c.powers {
		if p == nil && c.SetPower(id, newp) {
			return true
		}
	}
	return false
}

func (c *testcreature) GetPowers() [4]Power {
	return c.powers
}
func (c *testcreature) SetPower(id int, p Power) bool {
	if id >= 0 && id <= 3 {
		if c.powers[id].Unlearnable() {
			c.powers[id] = p
			return true
		}
	}
	return false
}


func (c *testcreature) setExp(exp int) {
	c.exp = exp
}

func (c *testcreature) SetExp(exp int) bool {
	if exp >= 0 {
		if lvl, ok := c.GetLevelForExp(exp); ok {
			c.setExp(exp)
			c.setLevel(lvl)
			return true
		}
	}
	return false
}

func (c *testcreature) GetExp() int {
	return c.exp
}

func (c *testcreature) setLevel(lvl int) {
	c.lvl = lvl
}

func (c *testcreature) SetLevel(lvl int) bool {
	if lvl >= 1 && lvl <= MaxLevel {
		c.setLevel(lvl)
		exp, _ := c.GetExpForLevel(lvl)
		c.setExp(exp)
		return true
	}
	return false
}

func (c *testcreature) GetLevel() int {
	return c.lvl
}


func (c *testcreature) GetExpForLevel(lvl int) (int, bool) {
	if lvl >= 1 && lvl <= MaxLevel {
		if c.exp_req[lvl] == 0 {
			e_req_last, _ := c.GetExpForLevel(lvl - 1)
			c.exp_req[lvl] = 5*lvl*(lvl-1) + e_req_last
		}
		return c.exp_req[lvl], true
	}
	return 0, false
}

func (c *testcreature) GetLevelForExp(exp int) (int, bool) {
        if exp >= 0 {
                lvl := MaxLevel / 2
                diff := lvl / 2

                for diff > 0 {
                        switch exp_lower, _ := c.GetExpForLevel(lvl); {
                        case exp_lower >= exp:
                                lvl -= diff
                        case exp_lower <= exp:
                                exp_upper, _ := c.GetExpForLevel(lvl + 1)
                                if exp <= exp_upper {
                                        return lvl, true
                                }
                                lvl += diff
                        }

                        diff /= 2
                }
        }
        return 0, false
        
}


func (c *testcreature) GetElementList() ElementList {
	return c.elems
}


func (c *testcreature) BoxHp(hp int) int {
        switch {
        case hp < 0:
                hp = 0
        case hp > c.GetMaxHp():
                hp = c.GetMaxHp()
        }
        return hp
}

func (c *testcreature) BoxMaxHp(maxhp int) int {
	switch {
	case maxhp < 1:
		maxhp = 1
	case maxhp > MaxHp:
		maxhp = MaxHp
	}
	return maxhp
}
func (c *testcreature) BoxLevel(lvl int) int {
	switch {
	case lvl < 1:
		lvl = 1
	case lvl > MaxLevel:
		lvl = MaxLevel
	}
	return lvl
}
func (c *testcreature) BoxStat(stat int) int {
        switch {
        case stat < 1:
                stat = 1
        case stat > MaxStat:
                stat = MaxStat
        }
        return stat
}

func (c *testcreature) BoxAccuracy(acc float) float{
        switch {
        case acc < MinAccuracy:
                acc = MinAccuracy
        case acc > MaxAccuracy:
                acc = MaxAccuracy
        }
        return acc
}

func (c *testcreature) BoxDefense(def float) float {
        switch {
        case def < MinDefense:
                def = MinDefense
        case def > MaxDefense:
                def = MaxDefense
        }
        return def
}

func (c *testcreature) SyncAccess( foo func() ) {
        c.keyman.Atomic( foo )
}





/*
x  fx
1  10 10 10
2  30 20 10
3  60 30 10
4 100 40 10
*/
