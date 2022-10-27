// Copyright 2022 Steven Stern
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package systests

import (
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/systests/driver"
	"github.com/stretchr/testify/suite"
	syncatomic "sync/atomic"
	"testing"
	"time"
)

const (
	channel            = "aeron:udp?endpoint=localhost:24325"
	streamId           = 1001
	fragmentCountLimit = 10
	messageLength      = 200
)

type BufferClaimMessageTestSuite struct {
	suite.Suite
	mediaDriver *driver.MediaDriver
	connection  *aeron.Aeron
	errorSink   chan error
}

func (suite *BufferClaimMessageTestSuite) SetupSuite() {
	mediaDriver, err := driver.StartMediaDriver()
	suite.Require().NoError(err, "Couldn't start Media Driver")
	suite.mediaDriver = mediaDriver

	suite.errorSink = make(chan error, 100)
	connect, err := aeron.Connect(aeron.
		NewContext().
		AeronDir(suite.mediaDriver.TempDir).
		ErrorHandler(func(err error) { suite.errorSink <- err }))
	if err != nil {
		// Testify does not run TearDownSuite if SetupSuite fails.  We have to manually stop Media Driver.
		suite.mediaDriver.StopMediaDriver()
		suite.Require().NoError(err, "aeron couldn't connect")
	}
	suite.connection = connect
}

func (suite *BufferClaimMessageTestSuite) TearDownSuite() {
	suite.connection.Close()
	suite.mediaDriver.StopMediaDriver()
}

// TODO: Parameterize this test
func (suite *BufferClaimMessageTestSuite) TestShouldReceivePublishedMessageWithInterleavedAbort() {
	/*
		logging.SetLevel(logging.DEBUG, "aeron")
		logging.SetLevel(logging.DEBUG, "memmap")
		logging.SetLevel(logging.DEBUG, "driver")
		logging.SetLevel(logging.DEBUG, "counters")
		logging.SetLevel(logging.DEBUG, "logbuffers")
		logging.SetLevel(logging.DEBUG, "buffer")
		logging.SetLevel(logging.DEBUG, "rb")

	*/
	var fragmentCount *int32 = new(int32)
	fragmentHandler := func(*atomic.Buffer, int32, int32, *logbuffer.Header) {
		syncatomic.AddInt32(fragmentCount, 1)
	}

	var bufferClaim logbuffer.Claim
	arr := make([]byte, messageLength)
	srcBuffer := atomic.MakeBuffer(arr)

	subscription := <-suite.connection.AddSubscription(channel, streamId)
	suite.Require().NotNil(subscription)
	defer subscription.Close()
	publication := <-suite.connection.AddPublication(channel, streamId)
	suite.Require().NotNil(publication)
	defer publication.Close()

	suite.publishMessage(srcBuffer, *publication)

	for publication.TryClaim(messageLength, &bufferClaim) < 0 {
		// TODO: add a timeout
		time.Sleep(50 * time.Millisecond)
	}

	suite.publishMessage(srcBuffer, *publication)
	bufferClaim.Abort()

	expectedNumberOfFragments := 2
	numFragments := 0
	for numFragments < expectedNumberOfFragments {
		fragments := subscription.Poll(fragmentHandler, fragmentCountLimit)
		if fragments == 0 {
			// TODO: add a timeout
			time.Sleep(50 * time.Millisecond)
		}
		numFragments += fragments
	}

	suite.Assert().EqualValues(expectedNumberOfFragments, syncatomic.LoadInt32(fragmentCount))
}

func (suite *BufferClaimMessageTestSuite) publishMessage(srcBuffer *atomic.Buffer, publication aeron.Publication) {
	for publication.Offer(srcBuffer, 0, messageLength, nil) < 0 {
		// TODO: add a timeout
		time.Sleep(50 * time.Millisecond)
	}
}

func TestBufferClaimMessage(t *testing.T) {
	suite.Run(t, new(BufferClaimMessageTestSuite))
}
