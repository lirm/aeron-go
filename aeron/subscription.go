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
	sub.roundRobinIndex = 0
	sub.isClosed.Store(false)

	return sub
}

func (sub *Subscription) IsClosed() bool {
	return sub.isClosed.Load().(bool)
}

func (sub *Subscription) Close() error {
	if !sub.IsClosed() {
		sub.isClosed.Store(true)
		<-sub.conductor.releaseSubscription(sub.registrationId, sub.images.Load().([]*Image))
	}

	return nil
}

func (sub *Subscription) Images() []*Image {
	return sub.images.Load().([]*Image)
}

func (sub *Subscription) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	images := sub.images.Load().([]*Image)
	length := len(images)
	var fragmentsRead int = 0

	if length > 0 {
		var startingIndex int = sub.roundRobinIndex
		sub.roundRobinIndex++
		if startingIndex >= length {
			sub.roundRobinIndex = 0
			startingIndex = 0
		}

		for i := startingIndex; i < length && fragmentsRead < fragmentLimit; i++ {
			fragmentsRead += images[i].Poll(handler, fragmentLimit-fragmentsRead)
		}

		for i := 0; i < startingIndex && fragmentsRead < fragmentLimit; i++ {
			fragmentsRead += images[i].Poll(handler, fragmentLimit-fragmentsRead)
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

func (sub *Subscription) removeImage(correlationId int64) *Image {

	images := sub.images.Load().([]*Image)
	for ix, image := range images {
		if image.correlationId == correlationId {
			logger.Debugf("Removing image %v for subscription %d", image, sub.registrationId)

			images[ix] = images[len(images)-1]
			images[len(images)-1] = nil
			images = images[:len(images)-1]

			sub.images.Store(images)

			return image
		}
	}
	return nil
}

func (sub *Subscription) HasImages() bool {
	images := sub.images.Load().([]*Image)
	return len(images) > 0
}

func IsConnectedTo(sub *Subscription, pub *Publication) bool {
	images := sub.images.Load().([]*Image)
	if sub.channel == pub.channel && sub.streamId == pub.streamId {
		for _, image := range images {
			if image.sessionId == pub.sessionId {
				return true
			}
		}
	}

	return false
}
