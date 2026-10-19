package term

import (
	"fmt"
	"sync"
	"testing"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/util"
)

// Frame append cost with a concurrent consumer: the frame length written with full barriers
// (atomic.StoreInt32, XCHG on amd64) for both the negative length in the header and the final
// length, as before, versus the store-fence and release stores used now. The load mimics a service
// that does some work, sends a ~100 byte reply (2 MB term), does more work and appends a ~250 byte
// event (4 MB term). A consumer goroutine reads and cleans each publication like the driver does.
// The append cost is the difference from workOnly.
//
//	go test -run '^$' -bench AppendWithConsumer -count 8 ./aeron/logbuffer/term/

// fullBarrierAppendUnfragmentedMessage2 is AppendUnfragmentedMessage2 with both frame length stores
// done by PutInt32Ordered (atomic.StoreInt32).
func fullBarrierAppendUnfragmentedMessage2(appender *Appender,
	srcBufferOne *atomic.Buffer, srcOffsetOne int32, lengthOne int32,
	srcBufferTwo *atomic.Buffer, srcOffsetTwo int32, lengthTwo int32,
) (resultingOffset int64, termID int32) {
	frameLength := lengthOne + lengthTwo + logbuffer.DataFrameHeader.Length
	alignedLength := util.AlignInt32(frameLength, logbuffer.FrameAlignment)
	rawTail := appender.getAndAddRawTail(alignedLength)
	termLength := appender.termBuffer.Capacity()

	termID = logbuffer.TermID(rawTail)
	termOffset := rawTail & 0xFFFFFFFF
	resultingOffset = termOffset + int64(alignedLength)
	if resultingOffset > int64(termLength) {
		return handleEndOfLogCondition(termID, appender.termBuffer, int32(termOffset),
			&appender.headerWriter, termLength), termID
	}
	offset := int32(termOffset)
	dataOffset := offset + logbuffer.DataFrameHeader.Length
	appender.headerWriter.write(appender.termBuffer, offset, frameLength, termID)
	appender.termBuffer.PutInt32Ordered(offset, -frameLength)
	appender.termBuffer.PutBytes(dataOffset, srcBufferOne, srcOffsetOne, lengthOne)
	appender.termBuffer.PutBytes(dataOffset+lengthOne, srcBufferTwo, srcOffsetTwo, lengthTwo)
	reservedValue := DefaultReservedValueSupplier(appender.termBuffer, offset, frameLength)
	appender.termBuffer.PutInt64(offset+logbuffer.DataFrameHeader.ReservedValueFieldOffset, reservedValue)
	appender.termBuffer.PutInt32Ordered(offset, frameLength)
	return resultingOffset, termID
}

type benchmarkAppend func(appender *Appender, one *atomic.Buffer, lengthOne int32,
	two *atomic.Buffer, lengthTwo int32) (int64, int32)

var benchmarkVariants = []struct {
	name   string
	append benchmarkAppend
}{
	{"fullBarrier", func(appender *Appender, one *atomic.Buffer, lengthOne int32, two *atomic.Buffer, lengthTwo int32) (int64, int32) {
		return fullBarrierAppendUnfragmentedMessage2(appender, one, 0, lengthOne, two, 0, lengthTwo)
	}},
	{"release", func(appender *Appender, one *atomic.Buffer, lengthOne int32, two *atomic.Buffer, lengthTwo int32) (int64, int32) {
		return appender.AppendUnfragmentedMessage2(one, 0, lengthOne, two, 0, lengthTwo, DefaultReservedValueSupplier)
	}},
}

var benchmarkWorkState [16 << 10]uint64

// benchmarkWork imitates message processing: arithmetic and stores to its own state, so that the
// store buffer has something to drain.
func benchmarkWork(iterations int) {
	x := uint64(88172645463325252)
	for i := 0; i < iterations; i++ {
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
		benchmarkWorkState[x%uint64(len(benchmarkWorkState))] += x
	}
}

var benchmarkSink int64

func BenchmarkAppendWithConsumer(b *testing.B) {
	for _, iterations := range []int{100, 400} {
		b.Run(fmt.Sprintf("work=%d/workOnly", iterations), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkWork(iterations)
				benchmarkWork(iterations)
			}
		})
		for _, variant := range benchmarkVariants {
			b.Run(fmt.Sprintf("work=%d/%s", iterations, variant.name), func(b *testing.B) {
				reply := newTestPublication(b, 2<<20)
				event := newTestPublication(b, 4<<20)
				// The consumer touches every cache line of a frame, like a sender or a recorder.
				var sums [2]int64
				var consumers sync.WaitGroup
				for index, publication := range []*testPublication{reply, event} {
					consumers.Add(1)
					go func(index int, publication *testPublication) {
						defer consumers.Done()
						publication.consumer.run(func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
							for i := int32(0); i < length; i += 64 {
								sums[index] += int64(buffer.GetUInt8(offset + i))
							}
						})
					}(index, publication)
				}

				replyMessage := make([]byte, 100)
				eventMessage := make([]byte, 250)
				replyOne, replyTwo := atomic.MakeBuffer(replyMessage[:8]), atomic.MakeBuffer(replyMessage[8:])
				eventOne, eventTwo := atomic.MakeBuffer(eventMessage[:8]), atomic.MakeBuffer(eventMessage[8:])
				appendFrame := variant.append
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkWork(iterations)
					reply.offer(alignedFrameLength(100), func(appender *Appender) (int64, int32) {
						return appendFrame(appender, replyOne, 8, replyTwo, 92)
					})
					benchmarkWork(iterations)
					event.offer(alignedFrameLength(250), func(appender *Appender) (int64, int32) {
						return appendFrame(appender, eventOne, 8, eventTwo, 242)
					})
				}
				b.StopTimer()
				reply.consumer.stop.Store(true)
				event.consumer.stop.Store(true)
				consumers.Wait()
				benchmarkSink += sums[0] + sums[1]
			})
		}
	}
}
