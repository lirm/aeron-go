package cluster

import (
	"os"
	"time"

	"github.com/lirm/aeron-go/aeron/atomic"
	// "github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
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
	fly.StartTimestamp.Set(time.Now().Unix())

	return &ClusterMarkFile{
		file:      f,
		buffer:    b,
		flyweight: fly,
	}, nil
}

func (cmf *ClusterMarkFile) SignalReady() {
	// util.SemanticVersionCompose()
}
