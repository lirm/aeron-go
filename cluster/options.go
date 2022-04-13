package cluster

import (
	"time"

	"github.com/lirm/aeron-go/aeron/idlestrategy"
)

type Options struct {
	Timeout          time.Duration      // [runtime] How long to try sending/receiving control messages
	IdleStrategy     idlestrategy.Idler // [runtime] Idlestrategy for sending/receiving control messagesA
	RangeChecking    bool               // [runtime] archive protocol marshalling checks
	LogFragmentLimit int
	ClusterDir       string
}

func NewOptions() *Options {
	o := &Options{
		Timeout:          time.Second,
		IdleStrategy:     idlestrategy.NewDefaultBackoffIdleStrategy(),
		RangeChecking:    true,
		LogFragmentLimit: 50,
		ClusterDir:       "/tmp/aeron-cluster",
	}
	return o
}
