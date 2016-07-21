package idlestrategy

import "time"

type Idler interface {
	Idle(fragmentsRead int)
}

type Busy struct {
}

func (s Busy) Idle(fragmentsRead int) {

}

type Sleeping struct {
	SleepFor time.Duration
}

func (s Sleeping) Idle(fragmentsRead int) {
	if fragmentsRead == 0 {
		time.Sleep(s.SleepFor)
	}
}
