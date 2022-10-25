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
	"github.com/lirm/aeron-go/systests/driver"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ChannelValidationTestSuite struct {
	suite.Suite
	mediaDriver *driver.MediaDriver
	connection  *aeron.Aeron
	errorSink   chan error
}

func (suite *ChannelValidationTestSuite) SetupSuite() {
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

func (suite *ChannelValidationTestSuite) TearDownSuite() {
	suite.Assert().Nil(suite.nextError())
	suite.connection.Close()
	suite.mediaDriver.StopMediaDriver()
}

func (suite *ChannelValidationTestSuite) TestPublicationCantUseDifferentSoSndbufIfAlreadySetViaUri() {
	pub1 := suite.addPublication("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1000)
	defer pub1.Close()

	pub2 := suite.addPublication("aeron:udp?endpoint=localhost:9999|so-sndbuf=65536", 1001)
	suite.Assert().Nil(pub2)
	suite.Assert().NotNil(suite.nextError())

	pub3 := suite.addPublication("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1002)
	defer pub3.Close()
}

func (suite *ChannelValidationTestSuite) TestSubscriptionCantUseDifferentSoSndbufIfAlreadySetViaUri() {
	sub1 := suite.addSubscription("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1000)
	defer sub1.Close()

	sub2 := suite.addSubscription("aeron:udp?endpoint=localhost:9999|so-sndbuf=65536", 1001)
	suite.Assert().Nil(sub2)
	suite.Assert().NotNil(suite.nextError())

	sub3 := suite.addSubscription("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1002)
	defer sub3.Close()
}

func (suite *ChannelValidationTestSuite) addPublication(channel string, streamId int32) *aeron.Publication {
	return <-suite.connection.AddPublication(channel, streamId)
}

func (suite *ChannelValidationTestSuite) addSubscription(channel string, streamId int32) *aeron.Subscription {
	return <-suite.connection.AddSubscription(channel, streamId)
}

func (suite *ChannelValidationTestSuite) nextError() error {
	timeout := time.After(100 * time.Millisecond)
	select {
	case <-timeout:
		return nil
	case err := <-suite.errorSink:
		return err
	}
}

func TestChannelValidation(t *testing.T) {
	suite.Run(t, new(ChannelValidationTestSuite))
}
