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
	TRUE  int32 = 1
	FALSE int32 = 0
)

type Bool struct {
	val int32
}

func NewBool(val bool) *Bool {
	if val {
		return &Bool{val: TRUE}
	} else {
		return &Bool{val: FALSE}
	}
}

func (b *Bool) Get() bool {
	return atomic.LoadInt32(&b.val) == TRUE
}

func (b *Bool) Set(val bool) {
	if val {
		atomic.StoreInt32(&b.val, TRUE)
	} else {
		atomic.StoreInt32(&b.val, FALSE)
	}
}

func (b *Bool) CompareAndSet(oldVal, newVal bool) bool {
	var old, newer int32
	if oldVal {
		old = TRUE
	} else {
		old = FALSE
	}

	if newVal {
		newer = TRUE
	} else {
		newer = FALSE
	}

	return atomic.CompareAndSwapInt32(&b.val, old, newer)
}
