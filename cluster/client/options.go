package client

import (
	"time"

	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	RangeChecking      bool
	Loglevel           zapcore.Level // [runtime] via logging.SetLevel()
	IngressEndpoints   string
	IngressChannel     string
	IngressStreamId    int32
	EgressChannel      string
	EgressStreamId     int32
	IdleStrategy       idlestrategy.Idler
	IsIngressExclusive bool
	RetryBackoff       time.Duration // how long to wait before re-attempting to connect
}

func NewOptions() *Options {
	return &Options{
		RangeChecking:   true,
		Loglevel:        zapcore.WarnLevel,
		IngressStreamId: 101,
		EgressChannel:   "aeron:udp?alias=cluster-egress|endpoint=localhost:0",
		EgressStreamId:  102,
		IdleStrategy:    idlestrategy.NewDefaultBackoffIdleStrategy(),
		RetryBackoff:    5 * time.Second,
	}
}
