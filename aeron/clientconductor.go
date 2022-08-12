/*
Copyright 2016-2018 Stanislav Liberman
Copyright (C) 2022 Talos, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aeron

import (
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/broadcast"
	ctr "github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/driver"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logging"
)

var RegistrationStatus = struct {
	AwaitingMediaDriver   int
	RegisteredMediaDriver int
	ErroredMediaDriver    int
}{
	0,
	1,
	2,
}

const (
	keepaliveTimeoutNS = 500 * int64(time.Millisecond)
	resourceTimeoutNS  = 1000 * int64(time.Millisecond)

	// heartbeatTypeId is the type id of a heartbeat counter.
	heartbeatTypeId = int32(11)

	// registrationIdOffset is the offset in the key metadata for the registration id of the counter.
	heartheatRegistrationIdOffset = int32(0)
)

type publicationStateDefn struct {
	regID                    int64
	origRegID                int64
	timeOfRegistration       int64
	streamID                 int32
	sessionID                int32
	posLimitCounterID        int32
	channelStatusIndicatorID int32
	errorCode                int32
	status                   int
	channel                  string
	errorMessage             string
	buffers                  *logbuffer.LogBuffers
	publication              *Publication
}

func (pub *publicationStateDefn) Init(channel string, regID int64, streamID int32, now int64) *publicationStateDefn {
	pub.channel = channel
	pub.regID = regID
	pub.streamID = streamID
	pub.sessionID = -1
	pub.posLimitCounterID = -1
	pub.timeOfRegistration = now
	pub.status = RegistrationStatus.AwaitingMediaDriver

	return pub
}

type subscriptionStateDefn struct {
	regID              int64
	timeOfRegistration int64
	streamID           int32
	errorCode          int32
	status             int
	channel            string
	errorMessage       string
	subscription       *Subscription
}

func (sub *subscriptionStateDefn) Init(ch string, regID int64, sID int32, now int64) *subscriptionStateDefn {
	sub.channel = ch
	sub.regID = regID
	sub.streamID = sID
	sub.timeOfRegistration = now
	sub.status = RegistrationStatus.AwaitingMediaDriver

	return sub
}

type lingerResourse struct {
	lastTime int64
	resource io.Closer
}

type ClientConductor struct {
	pubs []*publicationStateDefn
	subs []*subscriptionStateDefn

	driverProxy *driver.Proxy

	counterValuesBuffer *atomic.Buffer
	counterReader       *ctr.Reader

	driverListenerAdapter *driver.ListenerAdapter

	adminLock sync.Mutex

	pendingCloses      map[int64]chan bool
	lingeringResources chan lingerResourse

	onNewPublicationHandler   NewPublicationHandler
	onNewSubscriptionHandler  NewSubscriptionHandler
	onAvailableImageHandler   AvailableImageHandler
	onUnavailableImageHandler UnavailableImageHandler
	errorHandler              func(error)

	running          atomic.Bool
	conductorRunning atomic.Bool
	driverActive     atomic.Bool

	timeOfLastKeepalive             int64
	timeOfLastCheckManagedResources int64
	timeOfLastDoWork                int64
	driverTimeoutNs                 int64
	interServiceTimeoutNs           int64
	publicationConnectionTimeoutNs  int64
	resourceLingerTimeoutNs         int64

	heartbeatTimestamp *ctr.AtomicCounter
}

// Init is the primary initialization method for ClientConductor
func (cc *ClientConductor) Init(driverProxy *driver.Proxy, bcast *broadcast.CopyReceiver,
	interServiceTo, driverTo, pubConnectionTo, lingerTo time.Duration, counters *ctr.MetaDataFlyweight) *ClientConductor {

	logger.Debugf("Initializing ClientConductor with: %v %v %d %d %d", driverProxy, bcast, interServiceTo,
		driverTo, pubConnectionTo)

	cc.driverProxy = driverProxy
	cc.running.Set(true)
	cc.driverActive.Set(true)
	cc.driverListenerAdapter = driver.NewAdapter(cc, bcast)
	cc.interServiceTimeoutNs = interServiceTo.Nanoseconds()
	cc.driverTimeoutNs = driverTo.Nanoseconds()
	cc.publicationConnectionTimeoutNs = pubConnectionTo.Nanoseconds()
	cc.resourceLingerTimeoutNs = lingerTo.Nanoseconds()

	cc.counterValuesBuffer = counters.ValuesBuf.Get()
	cc.counterReader = ctr.NewReader(counters.ValuesBuf.Get(), counters.MetaDataBuf.Get())

	cc.pendingCloses = make(map[int64]chan bool)
	cc.lingeringResources = make(chan lingerResourse, 1024)

	cc.pubs = make([]*publicationStateDefn, 0)
	cc.subs = make([]*subscriptionStateDefn, 0)

	return cc
}

// Close will terminate the Run() goroutine body and close all active publications and subscription. Run() can
// be restarted in a another goroutine.
func (cc *ClientConductor) Close() (err error) {
	logger.Debugf("Closing ClientConductor")

	now := time.Now().UnixNano()

	running := cc.running.Get()

	cc.closeAllResources(now)
	if running {
		cc.driverProxy.ClientClose()
	}

	timeoutDuration := 5 * time.Second
	timeout := time.Now().Add(timeoutDuration)
	for cc.conductorRunning.Get() && time.Now().Before(timeout) {
		time.Sleep(10 * time.Millisecond)
	}
	if cc.conductorRunning.Get() {
		msg := fmt.Sprintf("failed to stop conductor after %v", timeoutDuration)
		logger.Warning(msg)
		err = errors.New(msg)
	}

	logger.Debugf("Closed ClientConductor")
	return err
}

// Start begins the main execution loop of ClientConductor on a goroutine.
func (cc *ClientConductor) Start(idleStrategy idlestrategy.Idler) {
	cc.running.Set(true)
	go cc.run(idleStrategy)
}

// run is the main execution loop of ClientConductor.
func (cc *ClientConductor) run(idleStrategy idlestrategy.Idler) {
	now := time.Now().UnixNano()
	cc.timeOfLastKeepalive = now
	cc.timeOfLastCheckManagedResources = now
	cc.timeOfLastDoWork = now

	// Stay on the same thread for performance
	runtime.LockOSThread()

	// Clean exit from this particular go routine
	defer func() {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("Panic: %v", err)
			logger.Error(errStr)
			cc.onError(errors.New(errStr))
			cc.running.Set(false)
		}
		cc.conductorRunning.Set(false)
		logger.Infof("ClientConductor done")
	}()

	cc.conductorRunning.Set(true)
	for cc.running.Get() {
		idleStrategy.Idle(cc.doWork())
	}
}

func (cc *ClientConductor) doWork() (workCount int) {
	workCount += cc.driverListenerAdapter.ReceiveMessages()
	workCount += cc.onHeartbeatCheckTimeouts()
	return
}

func (cc *ClientConductor) verifyDriverIsActive() {
	if !cc.driverActive.Get() {
		log.Fatal("Driver is not active")
	}
}

// AddPublication sends the add publication command through the driver proxy
func (cc *ClientConductor) AddPublication(channel string, streamID int32) int64 {
	logger.Debugf("AddPublication: channel=%s, streamId=%d", channel, streamID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	regID := cc.driverProxy.AddPublication(channel, streamID)

	pubState := new(publicationStateDefn)
	pubState.Init(channel, regID, streamID, now)

	cc.pubs = append(cc.pubs, pubState)

	return regID
}

// AddExclusivePublication sends the add publication command through the driver proxy
func (cc *ClientConductor) AddExclusivePublication(channel string, streamID int32) int64 {
	logger.Debugf("AddExclusivePublication: channel=%s, streamId=%d", channel, streamID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	regID := cc.driverProxy.AddExclusivePublication(channel, streamID)

	pubState := new(publicationStateDefn)
	pubState.Init(channel, regID, streamID, now)

	cc.pubs = append(cc.pubs, pubState)

	return regID
}

func (cc *ClientConductor) FindPublication(regID int64) *Publication {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var publication *Publication
	for _, pub := range cc.pubs {
		if pub.regID == regID {
			if pub.publication != nil {
				publication = pub.publication
			} else {
				switch pub.status {
				case RegistrationStatus.AwaitingMediaDriver:
					waitForMediaDriver(pub.timeOfRegistration, cc)
				case RegistrationStatus.RegisteredMediaDriver:
					publication = NewPublication(pub.buffers)
					publication.conductor = cc
					publication.channel = pub.channel
					publication.regID = regID
					publication.originalRegID = pub.origRegID
					publication.streamID = pub.streamID
					publication.sessionID = pub.sessionID
					publication.pubLimit = NewPosition(cc.counterValuesBuffer, pub.posLimitCounterID)
					publication.channelStatusIndicatorID = pub.channelStatusIndicatorID
					pub.publication = publication

				case RegistrationStatus.ErroredMediaDriver:
					log.Fatalf("Error on %d: %d: %s", regID, pub.errorCode, pub.errorMessage)
				}
			}
			break
		}
	}

	return publication
}

func (cc *ClientConductor) releasePublication(regID int64) {
	logger.Debugf("ReleasePublication: regID=%d", regID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	pubcnt := len(cc.pubs)
	for i, pub := range cc.pubs {
		if pub != nil && pub.regID == regID {
			cc.driverProxy.RemovePublication(regID)

			cc.pubs[i] = cc.pubs[pubcnt-1]
			cc.pubs[pubcnt-1] = nil
			pubcnt--

			if pub.buffers.DecRef() == 0 {
				cc.lingeringResources <- lingerResourse{now, pub.buffers}
			}
		}
	}
	cc.pubs = cc.pubs[:pubcnt]
}

// AddSubscription sends the add subscription command through the driver proxy
func (cc *ClientConductor) AddSubscription(channel string, streamID int32) int64 {
	logger.Debugf("AddSubscription: channel=%s, streamId=%d", channel, streamID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	regID := cc.driverProxy.AddSubscription(channel, streamID)

	subState := new(subscriptionStateDefn)
	subState.Init(channel, regID, streamID, now)

	cc.subs = append(cc.subs, subState)

	return regID
}

func (cc *ClientConductor) FindSubscription(regID int64) *Subscription {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var subscription *Subscription
	for _, sub := range cc.subs {
		if sub.regID == regID {

			switch sub.status {
			case RegistrationStatus.AwaitingMediaDriver:
				waitForMediaDriver(sub.timeOfRegistration, cc)
			case RegistrationStatus.ErroredMediaDriver:
				errStr := fmt.Sprintf("Error on %d: %d: %s", regID, sub.errorCode, sub.errorMessage)
				cc.onError(errors.New(errStr))
				log.Fatalf(errStr)
			}

			subscription = sub.subscription
			break
		}
	}

	return subscription
}

func waitForMediaDriver(timeOfRegistration int64, cc *ClientConductor) {
	if now := time.Now().UnixNano(); now > (timeOfRegistration + cc.driverTimeoutNs) {
		errStr := fmt.Sprintf("No response from driver. started: %d, now: %d, to: %d",
			timeOfRegistration/time.Millisecond.Nanoseconds(),
			now/time.Millisecond.Nanoseconds(),
			cc.driverTimeoutNs/time.Millisecond.Nanoseconds())
		cc.onError(errors.New(errStr))
	}
}

func (cc *ClientConductor) releaseSubscription(regID int64, images []Image) {
	logger.Debugf("ReleaseSubscription: regID=%d", regID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	subcnt := len(cc.subs)
	for i, sub := range cc.subs {
		if sub != nil && sub.regID == regID {
			if logger.IsEnabledFor(logging.DEBUG) {
				logger.Debugf("Removing subscription: %d; %v", regID, images)
			}

			cc.driverProxy.RemoveSubscription(regID)

			cc.subs[i] = cc.subs[subcnt-1]
			cc.subs[subcnt-1] = nil
			subcnt--

			for i := range images {
				image := &images[i]
				if cc.onUnavailableImageHandler != nil {
					cc.onUnavailableImageHandler(image)
				}
				cc.lingeringResources <- lingerResourse{now, image}
			}
		}
	}
	cc.subs = cc.subs[:subcnt]
}

// AddDestination sends the add destination command through the driver proxy
func (cc *ClientConductor) AddDestination(registrationID int64, endpointChannel string) {
	logger.Debugf("AddDestination: regID=%d endpointChannel=%s", registrationID, endpointChannel)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	cc.driverProxy.AddDestination(registrationID, endpointChannel)
}

// RemoveDestination sends the remove destination command through the driver proxy
func (cc *ClientConductor) RemoveDestination(registrationID int64, endpointChannel string) {
	logger.Debugf("RemoveDestination: regID=%d endpointChannel=%s", registrationID, endpointChannel)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	cc.driverProxy.RemoveDestination(registrationID, endpointChannel)
}

// AddRcvDestination sends the add rcv destination command through the driver proxy
func (cc *ClientConductor) AddRcvDestination(registrationID int64, endpointChannel string) {
	logger.Debugf("AddRcvDestination: regID=%d endpointChannel=%s", registrationID, endpointChannel)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	cc.driverProxy.AddRcvDestination(registrationID, endpointChannel)
}

// RemoveRcvDestination sends the remove rcv destination command through the driver proxy
func (cc *ClientConductor) RemoveRcvDestination(registrationID int64, endpointChannel string) {
	logger.Debugf("RemoveRcvDestination: regID=%d endpointChannel=%s", registrationID, endpointChannel)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	cc.driverProxy.RemoveRcvDestination(registrationID, endpointChannel)
}

func (cc *ClientConductor) OnNewPublication(streamID int32, sessionID int32, posLimitCounterID int32,
	channelStatusIndicatorID int32, logFileName string, regID int64, origRegID int64) {

	logger.Debugf("OnNewPublication: streamId=%d, sessionId=%d, posLimitCounterID=%d, channelStatusIndicatorID=%d, logFileName=%s, correlationID=%d, regID=%d",
		streamID, sessionID, posLimitCounterID, channelStatusIndicatorID, logFileName, regID, origRegID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.pubs {
		if pubDef.regID == regID {
			pubDef.status = RegistrationStatus.RegisteredMediaDriver
			pubDef.sessionID = sessionID
			pubDef.posLimitCounterID = posLimitCounterID
			pubDef.channelStatusIndicatorID = channelStatusIndicatorID
			pubDef.buffers = logbuffer.Wrap(logFileName)
			pubDef.buffers.IncRef()
			pubDef.origRegID = origRegID

			logger.Debugf("Updated publication: %v", pubDef)

			if cc.onNewPublicationHandler != nil {
				cc.onNewPublicationHandler(pubDef.channel, streamID, sessionID, regID)
			}
		}
	}
}

// TODO Implement logic specific to exclusive publications
func (cc *ClientConductor) OnNewExclusivePublication(streamID int32, sessionID int32, posLimitCounterID int32,
	channelStatusIndicatorID int32, logFileName string, regID int64, origRegID int64) {

	logger.Debugf("OnNewExclusivePublication: streamId=%d, sessionId=%d, posLimitCounterID=%d, channelStatusIndicatorID=%d, logFileName=%s, correlationID=%d, regID=%d",
		streamID, sessionID, posLimitCounterID, channelStatusIndicatorID, logFileName, regID, origRegID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.pubs {
		if pubDef.regID == regID {
			pubDef.status = RegistrationStatus.RegisteredMediaDriver
			pubDef.sessionID = sessionID
			pubDef.posLimitCounterID = posLimitCounterID
			pubDef.channelStatusIndicatorID = channelStatusIndicatorID
			pubDef.buffers = logbuffer.Wrap(logFileName)
			pubDef.buffers.IncRef()
			pubDef.origRegID = origRegID

			logger.Debugf("Updated publication: %v", pubDef)

			if cc.onNewPublicationHandler != nil {
				cc.onNewPublicationHandler(pubDef.channel, streamID, sessionID, regID)
			}
		}
	}
}

func (cc *ClientConductor) OnAvailableCounter(correlationID int64, counterID int32) {
	logger.Debugf("OnAvailableCounter: correlationID=%d, counterID=%d",
		correlationID, counterID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	logger.Debug("OnAvailableCounter: Not supported yet")
}

func (cc *ClientConductor) OnUnavailableCounter(correlationID int64, counterID int32) {
	logger.Debugf("OnUnavailableCounter: correlationID=%d, counterID=%d",
		correlationID, counterID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	logger.Debug("OnUnavailableCounter: Not supported yet")
}

func (cc *ClientConductor) OnClientTimeout(clientID int64) {
	logger.Debugf("OnClientTimeout: clientID=%d", clientID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	if clientID == cc.driverProxy.ClientID() {
		errStr := fmt.Sprintf("OnClientTimeout for ClientID:%d", clientID)
		cc.onError(errors.New(errStr))
		cc.running.Set(false)
	}
}

func (cc *ClientConductor) OnSubscriptionReady(correlationID int64, channelStatusIndicatorID int32) {
	logger.Debugf("OnSubscriptionReady: correlationID=%d, channelStatusIndicatorID=%d",
		correlationID, channelStatusIndicatorID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {

		if sub.regID == correlationID {
			sub.status = RegistrationStatus.RegisteredMediaDriver
			sub.subscription = NewSubscription(cc, sub.channel, correlationID, sub.streamID, channelStatusIndicatorID)

			if cc.onNewSubscriptionHandler != nil {
				cc.onNewSubscriptionHandler(sub.channel, sub.streamID, correlationID)
			}
		}
	}

}

//go:norace
func (cc *ClientConductor) OnAvailableImage(streamID int32, sessionID int32, logFilename string, sourceIdentity string,
	subscriberPositionID int32, subsRegID int64, corrID int64) {
	logger.Debugf("OnAvailableImage: streamId=%d, sessionId=%d, logFilename=%s, sourceIdentity=%s, subsRegID=%d, corrID=%d",
		streamID, sessionID, logFilename, sourceIdentity, subsRegID, corrID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {

		// if sub.streamID == streamID && sub.subscription != nil {
		if sub.subscription != nil {
			// logger.Debugf("OnAvailableImage: sub.regID=%d subsRegID=%d corrID=%d %#v", sub.regID, subsRegID, corrID, sub)
			if sub.regID == subsRegID {

				image := NewImage(sessionID, corrID, logbuffer.Wrap(logFilename))
				image.subscriptionRegistrationID = sub.regID
				image.sourceIdentity = sourceIdentity
				image.subscriberPosition = NewPosition(cc.counterValuesBuffer, subscriberPositionID)
				image.exceptionHandler = cc.errorHandler
				logger.Debugf("OnAvailableImage: new image position: %v -> %d",
					image.subscriberPosition, image.subscriberPosition.get())

				sub.subscription.addImage(image)

				if nil != cc.onAvailableImageHandler {
					cc.onAvailableImageHandler(image)
				}
			}
		}
	}
}

func (cc *ClientConductor) OnUnavailableImage(corrID int64, subscriptionRegistrationID int64) {
	logger.Debugf("OnUnavailableImage: corrID=%d subscriptionRegistrationID=%d", corrID, subscriptionRegistrationID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {
		if sub.regID == subscriptionRegistrationID {
			if sub.subscription != nil {
				image := sub.subscription.removeImage(corrID)
				if cc.onUnavailableImageHandler != nil {
					cc.onUnavailableImageHandler(image)
				}
				cc.lingeringResources <- lingerResourse{time.Now().UnixNano(), image}
				runtime.KeepAlive(image)
			}
		}
	}
}

func (cc *ClientConductor) OnOperationSuccess(corrID int64) {
	logger.Debugf("OnOperationSuccess: correlationId=%d", corrID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

}

func (cc *ClientConductor) OnChannelEndpointError(corrID int64, errorMessage string) {
	logger.Debugf("OnChannelEndpointError: correlationID=%d, errorMessage=%s", corrID, errorMessage)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	statusIndicatorId := int32(corrID)

	for _, pubDef := range cc.pubs {
		if pubDef.publication != nil && pubDef.publication.ChannelStatusID() == statusIndicatorId {
			cc.onError(fmt.Errorf(errorMessage))
		}
	}

	for _, subDef := range cc.subs {
		if subDef.subscription != nil && subDef.subscription.ChannelStatusId() == statusIndicatorId {
			cc.onError(fmt.Errorf(errorMessage))
		}
	}
}

func (cc *ClientConductor) OnErrorResponse(corrID int64, errorCode int32, errorMessage string) {
	logger.Debugf("OnErrorResponse: correlationID=%d, errorCode=%d, errorMessage=%s", corrID, errorCode, errorMessage)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.pubs {
		if pubDef.regID == corrID {
			pubDef.status = RegistrationStatus.ErroredMediaDriver
			pubDef.errorCode = errorCode
			pubDef.errorMessage = errorMessage
			return
		}
	}

	for _, subDef := range cc.subs {
		if subDef.regID == corrID {
			subDef.status = RegistrationStatus.ErroredMediaDriver
			subDef.errorCode = errorCode
			subDef.errorMessage = errorMessage
		}
	}
}

func (cc *ClientConductor) onHeartbeatCheckTimeouts() int {
	var result int

	now := time.Now().UnixNano()

	if now > (cc.timeOfLastDoWork + cc.interServiceTimeoutNs) {
		cc.closeAllResources(now)

		err := fmt.Errorf("Timeout between service calls over %d ms (%d > %d + %d) (%d)",
			cc.interServiceTimeoutNs/time.Millisecond.Nanoseconds(),
			now/time.Millisecond.Nanoseconds(),
			cc.timeOfLastDoWork,
			cc.interServiceTimeoutNs/time.Millisecond.Nanoseconds(),
			(now-cc.timeOfLastDoWork)/time.Millisecond.Nanoseconds())
		cc.onError(err)
	}

	cc.timeOfLastDoWork = now

	if now > (cc.timeOfLastKeepalive + keepaliveTimeoutNS) {
		age := cc.driverProxy.TimeOfLastDriverKeepalive()*time.Millisecond.Nanoseconds() + cc.driverTimeoutNs
		if now > age {
			cc.driverActive.Set(false)
			err := fmt.Errorf("MediaDriver keepalive (ms): age=%d > timeout=%d",
				age,
				cc.driverTimeoutNs/time.Millisecond.Nanoseconds(),
			)
			cc.onError(err)
		}

		if cc.heartbeatTimestamp != nil {
			registrationID, ctrErr := cc.counterReader.GetKeyPartInt64(cc.heartbeatTimestamp.CounterId, heartheatRegistrationIdOffset)
			if ctrErr == nil && registrationID == cc.driverProxy.ClientID() {
				cc.heartbeatTimestamp.Set(now / time.Millisecond.Nanoseconds())
			} else {
				cc.closeAllResources(now)
				err := fmt.Errorf("client heartbeat timestamp not active")
				cc.onError(err)
			}
		} else {
			counterId := cc.counterReader.FindCounter(heartbeatTypeId, func(keyBuffer *atomic.Buffer) bool {
				return keyBuffer.GetInt64(heartheatRegistrationIdOffset) == cc.driverProxy.ClientID()
			})
			if counterId != ctr.NullCounterId {
				var ctrErr error
				if cc.heartbeatTimestamp, ctrErr = ctr.NewAtomicCounter(cc.counterReader, counterId); ctrErr != nil {
					logger.Warning("unable to allocate heartbeat counter %d", counterId)
				} else {
					cc.heartbeatTimestamp.Set(now / time.Millisecond.Nanoseconds())
				}
			}
		}

		cc.timeOfLastKeepalive = now
		result = 1
	}

	if now > (cc.timeOfLastCheckManagedResources + resourceTimeoutNS) {
		cc.onCheckManagedResources(now)
		cc.timeOfLastCheckManagedResources = now
		result = 1
	}

	return result
}

func (cc *ClientConductor) onCheckManagedResources(now int64) {
	moreToCheck := true
	for moreToCheck {
		select {
		case r := <-cc.lingeringResources:
			logger.Debugf("Resource to linger: %v", r)
			if cc.resourceLingerTimeoutNs < now-r.lastTime {
				res := r.resource
				logger.Debugf("lingering resource expired(%dms old): %v",
					(now-r.lastTime)/time.Millisecond.Nanoseconds(), res)
				if res != nil {
					err := res.Close()
					if err != nil {
						logger.Warningf("Failed to close lingering resource: %v", err)
						cc.onError(err)
					}
				}
			} else {
				// The assumption is that resources are queued in order
				moreToCheck = false
				// FIXME ..and we're breaking it here, but since there is no peek...
				cc.lingeringResources <- r
			}
		default:
			moreToCheck = false
		}
	}
}

func (cc *ClientConductor) isPublicationConnected(timeOfLastStatusMessage int64) bool {
	return time.Now().UnixNano() <= (timeOfLastStatusMessage*int64(time.Millisecond) + cc.publicationConnectionTimeoutNs)
}

func (cc *ClientConductor) CounterReader() *ctr.Reader {
	return cc.counterReader
}

func (cc *ClientConductor) closeAllResources(now int64) {
	var err error
	if cc.running.CompareAndSet(true, false) {
		for _, pub := range cc.pubs {
			if pub != nil && pub.publication != nil {
				err = pub.publication.Close()
				if err != nil {
					cc.onError(err)
				}
			}
		}
		cc.pubs = nil

		for _, sub := range cc.subs {
			if sub != nil && sub.subscription != nil {
				err = sub.subscription.Close()
				if err != nil {
					cc.onError(err)
				}
			}
		}
		cc.subs = nil
	}
}

func (cc *ClientConductor) onError(err error) {
	if cc.errorHandler != nil {
		logger.Error(err)
		cc.errorHandler(err)
	} else {
		log.Fatal(err)
	}
}
