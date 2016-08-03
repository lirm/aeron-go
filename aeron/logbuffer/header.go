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

package logbuffer

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"unsafe"
)

type Header struct {
	buffer              atomic.Buffer
	offset              int32
	initialTermId       int32
	positionBitsToShift int32
}

func (hdr *Header) Wrap(ptr unsafe.Pointer, length int32) *Header {
	hdr.buffer.Wrap(ptr, length)
	return hdr
}

func (hdr *Header) SetOffset(offset int32) *Header {
	hdr.offset = offset
	return hdr
}

func (hdr *Header) SetInitialTermId(initialTermId int32) *Header {
	hdr.initialTermId = initialTermId
	return hdr
}

func (hdr *Header) SetPositionBitsToShift(positionBitsToShift int32) *Header {
	hdr.positionBitsToShift = positionBitsToShift
	return hdr
}
