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
	// True value for atomic.Bool
	True int32 = 1
	// False value for atomic.Bool
	False int32 = 0
)

// Bool is an atomic boolean implementation used by aeron-go
type Bool struct {
	val int32
}

// Get returns the current state of the variable
func (b *Bool) Get() bool {
	return atomic.LoadInt32(&b.val) == True
}

// Set atomically sets the value of the variable
func (b *Bool) Set(val bool) {
	if val {
		atomic.StoreInt32(&b.val, True)
	} else {
		atomic.StoreInt32(&b.val, False)
	}
}

// CompareAndSet performs an atomic CAS operation on this variable
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
