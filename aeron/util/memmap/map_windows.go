// Copyright 2016 Stanislav Liberman
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

package memmap

import (
	"errors"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

//
// Blatantly borrowed from https://github.com/edsrzf/mmap-go/blob/master/mmap_windows.go
//

// We keep this map so that we can get back the original handle from the memory address.
var handleLock sync.Mutex
var handleMap = map[uintptr]syscall.Handle{}

func doMap(fd int, offset int64, length int) (*File, error) {

	maxSizeHigh := uint32((offset + int64(length)) >> 32)
	maxSizeLow := uint32((offset + int64(length)) & 0xFFFFFFFF)
	h, errno := syscall.CreateFileMapping(syscall.Handle(uintptr(fd)), nil, syscall.PAGE_READWRITE, maxSizeHigh, maxSizeLow, nil)
	if h == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}

	// Actually map a view of the data into memory. The view's size
	// is the length the user requested.
	fileOffsetHigh := uint32(offset >> 32)
	fileOffsetLow := uint32(offset & 0xFFFFFFFF)
	addr, errno := syscall.MapViewOfFile(h, syscall.FILE_MAP_WRITE, fileOffsetHigh, fileOffsetLow, uintptr(length))
	if addr == 0 {
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}

	handleLock.Lock()
	handleMap[addr] = h
	handleLock.Unlock()

	mmap := new(File)
	mmap.mmap = unsafe.Pointer(&addr)
	mmap.data = *(*[]byte)(mmap.mmap)
	mmap.size = length

	mmap.data[0] = 1 // test access

	return mmap, nil
}

// Close attempts to unmap the mapped memory region
func (mmap *File) Close() (err error) {
	addr := uintptr(mmap.mmap)

	handleLock.Lock()
	defer handleLock.Unlock()
	err = syscall.UnmapViewOfFile(addr)
	if err != nil {
		return err
	}

	handle, ok := handleMap[addr]
	if !ok {
		// should be impossible; we would've errored above
		return errors.New("unknown base address")
	}
	delete(handleMap, addr)

	e := syscall.CloseHandle(syscall.Handle(handle))
	return os.NewSyscallError("CloseHandle", e)
}
