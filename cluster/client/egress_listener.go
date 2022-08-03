package client

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
)

type EgressListener interface {
	OnConnect(cluster *AeronCluster)

	OnDisconnect(cluster *AeronCluster, details string)

	OnMessage(
		cluster *AeronCluster,
		timestamp int64,
		buffer *atomic.Buffer,
		offset int32,
		length int32,
		header *logbuffer.Header,
	)

	OnNewLeader(
		cluster *AeronCluster,
		leadershipTermId int64,
		leaderMemberId int32,
	)

	OnError(cluster *AeronCluster, details string)
}
