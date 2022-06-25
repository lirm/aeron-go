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

// ControlledFragmentAssembler that sits in a chain-of-responsibility pattern that reassembles fragmented messages
// so that the next handler in the chain only sees whole messages.
//
// Unfragmented messages are delegated without copy. Fragmented messages are copied to a temporary
// buffer for reassembly before delegation.
//
// The Header passed to the delegate on assembling a message will be that of the last fragment.
//
// Session based buffers will be allocated and grown as necessary based on the length of messages to be assembled.
type ControlledFragmentAssembler struct {
	delegate              term.ControlledFragmentHandler
	initialBufferLength   int32
	builderBySessionIdMap map[int32]*bytes.Buffer
}

// NewControlledFragmentAssembler constructs an adapter to reassemble message fragments and delegate on whole messages.
func NewControlledFragmentAssembler(delegate term.ControlledFragmentHandler, initialBufferLength int32) *ControlledFragmentAssembler {
	return &ControlledFragmentAssembler{
		delegate:              delegate,
		initialBufferLength:   initialBufferLength,
		builderBySessionIdMap: make(map[int32]*bytes.Buffer),
	}
}

// OnFragment reassembles and forwards whole messages to the delegate.
func (f *ControlledFragmentAssembler) OnFragment(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header) (action term.ControlledPollAction) {
	flags := header.Flags()
	action = term.ControlledPollActionContinue

	if (flags & unfragmented) == unfragmented {
		action = f.delegate(buffer, offset, length, header)
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
			if ok {
				if limit := builder.Len(); limit > 0 {
					buffer.WriteBytes(builder, offset, length)
					if (flags & endFrag) == endFrag {
						msgLength := builder.Len()
						action := f.delegate(
							atomic.MakeBuffer(builder.Bytes(), msgLength),
							int32(0),
							int32(msgLength),
							header)
						if action == term.ControlledPollActionAbort {
							builder.Truncate(limit)
						} else {
							builder.Reset()
						}
					}
				}
			}
		}
	}
	return
}
