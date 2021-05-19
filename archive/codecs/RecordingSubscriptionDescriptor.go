// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingSubscriptionDescriptor struct {
	ControlSessionId int64
	CorrelationId    int64
	SubscriptionId   int64
	StreamId         int32
	StrippedChannel  []uint8
}

func (r *RecordingSubscriptionDescriptor) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.SubscriptionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.StreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.StrippedChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.StrippedChannel); err != nil {
		return err
	}
	return nil
}

func (r *RecordingSubscriptionDescriptor) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.ControlSessionIdInActingVersion(actingVersion) {
		r.ControlSessionId = r.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.ControlSessionId); err != nil {
			return err
		}
	}
	if !r.CorrelationIdInActingVersion(actingVersion) {
		r.CorrelationId = r.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.CorrelationId); err != nil {
			return err
		}
	}
	if !r.SubscriptionIdInActingVersion(actingVersion) {
		r.SubscriptionId = r.SubscriptionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.SubscriptionId); err != nil {
			return err
		}
	}
	if !r.StreamIdInActingVersion(actingVersion) {
		r.StreamId = r.StreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.StreamId); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}

	if r.StrippedChannelInActingVersion(actingVersion) {
		var StrippedChannelLength uint32
		if err := _m.ReadUint32(_r, &StrippedChannelLength); err != nil {
			return err
		}
		if cap(r.StrippedChannel) < int(StrippedChannelLength) {
			r.StrippedChannel = make([]uint8, StrippedChannelLength)
		}
		r.StrippedChannel = r.StrippedChannel[:StrippedChannelLength]
		if err := _m.ReadBytes(_r, r.StrippedChannel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := r.RangeCheck(actingVersion, r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordingSubscriptionDescriptor) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.ControlSessionIdInActingVersion(actingVersion) {
		if r.ControlSessionId < r.ControlSessionIdMinValue() || r.ControlSessionId > r.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.ControlSessionId (%v < %v > %v)", r.ControlSessionIdMinValue(), r.ControlSessionId, r.ControlSessionIdMaxValue())
		}
	}
	if r.CorrelationIdInActingVersion(actingVersion) {
		if r.CorrelationId < r.CorrelationIdMinValue() || r.CorrelationId > r.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on r.CorrelationId (%v < %v > %v)", r.CorrelationIdMinValue(), r.CorrelationId, r.CorrelationIdMaxValue())
		}
	}
	if r.SubscriptionIdInActingVersion(actingVersion) {
		if r.SubscriptionId < r.SubscriptionIdMinValue() || r.SubscriptionId > r.SubscriptionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SubscriptionId (%v < %v > %v)", r.SubscriptionIdMinValue(), r.SubscriptionId, r.SubscriptionIdMaxValue())
		}
	}
	if r.StreamIdInActingVersion(actingVersion) {
		if r.StreamId < r.StreamIdMinValue() || r.StreamId > r.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on r.StreamId (%v < %v > %v)", r.StreamIdMinValue(), r.StreamId, r.StreamIdMaxValue())
		}
	}
	return nil
}

func RecordingSubscriptionDescriptorInit(r *RecordingSubscriptionDescriptor) {
	return
}

func (*RecordingSubscriptionDescriptor) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*RecordingSubscriptionDescriptor) SbeTemplateId() (templateId uint16) {
	return 23
}

func (*RecordingSubscriptionDescriptor) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingSubscriptionDescriptor) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*RecordingSubscriptionDescriptor) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingSubscriptionDescriptor) ControlSessionIdId() uint16 {
	return 1
}

func (*RecordingSubscriptionDescriptor) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSubscriptionDescriptor) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ControlSessionIdSinceVersion()
}

func (*RecordingSubscriptionDescriptor) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*RecordingSubscriptionDescriptor) ControlSessionIdMetaAttribute(meta int) string {
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

func (*RecordingSubscriptionDescriptor) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSubscriptionDescriptor) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSubscriptionDescriptor) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSubscriptionDescriptor) CorrelationIdId() uint16 {
	return 2
}

func (*RecordingSubscriptionDescriptor) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSubscriptionDescriptor) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CorrelationIdSinceVersion()
}

func (*RecordingSubscriptionDescriptor) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*RecordingSubscriptionDescriptor) CorrelationIdMetaAttribute(meta int) string {
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

func (*RecordingSubscriptionDescriptor) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSubscriptionDescriptor) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSubscriptionDescriptor) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSubscriptionDescriptor) SubscriptionIdId() uint16 {
	return 3
}

func (*RecordingSubscriptionDescriptor) SubscriptionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSubscriptionDescriptor) SubscriptionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SubscriptionIdSinceVersion()
}

func (*RecordingSubscriptionDescriptor) SubscriptionIdDeprecated() uint16 {
	return 0
}

func (*RecordingSubscriptionDescriptor) SubscriptionIdMetaAttribute(meta int) string {
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

func (*RecordingSubscriptionDescriptor) SubscriptionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSubscriptionDescriptor) SubscriptionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSubscriptionDescriptor) SubscriptionIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSubscriptionDescriptor) StreamIdId() uint16 {
	return 4
}

func (*RecordingSubscriptionDescriptor) StreamIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSubscriptionDescriptor) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StreamIdSinceVersion()
}

func (*RecordingSubscriptionDescriptor) StreamIdDeprecated() uint16 {
	return 0
}

func (*RecordingSubscriptionDescriptor) StreamIdMetaAttribute(meta int) string {
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

func (*RecordingSubscriptionDescriptor) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingSubscriptionDescriptor) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingSubscriptionDescriptor) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*RecordingSubscriptionDescriptor) StrippedChannelMetaAttribute(meta int) string {
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

func (*RecordingSubscriptionDescriptor) StrippedChannelSinceVersion() uint16 {
	return 0
}

func (r *RecordingSubscriptionDescriptor) StrippedChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StrippedChannelSinceVersion()
}

func (*RecordingSubscriptionDescriptor) StrippedChannelDeprecated() uint16 {
	return 0
}

func (RecordingSubscriptionDescriptor) StrippedChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (RecordingSubscriptionDescriptor) StrippedChannelHeaderLength() uint64 {
	return 4
}
