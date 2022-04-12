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

const serviceId = 0
const MarkFileUpdateIntervalMs = 1000
const ServiceStreamId = 104
const ConsensusModuleStreamId = 105

type ClusteredServiceAgent struct {
	a                        *aeron.Aeron
	ctx                      *aeron.Context
	opts                     *Options
	proxy                    *ConsensusModuleProxy
	reader                   *counters.Reader
	serviceAdapter           *ServiceAdapter
	logAdapter               *BoundedLogAdapter
	markFile                 *ClusterMarkFile
	activeLogEvent           *activeLogEvent
	cachedTimeMs             int64
	markFileUpdateDeadlineMs int64
	logPosition              int64
	clusterTime              int64
	memberId                 int32
	nextAckId                int64
	role                     Role
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
	serviceAdapter := &ServiceAdapter{
		marshaller:   codecs.NewSbeGoMarshaller(),
		options:      options,
		subscription: sub,
	}
	logAdapter := &BoundedLogAdapter{
		marshaller: codecs.NewSbeGoMarshaller(),
		options:    options,
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
		a:              a,
		opts:           options,
		serviceAdapter: serviceAdapter,
		logAdapter:     logAdapter,
		ctx:            ctx,
		proxy:          proxy,
		reader:         reader,
		markFile:       cmf,
		role:           Follower,
	}
	serviceAdapter.agent = agent

	cmf.flyweight.ArchiveStreamId.Set(10)
	cmf.flyweight.ServiceStreamId.Set(ServiceStreamId)
	cmf.flyweight.ConsensusModuleStreamId.Set(ConsensusModuleStreamId)
	cmf.flyweight.IngressStreamId.Set(-1)
	cmf.flyweight.MemberId.Set(-1)
	cmf.flyweight.ServiceId.Set(serviceId)
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

	ackId := agent.nextAckId
	agent.nextAckId++
	agent.proxy.ServiceAckRequest(
		logPosition,
		agent.clusterTime,
		ackId,
		/* TODO: use agent.a.ClientID()? */ -1,
		serviceId,
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
	agent.serviceAdapter.Poll()

	if agent.activeLogEvent != nil && agent.logAdapter.image == nil {
		event := agent.activeLogEvent
		agent.activeLogEvent = nil
		agent.joinActiveLog(event)
	}

	// TODO:
	//if (NULL_POSITION != terminationPosition && logPosition >= terminationPosition)
	//{
	//	terminate();
	//}
}

func (agent *ClusteredServiceAgent) DoWork() int {
	work := 0

	if agent.checkForClockTick() {
		agent.pollServiceAdapter()
	}

	if agent.logAdapter.image != nil {
		polled := agent.logAdapter.Poll()
		work += polled
		if polled == 0 && agent.logAdapter.IsDone() {
			agent.closeLog()
		}
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
	agent.logAdapter.maxLogPosition = logPosition
	event := &activeLogEvent{
		logPosition:    logPosition,
		maxLogPosition: maxLogPosition,
		memberId:       memberId,
		logSessionId:   logSessionId,
		logStreamId:    logStreamId,
		isStartup:      isStartup,
		role:           role,
		logChannel:     logChannel,
	}
	agent.activeLogEvent = event
}

type activeLogEvent struct {
	logPosition    int64
	maxLogPosition int64
	memberId       int32
	logSessionId   int32
	logStreamId    int32
	isStartup      bool
	role           Role
	logChannel     string
}

func (agent *ClusteredServiceAgent) joinActiveLog(event *activeLogEvent) {
	logSub := <-agent.a.AddSubscription(event.logChannel, event.logStreamId)
	img := agent.awaitImage(event.logSessionId, logSub)
	if img.Position() != agent.logPosition {
		fmt.Printf("joinActiveLog - image.position: %v expected: %v\n", img.Position(), agent.logPosition)
		// TODO: close logSub and return error
	}
	agent.logAdapter.image = img
	agent.logAdapter.maxLogPosition = event.maxLogPosition

	ackId := agent.nextAckId
	agent.nextAckId++
	agent.proxy.ServiceAckRequest(event.logPosition, agent.clusterTime, ackId, -1, serviceId)

	agent.memberId = event.memberId
	agent.markFile.flyweight.MemberId.Set(agent.memberId)

	agent.setRole(event.role)
}

func (agent *ClusteredServiceAgent) closeLog() {
	imageLogPos := agent.logAdapter.image.Position()
	if imageLogPos > agent.logPosition {
		agent.logPosition = imageLogPos
	}
	if err := agent.logAdapter.Close(); err != nil {
		fmt.Println("error closing log image: ", err)
	}
	agent.setRole(Follower)
}

func (agent *ClusteredServiceAgent) setRole(newRole Role) {
	if newRole != agent.role {
		agent.role = newRole
		// TODO: service.onRoleChange(newRole)
	}
}

func (agent *ClusteredServiceAgent) awaitImage(sessionId int32, subscription *aeron.Subscription) *aeron.Image {
	for {
		if img := subscription.ImageBySessionID(sessionId); img != nil {
			return img
		}
		agent.opts.IdleStrategy.Idle(0)
	}
}