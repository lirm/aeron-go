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
	"fmt"
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/driver"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"io"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var RegistrationStatus = struct {
	AWAITING_MEDIA_DRIVER   int
	REGISTERED_MEDIA_DRIVER int
	ERRORED_MEDIA_DRIVER    int
}{
	0,
	1,
	2,
}

const (
	KEEPALIVE_TIMEOUT_NS int64 = 500 * int64(time.Millisecond)
	RESOURCE_TIMEOUT_NS  int64 = 1000 * int64(time.Millisecond)
)

type PublicationStateDefn struct {
	channel            string
	registrationId     int64
	streamId           int32
	sessionId          int32
	posLimitCounterId  int32
	timeOfRegistration int64
	status             int
	errorCode          int32
	errorMessage       string
	buffers            *logbuffer.LogBuffers
	publication        *Publication
}

func (pub *PublicationStateDefn) Init(channel string, registrationId int64, streamId int32, now int64) *PublicationStateDefn {
	pub.channel = channel
	pub.registrationId = registrationId
	pub.streamId = streamId
	pub.sessionId = -1
	pub.posLimitCounterId = -1
	pub.timeOfRegistration = now
	pub.status = RegistrationStatus.AWAITING_MEDIA_DRIVER

	return pub
}

type SubscriptionStateDefn struct {
	channel            string
	registrationId     int64
	streamId           int32
	timeOfRegistration int64
	status             int
	errorCode          int32
	errorMessage       string
	subscriptionCache  *Subscription
	subscription       *Subscription
}

func (sub *SubscriptionStateDefn) Init(ch string, regId int64, sId int32, now int64) *SubscriptionStateDefn {
	sub.channel = ch
	sub.registrationId = regId
	sub.streamId = sId
	sub.timeOfRegistration = now
	sub.status = RegistrationStatus.AWAITING_MEDIA_DRIVER

	return sub
}

type LingerResourse struct {
	lastTime int64
	resource io.Closer
}

type ClientConductor struct {
	pubs []*PublicationStateDefn
	subs []*SubscriptionStateDefn

	driverProxy *driver.Proxy

	counterValuesBuffer *buffers.Atomic

	driverListenerAdapter *driver.ListenerAdapter

	adminLock sync.Mutex

	lingeringResources chan LingerResourse

	onNewPublicationHandler   NewPublicationHandler
	onNewSubscriptionHandler  NewSubscriptionHandler
	onAvailableImageHandler   AvailableImageHandler
	onUnavailableImageHandler UnavailableImageHandler
	errorHandler              func(error)

	running      atomic.Value
	driverActive atomic.Value

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
	cc.running.Store(true)
	cc.driverActive.Store(true)
	cc.driverListenerAdapter = driver.NewAdapter(cc, bcast)
	cc.interServiceTimeoutNs = interServiceTimeout.Nanoseconds()
	cc.driverTimeoutNs = driverTimeout.Nanoseconds()
	cc.publicationConnectionTimeoutNs = pubConnectionTimeout.Nanoseconds()
	cc.resourceLingerTimeoutNs = lingerTo.Nanoseconds()

	cc.lingeringResources = make(chan LingerResourse, 1024)

	return cc
}

func (cc *ClientConductor) Close() error {
	cc.running.Store(false)

	// TODO accumulate errors
	var err error

	for _, pub := range cc.pubs {
		err = pub.publication.Close()
		// In Go 1.7 should have the following line
		// runtime.KeepAlice(pub.publication)
	}
	cc.pubs = nil

	for _, sub := range cc.subs {
		err = sub.subscription.Close()
		// In Go 1.7 should have the following line
		// runtime.KeepAlice(sub.subscription)
	}
	cc.subs = nil

	return err
}

func (cc *ClientConductor) Run(idleStrategy idlestrategy.Idler) {

	cc.timeOfLastKeepalive = time.Now().UnixNano()
	cc.timeOfLastCheckManagedResources = time.Now().UnixNano()
	cc.timeOfLastDoWork = time.Now().UnixNano()

	// In Go 1.7 should have the following line
	// runtime.LockOSThread()

	for cc.running.Load().(bool) {
		workCount := cc.driverListenerAdapter.ReceiveMessages()
		workCount += cc.onHeartbeatCheckTimeouts()
		idleStrategy.Idle(workCount)
	}
}

func (cc *ClientConductor) verifyDriverIsActive() {
	if !cc.driverActive.Load().(bool) {
		log.Fatal("Driver is not active")
	}
}

func (cc *ClientConductor) AddPublication(channel string, streamId int32) int64 {
	logger.Debugf("AddPublication: channel=%s, streamId=%d", channel, streamId)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pub := range cc.pubs {
		if pub.streamId == streamId && pub.channel == channel {
			return pub.registrationId
		}
	}

	now := time.Now().UnixNano()

	registrationId := cc.driverProxy.AddPublication(channel, streamId)

	pubState := new(PublicationStateDefn)
	pubState.Init(channel, registrationId, streamId, now)

	cc.pubs = append(cc.pubs, pubState)

	return registrationId
}

func (cc *ClientConductor) FindPublication(registrationId int64) *Publication {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var publication *Publication = nil
	for _, pub := range cc.pubs {
		if pub.registrationId == registrationId {
			if pub.publication != nil {
				publication = pub.publication
			} else {
				switch pub.status {
				case RegistrationStatus.AWAITING_MEDIA_DRIVER:
					now := time.Now().UnixNano()
					if now > (pub.timeOfRegistration + cc.driverTimeoutNs) {
						logger.Errorf("No response from driver. started: %d, now: %d, to: %d",
							pub.timeOfRegistration, now/time.Millisecond.Nanoseconds(),
							cc.driverTimeoutNs/time.Millisecond.Nanoseconds())
						log.Panic(fmt.Sprintf("No response from driver on %v of %v", pub, cc.pubs))
					}
				case RegistrationStatus.REGISTERED_MEDIA_DRIVER:
					publication = NewPublication(pub.buffers)
					publication.conductor = cc
					publication.channel = pub.channel
					publication.registrationId = registrationId
					publication.streamId = pub.streamId
					publication.sessionId = pub.sessionId
					publication.publicationLimit = buffers.NewPosition(cc.counterValuesBuffer,
						pub.posLimitCounterId)

				case RegistrationStatus.ERRORED_MEDIA_DRIVER:
					log.Fatalf("Error on %d: %d: %s", registrationId, pub.errorCode, pub.errorMessage)
				}
			}
			break
		}
	}

	return publication
}

func (cc *ClientConductor) ReleasePublication(registrationId int64) {
	logger.Debugf("ReleasePublication: registrationId=%d", registrationId)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for i, pub := range cc.pubs {
		if pub.registrationId == registrationId {
			cc.driverProxy.RemovePublication(registrationId)

			cc.pubs[i] = cc.pubs[len(cc.pubs)-1]
			cc.pubs[len(cc.pubs)-1] = nil
			cc.pubs = cc.pubs[:len(cc.pubs)-1]
		}
	}
}

func (cc *ClientConductor) AddSubscription(channel string, streamId int32) int64 {
	logger.Debugf("AddSubscription: channel=%s, streamId=%d", channel, streamId)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	registrationId := cc.driverProxy.AddSubscription(channel, streamId)

	subState := new(SubscriptionStateDefn)
	subState.status = RegistrationStatus.AWAITING_MEDIA_DRIVER
	subState.registrationId = registrationId
	subState.channel = channel
	subState.streamId = streamId
	subState.timeOfRegistration = now

	cc.subs = append(cc.subs, subState)

	return registrationId
}

func (cc *ClientConductor) FindSubscription(registrationId int64) *Subscription {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var subscription *Subscription = nil
	for _, sub := range cc.subs {
		if sub.registrationId == registrationId {

			switch sub.status {
			case RegistrationStatus.AWAITING_MEDIA_DRIVER:
				now := time.Now().UnixNano()
				if now > (sub.timeOfRegistration + cc.driverTimeoutNs) {
					logger.Errorf("No response from driver. started: %d, now: %d, to: %d",
						sub.timeOfRegistration, now/time.Millisecond.Nanoseconds(),
						cc.driverTimeoutNs/time.Millisecond.Nanoseconds())
					log.Fatalf("No response on %v of %v", sub, cc.subs)
				}
			case RegistrationStatus.ERRORED_MEDIA_DRIVER:
				log.Fatalf("Error on %d: %d: %s", registrationId, sub.errorCode, sub.errorMessage)
			}

			subscription = sub.subscription
			break
		}
	}

	return subscription
}

func (cc *ClientConductor) ReleaseSubscription(registrationId int64, images []*Image) {
	logger.Debugf("ReleaseSubscription: registrationId=%d", registrationId)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano()

	for i, sub := range cc.subs {
		if sub.registrationId == registrationId {
			logger.Debugf("Removing subscription: %d; %v", registrationId, images)
			cc.driverProxy.RemoveSubscription(registrationId)

			cc.subs[i] = cc.subs[len(cc.subs)-1]
			cc.subs[len(cc.subs)-1] = nil
			cc.subs = cc.subs[:len(cc.subs)-1]

			for _, image := range images {
				if cc.onUnavailableImageHandler != nil {
					cc.onUnavailableImageHandler(image)
				}
				image.Close()
				cc.lingeringResources <- LingerResourse{now, *image}
			}
		}
	}
}

func (cc *ClientConductor) OnNewPublication(streamId int32, sessionId int32, positionLimitCounterId int32,
	logFileName string, registrationId int64) {
	logger.Debugf("OnNewPublication: streamId=%d, sessionId=%d, positionLimitCounterId=%d, logFileName=%s, registrationId=%d",
		streamId, sessionId, positionLimitCounterId, logFileName, registrationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.pubs {
		if pubDef.registrationId == registrationId {
			pubDef.status = RegistrationStatus.REGISTERED_MEDIA_DRIVER
			pubDef.sessionId = sessionId
			pubDef.posLimitCounterId = positionLimitCounterId
			pubDef.buffers = logbuffer.Wrap(logFileName)

			logger.Debugf("Updated publication: %v", pubDef)

			if cc.onNewPublicationHandler != nil {
				cc.onNewPublicationHandler(pubDef.channel, streamId, sessionId, registrationId)
			}
		}
	}
}

func (cc *ClientConductor) OnAvailableImage(streamId int32, sessionId int32, logFilename string,
	sourceIdentity string, subscriberPositionCount int, subscriberPositions []driver.SubscriberPosition,
	correlationId int64) {
	logger.Debugf("OnAvailableImage: streamId=%d, sessionId=%d, logFilename=%s, sourceIdentity=%s, subscriberPositions=%v, correlationId=%d",
		streamId, sessionId, logFilename, sourceIdentity, subscriberPositions, correlationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {
		if sub.streamId == streamId && sub.subscription != nil {
			if !sub.subscription.hasImage(sessionId) {
				for _, subPos := range subscriberPositions {
					if sub.registrationId == subPos.RegistrationId() {

						image := NewImage(sessionId, correlationId, logbuffer.Wrap(logFilename))
						image.subscriptionRegistrationId = sub.registrationId
						image.sourceIdentity = sourceIdentity
						image.subscriberPosition = buffers.NewPosition(cc.counterValuesBuffer,
							subPos.IndicatorId())
						image.exceptionHandler = cc.errorHandler
						logger.Debugf("OnAvailableImage: new image position: %v -> %d",
							image.subscriberPosition, image.subscriberPosition.Get())

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

func (cc *ClientConductor) OnUnavailableImage(streamId int32, correlationId int64) {
	logger.Debugf("OnUnavailableImage: streamId=%d, correlationId=%d", streamId, correlationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {
		if sub.streamId == streamId {
			image := sub.subscription.removeImage(correlationId)
			if nil != image {
				cc.lingeringResources <- LingerResourse{time.Now().UnixNano(), *image}
			}
		}
	}
}

func (cc *ClientConductor) OnOperationSuccess(correlationId int64) {
	logger.Debugf("OnOperationSuccess: correlationId=%d", correlationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subs {
		if sub.registrationId == correlationId && sub.status == RegistrationStatus.AWAITING_MEDIA_DRIVER {
			sub.status = RegistrationStatus.REGISTERED_MEDIA_DRIVER
			sub.subscription = NewSubscription(cc, sub.channel, correlationId, sub.streamId)
		}
	}
}

func (cc *ClientConductor) OnErrorResponse(correlationId int64, errorCode int32, errorMessage string) {
	logger.Debugf("OnErrorResponse: correlationId=%d, errorCode=%d, errorMessage=%s", correlationId, errorCode, errorMessage)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.pubs {
		if pubDef.registrationId == correlationId {
			pubDef.status = RegistrationStatus.ERRORED_MEDIA_DRIVER
			pubDef.errorCode = errorCode
			pubDef.errorMessage = errorMessage
			return
		}
	}

	for _, subDef := range cc.pubs {
		if subDef.registrationId == correlationId {
			subDef.status = RegistrationStatus.ERRORED_MEDIA_DRIVER
			subDef.errorCode = errorCode
			subDef.errorMessage = errorMessage
		}
	}
}

func (cc *ClientConductor) onInterServiceTimeout(now int64) {
	log.Printf("onInterServiceTimeout: now=%d", now)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	cc.Close()
}

func (cc *ClientConductor) onHeartbeatCheckTimeouts() int {
	var result int = 0

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

	if now > (cc.timeOfLastKeepalive + KEEPALIVE_TIMEOUT_NS) {
		cc.driverProxy.SendClientKeepalive()

		hbTime := cc.driverProxy.TimeOfLastDriverKeepalive() * time.Millisecond.Nanoseconds()
		if now > (hbTime + cc.driverTimeoutNs) {
			cc.driverActive.Store(false)

			log.Fatalf("Driver has been inactive for over %d ms",
				cc.driverTimeoutNs/time.Millisecond.Nanoseconds())
		}

		cc.timeOfLastKeepalive = now
		result = 1
	}

	if now > (cc.timeOfLastCheckManagedResources + RESOURCE_TIMEOUT_NS) {
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
					res.Close()
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
