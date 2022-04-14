package cluster

import (
	"fmt"
	"github.com/lirm/aeron-go/archive"
	"strconv"
	"strings"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/cluster/codecs"
)

const NullValue = -1
const NullPosition = -1
const serviceId = 0
const MarkFileUpdateIntervalMs = 1000
const ServiceStreamId = 104
const ConsensusModuleStreamId = 105
const SnapshotStreamId = 106

const RecordingPosCounterTypeId = 100

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
	timeUnit                 codecs.ClusterTimeUnitEnum
	memberId                 int32
	nextAckId                int64
	terminationPosition      int64
	isServiceActive          bool
	role                     Role
	service                  ClusteredService
	sessions                 map[int64]ClientSession
}

func NewClusteredServiceAgent(
	ctx *aeron.Context,
	options *Options,
	service ClusteredService,
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

	cmf, err := NewClusterMarkFile(options.ClusterDir + "/cluster-mark-service-0.dat")
	if err != nil {
		return nil, err
	}

	agent := &ClusteredServiceAgent{
		a:                   a,
		opts:                options,
		serviceAdapter:      serviceAdapter,
		logAdapter:          logAdapter,
		ctx:                 ctx,
		proxy:               proxy,
		reader:              reader,
		markFile:            cmf,
		role:                Follower,
		service:             service,
		logPosition:         NullPosition,
		terminationPosition: NullPosition,
		sessions:            map[int64]ClientSession{},
	}
	serviceAdapter.agent = agent
	logAdapter.agent = agent

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

func (agent *ClusteredServiceAgent) StartAndRun() {
	if err := agent.OnStart(); err != nil {
		panic(err)
	}
	for agent.isServiceActive {
		agent.opts.IdleStrategy.Idle(agent.DoWork())
	}
}

func (agent *ClusteredServiceAgent) OnStart() error {
	id := agent.awaitCommitPositionCounter( /* TODO: get real cluster_id */ 0)
	fmt.Println("commit position counter: ", id)
	return agent.recoverState()
}

func (agent *ClusteredServiceAgent) awaitCommitPositionCounter(clusterID int) int32 {
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
	agent.isServiceActive = true

	if leadershipTermID == -1 {
		agent.service.OnStart(agent, nil)
	} else {
		// TODO: load snapshot
	}

	return agent.proxy.ServiceAckRequest(
		logPosition,
		agent.clusterTime,
		agent.getAndIncrementNextAckId(),
		agent.a.ClientID(),
		serviceId,
	)
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

	if agent.terminationPosition != NullPosition && agent.logPosition >= agent.terminationPosition {
		if agent.logPosition > agent.terminationPosition {
			fmt.Printf("WARNING: service terminate: logPos=%d > terminationPos=%d\n", agent.logPosition, agent.terminationPosition)
		}
		agent.terminate()
	}
}

func (agent *ClusteredServiceAgent) terminate() {
	agent.isServiceActive = false
	agent.service.OnTerminate(agent)
	err := agent.proxy.ServiceAckRequest(
		agent.logPosition,
		agent.clusterTime,
		agent.getAndIncrementNextAckId(),
		NullValue,
		serviceId,
	)
	if err != nil {
		fmt.Println("WARNING: failed to send termination service ack: ", err)
	}

	agent.terminationPosition = NullPosition
	// TODO: throw new ClusterTerminationException(isTerminationExpected)
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

// TODO: move this to its own file please :)
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

	err := agent.proxy.ServiceAckRequest(
		event.logPosition,
		agent.clusterTime,
		agent.getAndIncrementNextAckId(),
		NullValue,
		serviceId,
	)
	if err != nil {
		fmt.Println("ERROR: failed to send join log service ack: ", err)
	}

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
		agent.service.OnRoleChange(newRole)
	}
}

func (agent *ClusteredServiceAgent) awaitImage(
	sessionId int32,
	subscription *aeron.Subscription,
) *aeron.Image {
	for {
		if img := subscription.ImageBySessionID(sessionId); img != nil {
			return img
		}
		agent.opts.IdleStrategy.Idle(0)
	}
}

func (agent *ClusteredServiceAgent) onSessionOpen() {
	// TODO: implement
}

func (agent *ClusteredServiceAgent) onSessionClose() {
	// TODO: implement
}

func (agent *ClusteredServiceAgent) onSessionMessage(
	logPosition int64,
	clusterSessionId int64,
	timestamp int64,
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	agent.logPosition = logPosition
	agent.clusterTime = timestamp
	clientSession := agent.sessions[clusterSessionId]
	agent.service.OnSessionMessage(
		clientSession,
		timestamp,
		buffer,
		offset,
		length,
		header,
	)
}

func (agent *ClusteredServiceAgent) onNewLeadershipTermEvent(
	leadershipTermId int64,
	logPosition int64,
	timestamp int64,
	termBaseLogPosition int64,
	leaderMemberId int32,
	logSessionId int32,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32) {
	//if (util.SemanticVersionMajor(ctx.appVersion()) != SemanticVersion.major(appVersion))
	//{
	//	ctx.errorHandler().onError(new ClusterException(
	//	"incompatible version: " + SemanticVersion.toString(ctx.appVersion()) +
	//	" log=" + SemanticVersion.toString(appVersion)));
	//	throw new AgentTerminationException();
	//}
	//sessionMessageHeaderEncoder.leadershipTermId(leadershipTermId)
	agent.logPosition = logPosition
	agent.clusterTime = timestamp
	agent.timeUnit = timeUnit

	agent.service.OnNewLeadershipTermEvent(
		leadershipTermId,
		logPosition,
		timestamp,
		termBaseLogPosition,
		leaderMemberId,
		logSessionId,
		codecs.ClusterTimeUnit.MILLIS,
		appVersion)
}

func (agent *ClusteredServiceAgent) onServiceAction(leadershipTermId int64, logPos int64, timestamp int64, action codecs.ClusterActionEnum) {
	agent.logPosition = logPos
	agent.clusterTime = timestamp
	if action == codecs.ClusterAction.SNAPSHOT {
		recordingId, err := agent.takeSnapshot(logPos, leadershipTermId)
		if err != nil {
			fmt.Println("ERROR: take snapshot failed: ", err)
			return
		}
		if err := agent.proxy.ServiceAckRequest(logPos, timestamp, agent.getAndIncrementNextAckId(), recordingId, serviceId); err != nil {
			fmt.Println("WARNING: take snapshot service ack failed: ", err)
		}
	}
}

func (agent *ClusteredServiceAgent) takeSnapshot(logPos int64, leadershipTermId int64) (int64, error) {
	options := archive.DefaultOptions()
	//options.RequestChannel = *examples.Config.RequestChannel
	//options.RequestStream = int32(*examples.Config.RequestStream)
	//options.ResponseChannel = *examples.Config.ResponseChannel
	//options.ResponseStream = int32(*examples.Config.ResponseStream)

	arch, err := archive.NewArchive(options, agent.ctx)
	if err != nil {
		return NullValue, err
	}
	defer arch.Close()

	pub, err := arch.AddRecordedPublication("aeron:ipc?alias=snapshot", SnapshotStreamId)
	if err != nil {
		return NullValue, err
	}
	defer pub.Close()

	recordingId, err := agent.awaitRecordingId(pub.SessionID())
	if err != nil {
		return 0, err
	}

	fmt.Printf("takeSnapshot - got recordingId: %d\n", recordingId)
	snapshotTaker := NewSnapshotTaker(agent.opts, pub)
	if err := snapshotTaker.MarkBegin(logPos, leadershipTermId, agent.timeUnit, agent.opts.AppVersion); err != nil {
		return 0, err
	}
	for _, session := range agent.sessions {
		if err := snapshotTaker.SnapshotSession(session); err != nil {
			return 0, err
		}
	}
	if err := snapshotTaker.MarkEnd(logPos, leadershipTermId, agent.timeUnit, agent.opts.AppVersion); err != nil {
		return 0, err
	}
	agent.checkForClockTick()
	agent.service.OnTakeSnapshot(pub)

	return recordingId, nil
}

func (agent *ClusteredServiceAgent) awaitRecordingId(sessionId int32) (int64, error) {
	counterId := int32(NullValue)
	var err error = nil
	for err == nil && counterId == NullValue {
		agent.reader.Scan(func(counter counters.Counter) {
			if err == nil && counter.TypeId == RecordingPosCounterTypeId {
				thisSessionId, thisErr := agent.reader.GetKeyPartInt32(counter.Id, 8)
				if thisErr != nil {
					err = thisErr
				} else if thisSessionId == sessionId {
					//fmt.Printf("*** awaitRecordingId - sessionId=%d counterId=%d type=%d value=%v label=%v\n",
					//	sessionId, counter.Id, counter.TypeId, counter.Value, counter.Label)
					counterId = counter.Id
				}
			}
		})
		if counterId == NullValue {
			agent.opts.IdleStrategy.Idle(0)
		}
	}
	if err != nil {
		return 0, err
	}
	return agent.reader.GetKeyInt64(counterId, 0)
}

func (agent *ClusteredServiceAgent) onServiceTerminationPosition(position int64) {
	agent.terminationPosition = position
}

func (agent *ClusteredServiceAgent) getAndIncrementNextAckId() int64 {
	ackId := agent.nextAckId
	agent.nextAckId++
	return ackId
}

// BEGIN CLUSTER IMPLEMENTATION

func (agent *ClusteredServiceAgent) Offer(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	reservedValueSupplier term.ReservedValueSupplier,
) int64 {
	// TODO: implement. Needed for client session
	return 0
}

func (agent *ClusteredServiceAgent) getClientSession(id int64) (ClientSession, bool) {
	session, ok := agent.sessions[id]
	return session, ok
}

func (agent *ClusteredServiceAgent) closeClientSession(id int64) (ClientSession, bool) {
	// TODO: implement
	return nil, false
}

func (agent *ClusteredServiceAgent) LogPosition() int64 {
	return agent.logPosition
}

func (agent *ClusteredServiceAgent) MemberId() int32 {
	return agent.memberId
}

func (agent *ClusteredServiceAgent) Role() Role {
	return agent.role
}

func (agent *ClusteredServiceAgent) Time() int64 {
	return agent.clusterTime
}

// END CLUSTER IMPLEMENTATION
