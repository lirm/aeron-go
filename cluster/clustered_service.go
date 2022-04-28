package cluster

import (
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type ClusteredService interface {
	OnStart(cluster Cluster, image *aeron.Image)

	OnSessionOpen(session ClientSession, timestamp int64)

	OnSessionClose(
		session ClientSession,
		timestamp int64,
		closeReason codecs.CloseReasonEnum,
	)

	OnSessionMessage(
		session ClientSession,
		timestamp int64,
		buffer *atomic.Buffer,
		offset int32,
		length int32,
		header *logbuffer.Header,
	)

	OnTimerEvent(correlationId, timestamp int64)

	OnTakeSnapshot(publication *aeron.Publication)

	OnRoleChange(role Role)

	OnTerminate(cluster Cluster)

	OnNewLeadershipTermEvent(
		leadershipTermId int64,
		logPosition int64,
		timestamp int64,
		termBaseLogPosition int64,
		leaderMemberId int32,
		logSessionId int32,
		timeUnit codecs.ClusterTimeUnitEnum,
		appVersion int32,
	)
}
