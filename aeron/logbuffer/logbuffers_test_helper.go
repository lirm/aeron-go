// Copyright 2016 Stanislav Liberman
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

package logbuffer

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"os"
	"unsafe"
)

const Filename string = "logbuffers.bin"

// NewTestingLogbuffer is a helper function to create a new memmap file and logbuffer to wrap it.
func NewTestingLogbuffer() (*LogBuffers, error) {
	var logLength = (TermMinLength * PartitionCount) + LogMetaDataLength
	mmap, err := memmap.NewFile(Filename, 0, int(logLength))
	if err != nil {
		return nil, err
	}
	basePtr := uintptr(mmap.GetMemoryPtr())
	ptr := unsafe.Pointer(basePtr + uintptr(int64(logLength-LogMetaDataLength)))
	buf := atomic.MakeBuffer(ptr, LogMetaDataLength)

	var meta LogBufferMetaData
	meta.Wrap(buf, 0)

	meta.TermLen.Set(1024)

	if err := mmap.Close(); err != nil {
		return nil, err
	}

	return Wrap(Filename), nil
}

// RemoveTestingLogbufferFile is meant to be called by defer immediately after creating it above.
func RemoveTestingLogbufferFile() error {
	return os.Remove(Filename)
}
