/*
Copyright 2016 Stanislav Liberman

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
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/driver"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"io"
	"log"
	"sync"
	"time"
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
	keepaliveTimeoutNS int64 = 500 * int64(time.Millisecond)
	resourceTimeoutNS  int64 = 1000 * int64(time.Millisecond)
)

type PublicationStateDefn struct {
	registrationID     int64
	timeOfRegistration int64
	streamID           int32
	sessionID          int32
	posLimitCounterID  int32
	errorCode          int32
	status             int
	channel            string
	errorMessage       string
	buffers            *logbuffer.LogBuffers
	publication        *Publication
}

func (pub *PublicationStateDefn) Init(channel string, registrationID int64, streamID int32, now int64) *PublicationStateDefn {
	pub.channel = channel
	pub.registrationID = registrationID
	pub.streamID = streamID
	pub.sessionID = -1
	pub.posLimitCounterID = -1
	pub.timeOfRegistration = now
	pub.status = RegistrationStatus.AwaitingMediaDriver

	return pub
}

type SubscriptionStateDefn struct {
	registrationID     int64
	timeOfRegistration int64
	streamID           int32
	errorCode          int32
	status             int
	channel            string
	errorMessage       string
	subscription       *Subscription
}

func (sub *SubscriptionStateDefn) Init(ch string, regID int64, sID int32, now int64) *SubscriptionStateDefn {
	sub.channel = ch
	sub.registrationID = regID
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
	pubs []*PublicationStateDefn
	subs []*SubscriptionStateDefn

	driverProxy *driver.Proxy

	counterValuesBuffer *atomic.Buffer

	driverListenerAdapter *driver.ListenerAdapter

	adminLock sync.Mutex

	pendingCloses      map[int64]chan bool
	lingeringResources chan lingerResourse

	onNewPublicationHandler   NewPublicationHandler
	onNewSubscriptionHandler  NewSubscriptionHandler
	onAvailableImageHandler   AvailableImageHandler
	onUnavailableImageHandler UnavailableImageHandler
	errorHandler              func(error)

	running      atomic.Bool
	driverActive atomic.Bool

	timeOfLastKeepalive             int64
	timeOfLastCheckManagedResources int64
	timeOfLastDoWork                int64
	driverTimeoutNs                 int64
	interServiceTimeoutNs           int64
	publicationConnectionTimeoutNs  int64
	resourceLingerTimeoutNs         int64
}

func (cc *ClientConductor) Init(driverProxy *driver.Proxy, bcast *broadcast.CopyReceiver,
	interServiceTimeout time.Duration, driverTimeout time.Duration, pubConnectionTimeout time.Duration,
	lingerTo time.Duration) *ClientConductor {

	logger.Debugf("Initializing ClientConductor with: %v %v %d %d %d", driverProxy, bcast, interServiceTimeout, driverTimeout,
		pubConnectionTimeout)

	cc.driverProxy = driverProxy
	cc.running.Set(true)
	cc.driverActive.Set(true)
	cc.driverListenerAdapter = driver.NewAdapter(cc, bcast)
	cc.interServiceTimeoutNs = interServiceTimeout.Nanoseconds()
	cc.driverTimeoutNs = driverTimeout.Nanoseconds()
	cc.publicationConnectionTimeoutNs = pubConnectionTimeout.Nanoseconds()
	cc.resourceLingerTimeoutNs = lingerTo.Nanoseconds()

	cc.pendingCloses = make(map[int64]chan bool)
	cc.lingeringResources = make(chan lingerResourse, 1024)

	return cc
}

func (cc *ClientConductor) Close() error {

	var err error
	if cc.running.CompareAndSet(true, false) {
		// TODO accumulate errors

		for _, pub := range cc.pubs {
			err = pub.publication.Close()
			// In Go 1.7 should have the following line
			// runtime.KeepAlive(pub.publication)
		}
		cc.pubs = nil

		for _, sub := range cc.subs {
			err = sub.subscription.Close()
			// In Go 1.7 should have the following line
			// runtime.KeepAlive(sub.subscription)
		}
		cc.subs = nil
	}

	return err
}

// Run is the main execution loop of ClientConductor.
func (cc *ClientConductor) Run(idleStrategy idlestrategy.Idler) {

	now := time.Now().UnixNano()
	cc.timeOfLastKeepalive = now
	cc.timeOfLastCheckManagedResources = now
	cc.timeOfLastDoWork = now

	// In Go 1.7 should have the following line
	// runtime.LockOSThread()

	// Clean exit from this particular go routine
	defer func() {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("Panic: %v", err)
			logger.Error(errStr)
			cc.errorHandler(errors.New(errStr))
			cc.running.Set(false)
		}
	}()

	for cc.running.Get() {
		workCount := cc.driverListenerAdapter.ReceiveMessages()
		workCount += cc.onHeartbeatCheckTimeouts()
		idleStrategy.Idle(workCount)
	}
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

	for _, pub := range cc.pubs {
		if pub.streamID == streamID && pub.channel == channel {
			return pub.registrationID
		}
	}

	now := time.Now().UnixNano()

	registrationID := cc.driverProxy.AddPublication(channel, streamID)

	pubState := new(PublicationStateDefn)
	pubState.Init(channel, registrationID, streamID, now)

	cc.pubs = append(cc.pubs, pubState)

	return registrationID
}

func (cc *ClientConductor) FindPublication(registrationID int64) *Publication {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var publication *Publication
	for _, pub := range cc.pubs {
		if pub.registrationID == registrationID {
			if pub.publication != nil {
				publication = pub.publication
			} else {
				switch pub.status {
				case RegistrationStatus.AwaitingMediaDriver:
					if now := time.Now().UnixNano(); now > (pub.timeOfRegistration + cc.driverTimeoutNs) {
						logger.Errorf("No response from driver. started: %d, now: %d, to: %d",
							pub.timeOfRegistration, now/time.Millisecond.Nanoseconds(),
							cc.driverTimeoutNs/time.Millisecond.Nanoseconds())
						log.Panic(fmt.Sprintf("No response from driver on %v of %v", pub, cc.pubs))
					}
				case RegistrationStatus.RegisteredMediaDriver:
					publication = NewPublication(pub.buffers)
					publication.conductor = cc
					publication.channel = pub.channel
					publication.registrationID = registrationID
					publication.streamID = pub.streamID
					publication.sessionID = pub.sessionID
					publication.publicationLimit = NewPosition(cc.counterValuesBuffer,
						pub.posLimitCounterID)

				case RegistrationStatus.ErroredMediaDriver:
					log.Fatalf("Error on %d: %d: %s", registrationID, pub.errorCode, pub.errorMessage)
				}
			}
			break
		}
	}

	return publication
}

func (cc *ClientConductor) releasePublication(registrationID int64) chan bool {
	logger.Debugf("ReleasePublication: registrationId=%d", registrationID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	ch := make(chan bool, 1)
	found := false

	for i, pub := range cc.pubs {
		if pub.registrationID == registrationID {
			corrID := cc.driverProxy.RemovePublication(registrationID)

			cc.pubs[i] = cc.pubs[len(cc.pubs)-1]
			cc.pubs[len(cc.pubs)-1] = nil
			cc.pubs = cc.pubs[:len(cc.pubs)-1]

			cc.pendingCloses[corrID] = ch
			found = true
		}
	}

	// Need to report if it's already been closed
	if !found {
		ch <- false
	}

	return ch
}

// AddSubscription sends the add subscription command through the driver proxy
func (cc *ClientConductor) AddSubscription(channel string, streamID int32) int64 {
	logger.Debugf("AddSubscription: channel=%s, streamId=%d", channel, streamID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	registrationID := cc.driverProxy.AddSubscription(channel, streamID)

	subState := new(SubscriptionStateDefn)
	subState.Init(channel, registrationID, streamID, now)

	cc.subs = append(cc.subs, subState)

	return registrationID
}

func (cc *ClientConductor) FindSubscription(registrationID int64) *Subscription {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var subscription *Subscription
	for _, sub := range cc.subs {
		if sub.registrationID == registrationID {

			switch sub.status {
			case RegistrationStatus.AwaitingMediaDriver:
				if now := time.Now().UnixNano(); now > (sub.timeOfRegistration + cc.driverTimeoutNs) {
					logger.Errorf("No response from driver. started: %d, now: %d, to: %d",
						sub.timeOfRegistration, now/time.Millisecond.Nanoseconds(),
						cc.driverTimeoutNs/time.Millisecond.Nanoseconds())
					log.Fatalf("No response on %v of %v", sub, cc.subs)
				}
			case RegistrationStatus.ErroredMediaDriver:
				errStr := fmt.Sprintf("Error on %d: %d: %s", registrationID, sub.errorCode, sub.errorMessage)
				cc.errorHandler(errors.New(errStr))
				log.Fatalf(errStr)
			}

			subscription = sub.subscription
			break
		}
	}

	return subscription
}

func (cc *ClientConductor) releaseSubscription(registrationID int64, images []*Image) chan bool {
	logger.Debugf("ReleaseSubscription: registrationId=%d", registrationID)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	ch := make(chan bool, 1)
	found := false

	for i, sub := range cc.subs {
		if sub.registrationID == registrationID {
			logger.Debugf("Removing subscription: %d; %v", registrationID, images)
			corrID := cc.driverProxy.RemoveSubscription(registrationID)

			cc.subs[i] = cc.subs[len(cc.subs)-1]
			cc.subs[len(cc.subs)-1] = nil
			cc.subs = cc.subs[:len(cc.subs)-1]

			for _, image := range images {
				if cc.onUnavailableImageHandler != nil {
					cc.onUnavailableImageHandler(image)
				}
				err := image.Close()
				if err != nil {
					logger.Warningf("Failed to close subscription: %d", registrationID)
					cc.errorHandler(err)
				}
				cc.lingeringResources <- lingerResourse{now, *image}
			}

			cc.pendingCloses[corrID] = ch
			found = true
		}
	}

	// Need to report if it's already been closed
	if !found {
		ch <- false
	}

	return ch
}

func (cc *ClientConductor) OnNewPublication(streamID int32, sessionID int32, positionLimitCounterID int32,
	logFileName string, registrationID int64) {
	logger.Debugf("OnNewPublication: streamId=%d, sessionId=%d, positionLimitCounterId=%d, logFileName=%s, registrationId=%d",
		streamID, sessionID, positionLimitCounterID, logFileName, registrationID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.pubs {
		if pubDef.registrationID == registrationID {
			pubDef.status = RegistrationStatus.RegisteredMediaDriver
			pubDef.sessionID = sessionID
			pubDef.posLimitCounterID = positionLimitCounterID
			pubDef.buffers = logbuffer.Wrap(logFileName)

			logger.Debugf("Updated publication: %v", pubDef)

			if cc.onNewPublicationHandler != nil {
				cc.onNewPublicationHandler(pubDef.channel, streamID, sessionID, registrationID)
			}
		}
	}
}

func (cc *ClientConductor) OnAvailableImage(streamID int32, sessionID int32, logFilename string,
	sourceIdentity string, subscriberPositionCount int, subscriberPositions []driver.SubscriberPosition,
	correlationID int64) {
	logger.Debugf("OnAvailableImage: streamId=%d, sessionId=%d, logFilename=%s, sourceIdentity=%s, subscriberPositions=%v, correlationId=%d",
		streamID, sessionID, logFilename, sourceIdentity, subscriberPositions, correlationID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {
		if sub.streamID == streamID && sub.subscription != nil {
			if !sub.subscription.hasImage(sessionID) {
				for _, subPos := range subscriberPositions {
					if sub.registrationID == subPos.RegistrationID() {

						image := NewImage(sessionID, correlationID, logbuffer.Wrap(logFilename))
						image.subscriptionRegistrationID = sub.registrationID
						image.sourceIdentity = sourceIdentity
						image.subscriberPosition = NewPosition(cc.counterValuesBuffer,
							subPos.IndicatorID())
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
	}
}

func (cc *ClientConductor) OnUnavailableImage(streamID int32, correlationID int64) {
	logger.Debugf("OnUnavailableImage: streamId=%d, correlationId=%d", streamID, correlationID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {
		if sub.streamID == streamID {
			if sub.subscription != nil {
				image := sub.subscription.removeImage(correlationID)
				// FIXME Howzzat?!
				if nil != image {
					cc.lingeringResources <- lingerResourse{time.Now().UnixNano(), *image}
				}
			}
		}
	}
}

func (cc *ClientConductor) OnOperationSuccess(corrID int64) {
	logger.Debugf("OnOperationSuccess: correlationId=%d", corrID)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {
		if sub.registrationID == corrID && sub.status == RegistrationStatus.AwaitingMediaDriver {
			sub.status = RegistrationStatus.RegisteredMediaDriver
			sub.subscription = NewSubscription(cc, sub.channel, corrID, sub.streamID)

			if cc.onNewSubscriptionHandler != nil {
				cc.onNewSubscriptionHandler(sub.channel, sub.streamID, corrID)
			}
		}
	}

	if cc.pendingCloses[corrID] != nil {
		cc.pendingCloses[corrID] <- true
		close(cc.pendingCloses[corrID])
		delete(cc.pendingCloses, corrID)
	}
}

func (cc *ClientConductor) OnErrorResponse(correlationID int64, errorCode int32, errorMessage string) {
	logger.Debugf("OnErrorResponse: correlationId=%d, errorCode=%d, errorMessage=%s", correlationID, errorCode, errorMessage)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.pubs {
		if pubDef.registrationID == correlationID {
			pubDef.status = RegistrationStatus.ErroredMediaDriver
			pubDef.errorCode = errorCode
			pubDef.errorMessage = errorMessage
			return
		}
	}

	for _, subDef := range cc.pubs {
		if subDef.registrationID == correlationID {
			subDef.status = RegistrationStatus.ErroredMediaDriver
			subDef.errorCode = errorCode
			subDef.errorMessage = errorMessage
		}
	}
}

func (cc *ClientConductor) onInterServiceTimeout(now int64) {
	log.Printf("onInterServiceTimeout: now=%d", now)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	err := cc.Close()
	if err != nil {
		logger.Warningf("Failed to close client conductor: %v", err)
		cc.errorHandler(err)
	}
}

func (cc *ClientConductor) onHeartbeatCheckTimeouts() int {
	var result int

	now := time.Now().UnixNano()

	if now > (cc.timeOfLastDoWork + cc.interServiceTimeoutNs) {
		cc.onInterServiceTimeout(now)

		log.Fatalf("Timeout between service calls over %d ms (%d > %d + %d) (%d)",
			cc.interServiceTimeoutNs/time.Millisecond.Nanoseconds(),
			now/time.Millisecond.Nanoseconds(),
			cc.timeOfLastDoWork,
			cc.interServiceTimeoutNs/time.Millisecond.Nanoseconds(),
			(now-cc.timeOfLastDoWork)/time.Millisecond.Nanoseconds())
	}

	cc.timeOfLastDoWork = now

	if now > (cc.timeOfLastKeepalive + keepaliveTimeoutNS) {
		cc.driverProxy.SendClientKeepalive()

		hbTime := cc.driverProxy.TimeOfLastDriverKeepalive() * time.Millisecond.Nanoseconds()
		if now > (hbTime + cc.driverTimeoutNs) {
			cc.driverActive.Set(false)

			log.Fatalf("Driver has been inactive for over %d ms",
				cc.driverTimeoutNs/time.Millisecond.Nanoseconds())
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
						cc.errorHandler(err)
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
