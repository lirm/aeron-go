// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SnapshotRecordingQuery struct {
	CorrelationId   int64
	RequestMemberId int32
}

func (s *SnapshotRecordingQuery) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.RequestMemberId); err != nil {
		return err
	}
	return nil
}

func (s *SnapshotRecordingQuery) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
			return err
		}
	}
	if !s.RequestMemberIdInActingVersion(actingVersion) {
		s.RequestMemberId = s.RequestMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.RequestMemberId); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapshotRecordingQuery) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	if s.RequestMemberIdInActingVersion(actingVersion) {
		if s.RequestMemberId < s.RequestMemberIdMinValue() || s.RequestMemberId > s.RequestMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on s.RequestMemberId (%v < %v > %v)", s.RequestMemberIdMinValue(), s.RequestMemberId, s.RequestMemberIdMaxValue())
		}
	}
	return nil
}

func SnapshotRecordingQueryInit(s *SnapshotRecordingQuery) {
	return
}

func (*SnapshotRecordingQuery) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*SnapshotRecordingQuery) SbeTemplateId() (templateId uint16) {
	return 72
}

func (*SnapshotRecordingQuery) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SnapshotRecordingQuery) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SnapshotRecordingQuery) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SnapshotRecordingQuery) CorrelationIdId() uint16 {
	return 1
}

func (*SnapshotRecordingQuery) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingQuery) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*SnapshotRecordingQuery) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingQuery) CorrelationIdMetaAttribute(meta int) string {
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

func (*SnapshotRecordingQuery) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotRecordingQuery) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotRecordingQuery) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotRecordingQuery) RequestMemberIdId() uint16 {
	return 2
}

func (*SnapshotRecordingQuery) RequestMemberIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingQuery) RequestMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.RequestMemberIdSinceVersion()
}

func (*SnapshotRecordingQuery) RequestMemberIdDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingQuery) RequestMemberIdMetaAttribute(meta int) string {
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

func (*SnapshotRecordingQuery) RequestMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*SnapshotRecordingQuery) RequestMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*SnapshotRecordingQuery) RequestMemberIdNullValue() int32 {
	return math.MinInt32
}
