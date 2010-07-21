package char

import (
	"curses"
	"time"
)

var win *curses.Window
var char int
var kc chan int

func GetLast() int {
	return char
}
func Get() int {
	return <-kc
}


func GetWithTimeout(timeout int64) int {
	t := time.NewTicker(timeout)
	defer t.Stop()
	select {
	case k := <-kc:
		return k
	case <-t.C:
	}
	return -1

}

func Print(y, x int, s string, v ...interface{}) { win.Addstr(x, y, s, 0, v) }
func Flush()                                     { win.Refresh() }

func Start() {
	kc = make(chan int)
	var err interface{}
	if win, err = curses.Initscr(); err == nil {
		curses.Noecho()
		curses.Curs_set(curses.CURS_HIDE)
		win.Keypad(true)
		go func(){ 
			for {
				char = win.Getch()
				kc <- char
				
			}
			
		}()

	}
	
}

func End() {
	curses.Endwin()
}
