// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type JoinCluster struct {
	LeadershipTermId int64
	MemberId         int32
}

func (j *JoinCluster) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := j.RangeCheck(j.SbeSchemaVersion(), j.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, j.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, j.MemberId); err != nil {
		return err
	}
	return nil
}

func (j *JoinCluster) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !j.LeadershipTermIdInActingVersion(actingVersion) {
		j.LeadershipTermId = j.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &j.LeadershipTermId); err != nil {
			return err
		}
	}
	if !j.MemberIdInActingVersion(actingVersion) {
		j.MemberId = j.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &j.MemberId); err != nil {
			return err
		}
	}
	if actingVersion > j.SbeSchemaVersion() && blockLength > j.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-j.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := j.RangeCheck(actingVersion, j.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (j *JoinCluster) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if j.LeadershipTermIdInActingVersion(actingVersion) {
		if j.LeadershipTermId < j.LeadershipTermIdMinValue() || j.LeadershipTermId > j.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on j.LeadershipTermId (%v < %v > %v)", j.LeadershipTermIdMinValue(), j.LeadershipTermId, j.LeadershipTermIdMaxValue())
		}
	}
	if j.MemberIdInActingVersion(actingVersion) {
		if j.MemberId < j.MemberIdMinValue() || j.MemberId > j.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on j.MemberId (%v < %v > %v)", j.MemberIdMinValue(), j.MemberId, j.MemberIdMaxValue())
		}
	}
	return nil
}

func JoinClusterInit(j *JoinCluster) {
	return
}

func (*JoinCluster) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*JoinCluster) SbeTemplateId() (templateId uint16) {
	return 74
}

func (*JoinCluster) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*JoinCluster) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*JoinCluster) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*JoinCluster) LeadershipTermIdId() uint16 {
	return 1
}

func (*JoinCluster) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (j *JoinCluster) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.LeadershipTermIdSinceVersion()
}

func (*JoinCluster) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*JoinCluster) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*JoinCluster) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*JoinCluster) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*JoinCluster) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*JoinCluster) MemberIdId() uint16 {
	return 2
}

func (*JoinCluster) MemberIdSinceVersion() uint16 {
	return 0
}

func (j *JoinCluster) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.MemberIdSinceVersion()
}

func (*JoinCluster) MemberIdDeprecated() uint16 {
	return 0
}

func (*JoinCluster) MemberIdMetaAttribute(meta int) string {
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

func (*JoinCluster) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*JoinCluster) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*JoinCluster) MemberIdNullValue() int32 {
	return math.MinInt32
}
