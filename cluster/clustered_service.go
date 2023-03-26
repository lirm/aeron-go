package cluster

import (
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type ClusteredService interface {
	/**
	 * Start event for the service where the service can perform any initialisation required and load snapshot state.
	 * The snapshot image can be null if no previous snapshot exists.
	 * <p>
	 * <b>Note:</b> As this is a potentially long-running operation the implementation should use
	 * {@link Cluster#idleStrategy()} and then occasionally call {@link org.agrona.concurrent.IdleStrategy#idle()} or
	 * {@link org.agrona.concurrent.IdleStrategy#idle(int)}, especially when polling the {@link Image} returns 0.
	 *
	 * @param cluster       with which the service can interact.
	 * @param snapshotImage from which the service can load its archived state which can be null when no snapshot.
	 */
	OnStart(cluster Cluster, image aeron.Image)

	/**
	 * A session has been opened for a client to the cluster.
	 *
	 * @param session   for the client which have been opened.
	 * @param timestamp at which the session was opened.
	 */
	OnSessionOpen(session ClientSession, timestamp int64)

	/**
	 * A session has been closed for a client to the cluster.
	 *
	 * @param session     that has been closed.
	 * @param timestamp   at which the session was closed.
	 * @param closeReason the session was closed.
	 */
	OnSessionClose(
		session ClientSession,
		timestamp int64,
		closeReason codecs.CloseReasonEnum,
	)

	/**
	 * A message has been received to be processed by a clustered service.
	 *
	 * @param session   for the client which sent the message. This can be null if the client was a service.
	 * @param timestamp for when the message was received.
	 * @param buffer    containing the message.
	 * @param offset    in the buffer at which the message is encoded.
	 * @param length    of the encoded message.
	 * @param header    aeron header for the incoming message.
	 */
	OnSessionMessage(
		session ClientSession,
		timestamp int64,
		buffer *atomic.Buffer,
		offset int32,
		length int32,
		header *logbuffer.Header,
	)

	/**
	 * A scheduled timer has expired.
	 *
	 * @param correlationId for the expired timer.
	 * @param timestamp     at which the timer expired.
	 */
	OnTimerEvent(correlationId, timestamp int64)

	/**
	 * The service should take a snapshot and store its state to the provided archive {@link Publication}.
	 * <p>
	 * <b>Note:</b> As this is a potentially long-running operation the implementation should use
	 * {@link Cluster#idleStrategy()} and then occasionally call {@link org.agrona.concurrent.IdleStrategy#idle()} or
	 * {@link org.agrona.concurrent.IdleStrategy#idle(int)},
	 * especially when the snapshot {@link ExclusivePublication} returns {@link Publication#BACK_PRESSURED}.
	 *
	 * @param publication to which the state should be recorded.
	 */
	OnTakeSnapshot(publication *aeron.Publication)

	/**
	 * Notify that the cluster node has changed role.
	 *
	 * @param role that the node has assumed.
	 */
	OnRoleChange(role Role)

	/**
	 * Called when the container is going to terminate.
	 *
	 * @param cluster with which the service can interact.
	 */
	OnTerminate(cluster Cluster)

	/**
	 * An election has been successful and a leader has entered a new term.
	 *
	 * @param leadershipTermId    identity for the new leadership term.
	 * @param logPosition         position the log has reached as the result of this message.
	 * @param timestamp           for the new leadership term.
	 * @param termBaseLogPosition position at the beginning of the leadership term.
	 * @param leaderMemberId      who won the election.
	 * @param logSessionId        session id for the publication of the log.
	 * @param timeUnit            for the timestamps in the coming leadership term.
	 * @param appVersion          for the application configured in the consensus module.
	 */
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
