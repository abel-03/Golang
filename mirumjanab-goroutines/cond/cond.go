package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
type Cond struct {
	L           Locker
	queue       []int
	signalCount chan int
}

// New returns a new Cond with Locker l.
func New(locker Locker) *Cond {
	return &Cond{
		L:    locker,
		queue:   make([]int, 0),
		signalCount: make(chan int, 1),
	}
}


// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
func (c *Cond) Wait() {
	n := 0
	if len(c.queue) != 0 {
		n = c.queue[len(c.queue)-1] + 1
	}
	c.queue = append(c.queue, n)
	c.L.Unlock()
	
	for {
		count := <-c.signalCount
		if count > 0 && c.queue[0] == n {
			c.L.Lock()
			c.queue = c.queue[1:]
			c.L.Unlock()
			if count != 1 {
				c.signalCount <- count - 1
			}
			break
		}
		c.signalCount <- count
	}
	defer c.L.Lock()

	
	
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal() {
	c.signalCount <- 1
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {
	c.signalCount <- len(c.queue)
}
