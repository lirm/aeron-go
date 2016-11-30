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

package memmap

import (
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// File is the wrapper around memory mapped file
type File struct {
	mmap unsafe.Pointer
	data []byte
	size int
}

var logger = logging.MustGetLogger("memmap")

// GetFileSize is a helper function to retrieve file size
func GetFileSize(filename string) int64 {
	file, err := os.Open(filename)
	if err != nil {
		logger.Error(err)
		return -1
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		logger.Fatal(err)
		return -1
	}

	return fi.Size()
}

// MapExisting is the factory method used to create memory maps for existing file. This function will fail if
// the file does not exist.
func MapExisting(filename string, offset int64, length int) (*File, error) {
	logger.Debugf("Will try to map existing %s, %d, %d", filename, offset, length)

	f, err := os.OpenFile(filename, syscall.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	if size == 0 {
		return nil, errors.New("zero size for existing file")
	}
	if size < 0 {
		return nil, fmt.Errorf("mmap: stat %q returned %d", filename, size)
	}
	if size != int64(size) {
		return nil, fmt.Errorf("mmap: file %q is too large", filename)
	}

	var mapSize int
	if length != 0 {
		mapSize = length
	} else {
		mapSize = int(size)
	}

	logger.Debugf("Mapping existing file: fd: %d, size: %d, offset: %d", f.Fd(), size, offset)
	mmap := new(File)
	mmap.data, err = syscall.Mmap(int(f.Fd()), offset, mapSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	mmap.mmap = unsafe.Pointer(&mmap.data[0])
	mmap.size = len(mmap.data)
	logger.Debugf("Mapped existing file @%v for %d", mmap.mmap, mmap.size)

	return mmap, err
}

// NewFile is a factory method to create a new memory mapped file with the specified capacity
func NewFile(filename string, offset int64, length int) (*File, error) {
	logger.Debugf("Will try to map new %s, %d, %d", filename, offset, length)

	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Seek(int64(length-1), 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write([]byte("\000"))
	if err != nil {
		log.Fatal(err)
	}

	mmap := new(File)
	mmap.data, err = syscall.Mmap(int(f.Fd()), offset, length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	mmap.mmap = unsafe.Pointer(&mmap.data[0])
	mmap.size = len(mmap.data)
	logger.Debugf("Mapped a new file @%v for %d", mmap.mmap, mmap.size)

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

// GetMemoryPtr return the pointer to the mapped region of the file
func (mmap *File) GetMemoryPtr() unsafe.Pointer {
	return mmap.mmap
}

// GetMemorySize return the size of the mapped region of the file
func (mmap *File) GetMemorySize() int {
	return mmap.size
}
