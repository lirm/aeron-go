package cluster

import (
	"github.com/corymonroe-coinbase/aeron-go/aeron/util"
	"os"
	"time"

	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	// "github.com/lirm/aeron-go/aeron/util"
	"github.com/corymonroe-coinbase/aeron-go/aeron/util/memmap"
)

const (
	HeaderLength      = 8 * 1024
	ErrorBufferLength = 1024 * 1024
)

type ClusterMarkFile struct {
	file      *memmap.File
	buffer    *atomic.Buffer
	flyweight *MarkFileHeaderFlyweight
}

func NewClusterMarkFile(filename string) (*ClusterMarkFile, error) {
	f, err := memmap.NewFile(filename, 0, HeaderLength+ErrorBufferLength)
	if err != nil {
		return nil, err
	}

	b := atomic.MakeBuffer(f.GetMemoryPtr(), f.GetMemorySize())

	fly := &MarkFileHeaderFlyweight{}
	fly.Wrap(b, 0)

	fly.CandidateTermId.Set(-1)
	fly.ComponentType.Set(2)
	fly.HeaderLength.Set(HeaderLength)
	fly.ErrorBufferLength.Set(ErrorBufferLength)
	fly.Pid.Set(int64(os.Getpid()))
	fly.StartTimestamp.Set(time.Now().UnixMilli())

	return &ClusterMarkFile{
		file:      f,
		buffer:    b,
		flyweight: fly,
	}, nil
}

func (cmf *ClusterMarkFile) SignalReady() {
	semanticVersion := util.SemanticVersionCompose(0, 3, 0)
	//cmf.flyweight.Version.Set(int32(semanticVersion))
	cmf.buffer.PutInt32Ordered(0, int32(semanticVersion))
}

func (cmf *ClusterMarkFile) SignalFailedStart() {
	cmf.flyweight.Version.Set(-1)
}

func (cmf *ClusterMarkFile) UpdateActivityTimestamp(timestamp int64) {
	cmf.flyweight.ActivityTimestamp.Set(timestamp)
}
