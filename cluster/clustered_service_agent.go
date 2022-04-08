package cluster

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/counters"
)

type ClusteredServiceAgent struct {
	a      *aeron.Aeron
	ctx    *aeron.Context
	proxy  *ConsensusModuleProxy
	reader *counters.Reader
	sub    *aeron.Subscription
}

func NewClusteredServiceAgent(
	ctx *aeron.Context,
	options *Options,
) (*ClusteredServiceAgent, error) {
	a, err := aeron.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer a.Close()

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

	counterFile, _, _ := counters.MapFile(ctx.CncFileName())
	reader := counters.NewReader(
		counterFile.ValuesBuf.Get(),
		counterFile.MetaDataBuf.Get(),
	)

	return &ClusteredServiceAgent{
		a:      a,
		ctx:    ctx,
		proxy:  proxy,
		sub:    sub,
		reader: reader,
	}, nil
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
		/* TODO: get from label? */ time.Now().Unix(),
		/* TODO: real value? */ 1,
		agent.a.ClientID(),
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
