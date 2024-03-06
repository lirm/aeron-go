package archive

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive/codecs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestControl_PollForResponse(t *testing.T) {
	t.Run("times out when nothing to poll", func(t *testing.T) {
		control, image := newTestControl(t)
		image.On("Poll", mock.Anything, mock.Anything).Return(0)
		correlations.Store(int64(1), control)
		id, err := control.PollForResponse(1, 0)
		assert.Zero(t, id)
        assert.EqualError(t, err, `timeout waiting for correlationID 1`)
	})

	t.Run("returns the error response", func(t *testing.T) {
		control, image := newTestControl(t)
		mockPollResponses(t, image,
			&codecs.ControlResponse{
				Code:             codecs.ControlResponseCode.ERROR,
				CorrelationId:    1,
				ErrorMessage:     []byte(`b0rk`),
				RelevantId:       3,
			},
		)
		correlations.Store(int64(1), control)
		id, err := control.PollForResponse(1, 0)
		assert.EqualValues(t, 3, id)
		assert.EqualError(t, err, `Control Response failure: b0rk`)
	})

	t.Run("discards all responses preceding the first error", func(t *testing.T) {
		control, image := newTestControl(t)
		mockPollResponses(t, image,
			&codecs.ControlResponse{Code: codecs.ControlResponseCode.OK},
			&codecs.ControlResponse{Code: codecs.ControlResponseCode.OK},
			&codecs.ControlResponse{
				Code:             codecs.ControlResponseCode.ERROR,
				CorrelationId:    1,
				ErrorMessage:     []byte(`b0rk`),
				RelevantId:       3,
			},
		)
		correlations.Store(int64(1), control)
		id, err := control.PollForResponse(1, 0)
		assert.EqualValues(t, 3, id)
		assert.EqualError(t, err, `Control Response failure: b0rk`)
	})

	t.Run("does not process messages after result", func(t *testing.T) {
		control, image := newTestControl(t)
		mockPollResponses(t, image,
			&codecs.ControlResponse{
				Code:             codecs.ControlResponseCode.ERROR,
				CorrelationId:    1,
				ErrorMessage:     []byte(`b0rk`),
				RelevantId:       3,
			},
			&codecs.ControlResponse{Code: codecs.ControlResponseCode.OK},
		)
		correlations.Store(int64(1), control)
		id, err := control.PollForResponse(1, 0)
		assert.EqualValues(t, 3, id)
		assert.EqualError(t, err, `Control Response failure: b0rk`)
		fragments := image.Poll(func(buffer *atomic.Buffer, offset, length int32, header *logbuffer.Header) {}, 1)
		assert.EqualValues(t, 1, fragments)
	})
}

func TestControl_PollForErrorResponse(t *testing.T) {
	t.Run("returns zero when nothing to poll", func(t *testing.T) {
		control, image := newTestControl(t)
		image.On("Poll", mock.Anything, mock.Anything).Return(0)
		cnt, err := control.PollForErrorResponse()
		assert.Zero(t, cnt)
		assert.NoError(t, err)
	})

	t.Run("returns the error response", func(t *testing.T) {
		control, image := newTestControl(t)
		mockPollResponses(t, image,
			&codecs.ControlResponse{
				Code:         codecs.ControlResponseCode.ERROR,
				ErrorMessage: []byte(`b0rk`),
			},
		)
		cnt, err := control.PollForErrorResponse()
		assert.EqualValues(t, 1, cnt)
		assert.EqualError(t, err, `PollForErrorResponse received a ControlResponse (correlationId:0 Code:ERROR error="b0rk"`)
	})

	t.Run("discards all queued responses unless error", func(t *testing.T) {
		control, image := newTestControl(t)
		mockPollResponses(t, image,
			&codecs.ControlResponse{Code: codecs.ControlResponseCode.OK},
			&codecs.ControlResponse{Code: codecs.ControlResponseCode.OK},
		)
		cnt, err := control.PollForErrorResponse()
		assert.EqualValues(t, 2, cnt)
		assert.NoError(t, err)
	})

	t.Run("does not process messages after error", func(t *testing.T) {
		control, image := newTestControl(t)
		mockPollResponses(t, image,
			&codecs.ControlResponse{
				Code:         codecs.ControlResponseCode.ERROR,
				ErrorMessage: []byte(`b0rk`),
			},
			&codecs.ControlResponse{Code: codecs.ControlResponseCode.OK},
		)
		cnt, err := control.PollForErrorResponse()
		assert.EqualValues(t, 1, cnt)
		assert.Error(t, err)
		cnt, err = control.PollForErrorResponse()
		assert.EqualValues(t, 1, cnt)
		assert.NoError(t, err)
	})
}

func mockPollResponses(t *testing.T, image *aeron.MockImage, responses ...encodable) {
	poll := image.On("Poll", mock.Anything, mock.Anything)
	poll.Run(func(args mock.Arguments) {
		fragmentCount := args.Get(1).(int)
		count := 0
		for ; count < fragmentCount; count++ {
			if len(responses) == 0 {
				break
			}
			buffer := encode(t, responses[0])
			responses = responses[1:]
			args.Get(0).(term.FragmentHandler)(buffer, 0, buffer.Capacity(), newTestHeader())
		}
		poll.Return(count)
	})
}

type encodable interface {
	SbeBlockLength() uint16
	SbeTemplateId() uint16
	SbeSchemaId() uint16
	SbeSchemaVersion() uint16
	Encode(*codecs.SbeGoMarshaller, io.Writer, bool) error
}

func encode(t *testing.T, data encodable) *atomic.Buffer {
	m := codecs.NewSbeGoMarshaller()
	buf := new(bytes.Buffer)
	header := codecs.MessageHeader{
		BlockLength: data.SbeBlockLength(),
		TemplateId:  data.SbeTemplateId(),
		SchemaId:    data.SbeSchemaId(),
		Version:     data.SbeSchemaVersion(),
	}
	if !assert.NoError(t, header.Encode(m, buf)) {
		return nil
	}
	if !assert.NoError(t, data.Encode(m, buf, false)) {
		return nil
	}
	return atomic.MakeBuffer(buf.Bytes())
}

func newTestControl(t *testing.T) (*Control, *aeron.MockImage) {
	image := aeron.NewMockImage(t)
	c := &Control{
		Subscription: newTestSub(image),
	}
	c.archive = &Archive{
		Listeners: &ArchiveListeners{},
		Options:   &Options{
			Timeout: 100*time.Millisecond,
			IdleStrategy: idlestrategy.Yielding{},
		},
	}
	c.fragmentAssembler = aeron.NewControlledFragmentAssembler(
		c.onFragment, aeron.DefaultFragmentAssemblyBufferLength,
	)
	return c, image
}

func newTestSub(image aeron.Image) *aeron.Subscription {
	images := aeron.NewImageList()
	images.Set([]aeron.Image{image})
	sub := aeron.NewSubscription(nil, "", 1, 2, 3, nil, nil)
	rsub := reflect.ValueOf(sub)
	rfimages := rsub.Elem().FieldByName("images")
	rfimages = reflect.NewAt(rfimages.Type(), unsafe.Pointer(rfimages.UnsafeAddr())).Elem()
	rfimages.Set(reflect.ValueOf(images))
	return sub
}

func newTestHeader() *logbuffer.Header {
	buffer := atomic.MakeBuffer(make([]byte, logbuffer.DataFrameHeader.Length))
	buffer.PutUInt8(logbuffer.DataFrameHeader.FlagsFieldOffset, 0xc0) // unfragmented
	return new(logbuffer.Header).Wrap(buffer.Ptr(), buffer.Capacity())
}
