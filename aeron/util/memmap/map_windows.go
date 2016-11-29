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
	mapper "github.com/edsrzf/mmap-go"
	"os"
	"unsafe"
)

// We keep this map so that we can get back the original handle from the memory address.
var memories = make(map[unsafe.Pointer]mapper.MMap)

func doMap(f *os.File, offset int64, length int) (*File, error) {

	mm, err := mapper.MapRegion(f, length, mapper.RDWR, 0, offset)
	if err != nil {
		return nil, err
	}

	mmap := new(File)
	mmap.data = mm
	mmap.mmap = unsafe.Pointer(&mm[0])
	mmap.size = length

	mmap.data[0] = 1 // test access

	memories[mmap.mmap] = mm

	return mmap, err
}

// Close attempts to unmap the mapped memory region
func (mmap *File) Close() (err error) {
	mm := memories[mmap.mmap]
	return mm.Unmap()
}
