// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClusterMembersExtendedResponse struct {
	CorrelationId  int64
	CurrentTimeNs  int64
	LeaderMemberId int32
	MemberId       int32
	ActiveMembers  []ClusterMembersExtendedResponseActiveMembers
	PassiveMembers []ClusterMembersExtendedResponsePassiveMembers
}
type ClusterMembersExtendedResponseActiveMembers struct {
	LeadershipTermId   int64
	LogPosition        int64
	TimeOfLastAppendNs int64
	MemberId           int32
	IngressEndpoint    []uint8
	ConsensusEndpoint  []uint8
	LogEndpoint        []uint8
	CatchupEndpoint    []uint8
	ArchiveEndpoint    []uint8
}
type ClusterMembersExtendedResponsePassiveMembers struct {
	LeadershipTermId   int64
	LogPosition        int64
	TimeOfLastAppendNs int64
	MemberId           int32
	IngressEndpoint    []uint8
	ConsensusEndpoint  []uint8
	LogEndpoint        []uint8
	CatchupEndpoint    []uint8
	ArchiveEndpoint    []uint8
}

func (c *ClusterMembersExtendedResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.CurrentTimeNs); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.LeaderMemberId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.MemberId); err != nil {
		return err
	}
	var ActiveMembersBlockLength uint16 = 28
	if err := _m.WriteUint16(_w, ActiveMembersBlockLength); err != nil {
		return err
	}
	var ActiveMembersNumInGroup uint16 = uint16(len(c.ActiveMembers))
	if err := _m.WriteUint16(_w, ActiveMembersNumInGroup); err != nil {
		return err
	}
	for _, prop := range c.ActiveMembers {
		if err := prop.Encode(_m, _w); err != nil {
			return err
		}
	}
	var PassiveMembersBlockLength uint16 = 28
	if err := _m.WriteUint16(_w, PassiveMembersBlockLength); err != nil {
		return err
	}
	var PassiveMembersNumInGroup uint16 = uint16(len(c.PassiveMembers))
	if err := _m.WriteUint16(_w, PassiveMembersNumInGroup); err != nil {
		return err
	}
	for _, prop := range c.PassiveMembers {
		if err := prop.Encode(_m, _w); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterMembersExtendedResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
			return err
		}
	}
	if !c.CurrentTimeNsInActingVersion(actingVersion) {
		c.CurrentTimeNs = c.CurrentTimeNsNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CurrentTimeNs); err != nil {
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
	if !c.MemberIdInActingVersion(actingVersion) {
		c.MemberId = c.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.MemberId); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.ActiveMembersInActingVersion(actingVersion) {
		var ActiveMembersBlockLength uint16
		if err := _m.ReadUint16(_r, &ActiveMembersBlockLength); err != nil {
			return err
		}
		var ActiveMembersNumInGroup uint16
		if err := _m.ReadUint16(_r, &ActiveMembersNumInGroup); err != nil {
			return err
		}
		if cap(c.ActiveMembers) < int(ActiveMembersNumInGroup) {
			c.ActiveMembers = make([]ClusterMembersExtendedResponseActiveMembers, ActiveMembersNumInGroup)
		}
		c.ActiveMembers = c.ActiveMembers[:ActiveMembersNumInGroup]
		for i := range c.ActiveMembers {
			if err := c.ActiveMembers[i].Decode(_m, _r, actingVersion, uint(ActiveMembersBlockLength)); err != nil {
				return err
			}
		}
	}

	if c.PassiveMembersInActingVersion(actingVersion) {
		var PassiveMembersBlockLength uint16
		if err := _m.ReadUint16(_r, &PassiveMembersBlockLength); err != nil {
			return err
		}
		var PassiveMembersNumInGroup uint16
		if err := _m.ReadUint16(_r, &PassiveMembersNumInGroup); err != nil {
			return err
		}
		if cap(c.PassiveMembers) < int(PassiveMembersNumInGroup) {
			c.PassiveMembers = make([]ClusterMembersExtendedResponsePassiveMembers, PassiveMembersNumInGroup)
		}
		c.PassiveMembers = c.PassiveMembers[:PassiveMembersNumInGroup]
		for i := range c.PassiveMembers {
			if err := c.PassiveMembers[i].Decode(_m, _r, actingVersion, uint(PassiveMembersBlockLength)); err != nil {
				return err
			}
		}
	}
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterMembersExtendedResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if c.CurrentTimeNsInActingVersion(actingVersion) {
		if c.CurrentTimeNs < c.CurrentTimeNsMinValue() || c.CurrentTimeNs > c.CurrentTimeNsMaxValue() {
			return fmt.Errorf("Range check failed on c.CurrentTimeNs (%v < %v > %v)", c.CurrentTimeNsMinValue(), c.CurrentTimeNs, c.CurrentTimeNsMaxValue())
		}
	}
	if c.LeaderMemberIdInActingVersion(actingVersion) {
		if c.LeaderMemberId < c.LeaderMemberIdMinValue() || c.LeaderMemberId > c.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.LeaderMemberId (%v < %v > %v)", c.LeaderMemberIdMinValue(), c.LeaderMemberId, c.LeaderMemberIdMaxValue())
		}
	}
	if c.MemberIdInActingVersion(actingVersion) {
		if c.MemberId < c.MemberIdMinValue() || c.MemberId > c.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.MemberId (%v < %v > %v)", c.MemberIdMinValue(), c.MemberId, c.MemberIdMaxValue())
		}
	}
	for _, prop := range c.ActiveMembers {
		if err := prop.RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	for _, prop := range c.PassiveMembers {
		if err := prop.RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	return nil
}

func ClusterMembersExtendedResponseInit(c *ClusterMembersExtendedResponse) {
	return
}

func (c *ClusterMembersExtendedResponseActiveMembers) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt64(_w, c.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.TimeOfLastAppendNs); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.MemberId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.IngressEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.IngressEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ConsensusEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ConsensusEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.LogEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.LogEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.CatchupEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.CatchupEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ArchiveEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ArchiveEndpoint); err != nil {
		return err
	}
	return nil
}

func (c *ClusterMembersExtendedResponseActiveMembers) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if !c.LeadershipTermIdInActingVersion(actingVersion) {
		c.LeadershipTermId = c.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.LeadershipTermId); err != nil {
			return err
		}
	}
	if !c.LogPositionInActingVersion(actingVersion) {
		c.LogPosition = c.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.LogPosition); err != nil {
			return err
		}
	}
	if !c.TimeOfLastAppendNsInActingVersion(actingVersion) {
		c.TimeOfLastAppendNs = c.TimeOfLastAppendNsNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.TimeOfLastAppendNs); err != nil {
			return err
		}
	}
	if !c.MemberIdInActingVersion(actingVersion) {
		c.MemberId = c.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.MemberId); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.IngressEndpointInActingVersion(actingVersion) {
		var IngressEndpointLength uint32
		if err := _m.ReadUint32(_r, &IngressEndpointLength); err != nil {
			return err
		}
		if cap(c.IngressEndpoint) < int(IngressEndpointLength) {
			c.IngressEndpoint = make([]uint8, IngressEndpointLength)
		}
		c.IngressEndpoint = c.IngressEndpoint[:IngressEndpointLength]
		if err := _m.ReadBytes(_r, c.IngressEndpoint); err != nil {
			return err
		}
	}

	if c.ConsensusEndpointInActingVersion(actingVersion) {
		var ConsensusEndpointLength uint32
		if err := _m.ReadUint32(_r, &ConsensusEndpointLength); err != nil {
			return err
		}
		if cap(c.ConsensusEndpoint) < int(ConsensusEndpointLength) {
			c.ConsensusEndpoint = make([]uint8, ConsensusEndpointLength)
		}
		c.ConsensusEndpoint = c.ConsensusEndpoint[:ConsensusEndpointLength]
		if err := _m.ReadBytes(_r, c.ConsensusEndpoint); err != nil {
			return err
		}
	}

	if c.LogEndpointInActingVersion(actingVersion) {
		var LogEndpointLength uint32
		if err := _m.ReadUint32(_r, &LogEndpointLength); err != nil {
			return err
		}
		if cap(c.LogEndpoint) < int(LogEndpointLength) {
			c.LogEndpoint = make([]uint8, LogEndpointLength)
		}
		c.LogEndpoint = c.LogEndpoint[:LogEndpointLength]
		if err := _m.ReadBytes(_r, c.LogEndpoint); err != nil {
			return err
		}
	}

	if c.CatchupEndpointInActingVersion(actingVersion) {
		var CatchupEndpointLength uint32
		if err := _m.ReadUint32(_r, &CatchupEndpointLength); err != nil {
			return err
		}
		if cap(c.CatchupEndpoint) < int(CatchupEndpointLength) {
			c.CatchupEndpoint = make([]uint8, CatchupEndpointLength)
		}
		c.CatchupEndpoint = c.CatchupEndpoint[:CatchupEndpointLength]
		if err := _m.ReadBytes(_r, c.CatchupEndpoint); err != nil {
			return err
		}
	}

	if c.ArchiveEndpointInActingVersion(actingVersion) {
		var ArchiveEndpointLength uint32
		if err := _m.ReadUint32(_r, &ArchiveEndpointLength); err != nil {
			return err
		}
		if cap(c.ArchiveEndpoint) < int(ArchiveEndpointLength) {
			c.ArchiveEndpoint = make([]uint8, ArchiveEndpointLength)
		}
		c.ArchiveEndpoint = c.ArchiveEndpoint[:ArchiveEndpointLength]
		if err := _m.ReadBytes(_r, c.ArchiveEndpoint); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterMembersExtendedResponseActiveMembers) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.LeadershipTermIdInActingVersion(actingVersion) {
		if c.LeadershipTermId < c.LeadershipTermIdMinValue() || c.LeadershipTermId > c.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on c.LeadershipTermId (%v < %v > %v)", c.LeadershipTermIdMinValue(), c.LeadershipTermId, c.LeadershipTermIdMaxValue())
		}
	}
	if c.LogPositionInActingVersion(actingVersion) {
		if c.LogPosition < c.LogPositionMinValue() || c.LogPosition > c.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on c.LogPosition (%v < %v > %v)", c.LogPositionMinValue(), c.LogPosition, c.LogPositionMaxValue())
		}
	}
	if c.TimeOfLastAppendNsInActingVersion(actingVersion) {
		if c.TimeOfLastAppendNs < c.TimeOfLastAppendNsMinValue() || c.TimeOfLastAppendNs > c.TimeOfLastAppendNsMaxValue() {
			return fmt.Errorf("Range check failed on c.TimeOfLastAppendNs (%v < %v > %v)", c.TimeOfLastAppendNsMinValue(), c.TimeOfLastAppendNs, c.TimeOfLastAppendNsMaxValue())
		}
	}
	if c.MemberIdInActingVersion(actingVersion) {
		if c.MemberId < c.MemberIdMinValue() || c.MemberId > c.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.MemberId (%v < %v > %v)", c.MemberIdMinValue(), c.MemberId, c.MemberIdMaxValue())
		}
	}
	for idx, ch := range c.IngressEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.IngressEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.ConsensusEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.ConsensusEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.LogEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.LogEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.CatchupEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.CatchupEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.ArchiveEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.ArchiveEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func ClusterMembersExtendedResponseActiveMembersInit(c *ClusterMembersExtendedResponseActiveMembers) {
	return
}

func (c *ClusterMembersExtendedResponsePassiveMembers) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt64(_w, c.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.TimeOfLastAppendNs); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.MemberId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.IngressEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.IngressEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ConsensusEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ConsensusEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.LogEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.LogEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.CatchupEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.CatchupEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ArchiveEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ArchiveEndpoint); err != nil {
		return err
	}
	return nil
}

func (c *ClusterMembersExtendedResponsePassiveMembers) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if !c.LeadershipTermIdInActingVersion(actingVersion) {
		c.LeadershipTermId = c.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.LeadershipTermId); err != nil {
			return err
		}
	}
	if !c.LogPositionInActingVersion(actingVersion) {
		c.LogPosition = c.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.LogPosition); err != nil {
			return err
		}
	}
	if !c.TimeOfLastAppendNsInActingVersion(actingVersion) {
		c.TimeOfLastAppendNs = c.TimeOfLastAppendNsNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.TimeOfLastAppendNs); err != nil {
			return err
		}
	}
	if !c.MemberIdInActingVersion(actingVersion) {
		c.MemberId = c.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.MemberId); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.IngressEndpointInActingVersion(actingVersion) {
		var IngressEndpointLength uint32
		if err := _m.ReadUint32(_r, &IngressEndpointLength); err != nil {
			return err
		}
		if cap(c.IngressEndpoint) < int(IngressEndpointLength) {
			c.IngressEndpoint = make([]uint8, IngressEndpointLength)
		}
		c.IngressEndpoint = c.IngressEndpoint[:IngressEndpointLength]
		if err := _m.ReadBytes(_r, c.IngressEndpoint); err != nil {
			return err
		}
	}

	if c.ConsensusEndpointInActingVersion(actingVersion) {
		var ConsensusEndpointLength uint32
		if err := _m.ReadUint32(_r, &ConsensusEndpointLength); err != nil {
			return err
		}
		if cap(c.ConsensusEndpoint) < int(ConsensusEndpointLength) {
			c.ConsensusEndpoint = make([]uint8, ConsensusEndpointLength)
		}
		c.ConsensusEndpoint = c.ConsensusEndpoint[:ConsensusEndpointLength]
		if err := _m.ReadBytes(_r, c.ConsensusEndpoint); err != nil {
			return err
		}
	}

	if c.LogEndpointInActingVersion(actingVersion) {
		var LogEndpointLength uint32
		if err := _m.ReadUint32(_r, &LogEndpointLength); err != nil {
			return err
		}
		if cap(c.LogEndpoint) < int(LogEndpointLength) {
			c.LogEndpoint = make([]uint8, LogEndpointLength)
		}
		c.LogEndpoint = c.LogEndpoint[:LogEndpointLength]
		if err := _m.ReadBytes(_r, c.LogEndpoint); err != nil {
			return err
		}
	}

	if c.CatchupEndpointInActingVersion(actingVersion) {
		var CatchupEndpointLength uint32
		if err := _m.ReadUint32(_r, &CatchupEndpointLength); err != nil {
			return err
		}
		if cap(c.CatchupEndpoint) < int(CatchupEndpointLength) {
			c.CatchupEndpoint = make([]uint8, CatchupEndpointLength)
		}
		c.CatchupEndpoint = c.CatchupEndpoint[:CatchupEndpointLength]
		if err := _m.ReadBytes(_r, c.CatchupEndpoint); err != nil {
			return err
		}
	}

	if c.ArchiveEndpointInActingVersion(actingVersion) {
		var ArchiveEndpointLength uint32
		if err := _m.ReadUint32(_r, &ArchiveEndpointLength); err != nil {
			return err
		}
		if cap(c.ArchiveEndpoint) < int(ArchiveEndpointLength) {
			c.ArchiveEndpoint = make([]uint8, ArchiveEndpointLength)
		}
		c.ArchiveEndpoint = c.ArchiveEndpoint[:ArchiveEndpointLength]
		if err := _m.ReadBytes(_r, c.ArchiveEndpoint); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterMembersExtendedResponsePassiveMembers) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.LeadershipTermIdInActingVersion(actingVersion) {
		if c.LeadershipTermId < c.LeadershipTermIdMinValue() || c.LeadershipTermId > c.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on c.LeadershipTermId (%v < %v > %v)", c.LeadershipTermIdMinValue(), c.LeadershipTermId, c.LeadershipTermIdMaxValue())
		}
	}
	if c.LogPositionInActingVersion(actingVersion) {
		if c.LogPosition < c.LogPositionMinValue() || c.LogPosition > c.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on c.LogPosition (%v < %v > %v)", c.LogPositionMinValue(), c.LogPosition, c.LogPositionMaxValue())
		}
	}
	if c.TimeOfLastAppendNsInActingVersion(actingVersion) {
		if c.TimeOfLastAppendNs < c.TimeOfLastAppendNsMinValue() || c.TimeOfLastAppendNs > c.TimeOfLastAppendNsMaxValue() {
			return fmt.Errorf("Range check failed on c.TimeOfLastAppendNs (%v < %v > %v)", c.TimeOfLastAppendNsMinValue(), c.TimeOfLastAppendNs, c.TimeOfLastAppendNsMaxValue())
		}
	}
	if c.MemberIdInActingVersion(actingVersion) {
		if c.MemberId < c.MemberIdMinValue() || c.MemberId > c.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.MemberId (%v < %v > %v)", c.MemberIdMinValue(), c.MemberId, c.MemberIdMaxValue())
		}
	}
	for idx, ch := range c.IngressEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.IngressEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.ConsensusEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.ConsensusEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.LogEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.LogEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.CatchupEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.CatchupEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	for idx, ch := range c.ArchiveEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.ArchiveEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func ClusterMembersExtendedResponsePassiveMembersInit(c *ClusterMembersExtendedResponsePassiveMembers) {
	return
}

func (*ClusterMembersExtendedResponse) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*ClusterMembersExtendedResponse) SbeTemplateId() (templateId uint16) {
	return 43
}

func (*ClusterMembersExtendedResponse) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClusterMembersExtendedResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterMembersExtendedResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClusterMembersExtendedResponse) CorrelationIdId() uint16 {
	return 1
}

func (*ClusterMembersExtendedResponse) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponse) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ClusterMembersExtendedResponse) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponse) CorrelationIdMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponse) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponse) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponse) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponse) CurrentTimeNsId() uint16 {
	return 2
}

func (*ClusterMembersExtendedResponse) CurrentTimeNsSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponse) CurrentTimeNsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CurrentTimeNsSinceVersion()
}

func (*ClusterMembersExtendedResponse) CurrentTimeNsDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponse) CurrentTimeNsMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponse) CurrentTimeNsMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponse) CurrentTimeNsMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponse) CurrentTimeNsNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponse) LeaderMemberIdId() uint16 {
	return 3
}

func (*ClusterMembersExtendedResponse) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponse) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeaderMemberIdSinceVersion()
}

func (*ClusterMembersExtendedResponse) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponse) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponse) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembersExtendedResponse) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembersExtendedResponse) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembersExtendedResponse) MemberIdId() uint16 {
	return 4
}

func (*ClusterMembersExtendedResponse) MemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponse) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.MemberIdSinceVersion()
}

func (*ClusterMembersExtendedResponse) MemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponse) MemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponse) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembersExtendedResponse) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembersExtendedResponse) MemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdId() uint16 {
	return 6
}

func (*ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeadershipTermIdSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponseActiveMembers) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponseActiveMembers) LogPositionId() uint16 {
	return 7
}

func (*ClusterMembersExtendedResponseActiveMembers) LogPositionSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogPositionSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) LogPositionDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponseActiveMembers) LogPositionMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponseActiveMembers) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponseActiveMembers) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsId() uint16 {
	return 8
}

func (*ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.TimeOfLastAppendNsSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponseActiveMembers) TimeOfLastAppendNsNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponseActiveMembers) MemberIdId() uint16 {
	return 9
}

func (*ClusterMembersExtendedResponseActiveMembers) MemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.MemberIdSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) MemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponseActiveMembers) MemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembersExtendedResponseActiveMembers) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembersExtendedResponseActiveMembers) MemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembersExtendedResponseActiveMembers) IngressEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) IngressEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) IngressEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.IngressEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) IngressEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponseActiveMembers) IngressEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponseActiveMembers) IngressEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponseActiveMembers) ConsensusEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) ConsensusEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) ConsensusEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ConsensusEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) ConsensusEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponseActiveMembers) ConsensusEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponseActiveMembers) ConsensusEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponseActiveMembers) LogEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) LogEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) LogEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) LogEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponseActiveMembers) LogEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponseActiveMembers) LogEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponseActiveMembers) CatchupEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) CatchupEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) CatchupEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CatchupEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) CatchupEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponseActiveMembers) CatchupEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponseActiveMembers) CatchupEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponseActiveMembers) ArchiveEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponseActiveMembers) ArchiveEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponseActiveMembers) ArchiveEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ArchiveEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponseActiveMembers) ArchiveEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponseActiveMembers) ArchiveEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponseActiveMembers) ArchiveEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdId() uint16 {
	return 16
}

func (*ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeadershipTermIdSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponsePassiveMembers) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogPositionId() uint16 {
	return 17
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogPositionSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogPositionSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogPositionDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogPositionMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsId() uint16 {
	return 18
}

func (*ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.TimeOfLastAppendNsSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersExtendedResponsePassiveMembers) TimeOfLastAppendNsNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersExtendedResponsePassiveMembers) MemberIdId() uint16 {
	return 19
}

func (*ClusterMembersExtendedResponsePassiveMembers) MemberIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.MemberIdSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) MemberIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponsePassiveMembers) MemberIdMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterMembersExtendedResponsePassiveMembers) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterMembersExtendedResponsePassiveMembers) MemberIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterMembersExtendedResponsePassiveMembers) IngressEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) IngressEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) IngressEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.IngressEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) IngressEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponsePassiveMembers) IngressEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponsePassiveMembers) IngressEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponsePassiveMembers) ConsensusEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) ConsensusEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) ConsensusEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ConsensusEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) ConsensusEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponsePassiveMembers) ConsensusEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponsePassiveMembers) ConsensusEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) LogEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) LogEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) LogEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponsePassiveMembers) LogEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponsePassiveMembers) LogEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponsePassiveMembers) CatchupEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) CatchupEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) CatchupEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CatchupEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) CatchupEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponsePassiveMembers) CatchupEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponsePassiveMembers) CatchupEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponsePassiveMembers) ArchiveEndpointMetaAttribute(meta int) string {
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

func (*ClusterMembersExtendedResponsePassiveMembers) ArchiveEndpointSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponsePassiveMembers) ArchiveEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ArchiveEndpointSinceVersion()
}

func (*ClusterMembersExtendedResponsePassiveMembers) ArchiveEndpointDeprecated() uint16 {
	return 0
}

func (ClusterMembersExtendedResponsePassiveMembers) ArchiveEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterMembersExtendedResponsePassiveMembers) ArchiveEndpointHeaderLength() uint64 {
	return 4
}

func (*ClusterMembersExtendedResponse) ActiveMembersId() uint16 {
	return 5
}

func (*ClusterMembersExtendedResponse) ActiveMembersSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponse) ActiveMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ActiveMembersSinceVersion()
}

func (*ClusterMembersExtendedResponse) ActiveMembersDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponseActiveMembers) SbeBlockLength() (blockLength uint) {
	return 28
}

func (*ClusterMembersExtendedResponseActiveMembers) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterMembersExtendedResponse) PassiveMembersId() uint16 {
	return 15
}

func (*ClusterMembersExtendedResponse) PassiveMembersSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersExtendedResponse) PassiveMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.PassiveMembersSinceVersion()
}

func (*ClusterMembersExtendedResponse) PassiveMembersDeprecated() uint16 {
	return 0
}

func (*ClusterMembersExtendedResponsePassiveMembers) SbeBlockLength() (blockLength uint) {
	return 28
}

func (*ClusterMembersExtendedResponsePassiveMembers) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}
