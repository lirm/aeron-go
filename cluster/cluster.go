package cluster

import "github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"

type Cluster interface {
	LogPosition() int64
	MemberId() int32
	Role() Role
	Time() int64
	IdleStrategy() idlestrategy.Idler
}
