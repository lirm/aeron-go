// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type FindLastMatchingRecordingRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	MinRecordingId   int64
	SessionId        int32
	StreamId         int32
	Channel          []uint8
}

func (f *FindLastMatchingRecordingRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := f.RangeCheck(f.SbeSchemaVersion(), f.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, f.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, f.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, f.MinRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, f.SessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, f.StreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(f.Channel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, f.Channel); err != nil {
		return err
	}
	return nil
}

func (f *FindLastMatchingRecordingRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !f.ControlSessionIdInActingVersion(actingVersion) {
		f.ControlSessionId = f.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &f.ControlSessionId); err != nil {
			return err
		}
	}
	if !f.CorrelationIdInActingVersion(actingVersion) {
		f.CorrelationId = f.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &f.CorrelationId); err != nil {
			return err
		}
	}
	if !f.MinRecordingIdInActingVersion(actingVersion) {
		f.MinRecordingId = f.MinRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &f.MinRecordingId); err != nil {
			return err
		}
	}
	if !f.SessionIdInActingVersion(actingVersion) {
		f.SessionId = f.SessionIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &f.SessionId); err != nil {
			return err
		}
	}
	if !f.StreamIdInActingVersion(actingVersion) {
		f.StreamId = f.StreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &f.StreamId); err != nil {
			return err
		}
	}
	if actingVersion > f.SbeSchemaVersion() && blockLength > f.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-f.SbeBlockLength()))
	}

	if f.ChannelInActingVersion(actingVersion) {
		var ChannelLength uint32
		if err := _m.ReadUint32(_r, &ChannelLength); err != nil {
			return err
		}
		if cap(f.Channel) < int(ChannelLength) {
			f.Channel = make([]uint8, ChannelLength)
		}
		f.Channel = f.Channel[:ChannelLength]
		if err := _m.ReadBytes(_r, f.Channel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := f.RangeCheck(actingVersion, f.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (f *FindLastMatchingRecordingRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if f.ControlSessionIdInActingVersion(actingVersion) {
		if f.ControlSessionId < f.ControlSessionIdMinValue() || f.ControlSessionId > f.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on f.ControlSessionId (%v < %v > %v)", f.ControlSessionIdMinValue(), f.ControlSessionId, f.ControlSessionIdMaxValue())
		}
	}
	if f.CorrelationIdInActingVersion(actingVersion) {
		if f.CorrelationId < f.CorrelationIdMinValue() || f.CorrelationId > f.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on f.CorrelationId (%v < %v > %v)", f.CorrelationIdMinValue(), f.CorrelationId, f.CorrelationIdMaxValue())
		}
	}
	if f.MinRecordingIdInActingVersion(actingVersion) {
		if f.MinRecordingId < f.MinRecordingIdMinValue() || f.MinRecordingId > f.MinRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on f.MinRecordingId (%v < %v > %v)", f.MinRecordingIdMinValue(), f.MinRecordingId, f.MinRecordingIdMaxValue())
		}
	}
	if f.SessionIdInActingVersion(actingVersion) {
		if f.SessionId < f.SessionIdMinValue() || f.SessionId > f.SessionIdMaxValue() {
			return fmt.Errorf("Range check failed on f.SessionId (%v < %v > %v)", f.SessionIdMinValue(), f.SessionId, f.SessionIdMaxValue())
		}
	}
	if f.StreamIdInActingVersion(actingVersion) {
		if f.StreamId < f.StreamIdMinValue() || f.StreamId > f.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on f.StreamId (%v < %v > %v)", f.StreamIdMinValue(), f.StreamId, f.StreamIdMaxValue())
		}
	}
	return nil
}

func FindLastMatchingRecordingRequestInit(f *FindLastMatchingRecordingRequest) {
	return
}

func (*FindLastMatchingRecordingRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*FindLastMatchingRecordingRequest) SbeTemplateId() (templateId uint16) {
	return 16
}

func (*FindLastMatchingRecordingRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*FindLastMatchingRecordingRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*FindLastMatchingRecordingRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*FindLastMatchingRecordingRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*FindLastMatchingRecordingRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (f *FindLastMatchingRecordingRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= f.ControlSessionIdSinceVersion()
}

func (*FindLastMatchingRecordingRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*FindLastMatchingRecordingRequest) ControlSessionIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*FindLastMatchingRecordingRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*FindLastMatchingRecordingRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*FindLastMatchingRecordingRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*FindLastMatchingRecordingRequest) CorrelationIdId() uint16 {
	return 2
}

func (*FindLastMatchingRecordingRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (f *FindLastMatchingRecordingRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= f.CorrelationIdSinceVersion()
}

func (*FindLastMatchingRecordingRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*FindLastMatchingRecordingRequest) CorrelationIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*FindLastMatchingRecordingRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*FindLastMatchingRecordingRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*FindLastMatchingRecordingRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*FindLastMatchingRecordingRequest) MinRecordingIdId() uint16 {
	return 3
}

func (*FindLastMatchingRecordingRequest) MinRecordingIdSinceVersion() uint16 {
	return 0
}

func (f *FindLastMatchingRecordingRequest) MinRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= f.MinRecordingIdSinceVersion()
}

func (*FindLastMatchingRecordingRequest) MinRecordingIdDeprecated() uint16 {
	return 0
}

func (*FindLastMatchingRecordingRequest) MinRecordingIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*FindLastMatchingRecordingRequest) MinRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*FindLastMatchingRecordingRequest) MinRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*FindLastMatchingRecordingRequest) MinRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*FindLastMatchingRecordingRequest) SessionIdId() uint16 {
	return 4
}

func (*FindLastMatchingRecordingRequest) SessionIdSinceVersion() uint16 {
	return 0
}

func (f *FindLastMatchingRecordingRequest) SessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= f.SessionIdSinceVersion()
}

func (*FindLastMatchingRecordingRequest) SessionIdDeprecated() uint16 {
	return 0
}

func (*FindLastMatchingRecordingRequest) SessionIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*FindLastMatchingRecordingRequest) SessionIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*FindLastMatchingRecordingRequest) SessionIdMaxValue() int32 {
	return math.MaxInt32
}

func (*FindLastMatchingRecordingRequest) SessionIdNullValue() int32 {
	return math.MinInt32
}

func (*FindLastMatchingRecordingRequest) StreamIdId() uint16 {
	return 5
}

func (*FindLastMatchingRecordingRequest) StreamIdSinceVersion() uint16 {
	return 0
}

func (f *FindLastMatchingRecordingRequest) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= f.StreamIdSinceVersion()
}

func (*FindLastMatchingRecordingRequest) StreamIdDeprecated() uint16 {
	return 0
}

func (*FindLastMatchingRecordingRequest) StreamIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*FindLastMatchingRecordingRequest) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*FindLastMatchingRecordingRequest) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*FindLastMatchingRecordingRequest) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*FindLastMatchingRecordingRequest) ChannelMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*FindLastMatchingRecordingRequest) ChannelSinceVersion() uint16 {
	return 0
}

func (f *FindLastMatchingRecordingRequest) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= f.ChannelSinceVersion()
}

func (*FindLastMatchingRecordingRequest) ChannelDeprecated() uint16 {
	return 0
}

func (FindLastMatchingRecordingRequest) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (FindLastMatchingRecordingRequest) ChannelHeaderLength() uint64 {
	return 4
}
