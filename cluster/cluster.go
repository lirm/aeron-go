package cluster

import "github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"

type Cluster interface {
	LogPosition() int64
	MemberId() int32
	Role() Role
	Time() int64
	IdleStrategy() idlestrategy.Idler

	// ScheduleTimer schedules a timer for a given deadline
	ScheduleTimer(correlationId int64, deadline int64) bool
	CancelTimer(correlationId int64) bool
}
