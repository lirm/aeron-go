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
	// occasionally call IdleStrategy.Idle() or IdleStrategy.Idle(int), especially when polling the Image returns 0.
	//
	// cluster the Cluster with which the service can interact.
	// snapshotImage the Image from which the service can load its archived state, which can be nil when there is no snapshot.
	OnStart(cluster Cluster, image aeron.Image)

	// OnSessionOpen notifies the clustered service that a session has been opened for a client to the cluster.
	//
	// session   for the client which have been opened.
	// timestamp at which the session was opened.
	OnSessionOpen(session ClientSession, timestamp int64)

	// OnSessionClose notifies the clustered service that a session has been closed for a client to the cluster.
	//
	// session     that has been closed.
	// timestamp   at which the session was closed.
	// closeReason the session was closed.
	OnSessionClose(
		session ClientSession,
		timestamp int64,
		closeReason codecs.CloseReasonEnum,
	)

	// OnSessionMessage notifies the clustered service that a message has been received to be processed by a clustered service.
	//
	// session   for the client which sent the message. This can be null if the client was a service.
	// timestamp for when the message was received.
	// buffer    containing the message.
	// offset    in the buffer at which the message is encoded.
	// length    of the encoded message.
	// header    aeron header for the incoming message.
	OnSessionMessage(
		session ClientSession,
		timestamp int64,
		buffer *atomic.Buffer,
		offset int32,
		length int32,
		header *logbuffer.Header,
	)

	// OnTimerEvent notifies the clustered service that a scheduled timer has expired.
	//
	// correlationId for the expired timer.
	// timestamp     at which the timer expired.
	OnTimerEvent(correlationId, timestamp int64)

	// OnTakeSnapshot instructs the clustered service to take a snapshot and store its state to the provided aeron archive Publication.
	//
	// Note: As this is a potentially long-running operation the implementation should use
	// Cluster#idleStrategy() and then occasionally call IdleStrategy#idle()
	// especially when the snapshot ExclusivePublication returns Publication#BACK_PRESSURED.
	//
	// publication to which the state should be recorded.
	OnTakeSnapshot(publication *aeron.Publication)

	// OnRoleChange notifies the clustered service that the cluster node has changed role.
	//
	// role that the node has assumed.
	OnRoleChange(role Role)

	// OnTerminate notifies the clustered service that the container is going to terminate.
	//
	// cluster with which the service can interact.
	OnTerminate(cluster Cluster)

	// OnNewLeadershipTermEvent notifies the clustered service that an election has been successful and a leader has entered a new term.
	//
	// leadershipTermId    identity for the new leadership term.
	// logPosition         position the log has reached as the result of this message.
	// timestamp           for the new leadership term.
	// termBaseLogPosition position at the beginning of the leadership term.
	// leaderMemberId      who won the election.
	// logSessionId        session id for the publication of the log.
	// timeUnit            for the timestamps in the coming leadership term.
	// appVersion          for the application configured in the consensus module.
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
