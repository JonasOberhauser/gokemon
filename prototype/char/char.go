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


func Start() {
	kc = make(chan int)
	win = curses.Stdwin
	
	go func(){ 
		for {
			char = win.Getch()
			kc <- char
			
		}
		
	}()

	
}