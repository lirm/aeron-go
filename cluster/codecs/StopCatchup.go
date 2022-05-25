// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StopCatchup struct {
	LeadershipTermId int64
	FollowerMemberId int32
}

func (s *StopCatchup) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.FollowerMemberId); err != nil {
		return err
	}
	return nil
}

func (s *StopCatchup) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.LeadershipTermIdInActingVersion(actingVersion) {
		s.LeadershipTermId = s.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LeadershipTermId); err != nil {
			return err
		}
	}
	if !s.FollowerMemberIdInActingVersion(actingVersion) {
		s.FollowerMemberId = s.FollowerMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.FollowerMemberId); err != nil {
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

func (s *StopCatchup) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.LeadershipTermIdInActingVersion(actingVersion) {
		if s.LeadershipTermId < s.LeadershipTermIdMinValue() || s.LeadershipTermId > s.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeadershipTermId (%v < %v > %v)", s.LeadershipTermIdMinValue(), s.LeadershipTermId, s.LeadershipTermIdMaxValue())
		}
	}
	if s.FollowerMemberIdInActingVersion(actingVersion) {
		if s.FollowerMemberId < s.FollowerMemberIdMinValue() || s.FollowerMemberId > s.FollowerMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on s.FollowerMemberId (%v < %v > %v)", s.FollowerMemberIdMinValue(), s.FollowerMemberId, s.FollowerMemberIdMaxValue())
		}
	}
	return nil
}

func StopCatchupInit(s *StopCatchup) {
	return
}

func (*StopCatchup) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*StopCatchup) SbeTemplateId() (templateId uint16) {
	return 57
}

func (*StopCatchup) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*StopCatchup) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*StopCatchup) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StopCatchup) LeadershipTermIdId() uint16 {
	return 1
}

func (*StopCatchup) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *StopCatchup) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*StopCatchup) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*StopCatchup) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*StopCatchup) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopCatchup) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopCatchup) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*StopCatchup) FollowerMemberIdId() uint16 {
	return 2
}

func (*StopCatchup) FollowerMemberIdSinceVersion() uint16 {
	return 0
}

func (s *StopCatchup) FollowerMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.FollowerMemberIdSinceVersion()
}

func (*StopCatchup) FollowerMemberIdDeprecated() uint16 {
	return 0
}

func (*StopCatchup) FollowerMemberIdMetaAttribute(meta int) string {
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

func (*StopCatchup) FollowerMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*StopCatchup) FollowerMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*StopCatchup) FollowerMemberIdNullValue() int32 {
	return math.MinInt32
}
