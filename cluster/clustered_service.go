package cluster

import (
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type ClusteredService interface {
	// OnStart is called to initialize the service and load snapshot state, where the snapshot image can be nil if no previous snapshot exists.
	//
	// Note: As this can potentially be a long-running operation, the implementation should use Cluster.IdleStrategy() and
	// occasionally call IdleStrategy.Idle(), especially when polling the Image returns 0
	OnStart(cluster Cluster, image aeron.Image)

	// OnSessionOpen notifies the clustered service that a session has been opened for a client to the cluster
	OnSessionOpen(session ClientSession, timestamp int64)

	// OnSessionClose notifies the clustered service that a session has been closed for a client to the cluster
	OnSessionClose(
		session ClientSession,
		timestamp int64,
		closeReason codecs.CloseReasonEnum,
	)

	// OnSessionMessage notifies the clustered service that a message has been received to be processed by a clustered service
	OnSessionMessage(
		session ClientSession,
		timestamp int64,
		buffer *atomic.Buffer,
		offset int32,
		length int32,
		header *logbuffer.Header,
	)

	// OnTimerEvent notifies the clustered service that a scheduled timer has expired
	OnTimerEvent(correlationId, timestamp int64)

	// OnTakeSnapshot instructs the clustered service to take a snapshot and store its state to the provided aeron archive Publication.
	//
	// Note: As this is a potentially long-running operation the implementation should use
	// Cluster.idleStrategy() and then occasionally call IdleStrategy.idle()
	// especially when the snapshot ExclusivePublication returns Publication.BACK_PRESSURED
	OnTakeSnapshot(publication *aeron.Publication)

	// OnRoleChange notifies the clustered service that the cluster node has changed role
	OnRoleChange(role Role)

	// OnTerminate notifies the clustered service that the container is going to terminate
	OnTerminate(cluster Cluster)

	// OnNewLeadershipTermEvent notifies the clustered service that an election has been successful and a leader has entered a new term
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
