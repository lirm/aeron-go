package term

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
	goatomic "sync/atomic"
	"testing"
	"time"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
)

// Message seq: the first 8 bytes hold seq, followed by bytes byte(seq*31+i). The length and the
// append method depend on seq so that a single run covers every path that publishes a frame length:
// AppendUnfragmentedMessage, AppendUnfragmentedMessage2, AppendFragmentedMessage2, Claim+Commit and
// padding at the end of a term.

const stressMaxPayloadLength int32 = 192

func stressMessageLength(seq int64) int32 {
	if seq%10 == 9 {
		return 200 + int32(seq*53%600)
	}
	return 9 + int32(seq*37%int64(stressMaxPayloadLength-9))
}

func fillStressMessage(message []byte, seq int64) {
	binary.LittleEndian.PutUint64(message, uint64(seq))
	for i := 8; i < len(message); i++ {
		message[i] = byte(seq*31 + int64(i))
	}
}

func checkStressMessage(message []byte, expectedSeq int64) error {
	if len(message) < 8 {
		return fmt.Errorf("message %d: length %d is less than 8", expectedSeq, len(message))
	}
	seq := int64(binary.LittleEndian.Uint64(message))
	if seq != expectedSeq {
		return fmt.Errorf("expected message %d, read %d", expectedSeq, seq)
	}
	if int32(len(message)) != stressMessageLength(seq) {
		return fmt.Errorf("message %d: length %d, expected %d", seq, len(message), stressMessageLength(seq))
	}
	for i := 8; i < len(message); i++ {
		if message[i] != byte(seq*31+int64(i)) {
			return fmt.Errorf("message %d: byte %d = %d, expected %d", seq, i, message[i], byte(seq*31+int64(i)))
		}
	}
	return nil
}

func requiredFragmentedLength(length int32) int32 {
	numMaxPayloads := length / stressMaxPayloadLength
	required := numMaxPayloads * alignedFrameLength(stressMaxPayloadLength)
	if remaining := length % stressMaxPayloadLength; remaining > 0 {
		required += alignedFrameLength(remaining)
	}
	return required
}

// TestAppenderConcurrentReader runs a writer and a term.Read reader on different goroutines. Once the
// reader sees a positive frame length it must see the whole frame: the header and every message byte.
func TestAppenderConcurrentReader(t *testing.T) {
	if runtime.GOMAXPROCS(0) < 2 {
		t.Skip("needs two threads: a writer and a reader")
	}
	messages := int64(1_000_000)
	if testing.Short() {
		messages = 100_000
	}
	// Minimum term length: the log rotates thousands of times per run, with padding in every term.
	publication := newTestPublication(t, logbuffer.TermMinLength)

	var lastSeq goatomic.Int64
	var failure goatomic.Pointer[error]
	var assembly []byte
	expectedSeq := int64(1)
	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		flags := header.Flags()
		if flags&beginFrag != 0 {
			assembly = assembly[:0]
		}
		assembly = append(assembly, buffer.GetBytesArray(offset, length)...)
		if flags&endFrag == 0 {
			return
		}
		if err := checkStressMessage(assembly, expectedSeq); err != nil {
			failure.CompareAndSwap(nil, &err)
			publication.consumer.stop.Store(true)
			return
		}
		lastSeq.Store(expectedSeq)
		expectedSeq++
	}
	var consumerDone sync.WaitGroup
	consumerDone.Add(1)
	go func() {
		defer consumerDone.Done()
		publication.consumer.run(handler)
	}()

	message := make([]byte, 1024)
	one := atomic.MakeBuffer(message[:8])
	two := atomic.MakeBuffer(message[8:])
	whole := atomic.MakeBuffer(message)
	var claim logbuffer.Claim
	for seq := int64(1); seq <= messages; seq++ {
		length := stressMessageLength(seq)
		fillStressMessage(message[:length], seq)
		var written bool
		switch seq % 10 {
		case 9:
			written = publication.offer(requiredFragmentedLength(length), func(appender *Appender) (int64, int32) {
				return appender.AppendFragmentedMessage2(one, 0, 8, two, 0, length-8,
					stressMaxPayloadLength, DefaultReservedValueSupplier)
			})
		case 6, 7:
			written = publication.offer(alignedFrameLength(length), func(appender *Appender) (int64, int32) {
				resultingOffset, termID := appender.Claim(length, &claim)
				if resultingOffset > 0 {
					claim.Buffer().PutBytes(logbuffer.DataFrameHeader.Length, whole, 0, length)
					claim.Commit()
				}
				return resultingOffset, termID
			})
		case 8:
			written = publication.offer(alignedFrameLength(length), func(appender *Appender) (int64, int32) {
				return appender.AppendUnfragmentedMessage(whole, 0, length, DefaultReservedValueSupplier)
			})
		default:
			written = publication.offer(alignedFrameLength(length), func(appender *Appender) (int64, int32) {
				return appender.AppendUnfragmentedMessage2(one, 0, 8, two, 0, length-8, DefaultReservedValueSupplier)
			})
		}
		if !written {
			break
		}
	}

	deadline := time.Now().Add(30 * time.Second)
	for lastSeq.Load() < messages && failure.Load() == nil && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	publication.consumer.stop.Store(true)
	consumerDone.Wait()
	if err := failure.Load(); err != nil {
		t.Fatal(*err)
	}
	if lastSeq.Load() != messages {
		t.Fatalf("read %d of %d messages", lastSeq.Load(), messages)
	}
}
