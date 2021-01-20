// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingDescriptor struct {
	ControlSessionId  int64
	CorrelationId     int64
	RecordingId       int64
	StartTimestamp    int64
	StopTimestamp     int64
	StartPosition     int64
	StopPosition      int64
	InitialTermId     int32
	SegmentFileLength int32
	TermBufferLength  int32
	MtuLength         int32
	SessionId         int32
	StreamId          int32
	StrippedChannel   []uint8
	OriginalChannel   []uint8
	SourceIdentity    []uint8
}

func (r *RecordingDescriptor) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, r.StartTimestamp); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.StopTimestamp); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.StartPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.StopPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.InitialTermId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.SegmentFileLength); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.TermBufferLength); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.MtuLength); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.SessionId); err != nil {
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
	if err := _m.WriteUint32(_w, uint32(len(r.OriginalChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.OriginalChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.SourceIdentity))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.SourceIdentity); err != nil {
		return err
	}
	return nil
}

func (r *RecordingDescriptor) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !r.StartTimestampInActingVersion(actingVersion) {
		r.StartTimestamp = r.StartTimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.StartTimestamp); err != nil {
			return err
		}
	}
	if !r.StopTimestampInActingVersion(actingVersion) {
		r.StopTimestamp = r.StopTimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.StopTimestamp); err != nil {
			return err
		}
	}
	if !r.StartPositionInActingVersion(actingVersion) {
		r.StartPosition = r.StartPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.StartPosition); err != nil {
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
	if !r.InitialTermIdInActingVersion(actingVersion) {
		r.InitialTermId = r.InitialTermIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.InitialTermId); err != nil {
			return err
		}
	}
	if !r.SegmentFileLengthInActingVersion(actingVersion) {
		r.SegmentFileLength = r.SegmentFileLengthNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.SegmentFileLength); err != nil {
			return err
		}
	}
	if !r.TermBufferLengthInActingVersion(actingVersion) {
		r.TermBufferLength = r.TermBufferLengthNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.TermBufferLength); err != nil {
			return err
		}
	}
	if !r.MtuLengthInActingVersion(actingVersion) {
		r.MtuLength = r.MtuLengthNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.MtuLength); err != nil {
			return err
		}
	}
	if !r.SessionIdInActingVersion(actingVersion) {
		r.SessionId = r.SessionIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.SessionId); err != nil {
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

	if r.OriginalChannelInActingVersion(actingVersion) {
		var OriginalChannelLength uint32
		if err := _m.ReadUint32(_r, &OriginalChannelLength); err != nil {
			return err
		}
		if cap(r.OriginalChannel) < int(OriginalChannelLength) {
			r.OriginalChannel = make([]uint8, OriginalChannelLength)
		}
		r.OriginalChannel = r.OriginalChannel[:OriginalChannelLength]
		if err := _m.ReadBytes(_r, r.OriginalChannel); err != nil {
			return err
		}
	}

	if r.SourceIdentityInActingVersion(actingVersion) {
		var SourceIdentityLength uint32
		if err := _m.ReadUint32(_r, &SourceIdentityLength); err != nil {
			return err
		}
		if cap(r.SourceIdentity) < int(SourceIdentityLength) {
			r.SourceIdentity = make([]uint8, SourceIdentityLength)
		}
		r.SourceIdentity = r.SourceIdentity[:SourceIdentityLength]
		if err := _m.ReadBytes(_r, r.SourceIdentity); err != nil {
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

func (r *RecordingDescriptor) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if r.StartTimestampInActingVersion(actingVersion) {
		if r.StartTimestamp < r.StartTimestampMinValue() || r.StartTimestamp > r.StartTimestampMaxValue() {
			return fmt.Errorf("Range check failed on r.StartTimestamp (%v < %v > %v)", r.StartTimestampMinValue(), r.StartTimestamp, r.StartTimestampMaxValue())
		}
	}
	if r.StopTimestampInActingVersion(actingVersion) {
		if r.StopTimestamp < r.StopTimestampMinValue() || r.StopTimestamp > r.StopTimestampMaxValue() {
			return fmt.Errorf("Range check failed on r.StopTimestamp (%v < %v > %v)", r.StopTimestampMinValue(), r.StopTimestamp, r.StopTimestampMaxValue())
		}
	}
	if r.StartPositionInActingVersion(actingVersion) {
		if r.StartPosition < r.StartPositionMinValue() || r.StartPosition > r.StartPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.StartPosition (%v < %v > %v)", r.StartPositionMinValue(), r.StartPosition, r.StartPositionMaxValue())
		}
	}
	if r.StopPositionInActingVersion(actingVersion) {
		if r.StopPosition < r.StopPositionMinValue() || r.StopPosition > r.StopPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.StopPosition (%v < %v > %v)", r.StopPositionMinValue(), r.StopPosition, r.StopPositionMaxValue())
		}
	}
	if r.InitialTermIdInActingVersion(actingVersion) {
		if r.InitialTermId < r.InitialTermIdMinValue() || r.InitialTermId > r.InitialTermIdMaxValue() {
			return fmt.Errorf("Range check failed on r.InitialTermId (%v < %v > %v)", r.InitialTermIdMinValue(), r.InitialTermId, r.InitialTermIdMaxValue())
		}
	}
	if r.SegmentFileLengthInActingVersion(actingVersion) {
		if r.SegmentFileLength < r.SegmentFileLengthMinValue() || r.SegmentFileLength > r.SegmentFileLengthMaxValue() {
			return fmt.Errorf("Range check failed on r.SegmentFileLength (%v < %v > %v)", r.SegmentFileLengthMinValue(), r.SegmentFileLength, r.SegmentFileLengthMaxValue())
		}
	}
	if r.TermBufferLengthInActingVersion(actingVersion) {
		if r.TermBufferLength < r.TermBufferLengthMinValue() || r.TermBufferLength > r.TermBufferLengthMaxValue() {
			return fmt.Errorf("Range check failed on r.TermBufferLength (%v < %v > %v)", r.TermBufferLengthMinValue(), r.TermBufferLength, r.TermBufferLengthMaxValue())
		}
	}
	if r.MtuLengthInActingVersion(actingVersion) {
		if r.MtuLength < r.MtuLengthMinValue() || r.MtuLength > r.MtuLengthMaxValue() {
			return fmt.Errorf("Range check failed on r.MtuLength (%v < %v > %v)", r.MtuLengthMinValue(), r.MtuLength, r.MtuLengthMaxValue())
		}
	}
	if r.SessionIdInActingVersion(actingVersion) {
		if r.SessionId < r.SessionIdMinValue() || r.SessionId > r.SessionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SessionId (%v < %v > %v)", r.SessionIdMinValue(), r.SessionId, r.SessionIdMaxValue())
		}
	}
	if r.StreamIdInActingVersion(actingVersion) {
		if r.StreamId < r.StreamIdMinValue() || r.StreamId > r.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on r.StreamId (%v < %v > %v)", r.StreamIdMinValue(), r.StreamId, r.StreamIdMaxValue())
		}
	}
	return nil
}

func RecordingDescriptorInit(r *RecordingDescriptor) {
	return
}

func (*RecordingDescriptor) SbeBlockLength() (blockLength uint16) {
	return 80
}

func (*RecordingDescriptor) SbeTemplateId() (templateId uint16) {
	return 22
}

func (*RecordingDescriptor) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingDescriptor) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*RecordingDescriptor) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingDescriptor) ControlSessionIdId() uint16 {
	return 1
}

func (*RecordingDescriptor) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ControlSessionIdSinceVersion()
}

func (*RecordingDescriptor) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) ControlSessionIdMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingDescriptor) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingDescriptor) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingDescriptor) CorrelationIdId() uint16 {
	return 2
}

func (*RecordingDescriptor) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CorrelationIdSinceVersion()
}

func (*RecordingDescriptor) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) CorrelationIdMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingDescriptor) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingDescriptor) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingDescriptor) RecordingIdId() uint16 {
	return 3
}

func (*RecordingDescriptor) RecordingIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.RecordingIdSinceVersion()
}

func (*RecordingDescriptor) RecordingIdDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) RecordingIdMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingDescriptor) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingDescriptor) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingDescriptor) StartTimestampId() uint16 {
	return 4
}

func (*RecordingDescriptor) StartTimestampSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) StartTimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StartTimestampSinceVersion()
}

func (*RecordingDescriptor) StartTimestampDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) StartTimestampMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) StartTimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingDescriptor) StartTimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingDescriptor) StartTimestampNullValue() int64 {
	return math.MinInt64
}

func (*RecordingDescriptor) StopTimestampId() uint16 {
	return 5
}

func (*RecordingDescriptor) StopTimestampSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) StopTimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StopTimestampSinceVersion()
}

func (*RecordingDescriptor) StopTimestampDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) StopTimestampMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) StopTimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingDescriptor) StopTimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingDescriptor) StopTimestampNullValue() int64 {
	return math.MinInt64
}

func (*RecordingDescriptor) StartPositionId() uint16 {
	return 6
}

func (*RecordingDescriptor) StartPositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) StartPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StartPositionSinceVersion()
}

func (*RecordingDescriptor) StartPositionDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) StartPositionMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) StartPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingDescriptor) StartPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingDescriptor) StartPositionNullValue() int64 {
	return math.MinInt64
}

func (*RecordingDescriptor) StopPositionId() uint16 {
	return 7
}

func (*RecordingDescriptor) StopPositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) StopPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StopPositionSinceVersion()
}

func (*RecordingDescriptor) StopPositionDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) StopPositionMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) StopPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingDescriptor) StopPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingDescriptor) StopPositionNullValue() int64 {
	return math.MinInt64
}

func (*RecordingDescriptor) InitialTermIdId() uint16 {
	return 8
}

func (*RecordingDescriptor) InitialTermIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) InitialTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.InitialTermIdSinceVersion()
}

func (*RecordingDescriptor) InitialTermIdDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) InitialTermIdMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) InitialTermIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptor) InitialTermIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptor) InitialTermIdNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptor) SegmentFileLengthId() uint16 {
	return 9
}

func (*RecordingDescriptor) SegmentFileLengthSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) SegmentFileLengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SegmentFileLengthSinceVersion()
}

func (*RecordingDescriptor) SegmentFileLengthDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) SegmentFileLengthMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) SegmentFileLengthMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptor) SegmentFileLengthMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptor) SegmentFileLengthNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptor) TermBufferLengthId() uint16 {
	return 10
}

func (*RecordingDescriptor) TermBufferLengthSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) TermBufferLengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.TermBufferLengthSinceVersion()
}

func (*RecordingDescriptor) TermBufferLengthDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) TermBufferLengthMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) TermBufferLengthMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptor) TermBufferLengthMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptor) TermBufferLengthNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptor) MtuLengthId() uint16 {
	return 11
}

func (*RecordingDescriptor) MtuLengthSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) MtuLengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.MtuLengthSinceVersion()
}

func (*RecordingDescriptor) MtuLengthDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) MtuLengthMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) MtuLengthMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptor) MtuLengthMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptor) MtuLengthNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptor) SessionIdId() uint16 {
	return 12
}

func (*RecordingDescriptor) SessionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) SessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SessionIdSinceVersion()
}

func (*RecordingDescriptor) SessionIdDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) SessionIdMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) SessionIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptor) SessionIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptor) SessionIdNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptor) StreamIdId() uint16 {
	return 13
}

func (*RecordingDescriptor) StreamIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StreamIdSinceVersion()
}

func (*RecordingDescriptor) StreamIdDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptor) StreamIdMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptor) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptor) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptor) StrippedChannelMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) StrippedChannelSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) StrippedChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StrippedChannelSinceVersion()
}

func (*RecordingDescriptor) StrippedChannelDeprecated() uint16 {
	return 0
}

func (RecordingDescriptor) StrippedChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (RecordingDescriptor) StrippedChannelHeaderLength() uint64 {
	return 4
}

func (*RecordingDescriptor) OriginalChannelMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) OriginalChannelSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) OriginalChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.OriginalChannelSinceVersion()
}

func (*RecordingDescriptor) OriginalChannelDeprecated() uint16 {
	return 0
}

func (RecordingDescriptor) OriginalChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (RecordingDescriptor) OriginalChannelHeaderLength() uint64 {
	return 4
}

func (*RecordingDescriptor) SourceIdentityMetaAttribute(meta int) string {
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

func (*RecordingDescriptor) SourceIdentitySinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptor) SourceIdentityInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SourceIdentitySinceVersion()
}

func (*RecordingDescriptor) SourceIdentityDeprecated() uint16 {
	return 0
}

func (RecordingDescriptor) SourceIdentityCharacterEncoding() string {
	return "US-ASCII"
}

func (RecordingDescriptor) SourceIdentityHeaderLength() uint64 {
	return 4
}
