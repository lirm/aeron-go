package cluster

import (
	"fmt"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/cluster/codecs"
)

const SnapshotTypeId = 2

type SnapshotTaker struct {
	marshaller  *codecs.SbeGoMarshaller // currently shared as we're not reentrant (but could be here)
	options     *Options
	publication *aeron.Publication
}

func NewSnapshotTaker(
	options *Options,
	publication *aeron.Publication,
) *SnapshotTaker {
	return &SnapshotTaker{
		marshaller:  codecs.NewSbeGoMarshaller(),
		options:     options,
		publication: publication,
	}
}

func (st *SnapshotTaker) MarkBegin(
	logPosition int64,
	leadershipTermId int64,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) error {
	return st.markSnapshot(logPosition, leadershipTermId, codecs.SnapshotMark.BEGIN, timeUnit, appVersion)
}

func (st *SnapshotTaker) MarkEnd(
	logPosition int64,
	leadershipTermId int64,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) error {
	return st.markSnapshot(logPosition, leadershipTermId, codecs.SnapshotMark.END, timeUnit, appVersion)
}

func (st *SnapshotTaker) markSnapshot(
	logPosition int64,
	leadershipTermId int64,
	mark codecs.SnapshotMarkEnum,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) error {
	bytes, err := codecs.SnapshotMarkerPacket(
		st.marshaller,
		st.options.RangeChecking,
		SnapshotTypeId,
		logPosition,
		leadershipTermId,
		0,
		mark,
		timeUnit,
		appVersion,
	)
	if err != nil {
		return err
	}
	if ret := st.offer(bytes); ret < 0 {
		return fmt.Errorf("SnapshotTaker.offer failed: %d", ret)
	}
	return nil
}

func (st *SnapshotTaker) SnapshotSession(session ClientSession) error {
	bytes, err := codecs.ClientSessionPacket(st.marshaller, st.options.RangeChecking,
		session.Id(), session.ResponseStreamId(), []byte(session.ResponseChannel()), []byte(""))
	if err != nil {
		return err
	}
	if ret := st.offer(bytes); ret < 0 {
		return fmt.Errorf("SnapshotTaker.offer failed: %d", ret)
	}
	return nil
}

// Offer to our request publication
func (st *SnapshotTaker) offer(bytes []byte) int64 {
	buffer := atomic.MakeBuffer(bytes)
	length := int32(len(bytes))
	start := time.Now()
	var ret int64
	for time.Since(start) < st.options.Timeout {
		ret = st.publication.Offer(buffer, 0, length, nil)
		switch ret {
		// Retry on these
		case aeron.NotConnected, aeron.BackPressured, aeron.AdminAction:
			st.options.IdleStrategy.Idle(0)
		// Fail or succeed on other values
		default:
			return ret
		}
	}
	// Give up, returning the last failure
	return ret
}
