package logbuffer

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/buffers"
	"testing"
)

var MMAP string = "mmap.bin"

func TestMmapBasics(t *testing.T) {
	t.Logf("Beginning test")
	mmap, err := CreateMemoryMappedFile(MMAP, 0, 8000)
	// t.Logf("mmap: %v, len:%d, err:%v", mmap.GetMemoryPtr(), mmap.GetMemorySize(), err)
	fmt.Printf("mmap: %v, len:%d, err:%v\n", mmap.GetMemoryPtr(), mmap.GetMemorySize(), err)

	buf := buffers.MakeAtomic(mmap.GetMemoryPtr(), mmap.GetMemorySize())
	// t.Logf("created atomic buffer: %v", buf)
	fmt.Printf("created atomic buffer: %v\n", buf)

	//    buf.Fill(0)

	num := buf.GetInt64Volatile(0)
	// t.Logf("1: Got %d stuff from file", num)
	fmt.Printf("1: Got %d stuff from file\n", num)

	buf.PutInt64Ordered(0, 42)
	num = buf.GetInt64Volatile(0)
	t.Logf("2: Got %d stuff from file\n", num)

}
