package cluster

const (
	SBEHeaderLength            = 8
	SessionMessageHeaderLength = 24
)

type Role int32

const (
	Follower  Role = 0
	Candidate Role = 1
	Leader    Role = 2
)

const (
	ClusterSchemaId                 = 111
	ClusterSchemaVersion            = 8
	SessionMessageHeaderTemplateId  = 1
	SessionEventTemplateId          = 2
	SessionCloseRequestTemplateId   = 4
	SessionKeepAliveTemplateId      = 5
	NewLeaderEventTemlateId         = 6
	ChallengeTemplateId             = 7
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
	RequestServiceAckId             = 108
)

const SessionMessageHdrBlockLength = 24
