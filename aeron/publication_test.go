/*
Copyright 2016 Stanislav Liberman

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

package aeron

import (
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/counters"
	"github.com/corymonroe-coinbase/aeron-go/aeron/driver"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logging"
	"github.com/corymonroe-coinbase/aeron-go/aeron/ringbuffer"
	"github.com/corymonroe-coinbase/aeron-go/aeron/util"
	"github.com/corymonroe-coinbase/aeron-go/aeron/util/memmap"
	"os"
	"testing"
	"time"
	"unsafe"
)

func prepareFile(t *testing.T) (string, *logbuffer.LogBuffers) {
	logging.SetLevel(logging.INFO, "logbuffers")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "aeron")

	logBufferName := "logbuffers.bin"
	var logLength = (logbuffer.TermMinLength * logbuffer.PartitionCount) + logbuffer.LogMetaDataLength
	mmap, err := memmap.NewFile(logBufferName, 0, int(logLength))
	if err != nil {
		t.Error(err.Error())
		return "", nil
	}
	basePtr := uintptr(mmap.GetMemoryPtr())
	ptr := unsafe.Pointer(basePtr + uintptr(int64(logLength-logbuffer.LogMetaDataLength)))
	buf := atomic.MakeBuffer(ptr, logbuffer.LogMetaDataLength)

	var meta logbuffer.LogBufferMetaData
	meta.Wrap(buf, 0)

	meta.TermLen.Set(1024)

	mmap.Close()

	return logBufferName, logbuffer.Wrap(logBufferName)
}

func prepareCnc(t *testing.T) (string, *counters.MetaDataFlyweight) {
	logging.SetLevel(logging.INFO, "logbuffers")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "aeron")

	counterFileName := "cnc.dat"
	mmap, err := memmap.NewFile(counterFileName, 0, 256*1024)
	if err != nil {
		t.Error(err.Error())
	}

	cncBuffer := atomic.MakeBuffer(mmap.GetMemoryPtr(), mmap.GetMemorySize())
	var meta counters.MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)
	meta.CncVersion.Set(counters.CurrentCncVersion)
	meta.ToDriverBufLen.Set(1024)
	meta.ToClientBufLen.Set(1024)

	// Rewrap to pick up the new length
	meta.Wrap(cncBuffer, 0)
	t.Logf("meta data: %v", meta)

	return counterFileName, &meta
}

func TestNumberOfZeros(t *testing.T) {
	if x := util.NumberOfTrailingZeroes(65536); x != 16 {
		t.Logf("trailing zeroes: %d", x)
		t.Fail()
	}
	if x := util.NumberOfTrailingZeroes(131072); x != 17 {
		t.Logf("trailing zeroes: %d", x)
		t.Fail()
	}
	if x := util.NumberOfTrailingZeroes(4194304); x != 22 {
		t.Logf("trailing zeroes: %d", x)
		t.Fail()
	}
}

func TestNewPublication(t *testing.T) {

	cncName, meta := prepareCnc(t)
	defer os.Remove(cncName)

	var proxy driver.Proxy
	var rb rb.ManyToOne
	buf := meta.ToDriverBuf.Get()
	if buf == nil {
		t.FailNow()
	}
	t.Logf("RingBuffer backing atomic.Buffer: %v", buf)
	rb.Init(buf)
	proxy.Init(&rb)

	var cc ClientConductor
	cc.Init(&proxy, nil, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, meta)
	defer cc.Close()

	lbName, lb := prepareFile(t)
	defer os.Remove(lbName)

	lb.Meta().MTULen.Set(8192)

	pub := NewPublication(lb)

	pub.conductor = &cc
	pub.channel = "aeron:ipc"
	pub.regID = 1
	pub.streamID = 10
	pub.sessionID = 100
	metaBuffer := lb.Buffer(logbuffer.PartitionCount)
	if nil == metaBuffer {
		t.Logf("meta: %v", metaBuffer)
	}

	counter := atomic.MakeBuffer(make([]byte, 256))
	pub.pubLimit = NewPosition(counter, 0)
	t.Logf("pub: %v", pub.pubLimit)
	if pub.pubLimit.get() != 0 {
		t.Fail()
	}

	srcBuffer := atomic.MakeBuffer(make([]byte, 256))

	milliEpoch := (time.Nanosecond * time.Duration(time.Now().UnixNano())) / time.Millisecond
	pos := pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	t.Logf("new pos: %d", pos)
	if pos != NotConnected {
		t.Errorf("Expected publication to not be connected at %d: %d", milliEpoch, pos)
	}

	//pub.metaData.TimeOfLastStatusMsg.Set(milliEpoch.Nanoseconds())
	//pos = pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	//t.Logf("new pos: %d", pos)
	//if pos != BackPressured {
	//	t.Errorf("Expected publication to be back pressured: %d", pos)
	//}
	//
	pub.pubLimit.set(1024)
	pos = pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	t.Logf("new pos: %d", pos)
	if pos != int64(srcBuffer.Capacity()+logbuffer.DataFrameHeader.Length) {
		t.Errorf("Expected publication to advance to position %d", srcBuffer.Capacity()+logbuffer.DataFrameHeader.Length)
	}
}

func TestPublication_Offer(t *testing.T) {
	cncName, meta := prepareCnc(t)
	defer os.Remove(cncName)

	var proxy driver.Proxy
	var rb rb.ManyToOne
	buf := meta.ToDriverBuf.Get()
	if buf == nil {
		t.FailNow()
	}
	t.Logf("RingBuffer backing atomic.Buffer: %v", buf)
	rb.Init(buf)
	proxy.Init(&rb)

	var cc ClientConductor
	cc.Init(&proxy, nil, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, meta)
	defer cc.Close()

	lbName, lb := prepareFile(t)
	defer os.Remove(lbName)

	lb.Meta().MTULen.Set(8192)

	pub := NewPublication(lb)

	pub.conductor = &cc
	pub.channel = "aeron:ipc"
	pub.regID = 1
	pub.streamID = 10
	pub.sessionID = 100
	metaBuffer := lb.Buffer(logbuffer.PartitionCount)
	if nil == metaBuffer {
		t.Logf("meta: %v", metaBuffer)
	}

	counter := atomic.MakeBuffer(make([]byte, 256))
	pub.pubLimit = NewPosition(counter, 0)
	t.Logf("pub: %v", pub.pubLimit)
	if pub.pubLimit.get() != 0 {
		t.Fail()
	}

	srcBuffer := atomic.MakeBuffer(make([]byte, 256))
	//milliEpoch := (time.Nanosecond * time.Duration(time.Now().UnixNano())) / time.Millisecond
	//pub.metaData.TimeOfLastStatusMsg.Set(milliEpoch.Nanoseconds())

	termBufLen := lb.Buffer(0).Capacity()
	t.Logf("Term buffer length: %d", termBufLen)

	frameLen := int64(srcBuffer.Capacity() + logbuffer.DataFrameHeader.Length)
	pub.pubLimit.set(int64(termBufLen))
	for i := int64(1); i <= int64(termBufLen)/frameLen; i++ {
		pos := pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
		//t.Logf("new pos: %d", pos)
		if pos != frameLen*i {
			t.Errorf("Expected publication to advance to position %d (have %d)", frameLen*i, pos)
		}
	}

	pos := pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	t.Logf("new pos: %d", pos)
	if pos != AdminAction {
		t.Errorf("Expected publication to trigger AdminAction (%d)", pos)
	}

}
