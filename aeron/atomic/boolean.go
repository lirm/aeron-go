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

package atomic

import (
	"sync/atomic"
)

const (
	True  int32 = 1
	False int32 = 0
)

type Bool struct {
	val int32
}

func NewBool(val bool) *Bool {
	if val {
		return &Bool{val: True}
	}
	return &Bool{val: False}
}

func (b *Bool) Get() bool {
	return atomic.LoadInt32(&b.val) == True
}

func (b *Bool) Set(val bool) {
	if val {
		atomic.StoreInt32(&b.val, True)
	} else {
		atomic.StoreInt32(&b.val, False)
	}
}

func (b *Bool) CompareAndSet(oldVal, newVal bool) bool {
	var old, newer int32
	if oldVal {
		old = True
	} else {
		old = False
	}

	if newVal {
		newer = True
	} else {
		newer = False
	}

	return atomic.CompareAndSwapInt32(&b.val, old, newer)
}
