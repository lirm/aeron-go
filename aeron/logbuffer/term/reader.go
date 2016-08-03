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

package term

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
)

type FragmentHandler func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header)

type ReadOutcome struct {
	offset        int32
	fragmentsRead int
}

func (outcome *ReadOutcome) Offset() int32 {
	return outcome.offset
}

func (outcome *ReadOutcome) FragmentsRead() int {
	return outcome.fragmentsRead
}

func Read(outcome *ReadOutcome, termBuffer *atomic.Buffer, termOffset int32,
	handler FragmentHandler, fragmentsLimit int, header *logbuffer.Header) {

	outcome.fragmentsRead = 0
	outcome.offset = termOffset

	capacity := termBuffer.Capacity()

	for outcome.fragmentsRead < fragmentsLimit {
		frameLength := logbuffer.FrameLengthVolatile(termBuffer, termOffset)
		if frameLength <= 0 {
			break
		}

		fragmentOffset := termOffset
		termOffset += util.AlignInt32(frameLength, logbuffer.FrameDescriptor.Alignment)

		if !logbuffer.IsPaddingFrame(termBuffer, fragmentOffset) {
			header.Wrap(termBuffer.Ptr(), termBuffer.Capacity())
			header.SetOffset(fragmentOffset)
			handler(termBuffer, fragmentOffset+logbuffer.DataFrameHeader.Length,
				frameLength-logbuffer.DataFrameHeader.Length, header)

			outcome.fragmentsRead++
		}

		if termOffset >= capacity {
			break
		}
	}

	outcome.offset = termOffset
}
