package main

import (
	"fmt"
	"time"
	"strings"
	"rand"
	"char"
	"curses"
)

var (
	GAME_SPEED_MOD = 2.0
)

var (
	mainQueue = TestQueue()
	mainCIn   = make(chan bool)
	mainCOut  = make(chan bool)
	mainQuit  = make(chan bool)
	mainP     = int64(1e7)
)

type World interface {
	GetSingleArena() Arena
	PauseCombat(bool)

	InCombat() bool
	InField() bool
}

type world struct {
	a                Arena
	pausedC, pausedF bool
}

func NewWorld() World {
	w := new(world)
	w.a = TestArena1()
	return w
}

func (w *world) GetSingleArena() Arena {
	return w.a
}

func (w *world) PauseCombat(b bool) {
	if b && (!w.pausedC) {
		select {
		case mainCIn <- true:
			w.pausedC = true
		}
	} else if !b && w.pausedC {
		select {
		case <-mainCOut:
			w.pausedC = false
		}
	}
}

func (w *world) PausedCombat() bool {
	return w.pausedC
}
func (w *world) PausedField() bool {
	return true
}

func (w *world) InCombat() bool {
	return true
}

func (w *world) InField() bool {
	return false
}

var MainWorld = NewWorld()

func mainloop() {
	for {
		select {
		case <-mainCIn:
			mainCOut <- true

		default:
			time.Sleep(mainP)
			mainQueue.Update(int64(float(mainP) * GAME_SPEED_MOD))

		}
	}
}


var (
	Whole = "%-43s\n%-43s\n%-43s\n%-43s\n\n\n\n\n\n\n\n%43s\n%43s\n%43s\n%43s\n%43s\n"
	Bar   = "-------------------------"
	Name  = fmt.Sprintf("%-19s", " %-10s L%3d")
	Hp    = "  HP %18s "
	Ap    = "  AP %18s "
	Ep    = "  EP %18s "
	//       012345678901234567
	HpVal = " %3d/%3d"
)

var (
	CombatWin curses.Window
	CombatPan curses.Panel
)

func viewloop() {

	for {

		time.Sleep(3e7)
		if MainWorld.InCombat() {
			a := MainWorld.GetSingleArena()
			c_self, _ := a.GetCombattant(0)
			c_other, _ := a.GetCombattant(1)

			var name string
			var l, e, hp, ap int
			var t_hp, v_hp, m_hp, e_curLvl, e_nextLvl, cur_e int
			var progress float
			var show, canLvl bool

			a.SyncAccess(func() {
				if show = c_self != nil; show {
					c := c_self
					name = c.GetName()
					l = c.GetLevel()
					e_curLvl, _ = c.GetExpForLevel(l)
					e_nextLvl, canLvl = c.GetExpForLevel(l + 1)
					cur_e = c.GetExp()
					t_hp = c.GetTimedHp()
					v_hp = c.GetHp()
					m_hp = c.GetMaxHp()
					progress = c.GetSleepProgress()
				}
			})

			bar_self := fmt.Sprintf(Bar)
			name_self := ""
			hp_self := ""
			ap_self := ""
			ep_self := ""

			if show {
				//hp_val:=fmt.Sprintf( HpVal, v_hp, m_hp )
				hp = int(float(18*t_hp)/float(m_hp) + 0.99)
				ap = int(18.0*progress + 0.5)

				if canLvl {
					e = int(
						float(18*(cur_e-e_curLvl))/
							float(e_nextLvl-e_curLvl) + 0.99)
				} else {
					e = 18
				}

				name_self = fmt.Sprintf(Name, name, l)
				hp_self = fmt.Sprintf(Hp, fmt.Sprintf("%s%s",
					strings.Repeat("=", hp),
					strings.Repeat("-", 18-hp)))
				ap_self = fmt.Sprintf(Ap, fmt.Sprintf("%s%s",
					strings.Repeat("=", ap),
					strings.Repeat("-", 18-ap)))
				ep_self = fmt.Sprintf(Ep, fmt.Sprintf("%s%s",
					strings.Repeat("=", e),
					strings.Repeat("-", 18-e)))

			}

			a.SyncAccess(func() {
				if show = c_other != nil; show {
					c := c_other
					name = c.GetName()
					l = c.GetLevel()
					t_hp = c.GetTimedHp()
					m_hp = c.GetMaxHp()
					progress = c.GetSleepProgress()

				}
			})

			bar_other := fmt.Sprintf(Bar)
			name_other := ""
			hp_other := ""
			ap_other := ""

			if show {
				hp = int(float(18*t_hp)/float(m_hp) + 0.99)
				ap = int(18.0*progress + 0.5)

				name_other = fmt.Sprintf(Name, name, l)
				hp_other = fmt.Sprintf(Hp, fmt.Sprintf("%s%s", strings.Repeat("=", hp), strings.Repeat("-", 18-hp)))
				ap_other = fmt.Sprintf(Ap, fmt.Sprintf("%s%s", strings.Repeat("=", ap), strings.Repeat("-", 18-ap)))
			}
	
			curses.Stdwin.AddstrAlign( 0, 1, Whole, curses.A_NORMAL, 
				name_other,
				hp_other,
				ap_other,
				bar_other,

				bar_self,
				name_self,
				hp_self,
				ap_self,
				ep_self,
			)

		}
		curses.UpdatePanels()
		curses.DoUpdate()

	}
}

var (
	s  = 0
	ca = 0
)

func sleep(c Combattant, foo func()) {
	if c.RegenerateAp() {
		ca += 1

		foo()
	} else {
		s += 1
		mainQueue.AddEvent(c.GetApSleep(), func() {
			sleep(c, foo)
		})
	}
}

const (
	SUPER_EFFECTIVE = 2.0
	INEFFECTIVE     = 0.5
)


func main() {

	testmode := true
	Darkness := TestElement("Darkness")
	Light := TestElement("Light")
	Metal := TestElement("Metal")
	Fire := TestElement("Fire")

	Light.SetDamageMod(Metal, INEFFECTIVE)
	Darkness.SetDamageMod(Metal, INEFFECTIVE)
	Darkness.SetDamageMod(Light, SUPER_EFFECTIVE)
	Fire.SetDamageMod(Metal, SUPER_EFFECTIVE)
	//Metal.SetDamageMod(Darkness, SUPER_EFFECTIVE )

	//  80 => +8
	//  90 => +9
	// 100 => +10

	ava := TestCreature("Avatark",
		25, 80, 45, 70,
		12,
		Darkness,
	)

	fla := TestCreature("Flamex",
		170, 5, 45, 0,
		12,
		Fire,
	)

	ste := TestCreature("Steerox",
		65, 40, 80, 35,
		12,
		Metal,
	)

	if testmode {
		ava.SyncAccess(func() {
			name := ava.GetName()
			if "Avatark" != name {
				fmt.Println(name, "!=", "Avatark")
			}
			name = "John"
			if ava.SetName(name); ava.GetName() != name {
				fmt.Println(ava.GetName(), "!=", name)
			}
			name = "Avatark"
			if ava.SetName(name); ava.GetName() != name {
				fmt.Println(ava.GetName(), "!=", name)
			}

			str := ava.GetStrength()
			if ava.SetStrength(str); ava.GetStrength() != str || ava.GetStrength() != 40 {
				fmt.Println(ava.GetStrength(), "!=", str)
			}
		})

	}

	a := MainWorld.GetSingleArena()

	ava_com := TestCombattant(ava)
	fla_com := TestCombattant(fla)
	ste_com := TestCombattant(ste)

	a.SyncAccess(func() {

		if testmode {

			if ava.GetName() != ava_com.GetName() {
				fmt.Println(ava.GetName(), "!=", ava_com.GetName())
			}
			if ste.GetName() != ste_com.GetName() {
				fmt.Println(ste.GetName(), "!=", ste_com.GetName())
			}

			ava_com.SetStrength(45)
			ste_com.SetStrength(65)

			if ava.GetStrength() != 25 {
				fmt.Println(ava.GetStrength(), "!=", 25)
			}
			if ste.GetStrength() != 65 {
				fmt.Println(ste.GetStrength(), "!=", 65)
			}

			if ava_com.GetStrength() != 45 {
				fmt.Println(ava_com.GetStrength(), "!=", 45)
			}
			if ste_com.GetStrength() != 65 {
				fmt.Println(ste_com.GetStrength(), "!=", 65)
			}

			ava_com.SetStrength(25)
			ste_com.SetStrength(65)

		}

		a.SetCombattant(0, fla_com)

		a.SetCombattant(1, ste_com)
		a.SetCombattant(0, ava_com)

		if testmode {

			first, _ := a.GetCombattant(0)
			if first != ava_com {
				fmt.Println(first.GetName(), "!=", ava_com.GetName())
			}
			second, _ := a.GetCombattant(1)
			if second != ste_com {
				fmt.Println(second.GetName(), "!=", ste_com.GetName())
			}

			f_mod := first.GetArenaSpeedMod(a)
			s_mod := second.GetArenaSpeedMod(a)

			if f_mod != 1 {
				fmt.Println(f_mod, "!=", 1)
			}
			if s_mod <= f_mod {
				fmt.Println(s_mod, "<=", f_mod)
			}

		}
	})
	
	if _, err := curses.Initscr(); err != nil {
		panic("Window could not be initialised")
	}
	
	char.Start()
	defer curses.Endwin()
	curses.Stdwin.Addstr(0, 0, "Start?", curses.A_NORMAL )
	
	curses.Noecho()
	curses.Curs_set(curses.CURS_HIDE)
	curses.Stdwin.Keypad(true)
	curses.Stdwin.Refresh()
	char.GetWithTimeout(5e9)

	go func() {
		time_c := time.Tick(1e8)
		for _ = range time_c {
			if char.GetLast() == 'q' {
				break
			}
		}
		mainCIn <- true
		mainQuit <- true
	}()

	c1, _ := a.GetCombattant(0)
	c2, _ := a.GetCombattant(1)

	var foo1, foo2 func()

	p1 := TestSpecPower("Dark Ray", 30, 1.0, 3.0, Darkness)
	p2 := TestPhysPower("Metal Claw", 30, 1.0, 3.0, Metal)
	p3 := TestPhysPower("Flame Fist", 20, 1.0, 2.0, Fire)

	ava.AddPower(p1)
	ste.AddPower(p2)
	fla.AddPower(p3)

	foo1 = func() {
		MainWorld.PauseCombat(true)

		if c1.GetHp() > 0 {
			for p := range c1.GetPowers() {
				if p != nil {
					p.Use()
				}
			}
			if c2.GetHp() > 0 {
				sleep(c1, foo1)
			} else {
				a.SetCombattant(1, nil)
			}
		}

		MainWorld.PauseCombat(false)
	}

	foo2 = func() {
		MainWorld.PauseCombat(true)

		if c2.GetHp() > 0 {
			for p := range c1.GetPowers() {
				if p != nil {
					p.Use()
				}
			}
			if c1.GetHp() > 0 {
				sleep(c2, foo2)
			} else {
				a.SetCombattant(0, nil)
			}
		}

		MainWorld.PauseCombat(false)
	}

	if c1.SetSleepTime(a, 3e9) {
		sleep(c1, foo1)
	}

	if c2.SetSleepTime(a, 3e9) {
		sleep(c2, foo2)
	}

	rand.Seed(time.Nanoseconds())
	go mainloop()
	go viewloop()

	<-mainQuit
}
