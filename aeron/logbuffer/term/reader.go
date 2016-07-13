package term

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
	"log"
)

type FragmentHandler func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header)

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

func Read(outcome *ReadOutcome, termBuffer *buffers.Atomic, termOffset int32,
	handler FragmentHandler, fragmentsLimit int, header *logbuffer.Header, exceptionHandler func(error)) {

	outcome.fragmentsRead = 0
	outcome.offset = termOffset

	capacity := termBuffer.Capacity()
	/*
	   try
	   {
	*/

	// log.Printf("outcome.fragmentsRead: %d, fragmentsLimit: %d, termOffset: %d, capacity: %d",
	// outcome.fragmentsRead, fragmentsLimit, termOffset, capacity)
	for outcome.fragmentsRead < fragmentsLimit {
		// log.Printf("Reading at %d from %v", termOffset, termBuffer)
		frameLength := logbuffer.FrameLengthVolatile(termBuffer, termOffset)
		if frameLength <= 0 {
			break
		}

		fragmentOffset := termOffset
		termOffset += util.AlignInt32(frameLength, logbuffer.FrameDescriptor.FRAME_ALIGNMENT)

		if !logbuffer.IsPaddingFrame(termBuffer, fragmentOffset) {
			//log.Printf("Reading from termOffset: %d, frameLength: %d, fragmentOffset: %d", termOffset, frameLength, fragmentOffset)
			header.Wrap(termBuffer.Ptr(), termBuffer.Capacity())
			header.SetOffset(fragmentOffset)
			handler(termBuffer, fragmentOffset+logbuffer.DataFrameHeader.LENGTH, frameLength-logbuffer.DataFrameHeader.LENGTH,
				header)

			outcome.fragmentsRead++
		} else {
			log.Printf("Got padding frame at %d", termOffset)
		}

		// log.Printf("got %d/%d fragments", outcome.fragmentsRead, fragmentsLimit)

		if termOffset >= capacity {
			break
		}
	}

	/*
	   }
	   catch (std::exception ex)
	   {
	       exceptionHandler(ex);
	   }
	*/
	outcome.offset = termOffset
}
