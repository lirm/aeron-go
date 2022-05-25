// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClusterMembers struct {
	MemberId       int32
	HighMemberId   int32
	ClusterMembers []uint8
}

func (c *ClusterMembers) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt32(_w, c.MemberId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.HighMemberId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ClusterMembers))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ClusterMembers); err != nil {
		return err
	}
	return nil
}

func (c *ClusterMembers) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.MemberIdInActingVersion(actingVersion) {
		c.MemberId = c.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.MemberId); err != nil {
			return err
		}
	}
	if !c.HighMemberIdInActingVersion(actingVersion) {
		c.HighMemberId = c.HighMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.HighMemberId); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.ClusterMembersInActingVersion(actingVersion) {
		var ClusterMembersLength uint32
		if err := _m.ReadUint32(_r, &ClusterMembersLength); err != nil {
			return err
		}
		if cap(c.ClusterMembers) < int(ClusterMembersLength) {
			c.ClusterMembers = make([]uint8, ClusterMembersLength)
		}
		c.ClusterMembers = c.ClusterMembers[:ClusterMembersLength]
		if err := _m.ReadBytes(_r, c.ClusterMembers); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterMembers) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.MemberIdInActingVersion(actingVersion) {
		if c.MemberId < c.MemberIdMinValue() || c.MemberId > c.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.MemberId (%v < %v > %v)", c.MemberIdMinValue(), c.MemberId, c.MemberIdMaxValue())
		}
	}
	if c.HighMemberIdInActingVersion(actingVersion) {
		if c.HighMemberId < c.HighMemberIdMinValue() || c.HighMemberId > c.HighMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.HighMemberId (%v < %v > %v)", c.HighMemberIdMinValue(), c.HighMemberId, c.HighMemberIdMaxValue())
		}
	}
	for idx, ch := range c.ClusterMembers {
		if ch > 127 {
			return fmt.Errorf("c.ClusterMembers[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func ClusterMembersInit(c *ClusterMembers) {
	return
}

func (*ClusterMembers) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*ClusterMembers) SbeTemplateId() (templateId uint16) {
	return 106
}

func (*ClusterMembers) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClusterMembers) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterMembers) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClusterMembers) MemberIdId() uint16 {
	return 1
}

func (*ClusterMembers) MemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembers) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.MemberIdSinceVersion()
}

func (*ClusterMembers) MemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembers) MemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembers) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembers) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembers) MemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembers) HighMemberIdId() uint16 {
	return 2
}

func (*ClusterMembers) HighMemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembers) HighMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.HighMemberIdSinceVersion()
}

func (*ClusterMembers) HighMemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembers) HighMemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembers) HighMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembers) HighMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembers) HighMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembers) ClusterMembersMetaAttribute(meta int) string {
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

func (*ClusterMembers) ClusterMembersSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembers) ClusterMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ClusterMembersSinceVersion()
}

func (*ClusterMembers) ClusterMembersDeprecated() uint16 {
	return 0
}

func (ClusterMembers) ClusterMembersCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembers) ClusterMembersHeaderLength() uint64 {
	return 4
}
