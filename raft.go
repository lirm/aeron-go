package main

import (
	"fmt"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type Service struct {
}

func (s *Service) OnStart(cluster cluster.Cluster, image *aeron.Image) {
	fmt.Printf("OnStart called\n")
}

func (s *Service) OnSessionOpen(session cluster.ClientSession, timestamp int64) {}

func (s *Service) OnSessionClose(
	session cluster.ClientSession,
	timestamp int64,
	closeReason codecs.CloseReasonEnum,
) {
}

func (s *Service) OnSessionMessage(
	session cluster.ClientSession,
	timestamp int64,
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	fmt.Printf("OnSessionMessage called: %s\n", string(buffer.GetBytesArray(offset, length)))
}

func (s *Service) OnTimerEvent(correlationId, timestamp int64) {}

func (s *Service) OnTakeSnapshot(publication *aeron.Publication) {}

func (s *Service) OnRoleChange(role cluster.Role) {
	fmt.Printf("OnRoleChange called: %v\n", role)
}

func (s *Service) OnTerminate(cluster cluster.Cluster) {}

func (s *Service) OnNewLeadershipTermEvent(
	leadershipTermId int64,
	logPosition int64,
	timestamp int64,
	termBaseLogPosition int64,
	leaderMemberId int32,
	logSessionId int32,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) {
	fmt.Printf("OnNewLeadershipTermEvent called: %d, %d, %d\n", leadershipTermId, logPosition, timestamp)
}

func main() {
	ctx := aeron.NewContext()
	opts := cluster.NewOptions()
	service := &Service{}
	agent, err := cluster.NewClusteredServiceAgent(ctx, opts, service)
	if err != nil {
		panic(err)
	}

	agent.OnStart()
	for {
		opts.IdleStrategy.Idle(agent.DoWork())
	}
}
