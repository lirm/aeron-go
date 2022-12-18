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
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/systests/driver"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

const (
	counterTypeId = 1101
	counterLabel  = "counter label"
)

type CounterTestSuite struct {
	suite.Suite
	mediaDriver *driver.MediaDriver
	clientA     *aeron.Aeron
	clientB     *aeron.Aeron
	keyBuffer   *atomic.Buffer
	labelBuffer *atomic.Buffer
}

func (s *CounterTestSuite) SetupTest() {
	mediaDriver, err := driver.StartMediaDriver()
	s.Require().NoError(err, "Couldn't start Media Driver")
	s.mediaDriver = mediaDriver

	clientA, errA := aeron.Connect(aeron.NewContext().AeronDir(s.mediaDriver.TempDir))
	clientB, errB := aeron.Connect(aeron.NewContext().AeronDir(s.mediaDriver.TempDir))
	if errA != nil || errB != nil {
		// Testify does not run TearDownTest if SetupTest fails.  We have to manually stop Media Driver.
		s.mediaDriver.StopMediaDriver()
		s.Require().NoError(errA, "aeron couldn't connect")
		s.Require().NoError(errB, "aeron couldn't connect")
	}
	s.clientA = clientA
	s.clientB = clientB

	s.keyBuffer = atomic.MakeBuffer(make([]byte, 8))
	s.labelBuffer = atomic.MakeBuffer([]byte(counterLabel))
}

func (s *CounterTestSuite) TearDownTest() {
	s.clientA.Close()
	s.clientB.Close()
	s.mediaDriver.StopMediaDriver()
}

func (s *CounterTestSuite) TestShouldBeAbleToAddCounter() {
	chanA := make(chan int32, 10)
	availableCounterHandlerClientA := NewMockAvailableCounterHandler(s.T())
	s.clientA.AddAvailableCounterHandler(availableCounterHandlerClientA)
	availableCounterHandlerClientA.On("Handle",
		mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		// Java has a verify(call, timeout).  Go does not.  This allows us to wait for the right call before asserting
		// that the right call happened.  As redundant as it is, this is probably the easiest way to test this
		// functionality.
		chanA <- args.Get(2).(int32)
	})

	chanB := make(chan int32, 10)
	availableCounterHandlerClientB := NewMockAvailableCounterHandler(s.T())
	s.clientB.AddAvailableCounterHandler(availableCounterHandlerClientB)
	availableCounterHandlerClientB.On("Handle",
		mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		chanB <- args.Get(2).(int32)
	})

	regId, err := s.clientA.AddCounter(
		counterTypeId,
		s.keyBuffer,
		0,
		s.keyBuffer.Capacity(),
		s.labelBuffer,
		0,
		int32(len(counterLabel)))
	s.Require().NoError(err)
	var counter *aeron.Counter
	for counter == nil && err == nil {
		counter, err = s.clientA.FindCounter(regId)
	}
	s.Require().NoError(err)

	s.Assert().False(counter.IsClosed())
	s.Assert().Equal(counter.RegistrationId(), s.clientA.CounterReader().GetCounterRegistrationId(counter.Id()))
	s.Assert().Equal(s.clientA.ClientID(), s.clientA.CounterReader().GetCounterOwnerId(counter.Id()))

	GetCounterFromChanOrFail(counter.Id(), chanA, s.Suite)
	GetCounterFromChanOrFail(counter.Id(), chanB, s.Suite)
	availableCounterHandlerClientA.AssertCalled(s.T(), "Handle",
		mock.Anything, counter.RegistrationId(), counter.Id())
	availableCounterHandlerClientB.AssertCalled(s.T(), "Handle",
		mock.Anything, counter.RegistrationId(), counter.Id())
}

func GetCounterFromChanOrFail(counterId int32, ch chan int32, s suite.Suite) {
	for {
		select {
		case id := <-ch:
			if id == counterId {
				return
			}
		case <-time.After(5 * time.Second):
			s.FailNow("Timed out waiting for avaiable counter handler call")
		}
	}
}

func (s *CounterTestSuite) TestShouldBeAbleToAddReadableCounterWithinHandler() {
	ch, handler := s.createReadableCounterHandler()
	s.clientB.AddAvailableCounterHandler(handler)

	regId, err := s.clientA.AddCounter(
		counterTypeId,
		s.keyBuffer,
		0,
		s.keyBuffer.Capacity(),
		s.labelBuffer,
		0,
		int32(len(counterLabel)))
	s.Require().NoError(err)
	var counter *aeron.Counter
	for counter == nil && err == nil {
		counter, err = s.clientA.FindCounter(regId)
	}
	s.Require().NoError(err)

	var readableCounter *counters.ReadableCounter
	select {
	case readableCounter = <-ch:
	case <-time.After(5 * time.Second):
		s.Fail("Timed out waiting for a ReadableCounter")
	}
	s.Assert().Equal(counters.RecordAllocated, readableCounter.State())
	s.Assert().Equal(counter.Id(), readableCounter.CounterId)
	s.Assert().Equal(counter.RegistrationId(), readableCounter.RegistrationId)
}

func (s *CounterTestSuite) TestShouldCloseReadableCounterOnUnavailableCounter() {
	ch, readableHandler := s.createReadableCounterHandler()
	s.clientB.AddAvailableCounterHandler(readableHandler)

	regId, err := s.clientA.AddCounter(
		counterTypeId,
		s.keyBuffer,
		0,
		s.keyBuffer.Capacity(),
		s.labelBuffer,
		0,
		int32(len(counterLabel)))
	s.Require().NoError(err)
	var counter *aeron.Counter
	for counter == nil && err == nil {
		counter, err = s.clientA.FindCounter(regId)
	}
	s.Require().NoError(err)

	var readableCounter *counters.ReadableCounter
	select {
	case readableCounter = <-ch:
	case <-time.After(5 * time.Second):
		s.Fail("Timed out waiting for a ReadableCounter")
	}
	unavailableHandler := s.createUnavailableCounterHandler(readableCounter)
	s.clientB.AddUnavailableCounterHandler(unavailableHandler)

	s.Require().False(readableCounter.IsClosed())
	s.Require().Equal(counters.RecordAllocated, readableCounter.State())

	counter.Close()

	start := time.Now()
	for !readableCounter.IsClosed() {
		if time.Now().After(start.Add(5 * time.Second)) {
			s.Fail("Timed out waiting for ReadableCounter to close")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *CounterTestSuite) TestShouldGetUnavailableCounterWhenOwningClientIsClosed() {
	ch, readableHandler := s.createReadableCounterHandler()
	s.clientB.AddAvailableCounterHandler(readableHandler)

	regId, err := s.clientA.AddCounter(
		counterTypeId,
		s.keyBuffer,
		0,
		s.keyBuffer.Capacity(),
		s.labelBuffer,
		0,
		int32(len(counterLabel)))
	s.Require().NoError(err)
	var counter *aeron.Counter
	for counter == nil && err == nil {
		counter, err = s.clientA.FindCounter(regId)
	}
	s.Require().NoError(err)

	var readableCounter *counters.ReadableCounter
	select {
	case readableCounter = <-ch:
	case <-time.After(5 * time.Second):
		s.FailNow("Timed out waiting for a ReadableCounter")
	}
	unavailableHandler := s.createUnavailableCounterHandler(readableCounter)
	s.clientB.AddUnavailableCounterHandler(unavailableHandler)

	s.Require().False(readableCounter.IsClosed())
	s.Require().Equal(counters.RecordAllocated, readableCounter.State())

	s.clientA.Close()

	start := time.Now()
	for !readableCounter.IsClosed() {
		if time.Now().After(start.Add(5 * time.Second)) {
			s.FailNow("Timed out waiting for ReadableCounter to close")
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type ReadableCounterHandler struct {
	suite *suite.Suite
	ch    chan *counters.ReadableCounter
}

func (r ReadableCounterHandler) Handle(reader *counters.Reader, registrationId int64, counterId int32) {
	if counterTypeId == reader.GetCounterTypeId(counterId) {
		counter, err := counters.NewReadableRegisteredCounter(reader, registrationId, counterId)
		r.suite.Require().NoError(err)
		r.ch <- counter
	}
}

func (s *CounterTestSuite) createReadableCounterHandler() (chan *counters.ReadableCounter, aeron.AvailableCounterHandler) {
	ch := make(chan *counters.ReadableCounter, 1)
	handler := ReadableCounterHandler{suite: &s.Suite, ch: ch}
	return ch, handler
}

type UnavailableCounterHandler struct {
	suite   *suite.Suite
	counter *counters.ReadableCounter
}

func (r UnavailableCounterHandler) Handle(_ *counters.Reader, registrationId int64, _ int32) {
	if r.counter.RegistrationId == registrationId {
		r.counter.Close()
	}
}

func (s *CounterTestSuite) createUnavailableCounterHandler(counter *counters.ReadableCounter) aeron.UnavailableCounterHandler {
	return UnavailableCounterHandler{suite: &s.Suite, counter: counter}
}

func TestCounter(t *testing.T) {
	suite.Run(t, new(CounterTestSuite))
}

// If AvailableCounterHandler changes, recreate the mock code with the below command.
// mockery --name=AvailableCounterHandler --inpackage --structname=MockAvailableCounterHandler --print

// Code generated by mockery v2.14.0. DO NOT EDIT.

// MockAvailableCounterHandler is an autogenerated mock type for the AvailableCounterHandler type
type MockAvailableCounterHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: countersReader, registrationId, counterId
func (_m *MockAvailableCounterHandler) Handle(countersReader *counters.Reader, registrationId int64, counterId int32) {
	_m.Called(countersReader, registrationId, counterId)
}

type mockConstructorTestingTNewMockAvailableCounterHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAvailableCounterHandler creates a new instance of MockAvailableCounterHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAvailableCounterHandler(t mockConstructorTestingTNewMockAvailableCounterHandler) *MockAvailableCounterHandler {
	mock := &MockAvailableCounterHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
