package cluster

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/cluster/codecs"
)

const MarkFileUpdateIntervalMs = 1000
const ServiceStreamId = 104
const ConsensusModuleStreamId = 105

type ClusteredServiceAgent struct {
	a                        *aeron.Aeron
	ctx                      *aeron.Context
	proxy                    *ConsensusModuleProxy
	reader                   *counters.Reader
	adapter                  *ServiceAdapter
	markFile                 *ClusterMarkFile
	cachedTimeMs             int64
	markFileUpdateDeadlineMs int64
}

func NewClusteredServiceAgent(
	ctx *aeron.Context,
	options *Options,
) (*ClusteredServiceAgent, error) {
	a, err := aeron.Connect(ctx)
	if err != nil {
		return nil, err
	}

	pub := <-a.AddPublication(
		// TODO: constify?
		"aeron:ipc?term-length=128k|alias=consensus-control",
		int32(ConsensusModuleStreamId),
	)
	proxy := NewConsensusModuleProxy(options, pub)

	sub := <-a.AddSubscription(
		// TODO: constify?
		"aeron:ipc?term-length=128k|alias=consensus-control",
		int32(ServiceStreamId),
	)
	adapter := &ServiceAdapter{
		marshaller:   codecs.NewSbeGoMarshaller(),
		options:      options,
		subscription: sub,
	}

	counterFile, _, _ := counters.MapFile(ctx.CncFileName())
	reader := counters.NewReader(
		counterFile.ValuesBuf.Get(),
		counterFile.MetaDataBuf.Get(),
	)

	cmf, err := NewClusterMarkFile("/tmp/aeron-cluster/cluster-mark-service-0.dat")
	if err != nil {
		return nil, err
	}

	agent := &ClusteredServiceAgent{
		a:        a,
		adapter:  adapter,
		ctx:      ctx,
		proxy:    proxy,
		reader:   reader,
		markFile: cmf,
	}
	adapter.agent = agent

	cmf.flyweight.ArchiveStreamId.Set(10)
	cmf.flyweight.ServiceStreamId.Set(ServiceStreamId)
	cmf.flyweight.ConsensusModuleStreamId.Set(ConsensusModuleStreamId)
	cmf.flyweight.IngressStreamId.Set(-1)
	cmf.flyweight.MemberId.Set(-1)
	cmf.flyweight.ServiceId.Set(0)
	cmf.flyweight.ClusterId.Set(0)

	cmf.UpdateActivityTimestamp(time.Now().UnixMilli())
	cmf.SignalReady()

	return agent, nil
}

func (agent *ClusteredServiceAgent) OnStart() {
	id := agent.awaitCommitPositionCounter( /* TODO: get real cluster_id */ 0)
	fmt.Println("commit position counter: ", id)
	agent.recoverState()
}

func (agent *ClusteredServiceAgent) awaitCommitPositionCounter(
	clusterID int,
) int32 {
	id := int32(-1)
	agent.reader.Scan(func(counter counters.Counter) {
		if counter.TypeId == /* TODO: constify? */ 203 {
			id = counter.Id
		}
	})

	return id
}

func (agent *ClusteredServiceAgent) recoverState() error {
	id, label := agent.awaitRecoveryCounter()
	fmt.Println("label: ", label)

	parts := strings.Split(label, " ")
	leadershipTermID, err := strconv.ParseInt(strings.Split(parts[2], "=")[1], 10, 64)
	if err != nil {
		return err
	}

	logPosition, err := strconv.ParseInt(strings.Split(parts[3], "=")[1], 10, 64)
	if err != nil {
		return err
	}

	fmt.Printf("leader term id: %d, log position: %d\n", leadershipTermID, logPosition)

	fmt.Println("recovery position counter: ", id)

	if leadershipTermID == -1 {
		// TODO: add service callback class
	} else {
		// TODO: load snapshot
	}

	agent.proxy.ServiceAckRequest(
		logPosition,
		/* TODO: get from label? */ 0,
		/* TODO: real value? */ 1,
		/* TODO: use agent.a.ClientID()? */ -1,
		/* TODO: real value? */ 0,
	)

	return nil
}

func (agent *ClusteredServiceAgent) awaitRecoveryCounter() (int32, string) {
	id := int32(-1)
	label := ""
	agent.reader.Scan(func(counter counters.Counter) {
		if counter.TypeId == /* TODO: constify? */ 204 {
			id = counter.Id
			label = counter.Label
		}
	})

	return id, label
}

func (agent *ClusteredServiceAgent) checkForClockTick() bool {
	nowMs := time.Now().UnixMilli()
	if agent.cachedTimeMs != nowMs {
		agent.cachedTimeMs = nowMs
		if nowMs > agent.markFileUpdateDeadlineMs {
			agent.markFileUpdateDeadlineMs = nowMs + MarkFileUpdateIntervalMs
			agent.markFile.UpdateActivityTimestamp(nowMs)
		}
		return true
	}
	return false
}

func (agent *ClusteredServiceAgent) pollServiceAdapter() {
	agent.adapter.Poll()
}

func (agent *ClusteredServiceAgent) DoWork() int {
	work := 0

	if agent.checkForClockTick() {
		agent.pollServiceAdapter()
	}

	return work
}

type Role int32

const (
	Follower  Role = 0
	Candidate      = 1
	Leader         = 2
)

func (agent *ClusteredServiceAgent) onJoinLog(
	logPosition int64,
	maxLogPosition int64,
	memberId int32,
	logSessionId int32,
	logStreamId int32,
	isStartup bool,
	role Role,
	logChannel string,
) {
	fmt.Println("join log called: ", logPosition, isStartup, role, logChannel)
}
