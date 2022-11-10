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
)

type ChannelValidationTestSuite struct {
	suite.Suite
	mediaDriver *driver.MediaDriver
	connection  *aeron.Aeron
}

func (suite *ChannelValidationTestSuite) SetupSuite() {
	mediaDriver, err := driver.StartMediaDriver()
	suite.Require().NoError(err, "Couldn't start Media Driver")
	suite.mediaDriver = mediaDriver

	connect, err := aeron.Connect(aeron.NewContext().AeronDir(suite.mediaDriver.TempDir))
	if err != nil {
		// Testify does not run TearDownSuite if SetupSuite fails.  We have to manually stop Media Driver.
		suite.mediaDriver.StopMediaDriver()
		suite.Require().NoError(err, "aeron couldn't connect")
	}
	suite.connection = connect
}

func (suite *ChannelValidationTestSuite) TearDownSuite() {
	suite.connection.Close()
	suite.mediaDriver.StopMediaDriver()
}

func (suite *ChannelValidationTestSuite) TestPublicationCantUseDifferentSoSndbufIfAlreadySetViaUri() {
	pub1, err := suite.addPublication("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1000)
	suite.Require().NoError(err)
	defer pub1.Close()

	pub2, err := suite.addPublication("aeron:udp?endpoint=localhost:9999|so-sndbuf=65536", 1001)
	suite.Assert().Nil(pub2)
	suite.Assert().Error(err)

	pub3, err := suite.addPublication("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1002)
	suite.Require().NoError(err)
	defer pub3.Close()
}

func (suite *ChannelValidationTestSuite) TestSubscriptionCantUseDifferentSoSndbufIfAlreadySetViaUri() {
	sub1, err := suite.addSubscription("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1000)
	suite.Require().NoError(err)
	defer sub1.Close()

	sub2, err := suite.addSubscription("aeron:udp?endpoint=localhost:9999|so-sndbuf=65536", 1001)
	suite.Assert().Nil(sub2)
	suite.Assert().Error(err)

	sub3, err := suite.addSubscription("aeron:udp?endpoint=localhost:9999|so-sndbuf=131072", 1002)
	suite.Require().NoError(err)
	defer sub3.Close()
}

func (suite *ChannelValidationTestSuite) addPublication(channel string, streamId int32) (*aeron.Publication, error) {
	return suite.connection.AddPublication(channel, streamId)
}

func (suite *ChannelValidationTestSuite) addSubscription(channel string, streamId int32) (*aeron.Subscription, error) {
	return suite.connection.AddSubscription(channel, streamId)
}

func TestChannelValidation(t *testing.T) {
	suite.Run(t, new(ChannelValidationTestSuite))
}
