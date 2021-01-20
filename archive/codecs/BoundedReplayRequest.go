// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type BoundedReplayRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	Position         int64
	Length           int64
	LimitCounterId   int32
	ReplayStreamId   int32
	ReplayChannel    []uint8
}

func (b *BoundedReplayRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := b.RangeCheck(b.SbeSchemaVersion(), b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, b.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.Position); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.Length); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, b.LimitCounterId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, b.ReplayStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(b.ReplayChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, b.ReplayChannel); err != nil {
		return err
	}
	return nil
}

func (b *BoundedReplayRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !b.ControlSessionIdInActingVersion(actingVersion) {
		b.ControlSessionId = b.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.ControlSessionId); err != nil {
			return err
		}
	}
	if !b.CorrelationIdInActingVersion(actingVersion) {
		b.CorrelationId = b.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.CorrelationId); err != nil {
			return err
		}
	}
	if !b.RecordingIdInActingVersion(actingVersion) {
		b.RecordingId = b.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.RecordingId); err != nil {
			return err
		}
	}
	if !b.PositionInActingVersion(actingVersion) {
		b.Position = b.PositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.Position); err != nil {
			return err
		}
	}
	if !b.LengthInActingVersion(actingVersion) {
		b.Length = b.LengthNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.Length); err != nil {
			return err
		}
	}
	if !b.LimitCounterIdInActingVersion(actingVersion) {
		b.LimitCounterId = b.LimitCounterIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &b.LimitCounterId); err != nil {
			return err
		}
	}
	if !b.ReplayStreamIdInActingVersion(actingVersion) {
		b.ReplayStreamId = b.ReplayStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &b.ReplayStreamId); err != nil {
			return err
		}
	}
	if actingVersion > b.SbeSchemaVersion() && blockLength > b.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}

	if b.ReplayChannelInActingVersion(actingVersion) {
		var ReplayChannelLength uint32
		if err := _m.ReadUint32(_r, &ReplayChannelLength); err != nil {
			return err
		}
		if cap(b.ReplayChannel) < int(ReplayChannelLength) {
			b.ReplayChannel = make([]uint8, ReplayChannelLength)
		}
		b.ReplayChannel = b.ReplayChannel[:ReplayChannelLength]
		if err := _m.ReadBytes(_r, b.ReplayChannel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := b.RangeCheck(actingVersion, b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (b *BoundedReplayRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if b.ControlSessionIdInActingVersion(actingVersion) {
		if b.ControlSessionId < b.ControlSessionIdMinValue() || b.ControlSessionId > b.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on b.ControlSessionId (%v < %v > %v)", b.ControlSessionIdMinValue(), b.ControlSessionId, b.ControlSessionIdMaxValue())
		}
	}
	if b.CorrelationIdInActingVersion(actingVersion) {
		if b.CorrelationId < b.CorrelationIdMinValue() || b.CorrelationId > b.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on b.CorrelationId (%v < %v > %v)", b.CorrelationIdMinValue(), b.CorrelationId, b.CorrelationIdMaxValue())
		}
	}
	if b.RecordingIdInActingVersion(actingVersion) {
		if b.RecordingId < b.RecordingIdMinValue() || b.RecordingId > b.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on b.RecordingId (%v < %v > %v)", b.RecordingIdMinValue(), b.RecordingId, b.RecordingIdMaxValue())
		}
	}
	if b.PositionInActingVersion(actingVersion) {
		if b.Position < b.PositionMinValue() || b.Position > b.PositionMaxValue() {
			return fmt.Errorf("Range check failed on b.Position (%v < %v > %v)", b.PositionMinValue(), b.Position, b.PositionMaxValue())
		}
	}
	if b.LengthInActingVersion(actingVersion) {
		if b.Length < b.LengthMinValue() || b.Length > b.LengthMaxValue() {
			return fmt.Errorf("Range check failed on b.Length (%v < %v > %v)", b.LengthMinValue(), b.Length, b.LengthMaxValue())
		}
	}
	if b.LimitCounterIdInActingVersion(actingVersion) {
		if b.LimitCounterId < b.LimitCounterIdMinValue() || b.LimitCounterId > b.LimitCounterIdMaxValue() {
			return fmt.Errorf("Range check failed on b.LimitCounterId (%v < %v > %v)", b.LimitCounterIdMinValue(), b.LimitCounterId, b.LimitCounterIdMaxValue())
		}
	}
	if b.ReplayStreamIdInActingVersion(actingVersion) {
		if b.ReplayStreamId < b.ReplayStreamIdMinValue() || b.ReplayStreamId > b.ReplayStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on b.ReplayStreamId (%v < %v > %v)", b.ReplayStreamIdMinValue(), b.ReplayStreamId, b.ReplayStreamIdMaxValue())
		}
	}
	return nil
}

func BoundedReplayRequestInit(b *BoundedReplayRequest) {
	return
}

func (*BoundedReplayRequest) SbeBlockLength() (blockLength uint16) {
	return 48
}

func (*BoundedReplayRequest) SbeTemplateId() (templateId uint16) {
	return 18
}

func (*BoundedReplayRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*BoundedReplayRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*BoundedReplayRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*BoundedReplayRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*BoundedReplayRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ControlSessionIdSinceVersion()
}

func (*BoundedReplayRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*BoundedReplayRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BoundedReplayRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BoundedReplayRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*BoundedReplayRequest) CorrelationIdId() uint16 {
	return 2
}

func (*BoundedReplayRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.CorrelationIdSinceVersion()
}

func (*BoundedReplayRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*BoundedReplayRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BoundedReplayRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BoundedReplayRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*BoundedReplayRequest) RecordingIdId() uint16 {
	return 3
}

func (*BoundedReplayRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.RecordingIdSinceVersion()
}

func (*BoundedReplayRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*BoundedReplayRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BoundedReplayRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BoundedReplayRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*BoundedReplayRequest) PositionId() uint16 {
	return 4
}

func (*BoundedReplayRequest) PositionSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) PositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.PositionSinceVersion()
}

func (*BoundedReplayRequest) PositionDeprecated() uint16 {
	return 0
}

func (*BoundedReplayRequest) PositionMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) PositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BoundedReplayRequest) PositionMaxValue() int64 {
	return math.MaxInt64
}

func (*BoundedReplayRequest) PositionNullValue() int64 {
	return math.MinInt64
}

func (*BoundedReplayRequest) LengthId() uint16 {
	return 5
}

func (*BoundedReplayRequest) LengthSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) LengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LengthSinceVersion()
}

func (*BoundedReplayRequest) LengthDeprecated() uint16 {
	return 0
}

func (*BoundedReplayRequest) LengthMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) LengthMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BoundedReplayRequest) LengthMaxValue() int64 {
	return math.MaxInt64
}

func (*BoundedReplayRequest) LengthNullValue() int64 {
	return math.MinInt64
}

func (*BoundedReplayRequest) LimitCounterIdId() uint16 {
	return 6
}

func (*BoundedReplayRequest) LimitCounterIdSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) LimitCounterIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LimitCounterIdSinceVersion()
}

func (*BoundedReplayRequest) LimitCounterIdDeprecated() uint16 {
	return 0
}

func (*BoundedReplayRequest) LimitCounterIdMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) LimitCounterIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*BoundedReplayRequest) LimitCounterIdMaxValue() int32 {
	return math.MaxInt32
}

func (*BoundedReplayRequest) LimitCounterIdNullValue() int32 {
	return math.MinInt32
}

func (*BoundedReplayRequest) ReplayStreamIdId() uint16 {
	return 7
}

func (*BoundedReplayRequest) ReplayStreamIdSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) ReplayStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ReplayStreamIdSinceVersion()
}

func (*BoundedReplayRequest) ReplayStreamIdDeprecated() uint16 {
	return 0
}

func (*BoundedReplayRequest) ReplayStreamIdMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) ReplayStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*BoundedReplayRequest) ReplayStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*BoundedReplayRequest) ReplayStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*BoundedReplayRequest) ReplayChannelMetaAttribute(meta int) string {
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

func (*BoundedReplayRequest) ReplayChannelSinceVersion() uint16 {
	return 0
}

func (b *BoundedReplayRequest) ReplayChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ReplayChannelSinceVersion()
}

func (*BoundedReplayRequest) ReplayChannelDeprecated() uint16 {
	return 0
}

func (BoundedReplayRequest) ReplayChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (BoundedReplayRequest) ReplayChannelHeaderLength() uint64 {
	return 4
}
