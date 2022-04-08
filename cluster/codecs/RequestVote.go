// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RequestVote struct {
	LogLeadershipTermId int64
	LogPosition         int64
	CandidateTermId     int64
	CandidateMemberId   int32
}

func (r *RequestVote) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.LogLeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.CandidateTermId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.CandidateMemberId); err != nil {
		return err
	}
	return nil
}

func (r *RequestVote) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.LogLeadershipTermIdInActingVersion(actingVersion) {
		r.LogLeadershipTermId = r.LogLeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.LogLeadershipTermId); err != nil {
			return err
		}
	}
	if !r.LogPositionInActingVersion(actingVersion) {
		r.LogPosition = r.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.LogPosition); err != nil {
			return err
		}
	}
	if !r.CandidateTermIdInActingVersion(actingVersion) {
		r.CandidateTermId = r.CandidateTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.CandidateTermId); err != nil {
			return err
		}
	}
	if !r.CandidateMemberIdInActingVersion(actingVersion) {
		r.CandidateMemberId = r.CandidateMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.CandidateMemberId); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := r.RangeCheck(actingVersion, r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (r *RequestVote) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.LogLeadershipTermIdInActingVersion(actingVersion) {
		if r.LogLeadershipTermId < r.LogLeadershipTermIdMinValue() || r.LogLeadershipTermId > r.LogLeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on r.LogLeadershipTermId (%v < %v > %v)", r.LogLeadershipTermIdMinValue(), r.LogLeadershipTermId, r.LogLeadershipTermIdMaxValue())
		}
	}
	if r.LogPositionInActingVersion(actingVersion) {
		if r.LogPosition < r.LogPositionMinValue() || r.LogPosition > r.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.LogPosition (%v < %v > %v)", r.LogPositionMinValue(), r.LogPosition, r.LogPositionMaxValue())
		}
	}
	if r.CandidateTermIdInActingVersion(actingVersion) {
		if r.CandidateTermId < r.CandidateTermIdMinValue() || r.CandidateTermId > r.CandidateTermIdMaxValue() {
			return fmt.Errorf("Range check failed on r.CandidateTermId (%v < %v > %v)", r.CandidateTermIdMinValue(), r.CandidateTermId, r.CandidateTermIdMaxValue())
		}
	}
	if r.CandidateMemberIdInActingVersion(actingVersion) {
		if r.CandidateMemberId < r.CandidateMemberIdMinValue() || r.CandidateMemberId > r.CandidateMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on r.CandidateMemberId (%v < %v > %v)", r.CandidateMemberIdMinValue(), r.CandidateMemberId, r.CandidateMemberIdMaxValue())
		}
	}
	return nil
}

func RequestVoteInit(r *RequestVote) {
	return
}

func (*RequestVote) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*RequestVote) SbeTemplateId() (templateId uint16) {
	return 51
}

func (*RequestVote) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*RequestVote) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*RequestVote) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RequestVote) LogLeadershipTermIdId() uint16 {
	return 1
}

func (*RequestVote) LogLeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (r *RequestVote) LogLeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.LogLeadershipTermIdSinceVersion()
}

func (*RequestVote) LogLeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*RequestVote) LogLeadershipTermIdMetaAttribute(meta int) string {
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

func (*RequestVote) LogLeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RequestVote) LogLeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RequestVote) LogLeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*RequestVote) LogPositionId() uint16 {
	return 2
}

func (*RequestVote) LogPositionSinceVersion() uint16 {
	return 0
}

func (r *RequestVote) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.LogPositionSinceVersion()
}

func (*RequestVote) LogPositionDeprecated() uint16 {
	return 0
}

func (*RequestVote) LogPositionMetaAttribute(meta int) string {
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

func (*RequestVote) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RequestVote) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RequestVote) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*RequestVote) CandidateTermIdId() uint16 {
	return 3
}

func (*RequestVote) CandidateTermIdSinceVersion() uint16 {
	return 0
}

func (r *RequestVote) CandidateTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CandidateTermIdSinceVersion()
}

func (*RequestVote) CandidateTermIdDeprecated() uint16 {
	return 0
}

func (*RequestVote) CandidateTermIdMetaAttribute(meta int) string {
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

func (*RequestVote) CandidateTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RequestVote) CandidateTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RequestVote) CandidateTermIdNullValue() int64 {
	return math.MinInt64
}

func (*RequestVote) CandidateMemberIdId() uint16 {
	return 4
}

func (*RequestVote) CandidateMemberIdSinceVersion() uint16 {
	return 0
}

func (r *RequestVote) CandidateMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CandidateMemberIdSinceVersion()
}

func (*RequestVote) CandidateMemberIdDeprecated() uint16 {
	return 0
}

func (*RequestVote) CandidateMemberIdMetaAttribute(meta int) string {
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

func (*RequestVote) CandidateMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RequestVote) CandidateMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RequestVote) CandidateMemberIdNullValue() int32 {
	return math.MinInt32
}
