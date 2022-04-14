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
	"fmt"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"testing"
)

const (
	MMAP = "mmap.bin"
)

func TestMmapBasics(t *testing.T) {
	t.Log("Beginning test")
	mmap, err := NewFile(MMAP, 0, 8000)
	defer mmap.Close()
	// t.Logf("mmap: %v, len:%d, err:%v", mmap.GetMemoryPtr(), mmap.GetMemorySize(), err)
	fmt.Printf("mmap: %v, len:%d, err:%v\n", mmap.GetMemoryPtr(), mmap.GetMemorySize(), err)

	buf := atomic.MakeBuffer(mmap.GetMemoryPtr(), mmap.GetMemorySize())
	// t.Logf("created atomic buffer: %v", buf)
	fmt.Printf("created atomic buffer: %v\n", buf)

	//    buf.Fill(0)

	num := buf.GetInt64Volatile(0)
	// t.Logf("1: Got %d stuff from file", num)
	fmt.Printf("1: Got %d stuff from file\n", num)

	var expected int64 = 42
	buf.PutInt64Ordered(0, expected)
	num = buf.GetInt64Volatile(0)
	t.Logf("2: Got %d stuff from file\n", num)
	if num != expected {
		t.Fail()
	}
}
