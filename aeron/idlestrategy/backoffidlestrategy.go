package idlestrategy

import (
	"fmt"
	"runtime"
	"time"
)

// BackoffIdleStrategy is an idling strategy for threads when they have no work to do.
// Spin for maxSpins, then yield with runtime.Gosched for maxYields,
// then park with time.Sleep on an exponential backoff to maxParkPeriodNs.
type BackoffIdleStrategy struct {
	// configured max spins, yield, and min / max park period
	maxSpins, maxYields, minParkPeriodNs, maxParkPeriodNs int64
	// current state
	state backoffState
	// current number of spins, yield, and park period.
	spins, yields, parkPeriodNs int64
}

// NewBackoffIdleStrategy returns a BackoffIdleStrategy with the given parameters.
func NewBackoffIdleStrategy(maxSpins, maxYields, minParkPeriodNs, maxParkPeriodNs int64) *BackoffIdleStrategy {
	return &BackoffIdleStrategy{
		maxSpins:        maxSpins,
		maxYields:       maxYields,
		minParkPeriodNs: minParkPeriodNs,
		maxParkPeriodNs: maxParkPeriodNs,
	}
}

// DefaultMaxSpins is the default number of times the strategy will spin without work before going to next state.
const DefaultMaxSpins = 10

// DefaultMaxYields is the default number of times the strategy will yield without work before going to next state.
const DefaultMaxYields = 20

// DefaultMinParkNs is the default interval the strategy will park the thread on entering the park state.
const DefaultMinParkNs = 1000

// DefaultMaxParkNs is the default interval the strategy will park the thread will expand interval to as a max.
const DefaultMaxParkNs = int64(1 * time.Millisecond)

// NewDefaultBackoffIdleStrategy returns a BackoffIdleStrategy using DefaultMaxSpins, DefaultMaxYields,
// DefaultMinParkNs, and DefaultMaxParkNs.
func NewDefaultBackoffIdleStrategy() *BackoffIdleStrategy {
	return NewBackoffIdleStrategy(DefaultMaxSpins, DefaultMaxYields, DefaultMinParkNs, DefaultMaxParkNs)
}

type backoffState int8

const (
	// Denotes a non-idle state.
	backoffNotIdle backoffState = iota
	// Denotes a spinning state.
	backoffSpinning
	// Denotes an yielding state.
	backoffYielding
	// Denotes a parking state.
	backoffParking
)

func (s *BackoffIdleStrategy) Idle(workCount int) {
	if workCount > 0 {
		s.reset()
	} else {
		s.idle()
	}
}

func (s *BackoffIdleStrategy) String() string {
	return fmt.Sprintf("BackoffIdleStrategy(MaxSpins:%d, MaxYields:%d, MinParkPeriodNs:%d, MaxParkPeriodNs:%d)",
		s.maxSpins, s.maxYields, s.minParkPeriodNs, s.maxParkPeriodNs)
}

func (s *BackoffIdleStrategy) reset() {
	s.spins = 0
	s.yields = 0
	s.parkPeriodNs = s.minParkPeriodNs
	s.state = backoffNotIdle
}

func (s *BackoffIdleStrategy) idle() {
	switch s.state {
	case backoffNotIdle:
		s.state = backoffSpinning
		s.spins++
	case backoffSpinning:
		// We should call procyield here, see https://golang.org/src/runtime/lock_futex.go
		s.spins++
		if s.spins > s.maxSpins {
			s.state = backoffYielding
			s.yields = 0
		}
	case backoffYielding:
		s.yields++
		if s.yields > s.maxYields {
			s.state = backoffParking
			s.parkPeriodNs = s.minParkPeriodNs
		} else {
			runtime.Gosched()
		}
	case backoffParking:
		time.Sleep(time.Nanosecond * time.Duration(s.parkPeriodNs))
		s.parkPeriodNs = s.parkPeriodNs << 1
		if s.parkPeriodNs > s.maxParkPeriodNs {
			s.parkPeriodNs = s.maxParkPeriodNs
		}
	}
}
