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

package logbuffer

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"testing"
)

func prepareFile(t *testing.T) *LogBuffers {
	logging.SetLevel(logging.INFO, "logbuffers")
	logging.SetLevel(logging.INFO, "memmap")
	fname := "logbuffers.bin"
	mmap, err := memmap.NewFile(fname, 0, 256*1024)
	if err != nil {
		t.Error(err.Error())
		return nil
	}
	mmap.Close()

	return Wrap(fname)
}

func TestWrap(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("Panic: %v", err)
			logger.Error(errStr)
			t.FailNow()
		}
	}()

	lb := prepareFile(t)
	if lb != nil {
		defer lb.Close()
	} else {
		t.Fail()
	}
}

func TestWrapFail(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("Panic: %v", err)
			t.Logf("caught panic: %v", errStr)
		}
	}()

	fname := "logbuffers.bin"
	mmap, err := memmap.NewFile(fname, 0, 16*1024*1024-1)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	mmap.Close()

	Wrap(fname)
	t.Fail()
}

func TestLogBuffers_Buffer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("Panic: %v", err)
			logger.Error(errStr)
			t.FailNow()
		}
	}()

	lb := prepareFile(t)
	if lb != nil {
		defer lb.Close()
	} else {
		t.Fail()
	}

	for i := 0; i <= PartitionCount; i++ {
		if nil == lb.Buffer(i) {
			t.Errorf("nil buffer %d", i)
		}
	}
}

func TestLogBuffers_BufferFail(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Logf("Panic: %v", err)
		}
	}()

	lb := prepareFile(t)
	if lb != nil {
		defer lb.Close()
	} else {
		t.Fail()
	}

	lb.Buffer(PartitionCount + 1)
	t.Fail()
}

func TestLogBuffers_Meta(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("Panic: %v", err)
			logger.Error(errStr)
			t.FailNow()
		}
	}()

	lb := prepareFile(t)
	if lb != nil {
		defer lb.Close()
	} else {
		t.Fail()
	}

	for i := 0; i <= PartitionCount; i++ {
		if nil == lb.Buffer(i) {
			t.Errorf("nil buffer %d", i)
		}
	}

	meta := lb.Meta()

	t.Logf("meta fly size: %d", meta.Size())

	t.Logf("active term count offset: %d", meta.ActiveTermCountOff.Get())
	t.Logf("initTermID: %d", meta.InitTermID.Get())
	t.Logf("CorrelationId: %d", meta.CorrelationId.Get())
	t.Logf("tailCounter0: %d", meta.TailCounter[0].Get())
	t.Logf("tailCounter1: %d", meta.TailCounter[1].Get())
	t.Logf("tailCounter2: %d", meta.TailCounter[2].Get())
	t.Logf("defaultFrameHdrLen: %d", meta.DefaultFrameHdrLen.Get())

	if meta.Size() != int(util.CacheLineLength*7) {
		t.Errorf("Actual size: %d vs %d", meta.Size(), util.CacheLineLength*7)
	}
}
