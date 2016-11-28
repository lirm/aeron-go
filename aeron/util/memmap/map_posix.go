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

// +build darwin linux freebsd

package memmap

import (
	"syscall"
	"unsafe"
)

func doMap(fd int, offset int64, length int) (*File, error) {
	bytes, err := syscall.Mmap(fd, offset, length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	mmap := new(File)
	mmap.data = bytes
	mmap.mmap = unsafe.Pointer(&mmap.data[0])
	mmap.size = len(mmap.data)

	return mmap, err
}

// Close attempts to unmap the mapped memory region
func (mmap *File) Close() error {
	err := syscall.Munmap(mmap.data)
	if err != nil {
		logger.Errorf("Error unmapping file: %v", err)
	} else {
		mmap.data = nil
		mmap.mmap = nil
		logger.Debugf("Unmapped: %v", mmap)
	}
	return err
}
