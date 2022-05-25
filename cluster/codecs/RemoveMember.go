// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RemoveMember struct {
	MemberId  int32
	IsPassive BooleanTypeEnum
}

func (r *RemoveMember) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt32(_w, r.MemberId); err != nil {
		return err
	}
	if err := r.IsPassive.Encode(_m, _w); err != nil {
		return err
	}
	return nil
}

func (r *RemoveMember) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.MemberIdInActingVersion(actingVersion) {
		r.MemberId = r.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.MemberId); err != nil {
			return err
		}
	}
	if r.IsPassiveInActingVersion(actingVersion) {
		if err := r.IsPassive.Decode(_m, _r, actingVersion); err != nil {
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

func (r *RemoveMember) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.MemberIdInActingVersion(actingVersion) {
		if r.MemberId < r.MemberIdMinValue() || r.MemberId > r.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on r.MemberId (%v < %v > %v)", r.MemberIdMinValue(), r.MemberId, r.MemberIdMaxValue())
		}
	}
	if err := r.IsPassive.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func RemoveMemberInit(r *RemoveMember) {
	return
}

func (*RemoveMember) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*RemoveMember) SbeTemplateId() (templateId uint16) {
	return 35
}

func (*RemoveMember) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*RemoveMember) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*RemoveMember) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RemoveMember) MemberIdId() uint16 {
	return 1
}

func (*RemoveMember) MemberIdSinceVersion() uint16 {
	return 0
}

func (r *RemoveMember) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.MemberIdSinceVersion()
}

func (*RemoveMember) MemberIdDeprecated() uint16 {
	return 0
}

func (*RemoveMember) MemberIdMetaAttribute(meta int) string {
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

func (*RemoveMember) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RemoveMember) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RemoveMember) MemberIdNullValue() int32 {
	return math.MinInt32
}

func (*RemoveMember) IsPassiveId() uint16 {
	return 2
}

func (*RemoveMember) IsPassiveSinceVersion() uint16 {
	return 0
}

func (r *RemoveMember) IsPassiveInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.IsPassiveSinceVersion()
}

func (*RemoveMember) IsPassiveDeprecated() uint16 {
	return 0
}

func (*RemoveMember) IsPassiveMetaAttribute(meta int) string {
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
