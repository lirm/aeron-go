package cluster

const (
	SBEHeaderLength            = 8
	SessionMessageHeaderLength = 24
)

type Role int32

const (
	Follower  Role = 0
	Candidate      = 1
	Leader         = 2
)

const (
	clusterSchemaId                 = 111
	clusterSchemaVersion            = 8
	sessionMessageHeaderTemplateId  = 1
	timerEventTemplateId            = 20
	sessionOpenTemplateId           = 21
	sessionCloseTemplateId          = 22
	clusterActionReqTemplateId      = 23
	newLeadershipTermTemplateId     = 24
	membershipChangeTemplateId      = 25
	scheduleTimerTemplateId         = 31
	cancelTimerTemplateId           = 32
	joinLogTemplateId               = 40
	serviceTerminationPosTemplateId = 42
	snapshotMarkerTemplateId        = 100
	clientSessionTemplateId         = 102
)
