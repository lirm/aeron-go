package main

import (
	"fmt"
	"os"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type EchoService struct {
	cluster      cluster.Cluster
	messageCount int32
}

func (s *EchoService) OnStart(cluster cluster.Cluster, image *aeron.Image) {
	s.cluster = cluster
	if image == nil {
		fmt.Printf("OnStart with no image\n")
	} else {
		cnt := image.Poll(func(buf *atomic.Buffer, offset int32, length int32, hdr *logbuffer.Header) error {
			if length == 4 && s.messageCount == 0 {
				s.messageCount = buf.GetInt32(offset)
			} else {
				// TODO: Return a proper error
				fmt.Printf("WARNING: unexpected snapshot message - pos=%d offset=%d length=%d\n",
					hdr.Position(), offset, length)
			}
			return nil
		}, 100)
		fmt.Printf("OnStart with image - snapshotMsgCnt=%d messageCount=%d\n", cnt, s.messageCount)
	}
}

func (s *EchoService) OnSessionOpen(session cluster.ClientSession, timestamp int64) {
	fmt.Printf("OnSessionOpen - sessionId=%d timestamp=%v\n", session.Id(), timestamp)
}

func (s *EchoService) OnSessionClose(
	session cluster.ClientSession,
	timestamp int64,
	reason codecs.CloseReasonEnum,
) {
	fmt.Printf("OnSessionClose - sessionId=%d timestamp=%v reason=%v\n", session.Id(), timestamp, reason)
}

func (s *EchoService) OnSessionMessage(
	session cluster.ClientSession,
	timestamp int64,
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	s.messageCount++
	var result int64
	for offerCnt := 1; ; offerCnt++ {
		result = session.Offer(buffer, offset, length, nil)
		if result >= 0 {
			return
		} else if result == aeron.BackPressured || result == aeron.AdminAction {
			s.cluster.IdleStrategy().Idle(0)
		} else {
			fmt.Printf("WARNING: OnSessionMessage offer failed - sessionId=%d time=%d pos=%d len=%d offerCnt=%d result=%v\n",
				session.Id(), timestamp, header.Position(), length, offerCnt, result)
		}
	}
}

func (s *EchoService) OnTimerEvent(correlationId, timestamp int64) {
	fmt.Printf("OnTimerEvent - correlationId=%d timestamp=%v\n", correlationId, timestamp)
}

func (s *EchoService) OnTakeSnapshot(publication *aeron.Publication) {
	fmt.Printf("OnTakeSnapshot - streamId=%d sessionId=%d messageCount=%d\n",
		publication.StreamID(), publication.SessionID(), s.messageCount)
	buf := atomic.MakeBuffer(make([]byte, 4))
	buf.PutInt32(0, s.messageCount)
	for {
		result := publication.Offer(buf, 0, buf.Capacity(), nil)
		if result >= 0 {
			return
		} else if result == aeron.BackPressured || result == aeron.AdminAction {
			s.cluster.IdleStrategy().Idle(0)
		} else {
			fmt.Printf("WARNING: OnTakeSnapshot offer failed - result=%v\n", result)
		}
	}
}

func (s *EchoService) OnRoleChange(role cluster.Role) {
	fmt.Printf("OnRoleChange - role=%v\n", role)
}

func (s *EchoService) OnTerminate(cluster cluster.Cluster) {
	fmt.Printf("OnTerminate - role=%v logPos=%d\n", cluster.Role(), cluster.LogPosition())
}

func (s *EchoService) OnNewLeadershipTermEvent(
	leadershipTermId int64,
	logPosition int64,
	timestamp int64,
	termBaseLogPosition int64,
	leaderMemberId int32,
	logSessionId int32,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) {
	fmt.Printf("OnNewLeadershipTermEvent - leaderTermId=%d logPos=%d time=%d termBase=%d leaderId=%d logSessionId=%d timeUnit=%v appVer=%d\n",
		leadershipTermId, logPosition, timestamp, termBaseLogPosition, leaderMemberId, logSessionId, timeUnit, appVersion)
}

func main() {
	ctx := aeron.NewContext()
	if aeronDir := os.Getenv("AERON_DIR"); aeronDir != "" {
		ctx.AeronDir(aeronDir)
		fmt.Println("aeron dir: ", aeronDir)
	} else if _, err := os.Stat("/dev/shm"); err == nil {
		path := fmt.Sprintf("/dev/shm/aeron-%s", aeron.UserName)
		ctx.AeronDir(path)
		fmt.Println("aeron dir: ", path)
	}

	opts := cluster.NewOptions()
	if idleStr := os.Getenv("NO_OP_IDLE"); idleStr != "" {
		opts.IdleStrategy = &idlestrategy.Busy{}
	}
	if opts.ClusterDir = os.Getenv("CLUSTER_DIR"); opts.ClusterDir == "" {
		opts.ClusterDir = "/tmp/aeron-go-poc/cluster"
	}

	service := &EchoService{}
	agent, err := cluster.NewClusteredServiceAgent(ctx, opts, service)
	if err != nil {
		panic(err)
	}

	agent.StartAndRun()
}
