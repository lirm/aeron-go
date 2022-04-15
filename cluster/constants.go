package cluster

const (
	SBEHeaderLength            = 8
	SessionMessageHeaderLength = SBEHeaderLength + 24
)

const (
	clusterSchemaId                 = 111
	sessionMessageHeaderTemplateId  = 1
	timerEventTemplateId            = 20
	sessionOpenTemplateId           = 21
	sessionCloseTemplateId          = 22
	clusterActionReqTemplateId      = 23
	newLeadershipTermTemplateId     = 24
	membershipChangeTemplateId      = 25
	joinLogTemplateId               = 40
	serviceTerminationPosTemplateId = 42
	snapshotMarkerTemplateId        = 100
	clientSessionTemplateId         = 102
)
