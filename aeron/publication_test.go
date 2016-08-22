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
	//"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"testing"
)

func prepareFile(t *testing.T) *logbuffer.LogBuffers {
	logging.SetLevel(logging.INFO, "logbuffers")
	logging.SetLevel(logging.INFO, "memmap")
	logging.SetLevel(logging.INFO, "aeron")
	fname := "logbuffers.bin"
	mmap, err := memmap.NewFile(fname, 0, 256*1024)
	if err != nil {
		t.Error(err.Error())
		return nil
	}
	mmap.Close()

	return logbuffer.Wrap(fname)
}

func TestNewPublication(t *testing.T) {
	lb := prepareFile(t)

	pub := NewPublication(lb)

	//pub.conductor = cc
	pub.channel = "aeron:ipc"
	pub.regID = 1
	pub.streamID = 10
	pub.sessionID = 100
	metaBuffer := lb.Buffer(logbuffer.PartitionCount)
	if nil == metaBuffer {
		t.Logf("meta: %v", metaBuffer)
	}

	pub.pubLimit = NewPosition(metaBuffer, 0)
	t.Logf("pub: %v", pub.pubLimit)
	if pub.pubLimit.get() != 0 {
		t.Fail()
	}

	//srcBuffer := atomic.MakeBuffer(make([]byte, 256))

	// FIXME how can we make it work without creating the whole ClientConductor
	//pos := pub.Offer(srcBuffer, 0, srcBuffer.Capacity(), nil)
	//t.Logf("new pos: %d", pos)
}
