package cluster

import (
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/flyweight"
)

type MarkFileHeaderFlyweight struct {
	flyweight.FWBase

	Version           flyweight.Int32Field
	ComponentType     flyweight.Int32Field
	ActivityTimestamp flyweight.Int64Field
	StartTimestamp    flyweight.Int64Field

	Pid                     flyweight.Int64Field
	CandidateTermId         flyweight.Int64Field
	ArchiveStreamId         flyweight.Int32Field
	ServiceStreamId         flyweight.Int32Field
	ConsensusModuleStreamId flyweight.Int32Field
	IngressStreamId         flyweight.Int32Field
	MemberId                flyweight.Int32Field
	ServiceId               flyweight.Int32Field

	HeaderLength      flyweight.Int32Field
	ErrorBufferLength flyweight.Int32Field
	ClusterId         flyweight.Int32Field
}

func (f *MarkFileHeaderFlyweight) Wrap(
	buf *atomic.Buffer,
	offset int,
) flyweight.Flyweight {
	pos := offset
	pos += f.Version.Wrap(buf, pos)
	pos += f.ComponentType.Wrap(buf, pos)
	pos += f.ActivityTimestamp.Wrap(buf, pos)
	pos += f.StartTimestamp.Wrap(buf, pos)

	pos += f.Pid.Wrap(buf, pos)
	pos += f.CandidateTermId.Wrap(buf, pos)
	pos += f.ArchiveStreamId.Wrap(buf, pos)
	pos += f.ServiceStreamId.Wrap(buf, pos)
	pos += f.ConsensusModuleStreamId.Wrap(buf, pos)
	pos += f.IngressStreamId.Wrap(buf, pos)
	pos += f.MemberId.Wrap(buf, pos)
	pos += f.ServiceId.Wrap(buf, pos)

	pos += f.HeaderLength.Wrap(buf, pos)
	pos += f.ErrorBufferLength.Wrap(buf, pos)
	pos += f.ClusterId.Wrap(buf, pos)

	f.SetSize(HeaderLength + ErrorBufferLength)
	return f
}
