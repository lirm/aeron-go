/*
Copyright 2016 Stanislav Liberman
Copyright 2022 Steven Stern

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
	"fmt"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
)

// Read will attempt to read the next frame from the term and invoke the callback if successful.
// Method will return a tuple of new term offset, number of fragments read, and an error reading
// fragments.  An error will interrupt reading before fragmentLimit, but any fragments read prior
// to the error will still be processed, and that number will be returned with the error.
//
//go:norace
func Read(termBuffer *atomic.Buffer, termOffset int32, handler FragmentHandler, fragmentLimit int,
	header *logbuffer.Header) (int32, int, error) {

	capacity := termBuffer.Capacity()

	var fragmentsRead int
	var err error
	for fragmentsRead < fragmentLimit && termOffset < capacity && err == nil {
		frameLength := logbuffer.GetFrameLength(termBuffer, termOffset)
		if frameLength <= 0 {
			break
		}

		fragmentOffset := termOffset
		termOffset += util.AlignInt32(frameLength, logbuffer.FrameAlignment)

		if !logbuffer.IsPaddingFrame(termBuffer, fragmentOffset) {
			header.Wrap(termBuffer.Ptr(), termBuffer.Capacity())
			fragmentsRead++
			header.SetOffset(fragmentOffset)
			err = handler(termBuffer, fragmentOffset+logbuffer.DataFrameHeader.Length,
				frameLength-logbuffer.DataFrameHeader.Length, header)
		}
	}

	return termOffset, fragmentsRead, err
}

// BoundedRead will attempt to read frames from the term up to the specified offsetLimit.
// Method will return a tuple of new term offset and number of fragments read
func BoundedRead(termBuffer *atomic.Buffer, termOffset int32, offsetLimit int32, handler FragmentHandler,
	fragmentLimit int, header *logbuffer.Header) (int32, int, error) {

	var fragmentsRead int
	var err error
	for fragmentsRead < fragmentLimit && termOffset < offsetLimit && err == nil {
		frameLength := logbuffer.GetFrameLength(termBuffer, termOffset)
		if frameLength <= 0 {
			err = fmt.Errorf("invalid frameLength %d", frameLength)
			break
		}

		fragmentOffset := termOffset
		termOffset += util.AlignInt32(frameLength, logbuffer.FrameAlignment)

		if !logbuffer.IsPaddingFrame(termBuffer, fragmentOffset) {
			header.Wrap(termBuffer.Ptr(), termBuffer.Capacity())
			fragmentsRead++
			header.SetOffset(fragmentOffset)
			err = handler(termBuffer, fragmentOffset+logbuffer.DataFrameHeader.Length,
				frameLength-logbuffer.DataFrameHeader.Length, header)
		}
	}

	return termOffset, fragmentsRead, err
}
