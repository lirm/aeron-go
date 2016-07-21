package aeron

import (
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"sync/atomic"
)

type Subscription struct {
	conductor       *ClientConductor
	channel         string
	roundRobinIndex int
	registrationId  int64
	streamId        int32

	images atomic.Value

	isClosed atomic.Value
}

func NewSubscription(conductor *ClientConductor, channel string, registrationId int64, streamId int32) *Subscription {
	sub := new(Subscription)
	sub.images.Store(make([]*Image, 0))
	sub.conductor = conductor
	sub.channel = channel
	sub.registrationId = registrationId
	sub.streamId = streamId
	sub.isClosed.Store(false)

	return sub
}

func (sub *Subscription) IsClosed() bool {
	return sub.isClosed.Load().(bool)
}

func (sub *Subscription) Close() error {
	sub.isClosed.Store(true)
	sub.conductor.ReleaseSubscription(sub.registrationId, sub.images.Load().([]*Image))

	return nil
}

func (sub *Subscription) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	images := sub.images.Load().([]*Image)
	length := len(images)
	var fragmentsRead int = 0

	if length > 0 {
		var startingIndex int = sub.roundRobinIndex
		if startingIndex >= length {
			sub.roundRobinIndex = 0
			startingIndex = 0
		}

		var i int = startingIndex
		for fragmentsRead < fragmentLimit {
			fragmentsRead += images[i].Poll(handler, fragmentLimit-fragmentsRead)

			i++
			if i == length {
				i = 0
			}

			if i == startingIndex {
				break
			}
		}
	}

	return fragmentsRead
}

func (sub *Subscription) hasImage(sessionId int32) bool {
	images := sub.images.Load().([]*Image)
	for _, image := range images {
		if image.sessionId == sessionId {
			return true
		}
	}
	return false
}

func (sub *Subscription) addImage(image *Image) *[]*Image {

	images := sub.images.Load().([]*Image)

	sub.images.Store(append(images, image))

	return &images
}

func (sub *Subscription) removeImage(correlationId int64) (*Image, int) {

	var oldImage *Image
	var oldIx int

	images := sub.images.Load().([]*Image)
	for ix, image := range images {
		if image.correlationId == correlationId {
			logger.Debugf("Removing image %v for subscription %d", image, sub.registrationId)
			oldImage = image
			oldIx = ix

			images[ix] = images[len(images)-1]
			images[len(images)-1] = nil
			images = images[:len(images)-1]

			sub.images.Store(images)

			break
		}
	}

	return oldImage, oldIx
}
