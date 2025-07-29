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
	"fmt"
	"testing"
	"time"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/testdata"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	StreamId1                          = int32(1002)
	CorrelationId1                     = int64(2000)
	SessionId1                         = int32(13)
	SourceInfo                         = "127.0.0.1:40789"
	SubscriptionPositionId             = int32(2)
	SubscriptionPositionRegistrationId = int64(4001)
	CounterTypeId                      = int32(102)
	CounterLabel                       = "counter label"
	CounterRegId                       = int64(42)
	CounterRegId2                      = int64(43)
	CounterId                          = int32(3)
)

type ClientConductorTestSuite struct {
	suite.Suite
	cc          ClientConductor
	driverProxy *MockDriverProxy
}

func (c *ClientConductorTestSuite) SetupTest() {
	metaBuf := atomic.MakeBuffer(make([]byte, 256*1024))
	meta := counters.InitAndWrapMetaData(metaBuf, 0, 0, 0, 24*1024, 24*1024, 0)
	c.driverProxy = new(MockDriverProxy)
	c.cc.Init(c.driverProxy, nil, time.Millisecond*100, time.Millisecond*100,
		time.Millisecond*100, time.Millisecond*100, meta)
	c.driverProxy.On("NextCorrelationID").Return(int64(123))
}

func (c *ClientConductorTestSuite) TestAddPublicationShouldNotifyMediaDriver() {
	c.driverProxy.On("AddPublication", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddPublication(Channel, StreamId1)
	c.Assert().NoError(err)
	c.Assert().Equal(reg, CorrelationId1)
}

func (c *ClientConductorTestSuite) TestAddPublicationShouldTimeoutWithoutReadyMessage() {
	c.driverProxy.On("AddPublication", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddPublication(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, CorrelationId1)

	// First a temporary error.
	pub, err := c.cc.FindPublication(CorrelationId1)
	c.Require().NoError(err)
	c.Assert().Nil(pub)

	// Then a permanent error.
	oldTimeout := c.cc.driverTimeoutNs
	c.cc.driverTimeoutNs = 1
	defer func() { c.cc.driverTimeoutNs = oldTimeout }()
	pub, err = c.cc.FindPublication(CorrelationId1)
	c.Assert().Error(err)
	c.Assert().Nil(pub)
}

func (c *ClientConductorTestSuite) TestShouldFailToAddPublicationOnMediaDriverError() {
	c.driverProxy.On("AddPublication", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddPublication(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, CorrelationId1)

	c.cc.OnErrorResponse(CorrelationId1, 1, "error")
	pub, err := c.cc.FindPublication(CorrelationId1)
	c.Assert().Nil(pub)
	c.Assert().Error(err)
}

func (c *ClientConductorTestSuite) TestAddSubscriptionShouldNotifyMediaDriver() {
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Assert().NoError(err)
	c.Assert().Equal(reg, CorrelationId1)
}

func (c *ClientConductorTestSuite) TestAddSubscriptionShouldTimeoutWithoutOperationSuccessful() {
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, CorrelationId1)

	// First a temporary error.
	sub, err := c.cc.FindSubscription(CorrelationId1)
	c.Assert().NoError(err)
	c.Assert().Nil(sub)

	// Then a permanent error.
	oldTimeout := c.cc.driverTimeoutNs
	c.cc.driverTimeoutNs = 1
	defer func() { c.cc.driverTimeoutNs = oldTimeout }()
	sub, err = c.cc.FindSubscription(CorrelationId1)
	c.Assert().Error(err)
	c.Assert().Nil(sub)
}

func (c *ClientConductorTestSuite) TestShouldFailToAddSubscriptionOnMediaDriverError() {
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, CorrelationId1)

	c.cc.OnErrorResponse(CorrelationId1, 1, "error")
	sub, err := c.cc.FindSubscription(CorrelationId1)
	c.Assert().Nil(sub)
	c.Assert().Error(err)
}

func (c *ClientConductorTestSuite) TestClientNotifiedOfNewAndInactiveImagesWithDefaultHandler() {
	var availableHandler = FakeImageHandler{}
	var unavailableHandler = FakeImageHandler{}
	c.cc.onAvailableImageHandler = availableHandler.Handle
	c.cc.onUnavailableImageHandler = unavailableHandler.Handle
	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddSubscription(Channel, StreamId1)
	c.Require().NoError(err)
	c.Require().Equal(reg, CorrelationId1)

	c.cc.OnSubscriptionReady(CorrelationId1, -1)
	sub, err := c.cc.FindSubscription(CorrelationId1)
	c.Require().NoError(err)
	c.Require().NotNil(sub)

	var image MockImage
	c.cc.imageFactory = func(_ int32, _ int64, _ string, _ int64, _ string, _ *atomic.Buffer, _ int32) Image {
		return &image
	}

	c.Require().Equal(sub.ImageCount(), 0)
	c.cc.OnAvailableImage(StreamId1, SessionId1, fmt.Sprintf("%d-log", SessionId1), SourceInfo,
		SubscriptionPositionId, CorrelationId1, CorrelationId1)
	c.Require().Equal(sub.ImageCount(), 1)
	c.Require().Equal(availableHandler.GetHandledImage(), &image)
	c.Require().Nil(availableHandler.GetHandledImage())
	c.Require().Nil(unavailableHandler.GetHandledImage())

	image.On("CorrelationID").Return(CorrelationId1)
	image.On("Close").Return(nil)
	c.cc.OnUnavailableImage(CorrelationId1, sub.RegistrationID())
	c.Require().Equal(sub.ImageCount(), 0)
	c.Require().Equal(unavailableHandler.GetHandledImage(), &image)
	c.Require().Nil(unavailableHandler.GetHandledImage())
	c.Require().Nil(availableHandler.GetHandledImage())
}

func (c *ClientConductorTestSuite) TestClientNotifiedOfNewAndInactiveImagesWithSpecificHandler() {
	var failHandler = func(_ Image) {
		c.Fail("Default handler called instead of specified handler")
	}
	c.cc.onAvailableImageHandler = failHandler
	c.cc.onUnavailableImageHandler = failHandler

	var availableHandler = FakeImageHandler{}
	var unavailableHandler = FakeImageHandler{}

	c.driverProxy.On("AddSubscription", Channel, StreamId1).
		Return(CorrelationId1, nil)
	reg, err := c.cc.AddSubscriptionWithHandlers(Channel, StreamId1, availableHandler.Handle, unavailableHandler.Handle)
	c.Require().NoError(err)
	c.Require().Equal(reg, CorrelationId1)

	c.cc.OnSubscriptionReady(CorrelationId1, -1)
	sub, err := c.cc.FindSubscription(CorrelationId1)
	c.Require().NoError(err)
	c.Require().NotNil(sub)

	var image MockImage
	c.cc.imageFactory = func(_ int32, _ int64, _ string, _ int64, _ string, _ *atomic.Buffer, _ int32) Image {
		return &image
	}

	c.Require().Equal(sub.ImageCount(), 0)
	c.cc.OnAvailableImage(StreamId1, SessionId1, fmt.Sprintf("%d-log", SessionId1), SourceInfo,
		SubscriptionPositionId, CorrelationId1, CorrelationId1)
	c.Require().Equal(sub.ImageCount(), 1)
	c.Require().Equal(availableHandler.GetHandledImage(), &image)
	c.Require().Nil(availableHandler.GetHandledImage())
	c.Require().Nil(unavailableHandler.GetHandledImage())

	image.On("CorrelationID").Return(CorrelationId1)
	image.On("Close").Return(nil)
	c.cc.OnUnavailableImage(CorrelationId1, sub.RegistrationID())
	c.Require().Equal(sub.ImageCount(), 0)
	c.Require().Equal(unavailableHandler.GetHandledImage(), &image)
	c.Require().Nil(unavailableHandler.GetHandledImage())
	c.Require().Nil(availableHandler.GetHandledImage())
}

func (c *ClientConductorTestSuite) TestShouldIgnoreUnknownNewImage() {
	var failHandler = func(_ Image) {
		c.Fail("Unknown image should not trigger any handler")
	}
	c.cc.onAvailableImageHandler = failHandler
	c.cc.onUnavailableImageHandler = failHandler

	c.cc.OnAvailableImage(StreamId1, SessionId1, fmt.Sprintf("%d-log", SessionId1), SourceInfo,
		SubscriptionPositionId, SubscriptionPositionRegistrationId, CorrelationId1)
}

func (c *ClientConductorTestSuite) TestShouldIgnoreUnknownInactiveImage() {
	var failHandler = func(_ Image) {
		c.Fail("Unknown image should not trigger any handler")
	}
	c.cc.onAvailableImageHandler = failHandler
	c.cc.onUnavailableImageHandler = failHandler

	c.cc.OnUnavailableImage(CorrelationId1, SubscriptionPositionRegistrationId)
}

func (c *ClientConductorTestSuite) TestShouldReturnNullForUnknownCounter() {
	counter, err := c.cc.FindCounter(100)
	c.Assert().Nil(counter)
	c.Assert().Error(err)
}

// Also indirectly tests C++'s shouldSendAddCounterToDriver.  We mock out driverProxy.  They check that driverProxy
// wrote the correct data to a test ring buffer.
func (c *ClientConductorTestSuite) TestShouldReturnNullForCounterWithoutOnAvailableCounter() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Return(CounterRegId, nil)

	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)
	c.Require().NoError(err)
	c.Require().Equal(id, CounterRegId)

	counter, err := c.cc.FindCounter(id)
	c.Assert().Nil(counter)
	c.Assert().NoError(err)
}

func (c *ClientConductorTestSuite) TestShouldReturnCounterAfterOnAvailableCounter() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Return(CounterRegId, nil)
	handler := testdata.NewMockAvailableCounterHandler(c.T())
	handler.On("Handle", mock.Anything, CounterRegId, CounterId).Once()
	c.cc.AddAvailableCounterHandler(handler)

	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)
	c.Require().NoError(err)
	c.Require().Equal(id, CounterRegId)

	c.cc.OnAvailableCounter(CounterRegId, CounterId)

	counter, err := c.cc.FindCounter(id)
	c.Require().NotNil(counter)
	c.Require().NoError(err)
	c.Assert().Equal(counter.RegistrationId(), CounterRegId)
	c.Assert().Equal(counter.counterId, CounterId)
}

// C++ uses a shared_ptr's implicit destructor to call Close().  Here we just call Close() manually.
func (c *ClientConductorTestSuite) TestShouldReleaseCounterAfterGoingOutOfScope() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Return(CounterRegId, nil)

	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)
	c.Require().NoError(err)
	c.Require().Equal(id, CounterRegId)

	c.cc.OnAvailableCounter(CounterRegId, CounterId)

	counter, err := c.cc.FindCounter(id)
	c.Require().NotNil(counter)
	c.Require().NoError(err)

	c.driverProxy.On("RemoveCounter", CounterRegId).
		Return(int64(123), nil)
	counter.Close()

	counter, err = c.cc.FindCounter(id)
	c.Require().Nil(counter)
	c.Require().Error(err)
}

func (c *ClientConductorTestSuite) TestShouldReturnDifferentIdsForDuplicateAddCounter() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Once().Return(CounterRegId, nil)
	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)
	c.Require().NoError(err)
	c.Require().Equal(id, CounterRegId)

	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Once().Return(CounterRegId2, nil)
	id, err = c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)
	c.Require().NoError(err)
	c.Require().Equal(id, CounterRegId2)
}

func (c *ClientConductorTestSuite) TestShouldReturnSameFindCounterAfterOnAvailableCounter() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Return(CounterRegId, nil)

	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)

	c.cc.OnAvailableCounter(CounterRegId, CounterId)

	counter1, err := c.cc.FindCounter(id)
	c.Require().NotNil(counter1)
	c.Require().NoError(err)
	counter2, err := c.cc.FindCounter(id)
	c.Require().NotNil(counter2)
	c.Require().NoError(err)
	c.Require().Equal(counter1, counter2)
}

func (c *ClientConductorTestSuite) TestShouldReturnDifferentCounterAfterOnAvailableCounter() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Once().Return(CounterRegId, nil)
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Once().Return(CounterRegId2, nil)
	handler := testdata.NewMockAvailableCounterHandler(c.T())
	handler.On("Handle", mock.Anything, CounterRegId, CounterId).Once()
	handler.On("Handle", mock.Anything, CounterRegId2, CounterId).Once()
	c.cc.AddAvailableCounterHandler(handler)

	id1, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)
	c.Require().Equal(id1, CounterRegId)
	id2, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)
	c.Require().Equal(id2, CounterRegId2)

	c.cc.OnAvailableCounter(CounterRegId, CounterId)
	c.cc.OnAvailableCounter(CounterRegId2, CounterId)

	counter1, err := c.cc.FindCounter(id1)
	c.Require().NotNil(counter1)
	c.Require().NoError(err)
	counter2, err := c.cc.FindCounter(id2)
	c.Require().NotNil(counter2)
	c.Require().NoError(err)
	c.Require().NotEqual(counter1, counter2)
}

func (c *ClientConductorTestSuite) TestShouldNotFindCounterOnAvailableCounterForUnknownCorrelationId() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Return(CounterRegId, nil)
	handler := testdata.NewMockAvailableCounterHandler(c.T())
	handler.On("Handle", mock.Anything, CounterRegId2, CounterId).Once()
	c.cc.AddAvailableCounterHandler(handler)

	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)

	c.cc.OnAvailableCounter(CounterRegId2, CounterId)

	counter1, err := c.cc.FindCounter(id)
	c.Require().Nil(counter1)
	c.Require().NoError(err)
}

func (c *ClientConductorTestSuite) TestShouldTimeoutAddCounterWithoutOnAvailableCounter() {
	c.cc.driverTimeoutNs = 1
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Return(CounterRegId, nil)
	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)

	counter, err := c.cc.FindCounter(id)
	c.Require().Nil(counter)
	c.Require().Error(err)
}

func (c *ClientConductorTestSuite) TestShouldErrorOnFindWhenReceivingErrorResponseOnAddCounter() {
	c.driverProxy.On("AddCounterByLabel", CounterTypeId, CounterLabel).
		Return(CounterRegId, nil)
	id, err := c.cc.AddCounterByLabel(CounterTypeId, CounterLabel)

	c.cc.OnErrorResponse(id, 123, "can't add counter")

	counter, err := c.cc.FindCounter(id)
	c.Require().Nil(counter)
	c.Require().ErrorContains(err, "can't add counter")
}

func (c *ClientConductorTestSuite) TestShouldCallOnUnavailableCounter() {
	id := int64(101)
	handler := testdata.NewMockUnavailableCounterHandler(c.T())
	handler.On("Handle", mock.Anything, id, CounterId).Once()
	c.cc.AddUnavailableCounterHandler(handler)

	c.cc.OnUnavailableCounter(id, CounterId)
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

func (f *FakeImageHandler) GetHandledImage() Image {
	if len(f.images) == 0 {
		return nil
	}
	ret := f.images[0]
	f.images = f.images[1:]
	return ret
}

// Everything below is auto generated by mockery using this command:
// mockery --name=DriverProxy --inpackage --structname=MockDriverProxy --print

// Code generated by mockery v2.14.0. DO NOT EDIT.

// MockDriverProxy is an autogenerated mock type for the DriverProxy type
type MockDriverProxy struct {
	mock.Mock
}

// AddCounter provides a mock function with given fields: typeId, keyBuffer, keyOffset, keyLength, labelBuffer, labelOffset, labelLength
func (_m *MockDriverProxy) AddCounter(typeId int32, keyBuffer *atomic.Buffer, keyOffset int32, keyLength int32, labelBuffer *atomic.Buffer, labelOffset int32, labelLength int32) (int64, error) {
	ret := _m.Called(typeId, keyBuffer, keyOffset, keyLength, labelBuffer, labelOffset, labelLength)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int32, *atomic.Buffer, int32, int32, *atomic.Buffer, int32, int32) int64); ok {
		r0 = rf(typeId, keyBuffer, keyOffset, keyLength, labelBuffer, labelOffset, labelLength)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int32, *atomic.Buffer, int32, int32, *atomic.Buffer, int32, int32) error); ok {
		r1 = rf(typeId, keyBuffer, keyOffset, keyLength, labelBuffer, labelOffset, labelLength)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddCounterByLabel provides a mock function with given fields: typeId, label
func (_m *MockDriverProxy) AddCounterByLabel(typeId int32, label string) (int64, error) {
	ret := _m.Called(typeId, label)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int32, string) int64); ok {
		r0 = rf(typeId, label)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int32, string) error); ok {
		r1 = rf(typeId, label)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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

// NextCorrelationID provides a mock function with given fields:
func (_m *MockDriverProxy) NextCorrelationID() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// RemoveCounter provides a mock function with given fields: registrationId
func (_m *MockDriverProxy) RemoveCounter(registrationId int64) (int64, error) {
	ret := _m.Called(registrationId)

	var r0 int64
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(registrationId)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(registrationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
