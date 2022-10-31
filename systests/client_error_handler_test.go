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
	"errors"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/systests/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	streamId = 1001
	channel  = "aeron:ipc"
	sleep    = 50 * time.Millisecond
)

// Note: This is inspired by the Java test shouldHaveCorrectTermBufferLength in ClientErrorHandlerTest.  The name is a
// copy/paste bug.
func TestShouldUseCorrectErrorHandlers(t *testing.T) {
	mediaDriver, err := driver.StartMediaDriver()
	require.NoError(t, err, "Couldn't start Media Driver")
	defer mediaDriver.StopMediaDriver()

	errChanOne := make(chan error, 10)
	clientCtxOne := aeron.NewContext().AeronDir(mediaDriver.TempDir).
		ErrorHandler(func(err error) { errChanOne <- err })

	// Note: Java uses a RethrowingErrorHandler.  There is no Go equivalent.  Instead, this tests that the correct error
	// handler is used, which is also part of the Java test.
	errChanTwo := make(chan error, 10)
	clientCtxTwo := aeron.NewContext().AeronDir(mediaDriver.TempDir).
		ErrorHandler(func(err error) { require.Fail(t, "Wrong error handler called") }).
		SubscriberErrorHandler(func(err error) { errChanTwo <- err })

	aeronOne, err := aeron.Connect(clientCtxOne)
	require.NoError(t, err, "Couldn't connect to Media Driver")
	defer aeronOne.Close()

	aeronTwo, err := aeron.Connect(clientCtxTwo)
	require.NoError(t, err, "Couldn't connect to Media Driver")
	defer aeronTwo.Close()

	publication := <-aeronTwo.AddPublication(channel, streamId)
	require.NotNil(t, publication, "Couldn't AddPublication")
	defer publication.Close()

	subscriptionOne := <-aeronOne.AddSubscription(channel, streamId)
	require.NotNil(t, subscriptionOne, "Couldn't AddSubscription")
	defer subscriptionOne.Close()

	subscriptionTwo := <-aeronTwo.AddSubscription(channel, streamId)
	require.NotNil(t, subscriptionTwo, "Couldn't AddSubscription")
	defer subscriptionTwo.Close()

	for !subscriptionOne.IsConnected() {
		time.Sleep(sleep)
	}
	for !subscriptionTwo.IsConnected() {
		time.Sleep(sleep)
	}

	buffer := atomic.MakeBuffer(make([]byte, 100), 100)
	for publication.Offer(buffer, 0, 100, nil) < 0 {
		time.Sleep(sleep)
	}

	expectedError := errors.New("expected")
	handler := func(_ *atomic.Buffer, _ int32, _ int32, _ *logbuffer.Header) error {
		return expectedError
	}

	for subscriptionOne.Poll(handler, 1) == 0 {
		time.Sleep(sleep)
	}
	select {
	case err := <-errChanOne:
		assert.Equal(t, expectedError, err)
	default:
		assert.Fail(t, "Error not propagated")
	}

	for subscriptionTwo.Poll(handler, 1) == 0 {
		time.Sleep(sleep)
	}
	select {
	case err := <-errChanTwo:
		assert.Equal(t, expectedError, err)
	default:
		assert.Fail(t, "Error not propagated")
	}

	select {
	case err := <-errChanOne:
		assert.Failf(t, "Too many errors handled: %s", err.Error())
	case err := <-errChanTwo:
		assert.Failf(t, "Too many errors handled: %s", err.Error())
	default:
	}
}
