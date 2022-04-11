package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type ClusteredServiceAgent struct {
	a       *aeron.Aeron
	ctx     *aeron.Context
	proxy   *ConsensusModuleProxy
	reader  *counters.Reader
	adapter *ServiceAdapter
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
		int32(105),
	)
	proxy := NewConsensusModuleProxy(options, pub)

	sub := <-a.AddSubscription(
		// TODO: constify?
		"aeron:ipc?term-length=128k|alias=consensus-control",
		int32(104),
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

	agent := &ClusteredServiceAgent{
		a:       a,
		adapter: adapter,
		ctx:     ctx,
		proxy:   proxy,
		reader:  reader,
	}

	adapter.agent = agent

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
	return true
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
