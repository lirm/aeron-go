package cluster

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type Cluster interface {
	/**
	 * Position the log has reached in bytes as of the current message.
	 *
	 * @return position the log has reached in bytes as of the current message.
	 */
	LogPosition() int64

	/**
	 * The unique id for the hosting member of the cluster. Useful only for debugging purposes.
	 *
	 * @return unique id for the hosting member of the cluster.
	 */
	MemberId() int32

	/**
	 * The role the cluster node is playing.
	 *
	 * @return the role the cluster node is playing.
	 */
	Role() Role

	/**
	 * Cluster time as {@link #timeUnit()}s since 1 Jan 1970 UTC.
	 *
	 * @return time as {@link #timeUnit()}s since 1 Jan 1970 UTC.
	 */
	Time() int64

	/**
	 * The unit of time applied when timestamping and {@link #time()} operations.
	 *
	 * @return the unit of time applied when timestamping and {@link #time()} operations.
	 */
	TimeUnit() codecs.ClusterTimeUnitEnum

	/**
	 * {@link IdleStrategy} which should be used by the service when it experiences back-pressure on egress,
	 * closing sessions, making timer requests, or any long-running actions.
	 *
	 * @return the {@link IdleStrategy} which should be used by the service when it experiences back-pressure.
	 */
	IdleStrategy() idlestrategy.Idler

	/**
	 * Schedule a timer for a given deadline and provide a correlation id to identify the timer when it expires or
	 * for cancellation. This action is asynchronous and will race with the timer expiring.
	 * <p>
	 * If the correlationId is for an existing scheduled timer then it will be rescheduled to the new deadline. However,
	 * it is best to generate correl~~ationIds in a monotonic fashion and be aware of potential clashes with other
	 * services in the same cluster. Service isolation can be achieved by using the upper bits for service id.
	 * <p>
	 * Timers should only be scheduled or cancelled in the context of processing a
	 * {@link ClusteredService#onSessionMessage(ClientSession, long, DirectBuffer, int, int, Header)},
	 * {@link ClusteredService#onTimerEvent(long, long)},
	 * {@link ClusteredService#onSessionOpen(ClientSession, long)}, or
	 * {@link ClusteredService#onSessionClose(ClientSession, long, CloseReason)}.
	 * If applied to other events then they are not guaranteed to be reliable.
	 * <p>
	 * Callers of this method should loop until the method succeeds.
	 *
	 * <pre>{@code
	 * private Cluster cluster;
	 * // Lines omitted...
	 *
	 * cluster.idleStrategy().reset();
	 * while (!cluster.scheduleTimer(correlationId, deadline))
	 * {
	 *     cluster.idleStrategy().idle();
	 * }
	 * }</pre>
	 *
	 * The cluster's idle strategy must be used in the body of the loop to allow for the clustered service to be
	 * shutdown if required.
	 *
	 * @param correlationId to identify the timer when it expires. {@link Long#MAX_VALUE} not supported.
	 * @param deadline      time after which the timer will fire. {@link Long#MAX_VALUE} not supported.
	 * @return true if the event to schedule a timer request has been sent or false if back-pressure is applied.
	 * @see #cancelTimer(long)
	 */
	ScheduleTimer(correlationId int64, deadline int64) bool

	/**
	 * Cancel a previously scheduled timer. This action is asynchronous and will race with the timer expiring.
	 * <p>
	 * Timers should only be scheduled or cancelled in the context of processing a
	 * {@link ClusteredService#onSessionMessage(ClientSession, long, DirectBuffer, int, int, Header)},
	 * {@link ClusteredService#onTimerEvent(long, long)},
	 * {@link ClusteredService#onSessionOpen(ClientSession, long)}, or
	 * {@link ClusteredService#onSessionClose(ClientSession, long, CloseReason)}.
	 * If applied to other events then they are not guaranteed to be reliable.
	 * <p>
	 * Callers of this method should loop until the method succeeds, see {@link
	 * io.aeron.cluster.service.Cluster#scheduleTimer(long, long)} for an example.
	 *
	 * @param correlationId for the timer provided when it was scheduled. {@link Long#MAX_VALUE} not supported.
	 * @return true if the event to cancel request has been sent or false if back-pressure is applied.
	 * @see #scheduleTimer(long, long)
	 */
	CancelTimer(correlationId int64) bool

	/**
	 * Offer a message as ingress to the cluster for sequencing. This will happen efficiently over IPC to the
	 * consensus module and have the cluster session of as the negative value of the
	 * {@link io.aeron.cluster.service.ClusteredServiceContainer.Configuration#SERVICE_ID_PROP_NAME}.
	 * <p>
	 * Callers of this method should loop until the method succeeds.
	 *
	 * <pre>{@code
	 * private Cluster cluster;
	 * // Lines omitted...
	 *
	 * cluster.idleStrategy().reset();
	 * do
	 * {
	 *     final long position = cluster.offer(buffer, offset, length);
	 *     if (position > 0)
	 *     {
	 *         break;
	 *     }
	 *     else if (Publication.ADMIN_ACTION != position || Publication.BACK_PRESSURED != position)
	 *     {
	 *         throw new ClusterException("Internal offer failed: " + position);
	 *     }
	 *
	 *     cluster.idleStrategy.idle();
	 * }
	 * while (true);
	 * }</pre>
	 *
	 * The cluster's idle strategy must be used in the body of the loop to allow for the clustered service to be
	 * shutdown if required.
	 *
	 * @param buffer containing the message to be offered.
	 * @param offset in the buffer at which the encoded message begins.
	 * @param length in the buffer of the encoded message.
	 * @return positive value if successful.
	 * @see io.aeron.Publication#offer(DirectBuffer, int, int)
	 */
	Offer(*atomic.Buffer, int32, int32) int64
}
