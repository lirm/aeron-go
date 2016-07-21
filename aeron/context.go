package aeron

import (
	"github.com/lirm/aeron-go/aeron/counters"
	"os/user"
	"time"
)

type Context struct {
	aeronDir     string
	errorHandler func(error)

	newPublicationHandler   NewPublicationHandler
	newSubscriptionHandler  NewSubscriptionHandler
	availableImageHandler   AvailableImageHandler
	unavailableImageHandler UnavailableImageHandler

	mediaDriverTo           time.Duration
	resourceLingerTo        time.Duration
	publicationConnectionTo time.Duration
}

func NewContext() *Context {
	ctx := new(Context)

	ctx.aeronDir = "/tmp"
	ctx.errorHandler = func(err error) { logger.Error(err) }

	ctx.newPublicationHandler = func(string, int32, int32, int64) {}
	ctx.newSubscriptionHandler = func(string, int32, int64) {}
	ctx.availableImageHandler = func(*Image) {}
	ctx.unavailableImageHandler = func(*Image) {}

	ctx.mediaDriverTo = time.Millisecond * 500
	ctx.resourceLingerTo = time.Millisecond * 10
	ctx.publicationConnectionTo = time.Second * 5

	return ctx
}

func (ctx *Context) AeronDir(dir string) *Context {
	ctx.aeronDir = dir
	return ctx
}

func (ctx *Context) MediaDriverTimeout(to time.Duration) *Context {
	ctx.mediaDriverTo = to
	return ctx
}

func (ctx *Context) cncFileName() string {
	user, err := user.Current()
	var uName string = "unknown"
	if err != nil {
		logger.Warningf("Failed to get current user name: %v", err)
	}
	uName = user.Username
	return ctx.aeronDir + "/aeron-" + uName + "/" + counters.Descriptor.CNC_FILE
}
