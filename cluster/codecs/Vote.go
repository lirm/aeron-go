// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Vote struct {
	CandidateTermId     int64
	LogLeadershipTermId int64
	LogPosition         int64
	CandidateMemberId   int32
	FollowerMemberId    int32
	Vote                BooleanTypeEnum
}

func (v *Vote) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := v.RangeCheck(v.SbeSchemaVersion(), v.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, v.CandidateTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, v.LogLeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, v.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, v.CandidateMemberId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, v.FollowerMemberId); err != nil {
		return err
	}
	if err := v.Vote.Encode(_m, _w); err != nil {
		return err
	}
	return nil
}

func (v *Vote) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !v.CandidateTermIdInActingVersion(actingVersion) {
		v.CandidateTermId = v.CandidateTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &v.CandidateTermId); err != nil {
			return err
		}
	}
	if !v.LogLeadershipTermIdInActingVersion(actingVersion) {
		v.LogLeadershipTermId = v.LogLeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &v.LogLeadershipTermId); err != nil {
			return err
		}
	}
	if !v.LogPositionInActingVersion(actingVersion) {
		v.LogPosition = v.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &v.LogPosition); err != nil {
			return err
		}
	}
	if !v.CandidateMemberIdInActingVersion(actingVersion) {
		v.CandidateMemberId = v.CandidateMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &v.CandidateMemberId); err != nil {
			return err
		}
	}
	if !v.FollowerMemberIdInActingVersion(actingVersion) {
		v.FollowerMemberId = v.FollowerMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &v.FollowerMemberId); err != nil {
			return err
		}
	}
	if v.VoteInActingVersion(actingVersion) {
		if err := v.Vote.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if actingVersion > v.SbeSchemaVersion() && blockLength > v.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-v.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := v.RangeCheck(actingVersion, v.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vote) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if v.CandidateTermIdInActingVersion(actingVersion) {
		if v.CandidateTermId < v.CandidateTermIdMinValue() || v.CandidateTermId > v.CandidateTermIdMaxValue() {
			return fmt.Errorf("Range check failed on v.CandidateTermId (%v < %v > %v)", v.CandidateTermIdMinValue(), v.CandidateTermId, v.CandidateTermIdMaxValue())
		}
	}
	if v.LogLeadershipTermIdInActingVersion(actingVersion) {
		if v.LogLeadershipTermId < v.LogLeadershipTermIdMinValue() || v.LogLeadershipTermId > v.LogLeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on v.LogLeadershipTermId (%v < %v > %v)", v.LogLeadershipTermIdMinValue(), v.LogLeadershipTermId, v.LogLeadershipTermIdMaxValue())
		}
	}
	if v.LogPositionInActingVersion(actingVersion) {
		if v.LogPosition < v.LogPositionMinValue() || v.LogPosition > v.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on v.LogPosition (%v < %v > %v)", v.LogPositionMinValue(), v.LogPosition, v.LogPositionMaxValue())
		}
	}
	if v.CandidateMemberIdInActingVersion(actingVersion) {
		if v.CandidateMemberId < v.CandidateMemberIdMinValue() || v.CandidateMemberId > v.CandidateMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on v.CandidateMemberId (%v < %v > %v)", v.CandidateMemberIdMinValue(), v.CandidateMemberId, v.CandidateMemberIdMaxValue())
		}
	}
	if v.FollowerMemberIdInActingVersion(actingVersion) {
		if v.FollowerMemberId < v.FollowerMemberIdMinValue() || v.FollowerMemberId > v.FollowerMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on v.FollowerMemberId (%v < %v > %v)", v.FollowerMemberIdMinValue(), v.FollowerMemberId, v.FollowerMemberIdMaxValue())
		}
	}
	if err := v.Vote.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func VoteInit(v *Vote) {
	return
}

func (*Vote) SbeBlockLength() (blockLength uint16) {
	return 36
}

func (*Vote) SbeTemplateId() (templateId uint16) {
	return 52
}

func (*Vote) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*Vote) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*Vote) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*Vote) CandidateTermIdId() uint16 {
	return 1
}

func (*Vote) CandidateTermIdSinceVersion() uint16 {
	return 0
}

func (v *Vote) CandidateTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.CandidateTermIdSinceVersion()
}

func (*Vote) CandidateTermIdDeprecated() uint16 {
	return 0
}

func (*Vote) CandidateTermIdMetaAttribute(meta int) string {
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

func (*Vote) CandidateTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Vote) CandidateTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*Vote) CandidateTermIdNullValue() int64 {
	return math.MinInt64
}

func (*Vote) LogLeadershipTermIdId() uint16 {
	return 2
}

func (*Vote) LogLeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (v *Vote) LogLeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.LogLeadershipTermIdSinceVersion()
}

func (*Vote) LogLeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*Vote) LogLeadershipTermIdMetaAttribute(meta int) string {
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

func (*Vote) LogLeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Vote) LogLeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*Vote) LogLeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*Vote) LogPositionId() uint16 {
	return 3
}

func (*Vote) LogPositionSinceVersion() uint16 {
	return 0
}

func (v *Vote) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.LogPositionSinceVersion()
}

func (*Vote) LogPositionDeprecated() uint16 {
	return 0
}

func (*Vote) LogPositionMetaAttribute(meta int) string {
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

func (*Vote) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Vote) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*Vote) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*Vote) CandidateMemberIdId() uint16 {
	return 4
}

func (*Vote) CandidateMemberIdSinceVersion() uint16 {
	return 0
}

func (v *Vote) CandidateMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.CandidateMemberIdSinceVersion()
}

func (*Vote) CandidateMemberIdDeprecated() uint16 {
	return 0
}

func (*Vote) CandidateMemberIdMetaAttribute(meta int) string {
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

func (*Vote) CandidateMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*Vote) CandidateMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*Vote) CandidateMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*Vote) FollowerMemberIdId() uint16 {
	return 5
}

func (*Vote) FollowerMemberIdSinceVersion() uint16 {
	return 0
}

func (v *Vote) FollowerMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.FollowerMemberIdSinceVersion()
}

func (*Vote) FollowerMemberIdDeprecated() uint16 {
	return 0
}

func (*Vote) FollowerMemberIdMetaAttribute(meta int) string {
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

func (*Vote) FollowerMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*Vote) FollowerMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*Vote) FollowerMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*Vote) VoteId() uint16 {
	return 6
}

func (*Vote) VoteSinceVersion() uint16 {
	return 0
}

func (v *Vote) VoteInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.VoteSinceVersion()
}

func (*Vote) VoteDeprecated() uint16 {
	return 0
}

func (*Vote) VoteMetaAttribute(meta int) string {
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
