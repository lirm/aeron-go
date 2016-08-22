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
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/counters"
	"github.com/lirm/aeron-go/aeron/driver"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/ringbuffer"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"os"
	"testing"
	"time"
)

var logBufferName = "logbuffers.bin"

func prepareFile(t *testing.T) *logbuffer.LogBuffers {
	logging.SetLevel(logging.INFO, "logbuffers")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "aeron")
	mmap, err := memmap.NewFile(logBufferName, 0, 256*1024)
	if err != nil {
		t.Error(err.Error())
		return nil
	}
	mmap.Close()

	return logbuffer.Wrap(logBufferName)
}

func TestNewPublication(t *testing.T) {
	counterFileName := "cnc.dat"
	mmap, err := memmap.NewFile(counterFileName, 0, 256*1024)
	if err != nil {
		t.Error(err.Error())
	}
	defer mmap.Close()
	defer os.Remove(counterFileName)
	defer os.Remove(logBufferName)

	cncBuffer := atomic.MakeBuffer(mmap.GetMemoryPtr(), mmap.GetMemorySize())
	var meta counters.MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)
	meta.CncVersion.Set(counters.CurrentCncVersion)

	var proxy driver.Proxy
	var rb rb.ManyToOne
	rb.Init(meta.ToDriverBuf.Get())
	proxy.Init(&rb)

	var cc ClientConductor
	cc.Init(&proxy, nil, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100)
	defer cc.Close()

	lb := prepareFile(t)

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
		t.Errorf("Expected publication to not be connected: %d", pos)
	}

	pub.metaData.TimeOfLastStatusMsg.Set(milliEpoch.Nanoseconds())
	pos = pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	t.Logf("new pos: %d", pos)
	if pos != BackPressured {
		t.Errorf("Expected publication to be back pressured: %d", pos)
	}

	pub.pubLimit.set(1024)
	pos = pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	t.Logf("new pos: %d", pos)
	if pos != int64(srcBuffer.Capacity()+logbuffer.DataFrameHeader.Length) {
		t.Errorf("Expected publication to advance to position %d", srcBuffer.Capacity()+logbuffer.DataFrameHeader.Length)
	}
}
