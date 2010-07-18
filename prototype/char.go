package main

import "curses"

var id = 0
var char int
var inch chan bool

func GetLast() int {
        return char
}
func Get() int {
        <- inch
        return GetLast()
}

func init() {
        inch = make( chan bool )
        go func() {
                for{    
                        char = curses.Stdwin.Getch()
                        id += 1
                        inch <- true
                }
        }()
}
