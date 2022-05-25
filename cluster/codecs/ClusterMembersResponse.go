// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClusterMembersResponse struct {
	CorrelationId    int64
	LeaderMemberId   int32
	ActiveMembers    []uint8
	PassiveFollowers []uint8
}

func (c *ClusterMembersResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteUint32(_w, uint32(len(c.PassiveFollowers))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.PassiveFollowers); err != nil {
		return err
	}
	return nil
}

func (c *ClusterMembersResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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

	if c.PassiveFollowersInActingVersion(actingVersion) {
		var PassiveFollowersLength uint32
		if err := _m.ReadUint32(_r, &PassiveFollowersLength); err != nil {
			return err
		}
		if cap(c.PassiveFollowers) < int(PassiveFollowersLength) {
			c.PassiveFollowers = make([]uint8, PassiveFollowersLength)
		}
		c.PassiveFollowers = c.PassiveFollowers[:PassiveFollowersLength]
		if err := _m.ReadBytes(_r, c.PassiveFollowers); err != nil {
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

func (c *ClusterMembersResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	for idx, ch := range c.PassiveFollowers {
		if ch > 127 {
			return fmt.Errorf("c.PassiveFollowers[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func ClusterMembersResponseInit(c *ClusterMembersResponse) {
	return
}

func (*ClusterMembersResponse) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*ClusterMembersResponse) SbeTemplateId() (templateId uint16) {
	return 41
}

func (*ClusterMembersResponse) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClusterMembersResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterMembersResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClusterMembersResponse) CorrelationIdId() uint16 {
	return 1
}

func (*ClusterMembersResponse) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersResponse) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ClusterMembersResponse) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersResponse) CorrelationIdMetaAttribute(meta int) string {
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

func (*ClusterMembersResponse) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersResponse) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersResponse) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersResponse) LeaderMemberIdId() uint16 {
	return 2
}

func (*ClusterMembersResponse) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersResponse) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeaderMemberIdSinceVersion()
}

func (*ClusterMembersResponse) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersResponse) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembersResponse) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembersResponse) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembersResponse) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembersResponse) ActiveMembersMetaAttribute(meta int) string {
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

func (*ClusterMembersResponse) ActiveMembersSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersResponse) ActiveMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ActiveMembersSinceVersion()
}

func (*ClusterMembersResponse) ActiveMembersDeprecated() uint16 {
	return 0
}

func (ClusterMembersResponse) ActiveMembersCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersResponse) ActiveMembersHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersResponse) PassiveFollowersMetaAttribute(meta int) string {
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

func (*ClusterMembersResponse) PassiveFollowersSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersResponse) PassiveFollowersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.PassiveFollowersSinceVersion()
}

func (*ClusterMembersResponse) PassiveFollowersDeprecated() uint16 {
	return 0
}

func (ClusterMembersResponse) PassiveFollowersCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersResponse) PassiveFollowersHeaderLength() uint64 {
	return 4
}
