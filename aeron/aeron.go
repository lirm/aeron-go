package aeron

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/broadcast"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"log"
	"os/user"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/driver"
)

type NewPublicationHandler func(string, int32, int32, int64)

type NewSubscriptionHandler func(string, int32, int64)

type AvailableImageHandler func(*Image)

type UnavailableImageHandler func(*Image)

type Context struct {
	aeronDir     string
	errorHandler func(error)

	newPublicationHandler        NewPublicationHandler
	newSubscriptionHandler       NewSubscriptionHandler
	availableImageHandler        AvailableImageHandler
	unavailableImageHandler      UnavailableImageHandler
	mediaDriverTimeoutMs         int64
	resourceLingerTimeout        int64
	publicationConnectionTimeout int64
}

func (ctx *Context) AeronDir(dir string) *Context {
	ctx.aeronDir = dir
	return ctx
}

func (ctx *Context) MediaDriverTimeout(to int64) *Context {
	ctx.mediaDriverTimeoutMs = to
	return ctx
}

func (ctx *Context) cncFileName() string {
	user, err := user.Current()
	var uName string = "unknown"
	if err != nil {
		log.Printf("Failed to get current user name: %v", err)
	}
	uName = user.Username
	return ctx.aeronDir + "/aeron-" + uName + "/" + counters.Descriptor.CNC_FILE
}

type Aeron struct {
	context            *Context
	conductor          ClientConductor
	toDriverRingBuffer buffers.ManyToOneRingBuffer
	driverProxy        driver.Proxy

	toDriverAtomicBuffer  *buffers.Atomic
	toClientsAtomicBuffer *buffers.Atomic
	counterValuesBuffer   *buffers.Atomic

	cncBuffer *logbuffer.MemoryMappedFile

	toClientsBroadcastReceiver *broadcast.Receiver
	toClientsCopyReceiver      *broadcast.CopyReceiver

	/*
			    std::random_device m_randomDevice;
		    std::default_random_engine m_randomEngine;
		    std::uniform_int_distribution<std::int32_t> m_sessionIdDistribution;



		    SleepingIdleStrategy m_idleStrategy;
	*/
}

func Connect(context *Context) *Aeron {
	aeron := new(Aeron)
	aeron.context = context

	aeron.cncBuffer = mapCncFile(context)

	aeron.toDriverAtomicBuffer = counters.CreateToDriverBuffer(aeron.cncBuffer)
	aeron.toClientsAtomicBuffer = counters.CreateToClientsBuffer(aeron.cncBuffer)
	aeron.counterValuesBuffer = counters.CreateCounterValuesBuffer(aeron.cncBuffer)

	aeron.toDriverRingBuffer.Init(aeron.toDriverAtomicBuffer)

	aeron.driverProxy.Init(&aeron.toDriverRingBuffer)

	aeron.toClientsBroadcastReceiver = new(broadcast.Receiver)
	aeron.toClientsBroadcastReceiver.Init(aeron.toClientsAtomicBuffer)
	//log.Printf("aeron.toClientsBroadcastReceiver: %v\n", aeron.toClientsBroadcastReceiver)

	aeron.toClientsCopyReceiver = new(broadcast.CopyReceiver)
	aeron.toClientsCopyReceiver.Init(aeron.toClientsBroadcastReceiver)

	clientLivenessTimeout := counters.ClientLivenessTimeout(aeron.cncBuffer)
	//log.Printf("clientLivenessTimeout=%d\n", clientLivenessTimeout)

	aeron.conductor.Init(&aeron.driverProxy, aeron.toClientsCopyReceiver, clientLivenessTimeout, context.mediaDriverTimeoutMs*1000000)
	aeron.conductor.counterValuesBuffer = aeron.counterValuesBuffer
	//log.Printf("aeron.conductor: %v\n", aeron.conductor)

	go aeron.conductor.Run()

	return aeron
}

func mapCncFile(context *Context) *logbuffer.MemoryMappedFile {

	//log.Printf("Trying to map file: %s", context.cncFileName())
	cncBuffer, err := logbuffer.MapExistingMemoryMappedFile(context.cncFileName(), 0, 0)
	if err != nil {
		log.Fatal("Failed to map the file: " + context.cncFileName() + " with " + err.Error())
	}

	cncVer := counters.CncVersion(cncBuffer)
	//log.Printf("Mapped %s for ver %d", context.cncFileName(), cncVer)

	if counters.Descriptor.CNC_VERSION != cncVer {
		log.Fatal(
			fmt.Sprintf("aeron cnc file version not understood: version=%d", cncVer))
	}

	return cncBuffer
}

func (aeron *Aeron) AddPublication(channel string, streamId int32) int64 {
	return aeron.conductor.AddPublication(channel, streamId)
}

func (aeron *Aeron) FindPublication(registrationId int64) *Publication {
	return aeron.conductor.FindPublication(registrationId)
}

func (aeron *Aeron) AddSubscription(channel string, streamId int32) int64 {
	return aeron.conductor.AddSubscription(channel, streamId)
}

func (aeron *Aeron) FindSubscription(registrationId int64) *Subscription {
	return aeron.conductor.FindSubscription(registrationId)
}
