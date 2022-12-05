/*
Copyright 2016 Stanislav Liberman
Copyright 2022 Steven Stern

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

package systests

import (
	"flag"
	"fmt"
	"github.com/lirm/aeron-go/systests/driver"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logging"
)

var ExamplesConfig = struct {
	TestChannel  *string
	TestStreamID *int
	LoggingOn    *bool
}{
	flag.String("c", "aeron:ipc", "test channel"),
	flag.Int("s", 10, "streamId to use"),
	flag.Bool("l", false, "enable logging"),
}

var logger = logging.MustGetLogger("systests")

type SysTestSuite struct {
	suite.Suite
	mediaDriver *driver.MediaDriver
}

func (suite *SysTestSuite) SetupSuite() {
	mediaDriver, err := driver.StartMediaDriver()
	suite.Require().NoError(err, "Couldn't start Media Driver: ")
	suite.mediaDriver = mediaDriver
}

func (suite *SysTestSuite) TearDownSuite() {
	suite.mediaDriver.StopMediaDriver()
}

func (suite *SysTestSuite) send(n int, pub *aeron.Publication) {
	message := "this is a message"
	srcBuffer := atomic.MakeBuffer(([]byte)(message))

	for i := 0; i < n; i++ {
		timeoutAt := time.Now().Add(time.Second * 5)
		var v int64
		for v <= 0 {
			v = pub.Offer(srcBuffer, 0, int32(len(message)), nil)
			if time.Now().After(timeoutAt) {
				suite.Fail("Send timed out")
			}
			time.Sleep(time.Millisecond)
		}
	}
}

func (suite *SysTestSuite) receive(n int, sub *aeron.Subscription) {
	counter := 0
	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		logger.Debugf("    message: %s", string(buffer.GetBytesArray(offset, length)))
		counter++
	}
	var fragmentsRead atomic.Int
	for int(fragmentsRead.Get()) < n {
		timeoutAt := time.Now().Add(time.Second)
		for {
			recvd := sub.Poll(handler, 10)
			if recvd >= 1 {
				fragmentsRead.Add(int32(recvd))
				logger.Debugf("  have %d fragments", fragmentsRead)
				break
			}
			if time.Now().After(timeoutAt) {
				suite.Fail("Receive timed out")
			}
			time.Sleep(time.Millisecond)
		}
	}
	suite.Assert().EqualValues(fragmentsRead.Get(), n)
	suite.Assert().EqualValues(counter, n)
}

func (suite *SysTestSuite) subAndSend(n int, a *aeron.Aeron, pub *aeron.Publication) {
	sub, err := a.AddSubscription(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
	suite.Require().NoError(err)
	defer sub.Close()

	// This is basically a requirement since we need to wait
	for !aeron.IsConnectedTo(sub, pub) {
		time.Sleep(time.Millisecond)
	}

	suite.send(n, pub)
	suite.receive(n, sub)
}

func logtest(flag bool) {
	fmt.Printf("Logging: %t\n", flag)
	if flag {
		logging.SetLevel(logging.DEBUG, "aeron")
		logging.SetLevel(logging.DEBUG, "memmap")
		logging.SetLevel(logging.DEBUG, "driver")
		logging.SetLevel(logging.DEBUG, "counters")
		logging.SetLevel(logging.DEBUG, "logbuffers")
		logging.SetLevel(logging.DEBUG, "buffer")
	} else {
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.INFO, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")

	}
}

// TestAeronBasics will check for a simple send/receive scenario.
func (suite *SysTestSuite) TestAeronBasics() {
	logger.Debug("Started TestAeronBasics")

	a, err := aeron.Connect(aeron.NewContext().AeronDir(suite.mediaDriver.TempDir))
	if err != nil {
		suite.Failf("Failed to connect to driver: %s", err.Error())
	}
	defer a.Close()

	pub, err := a.AddPublication(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
	suite.Require().NoError(err)
	defer pub.Close()
	//logger.Debugf("Added publication: %v\n", pub)

	suite.subAndSend(1, a, pub)
}

// TestAeronSendMultipleMessages tests sending and receive multiple messages in a row.
func (suite *SysTestSuite) TestAeronSendMultipleMessages() {
	logger.Debug("Started TestAeronSendMultipleMessages")

	a, err := aeron.Connect(aeron.NewContext().AeronDir(suite.mediaDriver.TempDir))
	suite.Require().Nil(err, "Failed to connect to driver: %s", err)
	defer a.Close()

	for i := 0; i < 3; i++ {
		logger.Debugf("NextCorrelationID = %d", a.NextCorrelationID())
	}
	suite.Require().NotEqual(a.NextCorrelationID(), 0, "invalid zero NextCorrelationID")

	pub, err := a.AddPublication(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
	suite.Require().NoError(err)
	defer pub.Close()

	sub, err := a.AddSubscription(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
	suite.Require().NoError(err)
	defer sub.Close()

	// This is basically a requirement since we need to wait
	for !aeron.IsConnectedTo(sub, pub) {
		time.Sleep(time.Millisecond)
	}

	itCount := 100
	go suite.send(itCount, pub)
	suite.receive(itCount, sub)
}

// TestAeronSendMultiplePublications tests sending on multiple publications with a sigle
// subscription receiving. In IPC local mode this will end up with using the same Publication
// but it's a scenario nonetheless. As all systests this assumes a running media driver.
func (suite *SysTestSuite) NotTestedYet_TestAeronSendMultiplePublications() {
	logger.Debug("Started TestAeronSendMultiplePublications")

	//go func() {
	//	sigs := make(chan os.Signal, 1)
	//	signal.Notify(sigs, syscall.SIGQUIT)
	//	buf := make([]byte, 1<<20)
	//	for {
	//		<-sigs
	//		stacklen := runtime.Stack(buf, true)
	//		log.Printf("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
	//	}
	//}()

	a, err := aeron.Connect(aeron.NewContext())
	if err != nil {
		logger.Fatalf("Failed to connect to driver: %s\n", err.Error())
	}
	defer a.Close()

	pubCount := 10
	itCount := 100

	sub, err := a.AddSubscription(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
	suite.Require().NoError(err)
	defer sub.Close()

	pubs := make([]*aeron.Publication, pubCount)

	for i := 0; i < pubCount; i++ {
		pub, err := a.AddPublication(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
		suite.Require().NoError(err)
		defer pub.Close()

		pubs[i] = pub

		// This is basically a requirement since we need to wait
		for !aeron.IsConnectedTo(sub, pub) {
			time.Sleep(time.Millisecond)
		}
	}

	logger.Debugf(" ==> Got pubs %v", pubs)

	go suite.receive(itCount*pubCount, sub)

	time.Sleep(200 * time.Millisecond)

	// Send
	for i := 0; i < itCount; i++ {
		for pIx, p := range pubs {
			suite.send(1, p)
			logger.Debugf("sent %d to pubs[%d]", i, pIx)
			logger.Debugf("sent %d to pubs[%d]", i, pIx)
		}
	}

}

// TestAeronResubscribe test using different subscriptions with the same publication
func (suite *SysTestSuite) NotTestedYet_TestAeronResubscribe() {
	logger.Debug("Started TestAeronResubscribe")

	a, err := aeron.Connect(aeron.NewContext())
	suite.Require().NoError(err)
	defer a.Close()

	pub, err := a.AddPublication(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
	suite.Require().NoError(err)

	suite.subAndSend(1, a, pub)
	suite.subAndSend(1, a, pub)
}

// TestResubStress tests sending and receiving when creating a new subscription for each cycle
func (suite *SysTestSuite) NotTestedYet_TestResubStress() {
	logger.Debug("Started TestAeronResubscribe")

	a, err := aeron.Connect(aeron.NewContext())
	suite.Require().NoError(err)
	defer a.Close()

	pub, err := a.AddPublication(*ExamplesConfig.TestChannel, int32(*ExamplesConfig.TestStreamID))
	suite.Require().NoError(err)
	for i := 0; i < 100; i++ {
		suite.subAndSend(1, a, pub)
		logger.Debugf("bounce %d", i)
	}
}

// TestAeronClose simply tests explicit call to ctx.Close()
func testAeronClose() {
	logger.Debug("Started TestAeronClose")

	ctx := aeron.NewContext().MediaDriverTimeout(time.Second * 5)
	a, err := aeron.Connect(ctx)
	if err != nil {
		logger.Fatalf("Failed to connect to driver: %s\n", err.Error())
	}
	a.Close()
}

func TestSuiteMain(t *testing.T) {
	flag.Parse()
	logtest(*ExamplesConfig.LoggingOn)
	suite.Run(t, new(SysTestSuite))

	//testAeronClose()

	//testAeronResubscribe()

	//testAeronSendMultiplePublications()

	//testResubStress()
}
