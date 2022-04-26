package cluster

import (
	"fmt"
	"time"

	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/counters"
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer/term"
	"github.com/corymonroe-coinbase/aeron-go/aeron/util"
	"github.com/corymonroe-coinbase/aeron-go/archive"
	"github.com/corymonroe-coinbase/aeron-go/cluster/codecs"
)

const NullValue = -1
const NullPosition = -1
const serviceId = 0
const MarkFileUpdateIntervalMs = 1000

const (
	ReplayStreamId          = 103
	ServiceStreamId         = 104
	ConsensusModuleStreamId = 105
	SnapshotStreamId        = 106
)

const (
	RecordingPosCounterTypeId  = 100
	CommitPosCounterTypeId     = 203
	RecoveryStateCounterTypeId = 204
)

type ClusteredServiceAgent struct {
	a                        *aeron.Aeron
	ctx                      *aeron.Context
	opts                     *Options
	proxy                    *consensusModuleProxy
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
	commitPosition           *counters.ReadableCounter
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
	proxy := newConsensusModuleProxy(options, pub)

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
	proxy.idleStrategy = agent

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
	if err := agent.awaitCommitPositionCounter(); err != nil {
		return err
	}
	return agent.recoverState()
}

func (agent *ClusteredServiceAgent) awaitCommitPositionCounter() error {
	for {
		id := agent.reader.FindCounter(CommitPosCounterTypeId, func(keyBuffer *atomic.Buffer) bool {
			return keyBuffer.GetInt32(0) == agent.opts.ClusterId
		})
		if id != counters.NullCounterId {
			commitPos, err := counters.NewReadableCounter(agent.reader, id)
			fmt.Printf("found commit position counter - id=%d value=%d\n", id, commitPos.Get())
			agent.commitPosition = commitPos
			return err
		}
		agent.Idle(0)
	}
}

func (agent *ClusteredServiceAgent) recoverState() error {
	counterId, leadershipTermId := agent.awaitRecoveryCounter()
	fmt.Printf("found recovery counter - id=%d leadershipTermId=%d logPos=%d clusterTime=%d\n",
		counterId, leadershipTermId, agent.logPosition, agent.clusterTime)
	agent.isServiceActive = true

	if leadershipTermId == -1 {
		agent.service.OnStart(agent, nil)
	} else {
		serviceCount, err := agent.reader.GetKeyPartInt32(counterId, 28)
		if err != nil {
			return err
		}
		if serviceCount < 1 {
			return fmt.Errorf("invalid service count: %d", serviceCount)
		}
		snapshotRecId, err := agent.reader.GetKeyPartInt64(counterId, 32+(serviceId*util.SizeOfInt64))
		if err != nil {
			return err
		}
		if err := agent.loadSnapshot(snapshotRecId); err != nil {
			return err
		}
	}

	agent.proxy.serviceAckRequest(
		agent.logPosition,
		agent.clusterTime,
		agent.getAndIncrementNextAckId(),
		agent.a.ClientID(),
		serviceId,
	)
	return nil
}

func (agent *ClusteredServiceAgent) awaitRecoveryCounter() (int32, int64) {
	for {
		var leadershipTermId int64
		id := agent.reader.FindCounter(RecoveryStateCounterTypeId, func(keyBuffer *atomic.Buffer) bool {
			if keyBuffer.GetInt32(24) == agent.opts.ClusterId {
				leadershipTermId = keyBuffer.GetInt64(0)
				agent.logPosition = keyBuffer.GetInt64(8)
				agent.clusterTime = keyBuffer.GetInt64(16)
				return true
			}
			return false
		})
		if id != counters.NullCounterId {
			return id, leadershipTermId
		}
		agent.Idle(0)
	}
}

func (agent *ClusteredServiceAgent) loadSnapshot(recordingId int64) error {
	archOps := archive.DefaultOptions()
	arch, err := archive.NewArchive(archOps, agent.ctx)
	if err != nil {
		return err
	}
	defer closeArchive(arch)

	channel := "aeron:ipc?alias=snapshot-replay"
	replaySessionId, err := arch.StartReplay(recordingId, 0, NullValue, channel, ReplayStreamId)
	if err != nil {
		return err
	}
	subChannel, err := archive.AddSessionIdToChannel(channel, archive.ReplaySessionIdToStreamId(replaySessionId))
	if err != nil {
		return err
	}

	fmt.Printf("replaying snapshot - recId=%d sessionId=%d streamId=%d\n", recordingId, replaySessionId, ReplayStreamId)
	subscription := <-arch.AddSubscription(subChannel, ReplayStreamId)
	defer closeSubscription(subscription)

	img := agent.awaitImage(int32(replaySessionId), subscription)
	loader := newSnapshotLoader(agent, img)
	for !loader.isDone {
		agent.opts.IdleStrategy.Idle(loader.poll())
	}
	agent.timeUnit = loader.timeUnit
	agent.service.OnStart(agent, img)
	return nil
}

func (agent *ClusteredServiceAgent) addSessionFromSnapshot(session *ContainerClientSession) {
	agent.sessions[session.id] = session
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
	agent.proxy.serviceAckRequest(
		agent.logPosition,
		agent.clusterTime,
		agent.getAndIncrementNextAckId(),
		NullValue,
		serviceId,
	)
	agent.terminationPosition = NullPosition
}

func (agent *ClusteredServiceAgent) DoWork() int {
	work := 0

	if agent.checkForClockTick() {
		agent.pollServiceAdapter()
	}

	if agent.logAdapter.image != nil {
		polled := agent.logAdapter.Poll(agent.commitPosition.Get())
		work += polled
		if polled == 0 && agent.logAdapter.IsDone() {
			agent.closeLog()
		}
	}

	return work
}

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

	agent.proxy.serviceAckRequest(
		event.logPosition,
		agent.clusterTime,
		agent.getAndIncrementNextAckId(),
		NullValue,
		serviceId,
	)

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

func (agent *ClusteredServiceAgent) onSessionOpen(
	leadershipTermId int64,
	logPosition int64,
	clusterSessionId int64,
	timestamp int64,
	responseStreamId int32,
	responseChannel string,
	encodedPrincipal []byte,
) {
	agent.logPosition = logPosition
	agent.clusterTime = timestamp
	if _, ok := agent.sessions[clusterSessionId]; ok {
		fmt.Printf("clashing session id: %d, ignored\n", clusterSessionId)
	} else {
		session := NewContainerClientSession(
			clusterSessionId,
			responseStreamId,
			responseChannel,
			agent,
		)
		// TODO: looks like we only want to connect if this is the leader
		// currently always connecting

		agent.sessions[session.id] = session
		agent.service.OnSessionOpen(session, timestamp)
	}
}

func (agent *ClusteredServiceAgent) onSessionClose(
	leadershipTermId int64,
	logPosition int64,
	clusterSessionId int64,
	timestamp int64,
	closeReason codecs.CloseReasonEnum,
) {
	agent.logPosition = logPosition
	agent.clusterTime = timestamp

	if session, ok := agent.sessions[clusterSessionId]; ok {
		delete(agent.sessions, clusterSessionId)
		agent.service.OnSessionClose(session, timestamp, closeReason)
	} else {
		fmt.Printf("unknown session id: %d, ignored\n", clusterSessionId)
	}
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
		agent.proxy.serviceAckRequest(logPos, timestamp, agent.getAndIncrementNextAckId(), recordingId, serviceId)
	}
}

func (agent *ClusteredServiceAgent) onTimerEvent(
	logPosition int64,
	correlationId int64,
	timestamp int64,
) {
	agent.logPosition = logPosition
	agent.clusterTime = timestamp
	agent.service.OnTimerEvent(correlationId, timestamp)
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
	defer closeArchive(arch)

	pub, err := arch.AddRecordedPublication("aeron:ipc?alias=snapshot", SnapshotStreamId)
	if err != nil {
		return NullValue, err
	}
	defer closePublication(pub)

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
	start := time.Now()
	for time.Since(start) < agent.opts.Timeout {
		recId := int64(NullValue)
		counterId := agent.reader.FindCounter(RecordingPosCounterTypeId, func(keyBuffer *atomic.Buffer) bool {
			if keyBuffer.GetInt32(8) == sessionId {
				recId = keyBuffer.GetInt64(0)
				return true
			}
			return false
		})
		if counterId != NullValue {
			return recId, nil
		}
		agent.Idle(0)
	}
	return NullValue, fmt.Errorf("timed out waiting for recordingId for sessionId=%d", sessionId)
}

func (agent *ClusteredServiceAgent) onServiceTerminationPosition(position int64) {
	agent.terminationPosition = position
}

func (agent *ClusteredServiceAgent) getAndIncrementNextAckId() int64 {
	ackId := agent.nextAckId
	agent.nextAckId++
	return ackId
}

func (agent *ClusteredServiceAgent) Offer(
	clusterSessionId int64,
	publication *aeron.Publication,
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	reservedValueSupplier term.ReservedValueSupplier,
) int64 {
	if agent.role != Leader {
		return ClientSessionMockedOffer
	}

	bytes, err := codecs.SessionMessageHeaderPacket(
		codecs.NewSbeGoMarshaller(),
		agent.opts.RangeChecking,
		/* leaderTermId??? */ 0,
		clusterSessionId,
		agent.clusterTime,
	)
	if err != nil {
		fmt.Printf("error on ClusteredServiceAgent.Offer: %v\n", err)
		return 0
	}

	combined := append(bytes, buffer.GetBytesArray(offset, length)...)

	return publication.Offer(
		atomic.MakeBuffer(combined),
		0, int32(len(combined)),
		reservedValueSupplier,
	)
}

func (agent *ClusteredServiceAgent) getClientSession(id int64) (ClientSession, bool) {
	session, ok := agent.sessions[id]
	return session, ok
}

func (agent *ClusteredServiceAgent) closeClientSession(id int64) {
	if _, ok := agent.sessions[id]; ok {
		agent.proxy.closeSessionRequest(id)
	} else {
		fmt.Printf("unknown session id: %d, ignored\n", id)
	}
}

func closeArchive(arch *archive.Archive) {
	err := arch.Close()
	if err != nil {
		fmt.Println("error closing archive connection: ", err)
	}
}

func closeSubscription(sub *aeron.Subscription) {
	err := sub.Close()
	if err != nil {
		fmt.Printf("error closing subscription, streamId=%d channel=%s: %s\n", sub.StreamID(), sub.Channel(), err)
	}
}

func closePublication(pub *aeron.Publication) {
	err := pub.Close()
	if err != nil {
		fmt.Printf("error closing publication, streamId=%d channel=%s: %s\n", pub.StreamID(), pub.Channel(), err)
	}
}

func (agent *ClusteredServiceAgent) Idle(workCount int) {
	agent.opts.IdleStrategy.Idle(workCount)
	if workCount <= 0 {
		agent.checkForClockTick()
	}
}

// BEGIN CLUSTER IMPLEMENTATION

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

func (agent *ClusteredServiceAgent) IdleStrategy() idlestrategy.Idler {
	return agent
}

func (agent *ClusteredServiceAgent) ScheduleTimer(correlationId int64, deadline int64) bool {
	return agent.proxy.scheduleTimer(correlationId, deadline)
}

func (agent *ClusteredServiceAgent) CancelTimer(correlationId int64) bool {
	return agent.proxy.cancelTimer(correlationId)
}

// END CLUSTER IMPLEMENTATION
