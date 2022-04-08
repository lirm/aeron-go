// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionMessageHeader struct {
	LeadershipTermId int64
	ClusterSessionId int64
	Timestamp        int64
}

func (s *SessionMessageHeader) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.Timestamp); err != nil {
		return err
	}
	return nil
}

func (s *SessionMessageHeader) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.LeadershipTermIdInActingVersion(actingVersion) {
		s.LeadershipTermId = s.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LeadershipTermId); err != nil {
			return err
		}
	}
	if !s.ClusterSessionIdInActingVersion(actingVersion) {
		s.ClusterSessionId = s.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.ClusterSessionId); err != nil {
			return err
		}
	}
	if !s.TimestampInActingVersion(actingVersion) {
		s.Timestamp = s.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.Timestamp); err != nil {
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

func (s *SessionMessageHeader) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.LeadershipTermIdInActingVersion(actingVersion) {
		if s.LeadershipTermId < s.LeadershipTermIdMinValue() || s.LeadershipTermId > s.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeadershipTermId (%v < %v > %v)", s.LeadershipTermIdMinValue(), s.LeadershipTermId, s.LeadershipTermIdMaxValue())
		}
	}
	if s.ClusterSessionIdInActingVersion(actingVersion) {
		if s.ClusterSessionId < s.ClusterSessionIdMinValue() || s.ClusterSessionId > s.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ClusterSessionId (%v < %v > %v)", s.ClusterSessionIdMinValue(), s.ClusterSessionId, s.ClusterSessionIdMaxValue())
		}
	}
	if s.TimestampInActingVersion(actingVersion) {
		if s.Timestamp < s.TimestampMinValue() || s.Timestamp > s.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on s.Timestamp (%v < %v > %v)", s.TimestampMinValue(), s.Timestamp, s.TimestampMaxValue())
		}
	}
	return nil
}

func SessionMessageHeaderInit(s *SessionMessageHeader) {
	return
}

func (*SessionMessageHeader) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*SessionMessageHeader) SbeTemplateId() (templateId uint16) {
	return 1
}

func (*SessionMessageHeader) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SessionMessageHeader) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SessionMessageHeader) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionMessageHeader) LeadershipTermIdId() uint16 {
	return 1
}

func (*SessionMessageHeader) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *SessionMessageHeader) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*SessionMessageHeader) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*SessionMessageHeader) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*SessionMessageHeader) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionMessageHeader) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionMessageHeader) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionMessageHeader) ClusterSessionIdId() uint16 {
	return 2
}

func (*SessionMessageHeader) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (s *SessionMessageHeader) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ClusterSessionIdSinceVersion()
}

func (*SessionMessageHeader) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*SessionMessageHeader) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*SessionMessageHeader) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionMessageHeader) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionMessageHeader) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionMessageHeader) TimestampId() uint16 {
	return 3
}

func (*SessionMessageHeader) TimestampSinceVersion() uint16 {
	return 0
}

func (s *SessionMessageHeader) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TimestampSinceVersion()
}

func (*SessionMessageHeader) TimestampDeprecated() uint16 {
	return 0
}

func (*SessionMessageHeader) TimestampMetaAttribute(meta int) string {
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

func (*SessionMessageHeader) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionMessageHeader) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionMessageHeader) TimestampNullValue() int64 {
	return math.MinInt64
}
