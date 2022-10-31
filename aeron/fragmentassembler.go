// Copyright 2016 Stanislav Liberman
// Copyright 2022 Talos, Inc.
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

package aeron

import (
	"bytes"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
)

const (
	DefaultFragmentAssemblyBufferLength = int32(4096)

	beginFrag    uint8 = 0x80
	endFrag      uint8 = 0x40
	unfragmented uint8 = 0x80 | 0x40
)

// FragmentAssembler that sits in a chain-of-responsibility pattern that reassembles fragmented messages
// so that the next handler in the chain only sees whole messages.
//
// Unfragmented messages are delegated without copy. Fragmented messages are copied to a temporary
// buffer for reassembly before delegation.
//
// The Header passed to the delegate on assembling a message will be that of the last fragment.
//
// Session based buffers will be allocated and grown as necessary based on the length of messages to be assembled.
type FragmentAssembler struct {
	delegate              term.FragmentHandler
	initialBufferLength   int32
	builderBySessionIdMap map[int32]*bytes.Buffer
}

// NewFragmentAssembler constructs an adapter to reassemble message fragments and delegate on whole messages.
func NewFragmentAssembler(delegate term.FragmentHandler, initialBufferLength int32) *FragmentAssembler {
	return &FragmentAssembler{
		delegate:              delegate,
		initialBufferLength:   initialBufferLength,
		builderBySessionIdMap: make(map[int32]*bytes.Buffer),
	}
}

// Clear removes all existing session buffers.
func (f *FragmentAssembler) Clear() {
	for k := range f.builderBySessionIdMap {
		delete(f.builderBySessionIdMap, k)
	}
}

// OnFragment reassembles and forwards whole messages to the delegate.
func (f *FragmentAssembler) OnFragment(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header) error {
	flags := header.Flags()
	if (flags & unfragmented) == unfragmented {
		f.delegate(buffer, offset, length, header)
	} else {
		if (flags & beginFrag) == beginFrag {
			builder, ok := f.builderBySessionIdMap[header.SessionId()]
			if !ok {
				builder = &bytes.Buffer{}
				f.builderBySessionIdMap[header.SessionId()] = builder
			}
			builder.Reset()
			buffer.WriteBytes(builder, offset, length)
		} else {
			builder, ok := f.builderBySessionIdMap[header.SessionId()]
			if ok && builder.Len() != 0 {
				buffer.WriteBytes(builder, offset, length)
				if (flags & endFrag) == endFrag {
					msgLength := builder.Len()
					f.delegate(
						atomic.MakeBuffer(builder.Bytes(), msgLength),
						int32(0),
						int32(msgLength),
						header)
					builder.Reset()
				}
			}
		}
	}
	return nil
}
