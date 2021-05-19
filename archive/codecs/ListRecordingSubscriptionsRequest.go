// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ListRecordingSubscriptionsRequest struct {
	ControlSessionId  int64
	CorrelationId     int64
	PseudoIndex       int32
	SubscriptionCount int32
	ApplyStreamId     BooleanTypeEnum
	StreamId          int32
	Channel           []uint8
}

func (l *ListRecordingSubscriptionsRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := l.RangeCheck(l.SbeSchemaVersion(), l.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, l.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, l.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, l.PseudoIndex); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, l.SubscriptionCount); err != nil {
		return err
	}
	if err := l.ApplyStreamId.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, l.StreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(l.Channel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, l.Channel); err != nil {
		return err
	}
	return nil
}

func (l *ListRecordingSubscriptionsRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !l.ControlSessionIdInActingVersion(actingVersion) {
		l.ControlSessionId = l.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &l.ControlSessionId); err != nil {
			return err
		}
	}
	if !l.CorrelationIdInActingVersion(actingVersion) {
		l.CorrelationId = l.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &l.CorrelationId); err != nil {
			return err
		}
	}
	if !l.PseudoIndexInActingVersion(actingVersion) {
		l.PseudoIndex = l.PseudoIndexNullValue()
	} else {
		if err := _m.ReadInt32(_r, &l.PseudoIndex); err != nil {
			return err
		}
	}
	if !l.SubscriptionCountInActingVersion(actingVersion) {
		l.SubscriptionCount = l.SubscriptionCountNullValue()
	} else {
		if err := _m.ReadInt32(_r, &l.SubscriptionCount); err != nil {
			return err
		}
	}
	if l.ApplyStreamIdInActingVersion(actingVersion) {
		if err := l.ApplyStreamId.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !l.StreamIdInActingVersion(actingVersion) {
		l.StreamId = l.StreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &l.StreamId); err != nil {
			return err
		}
	}
	if actingVersion > l.SbeSchemaVersion() && blockLength > l.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-l.SbeBlockLength()))
	}

	if l.ChannelInActingVersion(actingVersion) {
		var ChannelLength uint32
		if err := _m.ReadUint32(_r, &ChannelLength); err != nil {
			return err
		}
		if cap(l.Channel) < int(ChannelLength) {
			l.Channel = make([]uint8, ChannelLength)
		}
		l.Channel = l.Channel[:ChannelLength]
		if err := _m.ReadBytes(_r, l.Channel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := l.RangeCheck(actingVersion, l.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (l *ListRecordingSubscriptionsRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if l.ControlSessionIdInActingVersion(actingVersion) {
		if l.ControlSessionId < l.ControlSessionIdMinValue() || l.ControlSessionId > l.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on l.ControlSessionId (%v < %v > %v)", l.ControlSessionIdMinValue(), l.ControlSessionId, l.ControlSessionIdMaxValue())
		}
	}
	if l.CorrelationIdInActingVersion(actingVersion) {
		if l.CorrelationId < l.CorrelationIdMinValue() || l.CorrelationId > l.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on l.CorrelationId (%v < %v > %v)", l.CorrelationIdMinValue(), l.CorrelationId, l.CorrelationIdMaxValue())
		}
	}
	if l.PseudoIndexInActingVersion(actingVersion) {
		if l.PseudoIndex < l.PseudoIndexMinValue() || l.PseudoIndex > l.PseudoIndexMaxValue() {
			return fmt.Errorf("Range check failed on l.PseudoIndex (%v < %v > %v)", l.PseudoIndexMinValue(), l.PseudoIndex, l.PseudoIndexMaxValue())
		}
	}
	if l.SubscriptionCountInActingVersion(actingVersion) {
		if l.SubscriptionCount < l.SubscriptionCountMinValue() || l.SubscriptionCount > l.SubscriptionCountMaxValue() {
			return fmt.Errorf("Range check failed on l.SubscriptionCount (%v < %v > %v)", l.SubscriptionCountMinValue(), l.SubscriptionCount, l.SubscriptionCountMaxValue())
		}
	}
	if err := l.ApplyStreamId.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if l.StreamIdInActingVersion(actingVersion) {
		if l.StreamId < l.StreamIdMinValue() || l.StreamId > l.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on l.StreamId (%v < %v > %v)", l.StreamIdMinValue(), l.StreamId, l.StreamIdMaxValue())
		}
	}
	return nil
}

func ListRecordingSubscriptionsRequestInit(l *ListRecordingSubscriptionsRequest) {
	return
}

func (*ListRecordingSubscriptionsRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*ListRecordingSubscriptionsRequest) SbeTemplateId() (templateId uint16) {
	return 17
}

func (*ListRecordingSubscriptionsRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ListRecordingSubscriptionsRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ListRecordingSubscriptionsRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ListRecordingSubscriptionsRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*ListRecordingSubscriptionsRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingSubscriptionsRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.ControlSessionIdSinceVersion()
}

func (*ListRecordingSubscriptionsRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingSubscriptionsRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ListRecordingSubscriptionsRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingSubscriptionsRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingSubscriptionsRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingSubscriptionsRequest) CorrelationIdId() uint16 {
	return 2
}

func (*ListRecordingSubscriptionsRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingSubscriptionsRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.CorrelationIdSinceVersion()
}

func (*ListRecordingSubscriptionsRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingSubscriptionsRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ListRecordingSubscriptionsRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingSubscriptionsRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingSubscriptionsRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingSubscriptionsRequest) PseudoIndexId() uint16 {
	return 3
}

func (*ListRecordingSubscriptionsRequest) PseudoIndexSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingSubscriptionsRequest) PseudoIndexInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.PseudoIndexSinceVersion()
}

func (*ListRecordingSubscriptionsRequest) PseudoIndexDeprecated() uint16 {
	return 0
}

func (*ListRecordingSubscriptionsRequest) PseudoIndexMetaAttribute(meta int) string {
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

func (*ListRecordingSubscriptionsRequest) PseudoIndexMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ListRecordingSubscriptionsRequest) PseudoIndexMaxValue() int32 {
	return math.MaxInt32
}

func (*ListRecordingSubscriptionsRequest) PseudoIndexNullValue() int32 {
	return math.MinInt32
}

func (*ListRecordingSubscriptionsRequest) SubscriptionCountId() uint16 {
	return 4
}

func (*ListRecordingSubscriptionsRequest) SubscriptionCountSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingSubscriptionsRequest) SubscriptionCountInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.SubscriptionCountSinceVersion()
}

func (*ListRecordingSubscriptionsRequest) SubscriptionCountDeprecated() uint16 {
	return 0
}

func (*ListRecordingSubscriptionsRequest) SubscriptionCountMetaAttribute(meta int) string {
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

func (*ListRecordingSubscriptionsRequest) SubscriptionCountMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ListRecordingSubscriptionsRequest) SubscriptionCountMaxValue() int32 {
	return math.MaxInt32
}

func (*ListRecordingSubscriptionsRequest) SubscriptionCountNullValue() int32 {
	return math.MinInt32
}

func (*ListRecordingSubscriptionsRequest) ApplyStreamIdId() uint16 {
	return 5
}

func (*ListRecordingSubscriptionsRequest) ApplyStreamIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingSubscriptionsRequest) ApplyStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.ApplyStreamIdSinceVersion()
}

func (*ListRecordingSubscriptionsRequest) ApplyStreamIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingSubscriptionsRequest) ApplyStreamIdMetaAttribute(meta int) string {
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

func (*ListRecordingSubscriptionsRequest) StreamIdId() uint16 {
	return 6
}

func (*ListRecordingSubscriptionsRequest) StreamIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingSubscriptionsRequest) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.StreamIdSinceVersion()
}

func (*ListRecordingSubscriptionsRequest) StreamIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingSubscriptionsRequest) StreamIdMetaAttribute(meta int) string {
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

func (*ListRecordingSubscriptionsRequest) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ListRecordingSubscriptionsRequest) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ListRecordingSubscriptionsRequest) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ListRecordingSubscriptionsRequest) ChannelMetaAttribute(meta int) string {
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

func (*ListRecordingSubscriptionsRequest) ChannelSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingSubscriptionsRequest) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.ChannelSinceVersion()
}

func (*ListRecordingSubscriptionsRequest) ChannelDeprecated() uint16 {
	return 0
}

func (ListRecordingSubscriptionsRequest) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ListRecordingSubscriptionsRequest) ChannelHeaderLength() uint64 {
	return 4
}
