// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClusterMembersChange struct {
	CorrelationId  int64
	LeaderMemberId int32
	ActiveMembers  []uint8
	PassiveMembers []uint8
}

func (c *ClusterMembersChange) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.LeaderMemberId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ActiveMembers))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ActiveMembers); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.PassiveMembers))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.PassiveMembers); err != nil {
		return err
	}
	return nil
}

func (c *ClusterMembersChange) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
			return err
		}
	}
	if !c.LeaderMemberIdInActingVersion(actingVersion) {
		c.LeaderMemberId = c.LeaderMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.LeaderMemberId); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.ActiveMembersInActingVersion(actingVersion) {
		var ActiveMembersLength uint32
		if err := _m.ReadUint32(_r, &ActiveMembersLength); err != nil {
			return err
		}
		if cap(c.ActiveMembers) < int(ActiveMembersLength) {
			c.ActiveMembers = make([]uint8, ActiveMembersLength)
		}
		c.ActiveMembers = c.ActiveMembers[:ActiveMembersLength]
		if err := _m.ReadBytes(_r, c.ActiveMembers); err != nil {
			return err
		}
	}

	if c.PassiveMembersInActingVersion(actingVersion) {
		var PassiveMembersLength uint32
		if err := _m.ReadUint32(_r, &PassiveMembersLength); err != nil {
			return err
		}
		if cap(c.PassiveMembers) < int(PassiveMembersLength) {
			c.PassiveMembers = make([]uint8, PassiveMembersLength)
		}
		c.PassiveMembers = c.PassiveMembers[:PassiveMembersLength]
		if err := _m.ReadBytes(_r, c.PassiveMembers); err != nil {
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

func (c *ClusterMembersChange) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if c.LeaderMemberIdInActingVersion(actingVersion) {
		if c.LeaderMemberId < c.LeaderMemberIdMinValue() || c.LeaderMemberId > c.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.LeaderMemberId (%v < %v > %v)", c.LeaderMemberIdMinValue(), c.LeaderMemberId, c.LeaderMemberIdMaxValue())
		}
	}
	for idx, ch := range c.ActiveMembers {
		if ch > 127 {
			return fmt.Errorf("c.ActiveMembers[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.PassiveMembers {
		if ch > 127 {
			return fmt.Errorf("c.PassiveMembers[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func ClusterMembersChangeInit(c *ClusterMembersChange) {
	return
}

func (*ClusterMembersChange) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*ClusterMembersChange) SbeTemplateId() (templateId uint16) {
	return 71
}

func (*ClusterMembersChange) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClusterMembersChange) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterMembersChange) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClusterMembersChange) CorrelationIdId() uint16 {
	return 1
}

func (*ClusterMembersChange) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersChange) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ClusterMembersChange) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersChange) CorrelationIdMetaAttribute(meta int) string {
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

func (*ClusterMembersChange) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersChange) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersChange) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersChange) LeaderMemberIdId() uint16 {
	return 2
}

func (*ClusterMembersChange) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersChange) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeaderMemberIdSinceVersion()
}

func (*ClusterMembersChange) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersChange) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembersChange) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembersChange) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembersChange) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembersChange) ActiveMembersMetaAttribute(meta int) string {
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

func (*ClusterMembersChange) ActiveMembersSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersChange) ActiveMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ActiveMembersSinceVersion()
}

func (*ClusterMembersChange) ActiveMembersDeprecated() uint16 {
	return 0
}

func (ClusterMembersChange) ActiveMembersCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersChange) ActiveMembersHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersChange) PassiveMembersMetaAttribute(meta int) string {
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

func (*ClusterMembersChange) PassiveMembersSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersChange) PassiveMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.PassiveMembersSinceVersion()
}

func (*ClusterMembersChange) PassiveMembersDeprecated() uint16 {
	return 0
}

func (ClusterMembersChange) PassiveMembersCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersChange) PassiveMembersHeaderLength() uint64 {
	return 4
}
