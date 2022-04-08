// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionCloseEvent struct {
	LeadershipTermId int64
	ClusterSessionId int64
	Timestamp        int64
	CloseReason      CloseReasonEnum
}

func (s *SessionCloseEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := s.CloseReason.Encode(_m, _w); err != nil {
		return err
	}
	return nil
}

func (s *SessionCloseEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if s.CloseReasonInActingVersion(actingVersion) {
		if err := s.CloseReason.Decode(_m, _r, actingVersion); err != nil {
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

func (s *SessionCloseEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if err := s.CloseReason.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func SessionCloseEventInit(s *SessionCloseEvent) {
	return
}

func (*SessionCloseEvent) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*SessionCloseEvent) SbeTemplateId() (templateId uint16) {
	return 22
}

func (*SessionCloseEvent) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SessionCloseEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SessionCloseEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionCloseEvent) LeadershipTermIdId() uint16 {
	return 1
}

func (*SessionCloseEvent) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *SessionCloseEvent) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*SessionCloseEvent) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*SessionCloseEvent) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*SessionCloseEvent) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionCloseEvent) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionCloseEvent) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionCloseEvent) ClusterSessionIdId() uint16 {
	return 2
}

func (*SessionCloseEvent) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (s *SessionCloseEvent) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ClusterSessionIdSinceVersion()
}

func (*SessionCloseEvent) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*SessionCloseEvent) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*SessionCloseEvent) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionCloseEvent) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionCloseEvent) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionCloseEvent) TimestampId() uint16 {
	return 3
}

func (*SessionCloseEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (s *SessionCloseEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TimestampSinceVersion()
}

func (*SessionCloseEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*SessionCloseEvent) TimestampMetaAttribute(meta int) string {
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

func (*SessionCloseEvent) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionCloseEvent) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionCloseEvent) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*SessionCloseEvent) CloseReasonId() uint16 {
	return 4
}

func (*SessionCloseEvent) CloseReasonSinceVersion() uint16 {
	return 0
}

func (s *SessionCloseEvent) CloseReasonInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CloseReasonSinceVersion()
}

func (*SessionCloseEvent) CloseReasonDeprecated() uint16 {
	return 0
}

func (*SessionCloseEvent) CloseReasonMetaAttribute(meta int) string {
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
