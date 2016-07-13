package aeron

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"log"
	"sync"
	"time"
	"github.com/lirm/aeron-go/aeron/driver"
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
	KEEPALIVE_TIMEOUT_MS int64 = 500
	RESOURCE_TIMEOUT_MS  int64 = 1000
)

type PublicationStateDefn struct {
	channel                string
	registrationId         int64
	streamId               int32
	sessionId              int32
	positionLimitCounterId int32
	timeOfRegistration     int64
	status                 int
	errorCode              int32
	errorMessage           string
	buffers                *logbuffer.LogBuffers
	publication            *Publication
}

func (pub *PublicationStateDefn) Init(channel string, registrationId int64, streamId int32, now int64) *PublicationStateDefn {
	pub.channel = channel
	pub.registrationId = registrationId
	pub.streamId = streamId
	pub.sessionId = -1
	pub.positionLimitCounterId = -1
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

type ClientConductor struct {
	publications  []*PublicationStateDefn
	subscriptions []*SubscriptionStateDefn

	driverProxy *driver.Proxy

	counterValuesBuffer *buffers.Atomic

	driverListenerAdapter *driver.ListenerAdapter

	adminLock sync.Mutex

	/*
		epoch_clock_t m_epochClock;
	*/
	onNewPublicationHandler   NewPublicationHandler
	onNewSubscriptionHandler  NewSubscriptionHandler
	onAvailableImageHandler   AvailableImageHandler
	onUnavailableImageHandler UnavailableImageHandler
	errorHandler              func(error)

	running                         bool
	driverActive                    bool // atomic
	timeOfLastKeepalive             int64
	timeOfLastCheckManagedResources int64
	timeOfLastDoWork                int64
	driverTimeoutMs                 int64
	interServiceTimeoutMs           int64
	publicationConnectionTimeoutMs  int64
}

func (cc *ClientConductor) Init(driverProxy *driver.Proxy, bcast *broadcast.CopyReceiver, interServiceTimeoutNs int64, driverTimeoutNs int64) *ClientConductor {
	cc.driverProxy = driverProxy
	cc.driverActive = true
	cc.driverListenerAdapter = driver.NewAdapter(cc, bcast)
	cc.interServiceTimeoutMs = interServiceTimeoutNs / 1000000
	cc.driverTimeoutMs = driverTimeoutNs / 1000000

	return cc
}

func (cc *ClientConductor) Run() {

	cc.timeOfLastKeepalive = time.Now().UnixNano() / 1000000
	cc.timeOfLastCheckManagedResources = time.Now().UnixNano() / 1000000
	cc.timeOfLastDoWork = time.Now().UnixNano() / 1000000

	// FIXME needs a boolean flag for graceful shutdown
	for {
		workCount := cc.driverListenerAdapter.ReceiveMessages()
		workCount += cc.onHeartbeatCheckTimeouts()
		// FIXME this needs to happen!
		// idleStrategy.idle(workCount)
		time.Sleep(time.Microsecond)
	}
}

func (cc *ClientConductor) verifyDriverIsActive() {
	if !cc.driverActive {
		log.Fatal("Driver is not active")
	}
}

func (cc *ClientConductor) AddPublication(channel string, streamId int32) int64 {

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pub := range cc.publications {
		if pub.streamId == streamId && pub.channel == channel {
			return pub.registrationId
		}
	}

	now := time.Now().UnixNano() / 1000000

	registrationId := cc.driverProxy.AddPublication(channel, streamId)

	pubState := new(PublicationStateDefn)
	pubState.Init(channel, registrationId, streamId, now)

	cc.publications = append(cc.publications, pubState)

	return registrationId
}

func (cc *ClientConductor) FindPublication(registrationId int64) *Publication {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var publication *Publication = nil
	for _, pub := range cc.publications {
		if pub.registrationId == registrationId {
			if pub.publication != nil {
				publication = pub.publication
			} else {
				switch pub.status {
				case RegistrationStatus.AWAITING_MEDIA_DRIVER:
					now := time.Now().UnixNano() / 1000000
					if now > (pub.timeOfRegistration + cc.driverTimeoutMs) {
						log.Panic(fmt.Sprintf("No response on %v of %v", pub, cc.publications))
						panic(fmt.Sprintf("No response from driver in %d ms", cc.driverTimeoutMs))
					}
				case RegistrationStatus.REGISTERED_MEDIA_DRIVER:
					publication = new(Publication)
					publication.Init(pub.buffers)
					publication.conductor = cc
					publication.channel = pub.channel
					publication.registrationId = registrationId
					publication.streamId = pub.streamId
					publication.sessionId = pub.sessionId
					publication.publicationLimit = buffers.Position{
						cc.counterValuesBuffer,
						pub.positionLimitCounterId,
						counters.Offset(pub.positionLimitCounterId)}

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

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subscriptions {
		if sub.registrationId == registrationId {
			cc.driverProxy.RemovePublication(registrationId)

			// TODO remove from array
		}
	}
}

func (cc *ClientConductor) AddSubscription(channel string, streamId int32) int64 {

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	now := time.Now().UnixNano() / 1000000

	registrationId := cc.driverProxy.AddSubscription(channel, streamId)

	subState := new(SubscriptionStateDefn)
	subState.status = RegistrationStatus.AWAITING_MEDIA_DRIVER
	subState.registrationId = registrationId
	subState.channel = channel
	subState.streamId = streamId
	subState.timeOfRegistration = now

	cc.subscriptions = append(cc.subscriptions, subState)

	return registrationId
}

func (cc *ClientConductor) FindSubscription(registrationId int64) *Subscription {

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	var subscription *Subscription = nil
	for _, sub := range cc.subscriptions {
		if sub.registrationId == registrationId {

			switch sub.status {
			case RegistrationStatus.AWAITING_MEDIA_DRIVER:
				now := time.Now().UnixNano() / 1000000
				if now > (sub.timeOfRegistration + cc.driverTimeoutMs) {
					log.Fatalf("No response on %v of %v", sub, cc.publications)
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

func (cc *ClientConductor) ReleaseSubscription(registrationId int64, images []*Image, imagesLength int) {
	log.Printf("ReleaseSubscription: registrationId=%d", registrationId)

	cc.verifyDriverIsActive()

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subscriptions {
		if sub.registrationId == registrationId {
			cc.driverProxy.RemoveSubscription(registrationId)

			// TODO remove from array

			if cc.onUnavailableImageHandler != nil {
				for _, image := range images {
					cc.onUnavailableImageHandler(image)
				}
			}

			//lingerResources(m_epochClock(), images, imagesLength);
		}
	}
}

func (cc *ClientConductor) OnNewPublication(streamId int32, sessionId int32, positionLimitCounterId int32,
	logFileName string, registrationId int64) {
	log.Printf("OnNewPublication: streamId=%d, sessionId=%d, positionLimitCounterId=%d, logFileName=%s, registrationId=%d",
		streamId, sessionId, positionLimitCounterId, logFileName, registrationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.publications {
		if pubDef.registrationId == registrationId {
			pubDef.status = RegistrationStatus.REGISTERED_MEDIA_DRIVER
			pubDef.sessionId = sessionId
			pubDef.positionLimitCounterId = positionLimitCounterId
			pubDef.buffers = logbuffer.Wrap(logFileName)

			log.Printf("Updated publication: %v", pubDef)

			if cc.onNewPublicationHandler != nil {
				cc.onNewPublicationHandler(pubDef.channel, streamId, sessionId, registrationId)
			}
		}
	}
}

func (cc *ClientConductor) OnAvailableImage(streamId int32, sessionId int32, logFilename string,
	sourceIdentity string, subscriberPositionCount int32, subscriberPositions []driver.SubscriberPosition,
	correlationId int64) {
	log.Printf("OnAvailableImage: streamId=%d, sessionId=%d, logFilename=%s, sourceIdentity=%s, subscriberPositions=%v, correlationId=%d",
		streamId, sessionId, logFilename, sourceIdentity, subscriberPositions, correlationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subscriptions {
		if sub.streamId == streamId && sub.subscription != nil {
			if !sub.subscription.hasImage(sessionId) {
				for _, subPos := range subscriberPositions {
					if sub.registrationId == subPos.RegistrationId() {

						image := NewImage(sessionId, correlationId, logbuffer.Wrap(logFilename))
						image.subscriptionRegistrationId = sub.registrationId
						image.sourceIdentity = sourceIdentity
						image.subscriberPosition = buffers.Position{
							cc.counterValuesBuffer,
							subPos.IndicatorId(),
							counters.Offset(subPos.IndicatorId())}
						image.exceptionHandler = cc.errorHandler

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
	log.Printf("OnUnavailableImage: streamId=%d, correlationId=%d", streamId, correlationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subscriptions {
		if sub.streamId == streamId {
			sub.subscription.removeImage(correlationId)
		}
	}
}

func (cc *ClientConductor) OnOperationSuccess(correlationId int64) {
	log.Printf("OnOperationSuccess: correlationId=%d", correlationId)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, sub := range cc.subscriptions {
		if sub.registrationId == correlationId && sub.status == RegistrationStatus.AWAITING_MEDIA_DRIVER {
			sub.status = RegistrationStatus.REGISTERED_MEDIA_DRIVER
			sub.subscription = new(Subscription)
			sub.subscription.conductor = cc
			sub.subscription.channel = sub.channel
			sub.subscription.registrationId = correlationId
			sub.subscription.streamId = sub.streamId
		}
	}
}

func (cc *ClientConductor) OnErrorResponse(correlationId int64, errorCode int32, errorMessage string) {
	log.Printf("OnErrorResponse: correlationId=%d, errorCode=%d, errorMessage=%s", correlationId, errorCode, errorMessage)

	cc.adminLock.Lock()
	defer cc.adminLock.Unlock()

	for _, pubDef := range cc.publications {
		if pubDef.registrationId == correlationId {
			pubDef.status = RegistrationStatus.ERRORED_MEDIA_DRIVER
			pubDef.errorCode = errorCode
			pubDef.errorMessage = errorMessage
			return
		}
	}

	for _, subDef := range cc.publications {
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
	/*
	   std::for_each(cc.publications.begin(), cc.publications.end(),
	       [&](PublicationStateDefn& entry)
	       {
	           std::shared_ptr<Publication> pub = entry.m_publication.lock();

	           if (nullptr != pub)
	           {
	               pub->close();
	           }
	       });

	   cc.publications.clear();

	   std::for_each(cc.subscriptions.begin(), cc.subscriptions.end(),
	       [&](SubscriptionStateDefn& entry)
	       {
	           std::shared_ptr<Subscription> sub = entry.m_subscription.lock();

	           if (nullptr != sub)
	           {
	               std::pair<Image *, int> removeResult = sub->removeAndCloseAllImages();
	               Image* images = removeResult.first;
	               const int imagesLength = removeResult.second;
	           }
	       });

	   cc.subscriptions.clear();
	*/
}

func (cc *ClientConductor) onHeartbeatCheckTimeouts() int {
	var result int = 0

	now := time.Now().UnixNano() / 1000000

	if now > (cc.timeOfLastDoWork + cc.interServiceTimeoutMs) {
		cc.onInterServiceTimeout(now)

		log.Fatalf("Timeout between service calls over %d ms (%d > %d + %d) (%d)",
			cc.interServiceTimeoutMs,
			now,
			cc.timeOfLastDoWork,
			cc.interServiceTimeoutMs,
			now-cc.timeOfLastDoWork)
	}

	cc.timeOfLastDoWork = now

	if now > (cc.timeOfLastKeepalive + KEEPALIVE_TIMEOUT_MS) {
		cc.driverProxy.SendClientKeepalive()

		hbTime := cc.driverProxy.TimeOfLastDriverKeepalive()
		if now > (hbTime + cc.driverTimeoutMs) {
			cc.driverActive = false

			log.Fatalf("Driver has been inactive for over %d ms", cc.driverTimeoutMs)
		}

		cc.timeOfLastKeepalive = now
		result = 1
	}

	if now > (cc.timeOfLastCheckManagedResources + RESOURCE_TIMEOUT_MS) {
		cc.timeOfLastCheckManagedResources = now
		result = 1
	}

	return result
}

func (cc *ClientConductor) isPublicationConnected(timeOfLastStatusMessage int64) bool {
	now := time.Now().UnixNano() / 1000000
	return now <= (timeOfLastStatusMessage + cc.publicationConnectionTimeoutMs)
}
