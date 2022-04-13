package cluster

type Cluster interface {
	LogPosition() int64
	MemberId() int32
	Role() Role
	Time() int64
}
