package main

import (
        "time"
        )

type Combattant interface {
	GetName() string

	SetSleepTime(Arena, int64) bool
        SetSleepRemaining(Arena, int64) bool
        GetArenaSpeedMod(Arena) float
        
        SetSleepTimeHard(int64) bool
        SetSleepRemainingHard(int64) bool
        GetSleepRemaining() int64
	GetSpeedMod() float
	  
	
        GetSleepProgress() float
        SetSleepProgress(float) bool
	RegenerateAp() bool
	GetApSleep() int64

        GetSTAB(Element) float
        GetPhysicalDMG(int, Combattant) int // How much damage does a physical move with a certain power do?
        GetSpecialDMG(int, Combattant) int  // How much damage does a special move with a certain power do?    
        
        DealDamage(int, Element) int
        HealDamage(int) int
        
	GetMaxHp() int
	GetHp() int
	SetHp(int) bool
        GetTimedHp() int
	
	GetPowers() [4]Power
	GetElementList() ElementList

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

	SetAccuracy(float) bool
	SetDefense(float) bool
	
	
        GetExp() int
        GetLevel() int
        GetExpForLevel(int) (int, bool)

	
	SyncAccess(func())
}


type testcombattant struct {
	Creature
        timed_hp                   int
	ap, max_ap                 uint16
	s_mod, a_mod, v_mod, e_mod int
	def_mod, acc_mod           float
        k                          KeyMan
}

func TestCombattant(cr Creature) Combattant {
	c := new(testcombattant)
	c.Creature = cr
	c.ap = 0
	c.max_ap = 0
        c.timed_hp = c.GetHp()
	c.k = TestKeyMan()
	return c
}

// A modifier for time intervals. A faster gokemon will return a lower time interval.
func (c *testcombattant) GetSpeedMod() float {
	return 1.0 / float( c.GetAgility() + c.GetLevel() )
}

// fastest.GetArenaSpeedMod( arena ) = 1
// returns a value >= 1
func (c *testcombattant) GetArenaSpeedMod( a Arena ) float {
        var fastest Combattant
        fastest = c
        for com := range a.Iter() {
            if com.GetSpeedMod() < fastest.GetSpeedMod() {
                    fastest = com
            }
        }
        // GetSpeedMod : faster <=> lower value
        //   => c.GetSpeedMod() >= fastest.GetSpeedMod()
        return c.GetSpeedMod() / fastest.GetSpeedMod()
}


// Sleep some time, depending on speed
func (c *testcombattant) SetSleepTime(a Arena, amt int64) bool {
        return c.SetSleepTimeHard( int64( float( amt ) * c.GetArenaSpeedMod(a) ) )
}

// Sleep amt nanoseconds
func (c *testcombattant) SetSleepTimeHard(amt int64 ) bool {
        ap := uint16( amt / c.GetApSleep() )
        if 0 <= ap && ap <= 65535 {
                c.ap = 0
                c.max_ap = ap
                return true
        }
        return false
}

func (c *testcombattant) GetSleepRemaining() int64 {
        return int64(c.max_ap-c.ap) * c.GetApSleep()
}

func (c *testcombattant) SetSleepRemaining(a Arena, amt int64 ) bool {
        return c.SetSleepRemainingHard( int64( float( amt ) * c.GetArenaSpeedMod(a) ) )
}

func (c *testcombattant) SetSleepRemainingHard(amt int64 ) bool {
        ap := c.ap + uint16(amt / c.GetApSleep())
        if 0 <= ap && ap <= 65535 {
                c.max_ap = ap
                return true
        }
        return false
}

func (c *testcombattant) GetSleepProgress() float {
	return float(c.ap) / float(c.max_ap)
}
func (c *testcombattant) SetSleepProgress( p float ) bool {
        if p < 0 || p > 1 {
                return false 
        }
        c.ap = uint16( float( c.max_ap ) * p )
        return true
}

// Regenerates Ap and then returns true if Ap are fully regenerated 
func (c *testcombattant) RegenerateAp() bool {
        if c.ap >= c.max_ap {
                return true 
        }
        c.ap += 1
        return false
}


// ap in ns
func (c *testcombattant) getBaseApSleep() int64 {
        return 4e7
}


// fast 0.75 -> x1
// slow 1 -> x1.33
// fast.Sleep( 3e9 ) -> fast.Sleep( 3e9 )
// 3e9 / 4e7 = 75 AP
// slow.Sleep( 3e9 ) -> slow.Sleep( 4e9 )
// 4e9 / 4e7 = 100 AP
// -> different amount of AP, BUT:
// fast regenerates 25 AP per second
// slow regenerates 25 AP per second


// All Creatures regenerate the same amount of AP per second.
func (c *testcombattant) GetApSleep() int64 {
        return 4e7
}


func (c *testcombattant) GetSTAB ( elemType Element ) float {
        stab := 1.0
        for e := range c.GetElementList().Iter() { 
                if e == elemType {
                        // Apply STAB
                        stab *= 1.5
                }
        }
        return stab
}

func (c *testcombattant) GetPhysicalDMG( p int, t Combattant ) int {
        l := p * ( c.GetStrength() + c.GetLevel() )
        d := t.GetDefense() / 150.0
        mod := ( 100.0 + 0.75 * float( c.GetStrength() ) ) / 100.0
        return int( (float(l) * d + 2.0 ) * mod )
}

func (c *testcombattant) GetSpecialDMG( p int, t Combattant ) int {
        l := p * ( c.GetEnergy() + c.GetLevel() )
        d := t.GetDefense() / 150.0
        mod := ( 100.0 + 0.75 * float( c.GetEnergy() ) ) / 100.0
        return int( (float(l) * d + 2.0 ) * mod )
}

func (c *testcombattant) GetTimedHp() int {
        return c.timed_hp
}

func sgn (v int) int {
        if v >= 0 {
                return 1
        }
        return -1
}

func (c *testcombattant) adjustTimedHp() {
        go c.k.Atomic( func(){
                intervals := 18 // intervals for full hp
                t := int64( 5e9 ) // time for full hp
                var diff_hp int // how large is the difference between timed and real?
                var direction int // -1 or +1
                
                dT := t / int64( intervals ) // time between two refreshes
                dHp := c.GetMaxHp() / intervals // hp refreshed
                
                c.SyncAccess( func(){
                        diff_hp = c.timed_hp - c.GetHp()
                        direction = sgn( diff_hp )
                } )
                
                for {
                        if direction * diff_hp <= dHp {
                                break
                        }
                        c.SyncAccess( func(){
                                c.timed_hp -= direction * dHp
                        })
                        
                        
                        time.Sleep( dT )
                        
                        
                        c.SyncAccess( func(){
                                dHp = c.GetMaxHp() / intervals 
                                diff_hp = c.timed_hp - c.GetHp()
                                direction = sgn( diff_hp )
                        } )
                }
                
                c.SyncAccess( func(){
                        c.timed_hp -= diff_hp
                })
                
        } )
}

func (c *testcombattant) DealDamage(dam int, elemType Element) int {
        c.SyncAccess( func(){
                dam = c.Creature.DealDamage(dam, elemType)
        })
        
        c.adjustTimedHp()
        
        return dam
}

func (c *testcombattant) HealDamage(dam int) int {
        c.SyncAccess( func(){
                c.timed_hp = c.GetHp()
                dam = c.Creature.HealDamage(dam)
        })
        
        c.adjustTimedHp()
        
        return dam
}

func (c *testcombattant) GetStrength() int {
	return c.Creature.GetStrength() + c.s_mod
}

func (c *testcombattant) GetAgility() int {
	return c.Creature.GetAgility() + c.a_mod
}

func (c *testcombattant) GetVitality() int {
	return c.Creature.GetVitality() + c.v_mod
}

func (c *testcombattant) GetEnergy() int {
	return c.Creature.GetEnergy() + c.e_mod
}

func (c *testcombattant) setStrengthMod(s_mod int) {
	c.s_mod = s_mod
}
func (c *testcombattant) setAgilityMod(a_mod int) {
	c.a_mod = a_mod
}
func (c *testcombattant) setVitalityMod(v_mod int) {
	c.v_mod = v_mod
}
func (c *testcombattant) setEnergyMod(e_mod int) {
	c.e_mod = e_mod
}

func (c *testcombattant) SetStrength(s int) bool {
	if s >= 0 && s <= MaxStat {
		c.setStrengthMod(s - c.Creature.GetStrength())
	}
	return false
}
func (c *testcombattant) SetAgility(a int) bool {
	if a >= 0 && a <= MaxStat {
		c.setAgilityMod(a - c.Creature.GetAgility())
	}
	return false
}
func (c *testcombattant) SetVitality(v int) bool {
	if v >= 0 && v <= MaxStat {
		c.setVitalityMod(v - c.Creature.GetVitality())
	}
	return false
}
func (c *testcombattant) SetEnergy(e int) bool {
	if e >= 0 && e <= MaxStat {
		c.setEnergyMod(e - c.Creature.GetEnergy())
	}
	return false
}

func (c *testcombattant) GetAccuracy() float {
	return c.Creature.GetAccuracy() + c.acc_mod
}
func (c *testcombattant) GetDefense() float {
	return c.Creature.GetAccuracy() + c.def_mod
}

func (c *testcombattant) SetAccuracy(acc float) bool {
	return false
}
func (c *testcombattant) SetDefense(def float) bool {
	return false
}
