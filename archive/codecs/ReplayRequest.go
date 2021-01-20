// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ReplayRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	Position         int64
	Length           int64
	ReplayStreamId   int32
	ReplayChannel    []uint8
}

func (r *ReplayRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, r.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.Position); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.Length); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.ReplayStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.ReplayChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.ReplayChannel); err != nil {
		return err
	}
	return nil
}

func (r *ReplayRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !r.RecordingIdInActingVersion(actingVersion) {
		r.RecordingId = r.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.RecordingId); err != nil {
			return err
		}
	}
	if !r.PositionInActingVersion(actingVersion) {
		r.Position = r.PositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.Position); err != nil {
			return err
		}
	}
	if !r.LengthInActingVersion(actingVersion) {
		r.Length = r.LengthNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.Length); err != nil {
			return err
		}
	}
	if !r.ReplayStreamIdInActingVersion(actingVersion) {
		r.ReplayStreamId = r.ReplayStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.ReplayStreamId); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}

	if r.ReplayChannelInActingVersion(actingVersion) {
		var ReplayChannelLength uint32
		if err := _m.ReadUint32(_r, &ReplayChannelLength); err != nil {
			return err
		}
		if cap(r.ReplayChannel) < int(ReplayChannelLength) {
			r.ReplayChannel = make([]uint8, ReplayChannelLength)
		}
		r.ReplayChannel = r.ReplayChannel[:ReplayChannelLength]
		if err := _m.ReadBytes(_r, r.ReplayChannel); err != nil {
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

func (r *ReplayRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if r.RecordingIdInActingVersion(actingVersion) {
		if r.RecordingId < r.RecordingIdMinValue() || r.RecordingId > r.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.RecordingId (%v < %v > %v)", r.RecordingIdMinValue(), r.RecordingId, r.RecordingIdMaxValue())
		}
	}
	if r.PositionInActingVersion(actingVersion) {
		if r.Position < r.PositionMinValue() || r.Position > r.PositionMaxValue() {
			return fmt.Errorf("Range check failed on r.Position (%v < %v > %v)", r.PositionMinValue(), r.Position, r.PositionMaxValue())
		}
	}
	if r.LengthInActingVersion(actingVersion) {
		if r.Length < r.LengthMinValue() || r.Length > r.LengthMaxValue() {
			return fmt.Errorf("Range check failed on r.Length (%v < %v > %v)", r.LengthMinValue(), r.Length, r.LengthMaxValue())
		}
	}
	if r.ReplayStreamIdInActingVersion(actingVersion) {
		if r.ReplayStreamId < r.ReplayStreamIdMinValue() || r.ReplayStreamId > r.ReplayStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on r.ReplayStreamId (%v < %v > %v)", r.ReplayStreamIdMinValue(), r.ReplayStreamId, r.ReplayStreamIdMaxValue())
		}
	}
	return nil
}

func ReplayRequestInit(r *ReplayRequest) {
	return
}

func (*ReplayRequest) SbeBlockLength() (blockLength uint16) {
	return 44
}

func (*ReplayRequest) SbeTemplateId() (templateId uint16) {
	return 6
}

func (*ReplayRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ReplayRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*ReplayRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ReplayRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*ReplayRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (r *ReplayRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ControlSessionIdSinceVersion()
}

func (*ReplayRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ReplayRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ReplayRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplayRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplayRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplayRequest) CorrelationIdId() uint16 {
	return 2
}

func (*ReplayRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (r *ReplayRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CorrelationIdSinceVersion()
}

func (*ReplayRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ReplayRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ReplayRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplayRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplayRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplayRequest) RecordingIdId() uint16 {
	return 3
}

func (*ReplayRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (r *ReplayRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.RecordingIdSinceVersion()
}

func (*ReplayRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*ReplayRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*ReplayRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplayRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplayRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplayRequest) PositionId() uint16 {
	return 4
}

func (*ReplayRequest) PositionSinceVersion() uint16 {
	return 0
}

func (r *ReplayRequest) PositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.PositionSinceVersion()
}

func (*ReplayRequest) PositionDeprecated() uint16 {
	return 0
}

func (*ReplayRequest) PositionMetaAttribute(meta int) string {
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

func (*ReplayRequest) PositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplayRequest) PositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplayRequest) PositionNullValue() int64 {
	return math.MinInt64
}

func (*ReplayRequest) LengthId() uint16 {
	return 5
}

func (*ReplayRequest) LengthSinceVersion() uint16 {
	return 0
}

func (r *ReplayRequest) LengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.LengthSinceVersion()
}

func (*ReplayRequest) LengthDeprecated() uint16 {
	return 0
}

func (*ReplayRequest) LengthMetaAttribute(meta int) string {
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

func (*ReplayRequest) LengthMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplayRequest) LengthMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplayRequest) LengthNullValue() int64 {
	return math.MinInt64
}

func (*ReplayRequest) ReplayStreamIdId() uint16 {
	return 6
}

func (*ReplayRequest) ReplayStreamIdSinceVersion() uint16 {
	return 0
}

func (r *ReplayRequest) ReplayStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ReplayStreamIdSinceVersion()
}

func (*ReplayRequest) ReplayStreamIdDeprecated() uint16 {
	return 0
}

func (*ReplayRequest) ReplayStreamIdMetaAttribute(meta int) string {
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

func (*ReplayRequest) ReplayStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ReplayRequest) ReplayStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ReplayRequest) ReplayStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ReplayRequest) ReplayChannelMetaAttribute(meta int) string {
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

func (*ReplayRequest) ReplayChannelSinceVersion() uint16 {
	return 0
}

func (r *ReplayRequest) ReplayChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ReplayChannelSinceVersion()
}

func (*ReplayRequest) ReplayChannelDeprecated() uint16 {
	return 0
}

func (ReplayRequest) ReplayChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ReplayRequest) ReplayChannelHeaderLength() uint64 {
	return 4
}
