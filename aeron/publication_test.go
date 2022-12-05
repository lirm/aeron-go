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
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func prepareCnc(t *testing.T) (string, *counters.MetaDataFlyweight) {
	counterFileName := "cnc.dat"
	mmap, err := memmap.NewFile(counterFileName, 0, 256*1024)
	require.NoError(t, err)

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
	assert.EqualValues(t, util.NumberOfTrailingZeroes(65536), 16)
	assert.EqualValues(t, util.NumberOfTrailingZeroes(131072), 17)
	assert.EqualValues(t, util.NumberOfTrailingZeroes(4194304), 22)
}

func TestNewPublication(t *testing.T) {
	cncName, meta := prepareCnc(t)
	defer require.NoError(t, os.Remove(cncName))

	var proxy driver.Proxy
	var rb rb.ManyToOne
	buf := meta.ToDriverBuf.Get()
	require.NotNil(t, buf)

	t.Logf("RingBuffer backing atomic.Buffer: %v", buf)
	rb.Init(buf)
	proxy.Init(&rb)

	var cc ClientConductor
	cc.Init(&proxy, nil, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, meta)
	defer require.NoError(t, cc.Close())

	lb, err := logbuffer.NewTestingLogbuffer()
	defer require.NoError(t, logbuffer.RemoveTestingLogbufferFile())
	require.NoError(t, err)

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
	require.EqualValues(t, pub.pubLimit.get(), 0)

	srcBuffer := atomic.MakeBuffer(make([]byte, 256))

	milliEpoch := (time.Nanosecond * time.Duration(time.Now().UnixNano())) / time.Millisecond
	pos := pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	t.Logf("new pos: %d", pos)
	assert.Equalf(t, pos, NotConnected, "Expected publication to not be connected at %d: %d", milliEpoch, pos)

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
	assert.EqualValuesf(t, pos, srcBuffer.Capacity()+logbuffer.DataFrameHeader.Length,
		"Expected publication to advance to position %d", srcBuffer.Capacity()+logbuffer.DataFrameHeader.Length)
}

func TestPublication_Offer(t *testing.T) {
	cncName, meta := prepareCnc(t)
	defer require.NoError(t, os.Remove(cncName))

	var proxy driver.Proxy
	var rb rb.ManyToOne
	buf := meta.ToDriverBuf.Get()
	require.NotNil(t, buf)
	t.Logf("RingBuffer backing atomic.Buffer: %v", buf)
	rb.Init(buf)
	proxy.Init(&rb)

	var cc ClientConductor
	cc.Init(&proxy, nil, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, time.Millisecond*100, meta)
	defer require.NoError(t, cc.Close())

	lb, err := logbuffer.NewTestingLogbuffer()
	defer require.NoError(t, logbuffer.RemoveTestingLogbufferFile())
	require.NoError(t, err)

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
	require.EqualValues(t, pub.pubLimit.get(), 0)

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

		assert.Equalf(t, pos, frameLen*i,
			"Expected publication to advance to position %d (have %d)", frameLen*i, pos)
	}

	pos := pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	t.Logf("new pos: %d", pos)
	assert.Equalf(t, pos, AdminAction,
		"Expected publication to trigger AdminAction (%d)", pos)
}
