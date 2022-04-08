// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type AppendPosition struct {
	LeadershipTermId int64
	LogPosition      int64
	FollowerMemberId int32
	Flags            uint8
}

func (a *AppendPosition) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := a.RangeCheck(a.SbeSchemaVersion(), a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, a.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, a.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, a.FollowerMemberId); err != nil {
		return err
	}
	if err := _m.WriteUint8(_w, a.Flags); err != nil {
		return err
	}
	return nil
}

func (a *AppendPosition) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !a.LeadershipTermIdInActingVersion(actingVersion) {
		a.LeadershipTermId = a.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.LeadershipTermId); err != nil {
			return err
		}
	}
	if !a.LogPositionInActingVersion(actingVersion) {
		a.LogPosition = a.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.LogPosition); err != nil {
			return err
		}
	}
	if !a.FollowerMemberIdInActingVersion(actingVersion) {
		a.FollowerMemberId = a.FollowerMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &a.FollowerMemberId); err != nil {
			return err
		}
	}
	if !a.FlagsInActingVersion(actingVersion) {
		a.Flags = a.FlagsNullValue()
	} else {
		if err := _m.ReadUint8(_r, &a.Flags); err != nil {
			return err
		}
	}
	if actingVersion > a.SbeSchemaVersion() && blockLength > a.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-a.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := a.RangeCheck(actingVersion, a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (a *AppendPosition) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if a.LeadershipTermIdInActingVersion(actingVersion) {
		if a.LeadershipTermId < a.LeadershipTermIdMinValue() || a.LeadershipTermId > a.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on a.LeadershipTermId (%v < %v > %v)", a.LeadershipTermIdMinValue(), a.LeadershipTermId, a.LeadershipTermIdMaxValue())
		}
	}
	if a.LogPositionInActingVersion(actingVersion) {
		if a.LogPosition < a.LogPositionMinValue() || a.LogPosition > a.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on a.LogPosition (%v < %v > %v)", a.LogPositionMinValue(), a.LogPosition, a.LogPositionMaxValue())
		}
	}
	if a.FollowerMemberIdInActingVersion(actingVersion) {
		if a.FollowerMemberId < a.FollowerMemberIdMinValue() || a.FollowerMemberId > a.FollowerMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on a.FollowerMemberId (%v < %v > %v)", a.FollowerMemberIdMinValue(), a.FollowerMemberId, a.FollowerMemberIdMaxValue())
		}
	}
	if a.FlagsInActingVersion(actingVersion) {
		if a.Flags < a.FlagsMinValue() || a.Flags > a.FlagsMaxValue() {
			return fmt.Errorf("Range check failed on a.Flags (%v < %v > %v)", a.FlagsMinValue(), a.Flags, a.FlagsMaxValue())
		}
	}
	return nil
}

func AppendPositionInit(a *AppendPosition) {
	return
}

func (*AppendPosition) SbeBlockLength() (blockLength uint16) {
	return 21
}

func (*AppendPosition) SbeTemplateId() (templateId uint16) {
	return 54
}

func (*AppendPosition) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*AppendPosition) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*AppendPosition) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*AppendPosition) LeadershipTermIdId() uint16 {
	return 1
}

func (*AppendPosition) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (a *AppendPosition) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.LeadershipTermIdSinceVersion()
}

func (*AppendPosition) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*AppendPosition) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*AppendPosition) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AppendPosition) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AppendPosition) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*AppendPosition) LogPositionId() uint16 {
	return 2
}

func (*AppendPosition) LogPositionSinceVersion() uint16 {
	return 0
}

func (a *AppendPosition) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.LogPositionSinceVersion()
}

func (*AppendPosition) LogPositionDeprecated() uint16 {
	return 0
}

func (*AppendPosition) LogPositionMetaAttribute(meta int) string {
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

func (*AppendPosition) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AppendPosition) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*AppendPosition) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*AppendPosition) FollowerMemberIdId() uint16 {
	return 3
}

func (*AppendPosition) FollowerMemberIdSinceVersion() uint16 {
	return 0
}

func (a *AppendPosition) FollowerMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.FollowerMemberIdSinceVersion()
}

func (*AppendPosition) FollowerMemberIdDeprecated() uint16 {
	return 0
}

func (*AppendPosition) FollowerMemberIdMetaAttribute(meta int) string {
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

func (*AppendPosition) FollowerMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*AppendPosition) FollowerMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*AppendPosition) FollowerMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*AppendPosition) FlagsId() uint16 {
	return 4
}

func (*AppendPosition) FlagsSinceVersion() uint16 {
	return 8
}

func (a *AppendPosition) FlagsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.FlagsSinceVersion()
}

func (*AppendPosition) FlagsDeprecated() uint16 {
	return 0
}

func (*AppendPosition) FlagsMetaAttribute(meta int) string {
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

func (*AppendPosition) FlagsMinValue() uint8 {
	return 0
}

func (*AppendPosition) FlagsMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*AppendPosition) FlagsNullValue() uint8 {
	return math.MaxUint8
}
