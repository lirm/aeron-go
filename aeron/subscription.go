package aeron

import (
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
)

type Subscription struct {
	conductor       *ClientConductor
	channel         string
	roundRobinIndex int
	registrationId  int64
	streamId        int32

	images []*Image

	imagesLength int

	isClosed bool
}

func (sub *Subscription) IsClosed() bool {
	return sub.isClosed
}

// int poll(const fragment_handler_t fragmentHandler, int fragmentLimit)
func (sub *Subscription) Poll(handler term.FragmentHandler, fragmentLimit int) int {

	// const int length = std::atomic_load(&m_imagesLength);
	// Image *images = std::atomic_load(&m_images);
	var fragmentsRead int = 0
	var length int = sub.imagesLength

	if length > 0 {
		var startingIndex int = sub.roundRobinIndex
		if startingIndex >= length {
			sub.roundRobinIndex = 0
			startingIndex = 0
		}

		var i int = startingIndex
		for fragmentsRead < fragmentLimit {
			fragmentsRead += sub.images[i].Poll(handler, fragmentLimit-fragmentsRead)

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
	for _, image := range sub.images {
		if image.sessionId == sessionId {
			return true
		}
	}
	return false
}

// FIXME make atomic
func (sub *Subscription) addImage(image *Image) *[]*Image {

	sub.images = append(sub.images, image)
	sub.imagesLength = len(sub.images)

	return &sub.images
}

func (sub *Subscription) removeImage(correlationId int64) (*Image, int) {

	var oldImage *Image
	var oldIx int

	for ix, image := range sub.images {
		if image.correlationId == correlationId {
			oldImage = image
			oldIx = ix
			/*

			Image * newArray = new Image[length - 1]

			for int i = 0, j = 0; i < length; i++ {
				if i != index {
					newArray[j++] = std::move(oldArray[i]);
				}
			}

			std::atomic_store(&m_imagesLength, length - 1);  // set length first. Don't go over end of new array on poll
			std::atomic_store(&m_images, newArray);
			*/

		}
	}

	return oldImage, oldIx
}
