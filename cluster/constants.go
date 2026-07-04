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
	ClusterSchemaId                 = 111
	ClusterSchemaVersion            = 16
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
	requestServiceAckTemplateId     = 108
)

// Flags for a cluster action request, mirroring io.aeron.cluster.ConsensusModule
const (
	clusterActionFlagsDefault         int32 = 0
	clusterActionFlagsStandbySnapshot int32 = 1
)

const SessionMessageHdrBlockLength = 24
