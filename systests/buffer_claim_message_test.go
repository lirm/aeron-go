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
	"testing"
	"time"
)

const (
	streamId           = 1001
	fragmentCountLimit = 10
	messageLength      = 200
)

type BufferClaimMessageTestSuite struct {
	suite.Suite
	mediaDriver *driver.MediaDriver
	connection  *aeron.Aeron
	errorSink   chan error
	channel     string
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

func (suite *BufferClaimMessageTestSuite) TestShouldReceivePublishedMessageWithInterleavedAbort() {
	fragmentCount := 0
	fragmentHandler := func(*atomic.Buffer, int32, int32, *logbuffer.Header) {
		fragmentCount++
	}

	var bufferClaim logbuffer.Claim
	arr := make([]byte, messageLength)
	srcBuffer := atomic.MakeBuffer(arr)

	subscription, err := suite.connection.AddSubscription(suite.channel, streamId)
	suite.Require().NoError(err)
	defer subscription.Close()
	publication, err := suite.connection.AddPublication(suite.channel, streamId)
	suite.Require().NoError(err)
	defer publication.Close()

	suite.publishMessage(srcBuffer, *publication)

	for publication.TryClaim(messageLength, &bufferClaim) < 0 {
		time.Sleep(50 * time.Millisecond)
	}

	suite.publishMessage(srcBuffer, *publication)

	bufferClaim.Abort()

	expectedNumberOfFragments := 2
	numFragments := 0
	for numFragments < expectedNumberOfFragments {
		fragments := subscription.Poll(fragmentHandler, fragmentCountLimit)
		if fragments == 0 {
			time.Sleep(50 * time.Millisecond)
		}
		numFragments += fragments
	}

	suite.Assert().EqualValues(expectedNumberOfFragments, fragmentCount)
}

func (suite *BufferClaimMessageTestSuite) TestShouldTransferReservedValue() {
	var bufferClaim logbuffer.Claim

	subscription, err := suite.connection.AddSubscription(suite.channel, streamId)
	suite.Require().NoError(err)
	defer subscription.Close()
	publication, err := suite.connection.AddPublication(suite.channel, streamId)
	suite.Require().NoError(err)
	defer publication.Close()

	for publication.TryClaim(messageLength, &bufferClaim) < 0 {
		time.Sleep(50 * time.Millisecond)
	}

	reservedValue := time.Now().UnixMilli()
	bufferClaim.SetReservedValue(reservedValue)
	bufferClaim.Commit()

	fragmentHandler := func(_ *atomic.Buffer, _ int32, length int32, header *logbuffer.Header) {
		suite.Assert().EqualValues(messageLength, length)
		suite.Assert().EqualValues(reservedValue, header.GetReservedValue())
	}

	for 1 != subscription.Poll(fragmentHandler, fragmentCountLimit) {
		time.Sleep(50 * time.Millisecond)
	}
}

func (suite *BufferClaimMessageTestSuite) publishMessage(srcBuffer *atomic.Buffer, publication aeron.Publication) {
	for publication.Offer(srcBuffer, 0, messageLength, nil) < 0 {
		time.Sleep(50 * time.Millisecond)
	}
}

func TestBufferClaimMessage(t *testing.T) {
	suite.Run(t, &BufferClaimMessageTestSuite{channel: "aeron:udp?endpoint=localhost:24325"})
	suite.Run(t, &BufferClaimMessageTestSuite{channel: aeron.IpcChannel})
}
