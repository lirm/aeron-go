package cluster

import (
	"time"

	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/archive"
)

type Options struct {
	Timeout          time.Duration      // [runtime] How long to try sending/receiving control messages
	IdleStrategy     idlestrategy.Idler // [runtime] Idlestrategy for sending/receiving control messagesA
	RangeChecking    bool               // [runtime] archive protocol marshalling checks
	LogFragmentLimit int
	ClusterDir       string
	ClusterId        int32
	AppVersion       int32
	ArchiveOptions   *archive.Options
}

func NewOptions() *Options {
	archiveOpts := archive.DefaultOptions()
	archiveOpts.RequestChannel = "aeron:ipc?alias=cluster-service-archive-ctrl-req|term-length=128k"
	archiveOpts.ResponseChannel = "aeron:ipc?alias=cluster-service-archive-ctrl-resp|term-length=128k"
	return &Options{
		Timeout:          time.Second * 5,
		IdleStrategy:     idlestrategy.NewDefaultBackoffIdleStrategy(),
		RangeChecking:    true,
		LogFragmentLimit: 50,
		ClusterDir:       "/tmp/aeron-cluster",
		ArchiveOptions:   archiveOpts,
	}
}
