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

package aeron

import (
	"errors"
	"fmt"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

const (
	StreamId1              = int32(1002)
	RegistrationId1        = int64(2000)
	SessionId1             = int32(13)
	SourceInfo             = "127.0.0.1:40789"
	SubscriptionPositionId = int32(2)
)

type ClientConductorTestSuite struct {
	suite.Suite
	cc          ClientConductor
	driverProxy MockDriverProxy
}

func (c *ClientConductorTestSuite) SetupTest() {
	var meta counters.MetaDataFlyweight
	c.cc.Init(&c.driverProxy, nil, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, &meta)
}

func (c *ClientConductorTestSuite) TearDownSuite() {

}

func (c *ClientConductorTestSuite) TestAddPublicationShouldNotifyMediaDriver() {
	c.driverProxy.On("AddPublication", Channel, StreamId1).
		Return(RegistrationId1, nil)
	reg, err := c.cc.AddPublication(Channel, StreamId1)
	c.Assert().NoError(err)
	c.Assert().Equal(reg, RegistrationId1)
}

func (c *ClientConductorTestSuite) TestAddPublicationShouldTimeoutWithoutReadyMessage() {
	c.driverProxy.On("AddPublication", Channel, StreamId1).
		Return(RegistrationId1, nil)
	reg, err := c.cc.AddPublication(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, RegistrationId1)

	// First a temporary error.
	pub, err := c.cc.FindPublication(RegistrationId1)
	c.Assert().True(errors.Is(err, TemporaryError))
	c.Assert().Nil(pub)

	// Then a permanent error.
	oldTimeout := c.cc.driverTimeoutNs
	c.cc.driverTimeoutNs = 1
	defer func() { c.cc.driverTimeoutNs = oldTimeout }()
	pub, err = c.cc.FindPublication(RegistrationId1)
	c.Assert().False(errors.Is(err, TemporaryError))
	c.Assert().Nil(pub)
}

func (c *ClientConductorTestSuite) TestShouldFailToAddPublicationOnMediaDriverError() {
	c.driverProxy.On("AddPublication", Channel, StreamId1).
		Return(RegistrationId1, nil)
	reg, err := c.cc.AddPublication(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, RegistrationId1)

	// TODO: Revamp CC's OnError so this is an actual error that bubbles back to FindPublication below.
	c.cc.OnErrorResponse(RegistrationId1, 1, "error")
	pub, err := c.cc.FindPublication(RegistrationId1)
	c.Assert().Nil(pub)
	c.Assert().Error(err)
	c.Assert().False(errors.Is(err, TemporaryError))
}

func (c *ClientConductorTestSuite) TestAddSubscriptionShouldNotifyMediaDriver() {
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(RegistrationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Assert().NoError(err)
	c.Assert().Equal(reg, RegistrationId1)
}

func (c *ClientConductorTestSuite) TestAddSubscriptionShouldTimeoutWithoutOperationSuccessful() {
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(RegistrationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, RegistrationId1)

	// First a temporary error.
	sub, err := c.cc.FindSubscription(RegistrationId1)
	c.Assert().True(errors.Is(err, TemporaryError))
	c.Assert().Nil(sub)

	// Then a permanent error.
	oldTimeout := c.cc.driverTimeoutNs
	c.cc.driverTimeoutNs = 1
	defer func() { c.cc.driverTimeoutNs = oldTimeout }()
	sub, err = c.cc.FindSubscription(RegistrationId1)
	c.Assert().False(errors.Is(err, TemporaryError))
	c.Assert().Nil(sub)
}

func (c *ClientConductorTestSuite) TestShouldFailToAddSubscriptionOnMediaDriverError() {
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(RegistrationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, RegistrationId1)

	// TODO: Revamp CC's OnError so this is an actual error that bubbles back to FindSubscription below.
	c.cc.OnErrorResponse(RegistrationId1, 1, "error")
	sub, err := c.cc.FindSubscription(RegistrationId1)
	c.Assert().Nil(sub)
	c.Assert().Error(err)
	c.Assert().False(errors.Is(err, TemporaryError))
}

// Java calls this test clientNotifiedOfNewImageShouldMapLogFile and uses a mock LogBufferFactory as a way to check that
// Images were created.  For us, that would involve mocking several followup calls to LogBuffer, so it's easier to just
// have a mock ImageFactory.
func (c *ClientConductorTestSuite) TestClientNotifiedOfNewImageWithDefaultHandler() {
	var defaultImageHandler = FakeImageHandler{}
	c.cc.defaultOnAvailableImageHandler = defaultImageHandler.Handle
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(RegistrationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, RegistrationId1)

	c.cc.OnSubscriptionReady(RegistrationId1, -1)
	sub, err := c.cc.FindSubscription(RegistrationId1)
	c.Require().NoError(err)
	c.Require().NotNil(sub)

	var image MockImage
	c.cc.imageFactory = func(_ int32, _ int64, _ string, _ int64, _ string, _ *atomic.Buffer, _ int32) Image {
		return &image
	}

	c.Require().Equal(sub.ImageCount(), 0)
	c.cc.OnAvailableImage(StreamId1, SessionId1, fmt.Sprintf("%d-log", SessionId1), SourceInfo,
		SubscriptionPositionId, RegistrationId1, RegistrationId1)
	c.Require().Equal(sub.ImageCount(), 1)
	c.Require().Equal(defaultImageHandler.GetNextImage(), &image)
	c.Require().Nil(defaultImageHandler.GetNextImage())
}

func TestClientConductor(t *testing.T) {
	suite.Run(t, new(ClientConductorTestSuite))
}

type FakeImageHandler struct {
	images []Image
}

func (f *FakeImageHandler) Handle(image Image) {
	if f.images == nil {
		f.images = make([]Image, 0)
	}
	f.images = append(f.images, image)
}

func (f *FakeImageHandler) GetNextImage() Image {
	if len(f.images) == 0 {
		return nil
	}
	ret := f.images[0]
	f.images = f.images[1:]
	return ret
}

// Everything below is auto generated by mockery using this command:
// mockery --name=ReceivingConductor --inpackage --structname=MockReceivingConductor --print

// MockDriverProxy is an autogenerated mock type for the DriverProxy type
type MockDriverProxy struct {
	mock.Mock
}

// AddDestination provides a mock function with given fields: registrationID, channel
func (_m *MockDriverProxy) AddDestination(registrationID int64, channel string) (int64, error) {
	ret := _m.Called(registrationID, channel)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, string) int64); ok {
		r0 = rf(registrationID, channel)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(registrationID, channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddExclusivePublication provides a mock function with given fields: channel, streamID
func (_m *MockDriverProxy) AddExclusivePublication(channel string, streamID int32) (int64, error) {
	ret := _m.Called(channel, streamID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, int32) int64); ok {
		r0 = rf(channel, streamID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int32) error); ok {
		r1 = rf(channel, streamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddPublication provides a mock function with given fields: channel, streamID
func (_m *MockDriverProxy) AddPublication(channel string, streamID int32) (int64, error) {
	ret := _m.Called(channel, streamID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, int32) int64); ok {
		r0 = rf(channel, streamID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int32) error); ok {
		r1 = rf(channel, streamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddRcvDestination provides a mock function with given fields: registrationID, channel
func (_m *MockDriverProxy) AddRcvDestination(registrationID int64, channel string) (int64, error) {
	ret := _m.Called(registrationID, channel)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, string) int64); ok {
		r0 = rf(registrationID, channel)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(registrationID, channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddSubscription provides a mock function with given fields: channel, streamID
func (_m *MockDriverProxy) AddSubscription(channel string, streamID int32) (int64, error) {
	ret := _m.Called(channel, streamID)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, int32) int64); ok {
		r0 = rf(channel, streamID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int32) error); ok {
		r1 = rf(channel, streamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientClose provides a mock function with given fields:
func (_m *MockDriverProxy) ClientClose() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClientID provides a mock function with given fields:
func (_m *MockDriverProxy) ClientID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// RemoveDestination provides a mock function with given fields: registrationID, channel
func (_m *MockDriverProxy) RemoveDestination(registrationID int64, channel string) (int64, error) {
	ret := _m.Called(registrationID, channel)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, string) int64); ok {
		r0 = rf(registrationID, channel)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(registrationID, channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemovePublication provides a mock function with given fields: registrationID
func (_m *MockDriverProxy) RemovePublication(registrationID int64) error {
	ret := _m.Called(registrationID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(registrationID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveRcvDestination provides a mock function with given fields: registrationID, channel
func (_m *MockDriverProxy) RemoveRcvDestination(registrationID int64, channel string) (int64, error) {
	ret := _m.Called(registrationID, channel)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64, string) int64); ok {
		r0 = rf(registrationID, channel)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, string) error); ok {
		r1 = rf(registrationID, channel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveSubscription provides a mock function with given fields: registrationID
func (_m *MockDriverProxy) RemoveSubscription(registrationID int64) error {
	ret := _m.Called(registrationID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(registrationID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TimeOfLastDriverKeepalive provides a mock function with given fields:
func (_m *MockDriverProxy) TimeOfLastDriverKeepalive() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

type mockConstructorTestingTNewMockDriverProxy interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDriverProxy creates a new instance of MockDriverProxy. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDriverProxy(t mockConstructorTestingTNewMockDriverProxy) *MockDriverProxy {
	mock := &MockDriverProxy{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
