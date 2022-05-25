// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionEvent struct {
	ClusterSessionId int64
	CorrelationId    int64
	LeadershipTermId int64
	LeaderMemberId   int32
	Code             EventCodeEnum
	Version          int32
	Detail           []uint8
}

func (s *SessionEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.LeaderMemberId); err != nil {
		return err
	}
	if err := s.Code.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.Version); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.Detail))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.Detail); err != nil {
		return err
	}
	return nil
}

func (s *SessionEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.ClusterSessionIdInActingVersion(actingVersion) {
		s.ClusterSessionId = s.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.ClusterSessionId); err != nil {
			return err
		}
	}
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
			return err
		}
	}
	if !s.LeadershipTermIdInActingVersion(actingVersion) {
		s.LeadershipTermId = s.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LeadershipTermId); err != nil {
			return err
		}
	}
	if !s.LeaderMemberIdInActingVersion(actingVersion) {
		s.LeaderMemberId = s.LeaderMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.LeaderMemberId); err != nil {
			return err
		}
	}
	if s.CodeInActingVersion(actingVersion) {
		if err := s.Code.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !s.VersionInActingVersion(actingVersion) {
		s.Version = s.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.Version); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	if s.DetailInActingVersion(actingVersion) {
		var DetailLength uint32
		if err := _m.ReadUint32(_r, &DetailLength); err != nil {
			return err
		}
		if cap(s.Detail) < int(DetailLength) {
			s.Detail = make([]uint8, DetailLength)
		}
		s.Detail = s.Detail[:DetailLength]
		if err := _m.ReadBytes(_r, s.Detail); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.ClusterSessionIdInActingVersion(actingVersion) {
		if s.ClusterSessionId < s.ClusterSessionIdMinValue() || s.ClusterSessionId > s.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ClusterSessionId (%v < %v > %v)", s.ClusterSessionIdMinValue(), s.ClusterSessionId, s.ClusterSessionIdMaxValue())
		}
	}
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	if s.LeadershipTermIdInActingVersion(actingVersion) {
		if s.LeadershipTermId < s.LeadershipTermIdMinValue() || s.LeadershipTermId > s.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeadershipTermId (%v < %v > %v)", s.LeadershipTermIdMinValue(), s.LeadershipTermId, s.LeadershipTermIdMaxValue())
		}
	}
	if s.LeaderMemberIdInActingVersion(actingVersion) {
		if s.LeaderMemberId < s.LeaderMemberIdMinValue() || s.LeaderMemberId > s.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeaderMemberId (%v < %v > %v)", s.LeaderMemberIdMinValue(), s.LeaderMemberId, s.LeaderMemberIdMaxValue())
		}
	}
	if err := s.Code.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if s.VersionInActingVersion(actingVersion) {
		if s.Version != s.VersionNullValue() && (s.Version < s.VersionMinValue() || s.Version > s.VersionMaxValue()) {
			return fmt.Errorf("Range check failed on s.Version (%v < %v > %v)", s.VersionMinValue(), s.Version, s.VersionMaxValue())
		}
	}
	for idx, ch := range s.Detail {
		if ch > 127 {
			return fmt.Errorf("s.Detail[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func SessionEventInit(s *SessionEvent) {
	s.Version = 0
	return
}

func (*SessionEvent) SbeBlockLength() (blockLength uint16) {
	return 36
}

func (*SessionEvent) SbeTemplateId() (templateId uint16) {
	return 2
}

func (*SessionEvent) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SessionEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SessionEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionEvent) ClusterSessionIdId() uint16 {
	return 1
}

func (*SessionEvent) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (s *SessionEvent) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ClusterSessionIdSinceVersion()
}

func (*SessionEvent) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*SessionEvent) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*SessionEvent) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionEvent) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionEvent) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionEvent) CorrelationIdId() uint16 {
	return 2
}

func (*SessionEvent) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *SessionEvent) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*SessionEvent) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*SessionEvent) CorrelationIdMetaAttribute(meta int) string {
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

func (*SessionEvent) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionEvent) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionEvent) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionEvent) LeadershipTermIdId() uint16 {
	return 3
}

func (*SessionEvent) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *SessionEvent) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*SessionEvent) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*SessionEvent) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*SessionEvent) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionEvent) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionEvent) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionEvent) LeaderMemberIdId() uint16 {
	return 4
}

func (*SessionEvent) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (s *SessionEvent) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeaderMemberIdSinceVersion()
}

func (*SessionEvent) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*SessionEvent) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*SessionEvent) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*SessionEvent) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*SessionEvent) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*SessionEvent) CodeId() uint16 {
	return 5
}

func (*SessionEvent) CodeSinceVersion() uint16 {
	return 0
}

func (s *SessionEvent) CodeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CodeSinceVersion()
}

func (*SessionEvent) CodeDeprecated() uint16 {
	return 0
}

func (*SessionEvent) CodeMetaAttribute(meta int) string {
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

func (*SessionEvent) VersionId() uint16 {
	return 6
}

func (*SessionEvent) VersionSinceVersion() uint16 {
	return 6
}

func (s *SessionEvent) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.VersionSinceVersion()
}

func (*SessionEvent) VersionDeprecated() uint16 {
	return 0
}

func (*SessionEvent) VersionMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "optional"
	}
	return ""
}

func (*SessionEvent) VersionMinValue() int32 {
	return 1
}

func (*SessionEvent) VersionMaxValue() int32 {
	return 16777215
}

func (*SessionEvent) VersionNullValue() int32 {
	return 0
}

func (*SessionEvent) DetailMetaAttribute(meta int) string {
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

func (*SessionEvent) DetailSinceVersion() uint16 {
	return 0
}

func (s *SessionEvent) DetailInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.DetailSinceVersion()
}

func (*SessionEvent) DetailDeprecated() uint16 {
	return 0
}

func (SessionEvent) DetailCharacterEncoding() string {
	return "US-ASCII"
}

func (SessionEvent) DetailHeaderLength() uint64 {
	return 4
}
