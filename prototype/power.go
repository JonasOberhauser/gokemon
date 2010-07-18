package main

import "rand"

type Power interface {
        GetName() string
        SetName(string) bool
        
        Use()
        
        Unlearnable() bool
}


type attackpower struct {
        name string
        power int                                                                                                    
        accuracy, cooldown float                                                                                                               
        elemType           Element     
}

func TestAttackPower( name string, power int, accuracy, cooldown float, elemType Element ) Power {
        p := new(attackpower)
        p.name = name
        p.power = power
        p.accuracy = accuracy
        p.cooldown = cooldown
        p.elemType = elemType
        return p
}


func (p *attackpower) GetName() string {
        return p.name
}

func (p *attackpower) SetName(string) bool {
        return false
}

func (p *attackpower) Use() {
        if MainWorld.InCombat() {
                a:=MainWorld.GetSingleArena()
                c, ok1 := a.GetCombattant(0)
                t, ok2 := a.GetCombattant(1)
                if ok1 && ok2 && c != nil && t != nil {
                        if c.GetSleepRemaining() > t.GetSleepRemaining() {
                                c,t = t,c
                        }
                        c.SyncAccess(func(){
                                if c.GetHp() > 0 {
                                        
                                        dmg := float(c.GetSpecialDMG( p.power, t)) * c.GetSTAB( p.elemType ) * (0.85 + 0.15*rand.Float())
                                        
                                        t.DealDamage( int(dmg), p.elemType )
                                        
                                        c.SetSleepTime( a, int64( p.cooldown * 1e9 ) )

                                }
                        })
                }
                
        } else if MainWorld.InField() {
                
        }
}

func (p *attackpower) Unlearnable() bool {
        return false
}

type testphyspower struct {
        *attackpower
}

func TestPhysPower(name string, power int, accuracy, cooldown float, elemType Element) Power {
        p := new(testphyspower)
        p.attackpower = TestAttackPower( name, power, accuracy, cooldown, elemType ).( *attackpower )
        return p
}

func (p *testphyspower) Use() {
        if MainWorld.InCombat() {
                a:=MainWorld.GetSingleArena()
                c, ok1 := a.GetCombattant(0)
                t, ok2 := a.GetCombattant(1)
                if ok1 && ok2 && c != nil && t != nil {
                        if c.GetSleepRemaining() > t.GetSleepRemaining() {
                                c,t = t,c
                        }
                        
                        if c.GetHp() > 0 {
                                        
                                dmg := float(c.GetPhysicalDMG( p.power, t)) * c.GetSTAB( p.elemType ) * (0.85 + 0.15*rand.Float())
                                
                                t.DealDamage( int(dmg), p.elemType )
                                
                                c.SetSleepTime( a, int64( p.cooldown * 1e9 ) )

                        }
                        
                }
                
        } else if MainWorld.InField() {
                
        }
}

type testspecpower struct {
        *attackpower
}

func TestSpecPower(name string, power int, accuracy, cooldown float, elemType Element) Power {
        p := new(testspecpower)
        p.attackpower = TestAttackPower( name, power, accuracy, cooldown, elemType ).( *attackpower )
        return p
}

type teststatpower struct {
}
