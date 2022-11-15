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

package logbuffer

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWrap(t *testing.T) {
	defer func() {
		err := recover()
		assert.Nilf(t, err, "Panic: %v", err)
	}()

	lb, err := NewTestingLogbuffer()
	defer assert.NoError(t, RemoveTestingLogbufferFile())
	require.NoError(t, err)
	assert.NoError(t, lb.Close())
}

func TestWrapFail(t *testing.T) {
	fname := "logbuffers.bin"
	mmap, err := memmap.NewFile(fname, 0, 16*1024*1024-1)
	defer assert.NoError(t, os.Remove(fname))
	require.NoError(t, err)
	assert.NoError(t, mmap.Close())

	assert.Panics(t, func() { Wrap(fname) })
}

func TestLogBuffers_Buffer(t *testing.T) {
	defer func() {
		err := recover()
		assert.Nilf(t, err, "Panic: %v", err)
	}()

	lb, err := NewTestingLogbuffer()
	defer assert.NoError(t, RemoveTestingLogbufferFile())
	require.NoError(t, err)
	defer assert.NoError(t, lb.Close())

	for i := 0; i <= PartitionCount; i++ {
		assert.NotNilf(t, lb.Buffer(i), "buffer %d", i)
	}
}

func TestLogBuffers_BufferFail(t *testing.T) {
	lb, err := NewTestingLogbuffer()
	defer assert.NoError(t, RemoveTestingLogbufferFile())
	require.NoError(t, err)
	defer assert.NoError(t, lb.Close())

	// Index is zero-based
	assert.Panics(t, func() { lb.Buffer(PartitionCount + 1) })
}

func TestHeader_ReservedValue(t *testing.T) {
	bytes := make([]byte, 1000)
	buffer := atomic.MakeBuffer(bytes)
	assert.Equal(t, buffer.Capacity(), int32(1000))

	var header Header
	header.Wrap(buffer.Ptr(), 1000)
	assert.Equal(t, header.GetReservedValue(), int64(0))

	header.SetReservedValue(123)
	assert.Equal(t, header.GetReservedValue(), int64(123))
}

func TestLogBuffers_Meta(t *testing.T) {
	defer func() {
		err := recover()
		assert.Nilf(t, err, "Panic: %v", err)
	}()

	lb, err := NewTestingLogbuffer()
	defer assert.NoError(t, RemoveTestingLogbufferFile())
	require.NoError(t, err)
	defer assert.NoError(t, lb.Close())

	for i := 0; i <= PartitionCount; i++ {
		assert.NotNilf(t, lb.Buffer(i), "nil buffer %d", i)
	}

	meta := lb.Meta()

	//t.Logf("meta fly size: %d", meta.Size())
	//t.Logf("active term count offset: %d", meta.ActiveTermCountOff.Get())
	//t.Logf("initTermID: %d", meta.InitTermID.Get())
	//t.Logf("CorrelationId: %d", meta.CorrelationId.Get())
	//t.Logf("tailCounter0: %d", meta.TailCounter[0].Get())
	//t.Logf("tailCounter1: %d", meta.TailCounter[1].Get())
	//t.Logf("tailCounter2: %d", meta.TailCounter[2].Get())
	//t.Logf("defaultFrameHdrLen: %d", meta.DefaultFrameHdrLen.Get())

	assert.Equal(t, meta.Size(), int(util.CacheLineLength*7))
}
