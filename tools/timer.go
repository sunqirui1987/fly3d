package tools

import "time"

type ITimer interface {
	Reset()
	Start()
	Stop()
}

type Timer struct {
	D time.Duration
	F func(args ...interface{})

	Args  interface{}
	max   int
	count int
}

func (t *Timer) Start() {
	t.count = 0
	go timerStart(t)
}

func (t *Timer) Stop() {
	t.count, t.max = 0, 0
}

func (t *Timer) Reset() {
	t.count = 0
}

func timerStart(t *Timer) {
	if t.max == 0 {
		t.F(t.Args)
	} else {
		for {
			if t.count == t.max {
				break
			} else if t.max < 0 || t.count < t.max {
				time.Sleep(t.D)
				t.count++
				t.F(t.Args)
			}

		}

	}
}
