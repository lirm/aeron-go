// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ReplicateRequest2 struct {
	ControlSessionId   int64
	CorrelationId      int64
	SrcRecordingId     int64
	DstRecordingId     int64
	StopPosition       int64
	ChannelTagId       int64
	SubscriptionTagId  int64
	SrcControlStreamId int32
	SrcControlChannel  []uint8
	LiveDestination    []uint8
	ReplicationChannel []uint8
}

func (r *ReplicateRequest2) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, r.SrcRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.DstRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.StopPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.ChannelTagId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.SubscriptionTagId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.SrcControlStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.SrcControlChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.SrcControlChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.LiveDestination))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.LiveDestination); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.ReplicationChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.ReplicationChannel); err != nil {
		return err
	}
	return nil
}

func (r *ReplicateRequest2) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !r.SrcRecordingIdInActingVersion(actingVersion) {
		r.SrcRecordingId = r.SrcRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.SrcRecordingId); err != nil {
			return err
		}
	}
	if !r.DstRecordingIdInActingVersion(actingVersion) {
		r.DstRecordingId = r.DstRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.DstRecordingId); err != nil {
			return err
		}
	}
	if !r.StopPositionInActingVersion(actingVersion) {
		r.StopPosition = r.StopPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.StopPosition); err != nil {
			return err
		}
	}
	if !r.ChannelTagIdInActingVersion(actingVersion) {
		r.ChannelTagId = r.ChannelTagIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.ChannelTagId); err != nil {
			return err
		}
	}
	if !r.SubscriptionTagIdInActingVersion(actingVersion) {
		r.SubscriptionTagId = r.SubscriptionTagIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.SubscriptionTagId); err != nil {
			return err
		}
	}
	if !r.SrcControlStreamIdInActingVersion(actingVersion) {
		r.SrcControlStreamId = r.SrcControlStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.SrcControlStreamId); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}

	if r.SrcControlChannelInActingVersion(actingVersion) {
		var SrcControlChannelLength uint32
		if err := _m.ReadUint32(_r, &SrcControlChannelLength); err != nil {
			return err
		}
		if cap(r.SrcControlChannel) < int(SrcControlChannelLength) {
			r.SrcControlChannel = make([]uint8, SrcControlChannelLength)
		}
		r.SrcControlChannel = r.SrcControlChannel[:SrcControlChannelLength]
		if err := _m.ReadBytes(_r, r.SrcControlChannel); err != nil {
			return err
		}
	}

	if r.LiveDestinationInActingVersion(actingVersion) {
		var LiveDestinationLength uint32
		if err := _m.ReadUint32(_r, &LiveDestinationLength); err != nil {
			return err
		}
		if cap(r.LiveDestination) < int(LiveDestinationLength) {
			r.LiveDestination = make([]uint8, LiveDestinationLength)
		}
		r.LiveDestination = r.LiveDestination[:LiveDestinationLength]
		if err := _m.ReadBytes(_r, r.LiveDestination); err != nil {
			return err
		}
	}

	if r.ReplicationChannelInActingVersion(actingVersion) {
		var ReplicationChannelLength uint32
		if err := _m.ReadUint32(_r, &ReplicationChannelLength); err != nil {
			return err
		}
		if cap(r.ReplicationChannel) < int(ReplicationChannelLength) {
			r.ReplicationChannel = make([]uint8, ReplicationChannelLength)
		}
		r.ReplicationChannel = r.ReplicationChannel[:ReplicationChannelLength]
		if err := _m.ReadBytes(_r, r.ReplicationChannel); err != nil {
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

func (r *ReplicateRequest2) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if r.SrcRecordingIdInActingVersion(actingVersion) {
		if r.SrcRecordingId < r.SrcRecordingIdMinValue() || r.SrcRecordingId > r.SrcRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SrcRecordingId (%v < %v > %v)", r.SrcRecordingIdMinValue(), r.SrcRecordingId, r.SrcRecordingIdMaxValue())
		}
	}
	if r.DstRecordingIdInActingVersion(actingVersion) {
		if r.DstRecordingId < r.DstRecordingIdMinValue() || r.DstRecordingId > r.DstRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.DstRecordingId (%v < %v > %v)", r.DstRecordingIdMinValue(), r.DstRecordingId, r.DstRecordingIdMaxValue())
		}
	}
	if r.StopPositionInActingVersion(actingVersion) {
		if r.StopPosition < r.StopPositionMinValue() || r.StopPosition > r.StopPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.StopPosition (%v < %v > %v)", r.StopPositionMinValue(), r.StopPosition, r.StopPositionMaxValue())
		}
	}
	if r.ChannelTagIdInActingVersion(actingVersion) {
		if r.ChannelTagId < r.ChannelTagIdMinValue() || r.ChannelTagId > r.ChannelTagIdMaxValue() {
			return fmt.Errorf("Range check failed on r.ChannelTagId (%v < %v > %v)", r.ChannelTagIdMinValue(), r.ChannelTagId, r.ChannelTagIdMaxValue())
		}
	}
	if r.SubscriptionTagIdInActingVersion(actingVersion) {
		if r.SubscriptionTagId < r.SubscriptionTagIdMinValue() || r.SubscriptionTagId > r.SubscriptionTagIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SubscriptionTagId (%v < %v > %v)", r.SubscriptionTagIdMinValue(), r.SubscriptionTagId, r.SubscriptionTagIdMaxValue())
		}
	}
	if r.SrcControlStreamIdInActingVersion(actingVersion) {
		if r.SrcControlStreamId < r.SrcControlStreamIdMinValue() || r.SrcControlStreamId > r.SrcControlStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SrcControlStreamId (%v < %v > %v)", r.SrcControlStreamIdMinValue(), r.SrcControlStreamId, r.SrcControlStreamIdMaxValue())
		}
	}
	return nil
}

func ReplicateRequest2Init(r *ReplicateRequest2) {
	return
}

func (*ReplicateRequest2) SbeBlockLength() (blockLength uint16) {
	return 60
}

func (*ReplicateRequest2) SbeTemplateId() (templateId uint16) {
	return 66
}

func (*ReplicateRequest2) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ReplicateRequest2) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ReplicateRequest2) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ReplicateRequest2) ControlSessionIdId() uint16 {
	return 1
}

func (*ReplicateRequest2) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ControlSessionIdSinceVersion()
}

func (*ReplicateRequest2) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest2) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest2) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest2) CorrelationIdId() uint16 {
	return 2
}

func (*ReplicateRequest2) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CorrelationIdSinceVersion()
}

func (*ReplicateRequest2) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) CorrelationIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest2) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest2) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest2) SrcRecordingIdId() uint16 {
	return 3
}

func (*ReplicateRequest2) SrcRecordingIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) SrcRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SrcRecordingIdSinceVersion()
}

func (*ReplicateRequest2) SrcRecordingIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) SrcRecordingIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) SrcRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest2) SrcRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest2) SrcRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest2) DstRecordingIdId() uint16 {
	return 4
}

func (*ReplicateRequest2) DstRecordingIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) DstRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.DstRecordingIdSinceVersion()
}

func (*ReplicateRequest2) DstRecordingIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) DstRecordingIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) DstRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest2) DstRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest2) DstRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest2) StopPositionId() uint16 {
	return 5
}

func (*ReplicateRequest2) StopPositionSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) StopPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StopPositionSinceVersion()
}

func (*ReplicateRequest2) StopPositionDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) StopPositionMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) StopPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest2) StopPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest2) StopPositionNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest2) ChannelTagIdId() uint16 {
	return 6
}

func (*ReplicateRequest2) ChannelTagIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) ChannelTagIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ChannelTagIdSinceVersion()
}

func (*ReplicateRequest2) ChannelTagIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) ChannelTagIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) ChannelTagIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest2) ChannelTagIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest2) ChannelTagIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest2) SubscriptionTagIdId() uint16 {
	return 7
}

func (*ReplicateRequest2) SubscriptionTagIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) SubscriptionTagIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SubscriptionTagIdSinceVersion()
}

func (*ReplicateRequest2) SubscriptionTagIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) SubscriptionTagIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) SubscriptionTagIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest2) SubscriptionTagIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest2) SubscriptionTagIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest2) SrcControlStreamIdId() uint16 {
	return 8
}

func (*ReplicateRequest2) SrcControlStreamIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) SrcControlStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SrcControlStreamIdSinceVersion()
}

func (*ReplicateRequest2) SrcControlStreamIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest2) SrcControlStreamIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) SrcControlStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ReplicateRequest2) SrcControlStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ReplicateRequest2) SrcControlStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ReplicateRequest2) SrcControlChannelMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) SrcControlChannelSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) SrcControlChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SrcControlChannelSinceVersion()
}

func (*ReplicateRequest2) SrcControlChannelDeprecated() uint16 {
	return 0
}

func (ReplicateRequest2) SrcControlChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ReplicateRequest2) SrcControlChannelHeaderLength() uint64 {
	return 4
}

func (*ReplicateRequest2) LiveDestinationMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) LiveDestinationSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) LiveDestinationInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.LiveDestinationSinceVersion()
}

func (*ReplicateRequest2) LiveDestinationDeprecated() uint16 {
	return 0
}

func (ReplicateRequest2) LiveDestinationCharacterEncoding() string {
	return "US-ASCII"
}

func (ReplicateRequest2) LiveDestinationHeaderLength() uint64 {
	return 4
}

func (*ReplicateRequest2) ReplicationChannelMetaAttribute(meta int) string {
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

func (*ReplicateRequest2) ReplicationChannelSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest2) ReplicationChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ReplicationChannelSinceVersion()
}

func (*ReplicateRequest2) ReplicationChannelDeprecated() uint16 {
	return 0
}

func (ReplicateRequest2) ReplicationChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ReplicateRequest2) ReplicationChannelHeaderLength() uint64 {
	return 4
}
