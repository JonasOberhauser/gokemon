package main


type Arena interface {
	GetCombattant(int) (Combattant, bool)
	SetCombattant(int, Combattant) bool
	Iter() <-chan Combattant
	SyncAccess(func())
}

type testarena1 struct {
	coms         []Combattant
	numComs      int
	keyman       KeyMan
}

func TestArena1() Arena {
	a := new(testarena1)
	a.numComs = 2
	a.coms = make([]Combattant, a.numComs)
	a.keyman = TestKeyMan()
	return a
}

func (a *testarena1) SetCombattant(id int, c Combattant) bool {
	if id >= 0 && id <= a.numComs {
                a.coms[id] = c
                return true
	}
	return false
}

func (a *testarena1) GetCombattant(id int) (Combattant, bool) {
	if id >= 0 && id <= a.numComs {
                return a.coms[id], true
        }
        return nil, false
}

func (a *testarena1) iterate( ch chan<- Combattant ) {
        for _, c := range a.coms {
                if c != nil {
                        ch <- c
                }
        }
        close(ch)
}

// Iterator for range clause
func (a *testarena1) Iter() <-chan Combattant {
        ch := make( chan Combattant )
        go a.iterate( ch )
        return ch
}

func (a *testarena1) SyncAccess( foo func() ) {
        a.keyman.Atomic( foo )
}
