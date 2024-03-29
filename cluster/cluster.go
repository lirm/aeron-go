package cluster

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type Cluster interface {

	// LogPosition returns the position the log has reached in bytes as of the current message
	LogPosition() int64

	// MemberId returns the unique id for the hosting member of the cluster. Useful only for debugging purposes
	MemberId() int32

	// Role returns the role the cluster node is playing
	Role() Role

	// Time returns the cluster time as time units since 1 Jan 1970 UTC
	Time() int64

	// TimeUnit returns the unit of time applied when timestamping and time operations
	TimeUnit() codecs.ClusterTimeUnitEnum

	// IdleStrategy returns the IdleStrategy which should be used by the service when it experiences back-pressure on egress,
	// closing sessions, making timer requests, or any long-running actions
	IdleStrategy() idlestrategy.Idler

	// ScheduleTimer schedules a timer for a given deadline and provide a correlation id to identify the timer when it expires or
	// for cancellation. This action is asynchronous and will race with the timer expiring.
	//
	// If the correlationId is for an existing scheduled timer then it will be rescheduled to the new deadline. However,
	// it is best to generate correllationIds in a monotonic fashion and be aware of potential clashes with other
	// services in the same cluster. Service isolation can be achieved by using the upper bits for service id.
	//
	// Timers should only be scheduled or cancelled in the context of processing
	// - onSessionMessage
	// - onTimerEvent
	// - onSessionOpen
	// - onSessionClose
	// If applied to other events then they are not guaranteed to be reliable.
	//
	// Callers of this method should loop until the method succeeds.
	//
	// The cluster's idle strategy must be used in the body of the loop to allow for the clustered service to be
	// shutdown if required.
	//
	// ScheduleTimer returns true if the event to schedule a timer request has been sent or false if back-pressure is applied
	ScheduleTimer(correlationId int64, deadline int64) bool

	// CancelTimer cancels a previously scheduled timer. This action is asynchronous and will race with the timer expiring.
	//
	// Timers should only be scheduled or cancelled in the context of processing
	// - onSessionMessage
	// - onTimerEvent
	// - onSessionOpen
	// - onSessionClose
	// If applied to other events then they are not guaranteed to be reliable.
	//
	// Callers of this method should loop until the method succeeds.
	//
	// CancelTimer   returns true if the event to cancel request has been sent or false if back-pressure is applied.
	CancelTimer(correlationId int64) bool

	// Offer a message as ingress to the cluster for sequencing. This will happen efficiently over IPC to the
	// consensus module and set the cluster session as the negative value of the cluster's ServiceID
	//
	// Callers of this method should loop until the method succeeds.
	//
	// The cluster's idle strategy must be used in the body of the loop to allow for the clustered service to be
	// shutdown if required
	Offer(*atomic.Buffer, int32, int32) int64
}
