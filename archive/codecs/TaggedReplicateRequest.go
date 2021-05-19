// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type TaggedReplicateRequest struct {
	ControlSessionId   int64
	CorrelationId      int64
	SrcRecordingId     int64
	DstRecordingId     int64
	ChannelTagId       int64
	SubscriptionTagId  int64
	SrcControlStreamId int32
	SrcControlChannel  []uint8
	LiveDestination    []uint8
}

func (t *TaggedReplicateRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := t.RangeCheck(t.SbeSchemaVersion(), t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, t.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.SrcRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.DstRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.ChannelTagId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.SubscriptionTagId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, t.SrcControlStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(t.SrcControlChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, t.SrcControlChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(t.LiveDestination))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, t.LiveDestination); err != nil {
		return err
	}
	return nil
}

func (t *TaggedReplicateRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !t.ControlSessionIdInActingVersion(actingVersion) {
		t.ControlSessionId = t.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.ControlSessionId); err != nil {
			return err
		}
	}
	if !t.CorrelationIdInActingVersion(actingVersion) {
		t.CorrelationId = t.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.CorrelationId); err != nil {
			return err
		}
	}
	if !t.SrcRecordingIdInActingVersion(actingVersion) {
		t.SrcRecordingId = t.SrcRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.SrcRecordingId); err != nil {
			return err
		}
	}
	if !t.DstRecordingIdInActingVersion(actingVersion) {
		t.DstRecordingId = t.DstRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.DstRecordingId); err != nil {
			return err
		}
	}
	if !t.ChannelTagIdInActingVersion(actingVersion) {
		t.ChannelTagId = t.ChannelTagIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.ChannelTagId); err != nil {
			return err
		}
	}
	if !t.SubscriptionTagIdInActingVersion(actingVersion) {
		t.SubscriptionTagId = t.SubscriptionTagIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.SubscriptionTagId); err != nil {
			return err
		}
	}
	if !t.SrcControlStreamIdInActingVersion(actingVersion) {
		t.SrcControlStreamId = t.SrcControlStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &t.SrcControlStreamId); err != nil {
			return err
		}
	}
	if actingVersion > t.SbeSchemaVersion() && blockLength > t.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-t.SbeBlockLength()))
	}

	if t.SrcControlChannelInActingVersion(actingVersion) {
		var SrcControlChannelLength uint32
		if err := _m.ReadUint32(_r, &SrcControlChannelLength); err != nil {
			return err
		}
		if cap(t.SrcControlChannel) < int(SrcControlChannelLength) {
			t.SrcControlChannel = make([]uint8, SrcControlChannelLength)
		}
		t.SrcControlChannel = t.SrcControlChannel[:SrcControlChannelLength]
		if err := _m.ReadBytes(_r, t.SrcControlChannel); err != nil {
			return err
		}
	}

	if t.LiveDestinationInActingVersion(actingVersion) {
		var LiveDestinationLength uint32
		if err := _m.ReadUint32(_r, &LiveDestinationLength); err != nil {
			return err
		}
		if cap(t.LiveDestination) < int(LiveDestinationLength) {
			t.LiveDestination = make([]uint8, LiveDestinationLength)
		}
		t.LiveDestination = t.LiveDestination[:LiveDestinationLength]
		if err := _m.ReadBytes(_r, t.LiveDestination); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := t.RangeCheck(actingVersion, t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (t *TaggedReplicateRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if t.ControlSessionIdInActingVersion(actingVersion) {
		if t.ControlSessionId < t.ControlSessionIdMinValue() || t.ControlSessionId > t.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on t.ControlSessionId (%v < %v > %v)", t.ControlSessionIdMinValue(), t.ControlSessionId, t.ControlSessionIdMaxValue())
		}
	}
	if t.CorrelationIdInActingVersion(actingVersion) {
		if t.CorrelationId < t.CorrelationIdMinValue() || t.CorrelationId > t.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on t.CorrelationId (%v < %v > %v)", t.CorrelationIdMinValue(), t.CorrelationId, t.CorrelationIdMaxValue())
		}
	}
	if t.SrcRecordingIdInActingVersion(actingVersion) {
		if t.SrcRecordingId < t.SrcRecordingIdMinValue() || t.SrcRecordingId > t.SrcRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on t.SrcRecordingId (%v < %v > %v)", t.SrcRecordingIdMinValue(), t.SrcRecordingId, t.SrcRecordingIdMaxValue())
		}
	}
	if t.DstRecordingIdInActingVersion(actingVersion) {
		if t.DstRecordingId < t.DstRecordingIdMinValue() || t.DstRecordingId > t.DstRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on t.DstRecordingId (%v < %v > %v)", t.DstRecordingIdMinValue(), t.DstRecordingId, t.DstRecordingIdMaxValue())
		}
	}
	if t.ChannelTagIdInActingVersion(actingVersion) {
		if t.ChannelTagId < t.ChannelTagIdMinValue() || t.ChannelTagId > t.ChannelTagIdMaxValue() {
			return fmt.Errorf("Range check failed on t.ChannelTagId (%v < %v > %v)", t.ChannelTagIdMinValue(), t.ChannelTagId, t.ChannelTagIdMaxValue())
		}
	}
	if t.SubscriptionTagIdInActingVersion(actingVersion) {
		if t.SubscriptionTagId < t.SubscriptionTagIdMinValue() || t.SubscriptionTagId > t.SubscriptionTagIdMaxValue() {
			return fmt.Errorf("Range check failed on t.SubscriptionTagId (%v < %v > %v)", t.SubscriptionTagIdMinValue(), t.SubscriptionTagId, t.SubscriptionTagIdMaxValue())
		}
	}
	if t.SrcControlStreamIdInActingVersion(actingVersion) {
		if t.SrcControlStreamId < t.SrcControlStreamIdMinValue() || t.SrcControlStreamId > t.SrcControlStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on t.SrcControlStreamId (%v < %v > %v)", t.SrcControlStreamIdMinValue(), t.SrcControlStreamId, t.SrcControlStreamIdMaxValue())
		}
	}
	return nil
}

func TaggedReplicateRequestInit(t *TaggedReplicateRequest) {
	return
}

func (*TaggedReplicateRequest) SbeBlockLength() (blockLength uint16) {
	return 52
}

func (*TaggedReplicateRequest) SbeTemplateId() (templateId uint16) {
	return 62
}

func (*TaggedReplicateRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*TaggedReplicateRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*TaggedReplicateRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*TaggedReplicateRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*TaggedReplicateRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.ControlSessionIdSinceVersion()
}

func (*TaggedReplicateRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*TaggedReplicateRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TaggedReplicateRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TaggedReplicateRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*TaggedReplicateRequest) CorrelationIdId() uint16 {
	return 2
}

func (*TaggedReplicateRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.CorrelationIdSinceVersion()
}

func (*TaggedReplicateRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*TaggedReplicateRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TaggedReplicateRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TaggedReplicateRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*TaggedReplicateRequest) SrcRecordingIdId() uint16 {
	return 3
}

func (*TaggedReplicateRequest) SrcRecordingIdSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) SrcRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.SrcRecordingIdSinceVersion()
}

func (*TaggedReplicateRequest) SrcRecordingIdDeprecated() uint16 {
	return 0
}

func (*TaggedReplicateRequest) SrcRecordingIdMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) SrcRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TaggedReplicateRequest) SrcRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TaggedReplicateRequest) SrcRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*TaggedReplicateRequest) DstRecordingIdId() uint16 {
	return 4
}

func (*TaggedReplicateRequest) DstRecordingIdSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) DstRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.DstRecordingIdSinceVersion()
}

func (*TaggedReplicateRequest) DstRecordingIdDeprecated() uint16 {
	return 0
}

func (*TaggedReplicateRequest) DstRecordingIdMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) DstRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TaggedReplicateRequest) DstRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TaggedReplicateRequest) DstRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*TaggedReplicateRequest) ChannelTagIdId() uint16 {
	return 5
}

func (*TaggedReplicateRequest) ChannelTagIdSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) ChannelTagIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.ChannelTagIdSinceVersion()
}

func (*TaggedReplicateRequest) ChannelTagIdDeprecated() uint16 {
	return 0
}

func (*TaggedReplicateRequest) ChannelTagIdMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) ChannelTagIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TaggedReplicateRequest) ChannelTagIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TaggedReplicateRequest) ChannelTagIdNullValue() int64 {
	return math.MinInt64
}

func (*TaggedReplicateRequest) SubscriptionTagIdId() uint16 {
	return 6
}

func (*TaggedReplicateRequest) SubscriptionTagIdSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) SubscriptionTagIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.SubscriptionTagIdSinceVersion()
}

func (*TaggedReplicateRequest) SubscriptionTagIdDeprecated() uint16 {
	return 0
}

func (*TaggedReplicateRequest) SubscriptionTagIdMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) SubscriptionTagIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TaggedReplicateRequest) SubscriptionTagIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TaggedReplicateRequest) SubscriptionTagIdNullValue() int64 {
	return math.MinInt64
}

func (*TaggedReplicateRequest) SrcControlStreamIdId() uint16 {
	return 7
}

func (*TaggedReplicateRequest) SrcControlStreamIdSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) SrcControlStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.SrcControlStreamIdSinceVersion()
}

func (*TaggedReplicateRequest) SrcControlStreamIdDeprecated() uint16 {
	return 0
}

func (*TaggedReplicateRequest) SrcControlStreamIdMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) SrcControlStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*TaggedReplicateRequest) SrcControlStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*TaggedReplicateRequest) SrcControlStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*TaggedReplicateRequest) SrcControlChannelMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) SrcControlChannelSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) SrcControlChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.SrcControlChannelSinceVersion()
}

func (*TaggedReplicateRequest) SrcControlChannelDeprecated() uint16 {
	return 0
}

func (TaggedReplicateRequest) SrcControlChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (TaggedReplicateRequest) SrcControlChannelHeaderLength() uint64 {
	return 4
}

func (*TaggedReplicateRequest) LiveDestinationMetaAttribute(meta int) string {
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

func (*TaggedReplicateRequest) LiveDestinationSinceVersion() uint16 {
	return 0
}

func (t *TaggedReplicateRequest) LiveDestinationInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.LiveDestinationSinceVersion()
}

func (*TaggedReplicateRequest) LiveDestinationDeprecated() uint16 {
	return 0
}

func (TaggedReplicateRequest) LiveDestinationCharacterEncoding() string {
	return "US-ASCII"
}

func (TaggedReplicateRequest) LiveDestinationHeaderLength() uint64 {
	return 4
}
