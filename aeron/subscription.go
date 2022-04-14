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
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer/term"
)

// Subscription is the object responsible for receiving messages from media driver. It is specific to a channel and
// stream ID combination.
type Subscription struct {
	conductor       *ClientConductor
	channel         string
	roundRobinIndex int
	registrationID  int64
	streamID        int32

	images *ImageList

	isClosed atomic.Bool
}

// NewSubscription is a factory method to create new subscription to be added to the media driver
func NewSubscription(conductor *ClientConductor, channel string, registrationID int64, streamID int32) *Subscription {
	sub := new(Subscription)
	sub.images = NewImageList()
	sub.conductor = conductor
	sub.channel = channel
	sub.registrationID = registrationID
	sub.streamID = streamID
	sub.roundRobinIndex = 0
	sub.isClosed.Set(false)

	return sub
}

// Channel returns the media address for delivery to the channel.
func (sub *Subscription) Channel() string {
	return sub.channel
}

// StreamID returns Stream identity for scoping within the channel media address.
func (sub *Subscription) StreamID() int32 {
	return sub.streamID
}

// IsClosed returns whether this subscription has been closed.
func (sub *Subscription) IsClosed() bool {
	return sub.isClosed.Get()
}

// Close will release all images in this subscription, send command to the driver and block waiting for response from
// the media driver. Images will be lingered by the ClientConductor.
func (sub *Subscription) Close() error {
	if sub.isClosed.CompareAndSet(false, true) {
		images := sub.images.Empty()
		sub.conductor.releaseSubscription(sub.registrationID, images)
	}

	return nil
}

// Poll is the primary receive mechanism on subscription.
func (sub *Subscription) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	img := sub.images.Get()
	length := len(img)
	var fragmentsRead int

	if length > 0 {
		startingIndex := sub.roundRobinIndex
		sub.roundRobinIndex++
		if startingIndex >= length {
			sub.roundRobinIndex = 0
			startingIndex = 0
		}

		for i := startingIndex; i < length && fragmentsRead < fragmentLimit; i++ {
			fragmentsRead += img[i].Poll(handler, fragmentLimit-fragmentsRead)
		}

		for i := 0; i < startingIndex && fragmentsRead < fragmentLimit; i++ {
			fragmentsRead += img[i].Poll(handler, fragmentLimit-fragmentsRead)
		}
	}

	return fragmentsRead
}

// PollWithContext as for Poll() but provides an integer argument for passing contextual information
func (sub *Subscription) PollWithContext(handler term.FragmentHandler, fragmentLimit int) int {

	img := sub.images.Get()
	length := len(img)
	var fragmentsRead int

	if length > 0 {
		startingIndex := sub.roundRobinIndex
		sub.roundRobinIndex++
		if startingIndex >= length {
			sub.roundRobinIndex = 0
			startingIndex = 0
		}

		for i := startingIndex; i < length && fragmentsRead < fragmentLimit; i++ {
			fragmentsRead += img[i].PollWithContext(handler, fragmentLimit-fragmentsRead)
		}

		for i := 0; i < startingIndex && fragmentsRead < fragmentLimit; i++ {
			fragmentsRead += img[i].PollWithContext(handler, fragmentLimit-fragmentsRead)
		}
	}

	return fragmentsRead
}

func (sub *Subscription) hasImage(sessionID int32) bool {
	img := sub.images.Get()
	for _, image := range img {
		if image.sessionID == sessionID {
			return true
		}
	}
	return false
}

func (sub *Subscription) addImage(image *Image) *[]Image {

	images := sub.images.Get()

	sub.images.Set(append(images, *image))

	return &images
}

func (sub *Subscription) removeImage(correlationID int64) *Image {

	img := sub.images.Get()
	for ix, image := range img {
		if image.correlationID == correlationID {
			logger.Debugf("Removing image %v for subscription %d", image, sub.registrationID)

			img[ix] = img[len(img)-1]
			img = img[:len(img)-1]

			sub.images.Set(img)

			return &image
		}
	}
	return nil
}

// RegistrationID returns the registration id.
func (sub *Subscription) RegistrationID() int64 {
	return sub.registrationID
}

// IsConnected returns if this subscription is connected by having at least one open publication Image.
func (sub *Subscription) IsConnected() bool {
	for _, image := range sub.images.Get() {
		if !image.IsClosed() {
			return true
		}
	}
	return false
}

// HasImages is a helper method checking whether this subscription has any images associated with it.
func (sub *Subscription) HasImages() bool {
	images := sub.images.Get()
	return len(images) > 0
}

// ImageCount count of images associated with this subscription.
func (sub *Subscription) ImageCount() int {
	images := sub.images.Get()
	return len(images)
}

// ImageBySessionId returns the associated with the given sessionId.
func (sub *Subscription) ImageBySessionID(sessionID int32) *Image {
	img := sub.images.Get()
	for _, image := range img {
		if image.sessionID == sessionID {
			return &image
		}
	}
	return nil
}

// IsConnectedTo is a helper function used primarily by tests, which is used within the same process to verify that
// subscription is connected to a specific publication.
func IsConnectedTo(sub *Subscription, pub *Publication) bool {
	img := sub.images.Get()
	if sub.channel == pub.channel && sub.streamID == pub.streamID {
		for _, image := range img {
			if image.sessionID == pub.sessionID {
				return true
			}
		}
	}

	return false
}
