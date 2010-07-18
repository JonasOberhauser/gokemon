package main

import (
        "fmt"
        )

type Event interface {
	Run()
	GetId() int64
}

type EventQueue interface {
	AddEvent(int64, func()) Event
	DropEvent(Event)
	Update(int64)
}

var (
        ekeyman = TestKeyMan()
        eeventid = int64( 0 )
        )
        
type event struct {
	foo func()
	id int64
}

func (e *event) Run() {
	e.foo()
}


func testEvent(foo func()) Event {
	e := new(event)
	e.foo = foo
	ekeyman.Atomic( func(){
	eeventid += 1 
	e.id = eeventid
        })
	return e
}

func (e *event) GetId() int64 {
        return e.id
}

type testqueue struct {
	heap *testheap
	k    KeyMan
        time int64
}

func TestQueue() EventQueue {
        q := new(testqueue)
        q.k = TestKeyMan()
        q.time = 0
        return q
}

func (q *testqueue) AddEvent( time int64, foo func() ) Event {
        e:=testEvent(foo)
        q.k.Atomic( func(){ 
                time = q.time + time
                if q.heap == nil {
                        q.heap = testHeap( time, e, nil )
                } else {
                        q.heap.add( time, e )
                } 
        })
        
        return e
}

func (q *testqueue) DropEvent( e Event ) {
        loadByE( e ).delete()
}

func (q *testqueue) Update( time int64 ) {
        q.k.Atomic( func(){
                for q.heap != nil && q.heap.valid && q.heap.time <= q.time + time {
                        e := q.heap.e
                        q.heap.delete()
                        go e.Run()
                }
                q.time += time
        })
}

var heapByEvent = make( map[int64] *testheap )

type testheap struct {
	time              int64
	e                 Event
	parent            *testheap
	h                 []*testheap
        nextadd           int
	valid             bool
}

func testHeap(t int64, e Event, parent *testheap ) *testheap {
        
	h := new(testheap)
	h.init( t, e )
	h.parent = parent
	h.nextadd = 0
	h.h = make([]*testheap,2)
	
	return h
}

func (h *testheap) bubble() {
	if h.parent == nil || h.parent.time < h.time {
		return
	}
	h.swap(h.parent)
	h.parent.bubble()
}

func (h *testheap) swap(o *testheap) {
	h.time, o.time = o.time, h.time
	h.e, o.e = o.e, h.e
        o.save_by_e()
        h.save_by_e()
}

func (h *testheap) sink() {
	for _, son := range h.h {
		if son != nil && son.valid && son.time < h.time {
			h.swap(son)
			son.sink()
			return
		}
	}
}

func (h *testheap) find() {
	h.bubble()
	h.sink()
}

func (h *testheap) save_by_e() {
        heapByEvent[ h.e.GetId() ] = h
}
func loadByE(e Event) *testheap {
        return heapByEvent[ e.GetId() ]
}

func (h *testheap) init( t int64, e Event ) {
        h.time = t
        h.e = e
        h.save_by_e()
        h.valid = true
}
                
func (h *testheap) add(t int64, e Event) {
        if h.valid {
                for i, son := range h.h {
                        
                        if son == nil {
                                nh := testHeap(t, e, h)
                                h.h[i] = nh
                                nh.bubble()
                                return
                        } else if !son.valid {
                                son.init( t, e ) 
                                son.bubble()
                                return
                        }
                }
                h.h[h.nextadd].add(t, e)
                
                h.nextadd += 1
                h.nextadd %= 2
        } else {
                h.init( t, e )
        }
}

// Returns a testheap without children.
func (h *testheap) getlowest() *testheap {
	for _, son := range h.h {
		if son != nil && son.valid {
			return son.getlowest()
		}
	}
	return h
}

func (h *testheap) delete() {
	lowest := h.getlowest()
	
	heapByEvent[ h.e.GetId() ] = nil, false
	lowest.valid = false
	
	lowest.swap(h)
	
	h.find()
}

func (h *testheap) String() string {
        var s string
        switch h.valid {
        case true: s = fmt.Sprint( h.time )
        case false: s = fmt.Sprint( "---" )
        }
                
        for _, son := range h.h {
                if son != nil {
                        s += fmt.Sprint( ", ", son )
                }
        }
        return "{" + s + "}"
}


func init() {
        go func() {
                testmode := true
                
                if testmode {
                        
                        e := make( []Event, 7 )
                        var h *testheap 
                        for i, _ := range e {
                                e[i] = testEvent( func(){} )
                        }
                        h = testHeap( 0, e[0], nil)
                        for i, _ := range e {
                                if i == 0 {
                                } else {
                                        h.add( int64(i), e[i] )
                                }
                        }
                        
                        if fmt.Sprint( h ) != "{0, {1, {3}, {5}}, {2, {4}, {6}}}" {
                                        fmt.Println( fmt.Sprint( h ), "!=", "{0, {1, {3}, {5}}, {2, {4}}}" )
                        }
                        h.delete()
                        if fmt.Sprint( h ) != "{1, {3, {---}, {5}}, {2, {4}, {6}}}" {
                                        fmt.Println( fmt.Sprint( h ), "!=", "{1, {3, {---}, {5}}, {2, {4}}}"  )
                        }
                        h.add( 0, e[0] )
                        if fmt.Sprint( h ) != "{0, {1, {3}, {5}}, {2, {4}, {6}}}" {
                                        fmt.Println( fmt.Sprint( h ), "!=", "{0, {1, {3}, {5}}, {2, {4}}}"  )
                        }
                        
                        
                }
        }()
}

