package term

import (
	"path/filepath"
	"runtime"
	goatomic "sync/atomic"
	"testing"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
)

// Harness for frame append tests: a three-term log in an mmap'ed file, a writer that rotates terms
// like aeron.Publication, and a consumer that cleans terms behind itself like a driver IPC
// publication (cleanBufferTo). The consumer spins on another core, so the term lines the writer
// stores to arrive from a foreign cache, as they do when the media driver reads and cleans them.

const testInitialTermID int32 = 100

func newTestLogBuffers(tb testing.TB, termLength int32) *logbuffer.LogBuffers {
	tb.Helper()
	fileName := filepath.Join(tb.TempDir(), "publication.logbuffer")
	logLength := int(termLength)*logbuffer.PartitionCount + int(logbuffer.LogMetaDataLength)
	file, err := memmap.NewFile(fileName, 0, logLength)
	if err != nil {
		tb.Fatal(err)
	}
	if err := file.Close(); err != nil {
		tb.Fatal(err)
	}
	logBuffers := logbuffer.Wrap(fileName)
	tb.Cleanup(func() { _ = logBuffers.Close() })

	meta := logBuffers.Meta()
	meta.TermLen.Set(termLength)
	meta.InitTermID.Set(testInitialTermID)
	// Tails of the inactive terms as in initialiseTailWithTermId in Java: RotateLog expects a term ID
	// PartitionCount less than the next one.
	for index := 0; index < logbuffer.PartitionCount; index++ {
		termID := testInitialTermID
		if index > 0 {
			termID = testInitialTermID - logbuffer.PartitionCount + int32(index)
		}
		meta.TailCounter[index].Set(int64(termID) << 32)
	}
	header := meta.DefaultFrameHeader.Get()
	header.PutInt32(logbuffer.DataFrameHeader.SessionIDFieldOffset, 7)
	header.PutInt32(logbuffer.DataFrameHeader.StreamIDFieldOffset, 1001)
	return logBuffers
}

type testPublication struct {
	meta       *logbuffer.LogBufferMetaData
	appenders  [logbuffer.PartitionCount]*Appender
	termLength int64
	consumer   *testConsumer
	// limit is a cached publication limit: the writer re-reads the consumer position only when it
	// hits the limit, like aeron.Publication with a pubLimit the driver does not update per frame.
	limit int64
}

func newTestPublication(tb testing.TB, termLength int32) *testPublication {
	logBuffers := newTestLogBuffers(tb, termLength)
	publication := &testPublication{meta: logBuffers.Meta(), termLength: int64(termLength)}
	for index := range publication.appenders {
		publication.appenders[index] = MakeAppender(logBuffers, index)
	}
	publication.consumer = &testConsumer{logBuffers: logBuffers, termLength: int64(termLength)}
	return publication
}

// offer appends a frame to the active term with write; when the term trips it rotates the log and
// retries, like Publication.Offer after AdminAction. The writer stays at most one term ahead of the
// consumer, which is the publication limit. It returns false if the consumer stopped and the frame
// was not written.
func (publication *testPublication) offer(requiredLength int32,
	write func(appender *Appender) (resultingOffset int64, termID int32)) bool {
	for {
		termCount := publication.meta.ActiveTermCountOff.Get()
		appender := publication.appenders[termCount%logbuffer.PartitionCount]
		termOffset := appender.RawTail() & 0xFFFFFFFF
		if termOffset > publication.termLength {
			termOffset = publication.termLength
		}
		position := int64(termCount)*publication.termLength + termOffset
		for position+int64(requiredLength) > publication.limit {
			if publication.consumer.stop.Load() {
				return false
			}
			publication.limit = publication.consumer.position.Load() + publication.termLength
			if position+int64(requiredLength) > publication.limit {
				runtime.Gosched()
			}
		}
		resultingOffset, termID := write(appender)
		if resultingOffset == AppenderTripped {
			logbuffer.RotateLog(publication.meta, termCount, termID)
			continue
		}
		if resultingOffset == AppenderFailed {
			panic("append past the end of the term")
		}
		return true
	}
}

type testConsumer struct {
	logBuffers *logbuffer.LogBuffers
	termLength int64
	position   goatomic.Int64
	stop       goatomic.Bool
}

// run reads frames with term.Read and zeroes the log one term behind itself until stop is set.
func (consumer *testConsumer) run(handler FragmentHandler) {
	const cleanChunk = 4096
	var position, cleanPosition int64
	var header logbuffer.Header
	header.SetInitialTermID(testInitialTermID)
	header.SetPositionBitsToShift(int32(util.NumberOfTrailingZeroes(uint32(consumer.termLength))))
	for !consumer.stop.Load() {
		termBuffer := consumer.logBuffers.Buffer(int((position / consumer.termLength) % logbuffer.PartitionCount))
		termOffset := int32(position % consumer.termLength)
		newOffset, _ := Read(termBuffer, termOffset, handler, 16, &header)
		if newOffset == termOffset {
			continue
		}
		position += int64(newOffset - termOffset)
		consumer.position.Store(position)
		for position-consumer.termLength-cleanPosition >= cleanChunk {
			cleanBuffer := consumer.logBuffers.Buffer(int((cleanPosition / consumer.termLength) % logbuffer.PartitionCount))
			cleanOffset := int32(cleanPosition % consumer.termLength)
			clearTerm(cleanBuffer, cleanOffset, cleanChunk)
			cleanPosition += cleanChunk
		}
	}
}

func clearTerm(buffer *atomic.Buffer, offset int32, length int32) {
	for i := int32(0); i < length; i += 8 {
		buffer.PutInt64(offset+i, 0)
	}
}

func alignedFrameLength(payloadLength int32) int32 {
	return util.AlignInt32(payloadLength+logbuffer.DataFrameHeader.Length, logbuffer.FrameAlignment)
}
