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

package aeron

import (
	"sync/atomic"
	"unsafe"
)

// ImageList is a helper class to manage list of images atomically without locks
type ImageList struct {
	img unsafe.Pointer
}

// NewImageList is a factory method for ImageList
func NewImageList() *ImageList {
	list := new(ImageList)

	var images []Image
	list.img = unsafe.Pointer(&images)

	return list
}

// Get returns a pointer to the underlying image array loaded atomically
func (l *ImageList) Get() []Image {
	return *(*[]Image)(atomic.LoadPointer(&l.img))
}

// Set atomically sets the reference to the underlying array
func (l *ImageList) Set(imgs []Image) {
	atomic.StorePointer(&l.img, unsafe.Pointer(&imgs))
}

// Empty is a convenience method to reset the contents of the list
func (l *ImageList) Empty() (oldList []Image) {
	oldList = l.Get()
	l.Set(make([]Image, 0))
	return
}
