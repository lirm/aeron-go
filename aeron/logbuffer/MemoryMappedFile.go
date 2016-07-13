package logbuffer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

type MemoryMappedFile struct {
	mmap unsafe.Pointer
	size int
}

func GetFileSize(filename string) int64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size()
}

func MapExistingMemoryMappedFile(filename string, offset int64, length int) (*MemoryMappedFile, error) {
	//log.Printf("Will try to map existing %s, %d, %d", filename, offset, length)

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
		// FIXME return error
		return nil, errors.New("zero size for existing file?!")
	}
	if size < 0 {
		// FIXME huh?
		return nil, fmt.Errorf("mmap: file %q has negative size", filename)
	}
	if size != int64(int(size)) {
		return nil, fmt.Errorf("mmap: file %q is too large", filename)
	}

	var mapSize int
	if length != 0 {
		mapSize = length
	} else {
		mapSize = int(size)
	}

	//log.Printf("Mapping existing file: fd: %d, size: %d, offset: %d", int(f.Fd()), int(size), offset)
	data, err := syscall.Mmap(int(f.Fd()), offset, mapSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	mmap := new(MemoryMappedFile)
	mmap.mmap = unsafe.Pointer(&data[0])
	mmap.size = len(data)
	//log.Printf("Mapped existing file @%v for %d", mmap.mmap, mmap.size)

	return mmap, nil
}

func CreateMemoryMappedFile(filename string, offset int64, length int) (*MemoryMappedFile, error) {
	//log.Printf("Will try to map new %s, %d, %d", filename, offset, length)

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

	data, err := syscall.Mmap(int(f.Fd()), offset, length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	mmap := new(MemoryMappedFile)
	mmap.mmap = unsafe.Pointer(&data[0])
	mmap.size = len(data)
	//log.Printf("Mapped a new file @%v for %d", mmap.mmap, mmap.size)

	return mmap, nil
}

func (mmap *MemoryMappedFile) GetMemoryPtr() unsafe.Pointer {
	return mmap.mmap
}

func (mmap *MemoryMappedFile) GetMemorySize() int {
	return mmap.size
}
