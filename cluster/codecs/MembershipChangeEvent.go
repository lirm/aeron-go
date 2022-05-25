// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type MembershipChangeEvent struct {
	LeadershipTermId int64
	LogPosition      int64
	Timestamp        int64
	LeaderMemberId   int32
	ClusterSize      int32
	ChangeType       ChangeTypeEnum
	MemberId         int32
	ClusterMembers   []uint8
}

func (m *MembershipChangeEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := m.RangeCheck(m.SbeSchemaVersion(), m.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, m.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, m.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, m.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, m.LeaderMemberId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, m.ClusterSize); err != nil {
		return err
	}
	if err := m.ChangeType.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, m.MemberId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(m.ClusterMembers))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, m.ClusterMembers); err != nil {
		return err
	}
	return nil
}

func (m *MembershipChangeEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !m.LeadershipTermIdInActingVersion(actingVersion) {
		m.LeadershipTermId = m.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &m.LeadershipTermId); err != nil {
			return err
		}
	}
	if !m.LogPositionInActingVersion(actingVersion) {
		m.LogPosition = m.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &m.LogPosition); err != nil {
			return err
		}
	}
	if !m.TimestampInActingVersion(actingVersion) {
		m.Timestamp = m.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &m.Timestamp); err != nil {
			return err
		}
	}
	if !m.LeaderMemberIdInActingVersion(actingVersion) {
		m.LeaderMemberId = m.LeaderMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &m.LeaderMemberId); err != nil {
			return err
		}
	}
	if !m.ClusterSizeInActingVersion(actingVersion) {
		m.ClusterSize = m.ClusterSizeNullValue()
	} else {
		if err := _m.ReadInt32(_r, &m.ClusterSize); err != nil {
			return err
		}
	}
	if m.ChangeTypeInActingVersion(actingVersion) {
		if err := m.ChangeType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !m.MemberIdInActingVersion(actingVersion) {
		m.MemberId = m.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &m.MemberId); err != nil {
			return err
		}
	}
	if actingVersion > m.SbeSchemaVersion() && blockLength > m.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-m.SbeBlockLength()))
	}

	if m.ClusterMembersInActingVersion(actingVersion) {
		var ClusterMembersLength uint32
		if err := _m.ReadUint32(_r, &ClusterMembersLength); err != nil {
			return err
		}
		if cap(m.ClusterMembers) < int(ClusterMembersLength) {
			m.ClusterMembers = make([]uint8, ClusterMembersLength)
		}
		m.ClusterMembers = m.ClusterMembers[:ClusterMembersLength]
		if err := _m.ReadBytes(_r, m.ClusterMembers); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := m.RangeCheck(actingVersion, m.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (m *MembershipChangeEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if m.LeadershipTermIdInActingVersion(actingVersion) {
		if m.LeadershipTermId < m.LeadershipTermIdMinValue() || m.LeadershipTermId > m.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on m.LeadershipTermId (%v < %v > %v)", m.LeadershipTermIdMinValue(), m.LeadershipTermId, m.LeadershipTermIdMaxValue())
		}
	}
	if m.LogPositionInActingVersion(actingVersion) {
		if m.LogPosition < m.LogPositionMinValue() || m.LogPosition > m.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on m.LogPosition (%v < %v > %v)", m.LogPositionMinValue(), m.LogPosition, m.LogPositionMaxValue())
		}
	}
	if m.TimestampInActingVersion(actingVersion) {
		if m.Timestamp < m.TimestampMinValue() || m.Timestamp > m.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on m.Timestamp (%v < %v > %v)", m.TimestampMinValue(), m.Timestamp, m.TimestampMaxValue())
		}
	}
	if m.LeaderMemberIdInActingVersion(actingVersion) {
		if m.LeaderMemberId < m.LeaderMemberIdMinValue() || m.LeaderMemberId > m.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on m.LeaderMemberId (%v < %v > %v)", m.LeaderMemberIdMinValue(), m.LeaderMemberId, m.LeaderMemberIdMaxValue())
		}
	}
	if m.ClusterSizeInActingVersion(actingVersion) {
		if m.ClusterSize < m.ClusterSizeMinValue() || m.ClusterSize > m.ClusterSizeMaxValue() {
			return fmt.Errorf("Range check failed on m.ClusterSize (%v < %v > %v)", m.ClusterSizeMinValue(), m.ClusterSize, m.ClusterSizeMaxValue())
		}
	}
	if err := m.ChangeType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if m.MemberIdInActingVersion(actingVersion) {
		if m.MemberId < m.MemberIdMinValue() || m.MemberId > m.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on m.MemberId (%v < %v > %v)", m.MemberIdMinValue(), m.MemberId, m.MemberIdMaxValue())
		}
	}
	for idx, ch := range m.ClusterMembers {
		if ch > 127 {
			return fmt.Errorf("m.ClusterMembers[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func MembershipChangeEventInit(m *MembershipChangeEvent) {
	return
}

func (*MembershipChangeEvent) SbeBlockLength() (blockLength uint16) {
	return 40
}

func (*MembershipChangeEvent) SbeTemplateId() (templateId uint16) {
	return 25
}

func (*MembershipChangeEvent) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*MembershipChangeEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*MembershipChangeEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*MembershipChangeEvent) LeadershipTermIdId() uint16 {
	return 1
}

func (*MembershipChangeEvent) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.LeadershipTermIdSinceVersion()
}

func (*MembershipChangeEvent) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*MembershipChangeEvent) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*MembershipChangeEvent) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*MembershipChangeEvent) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*MembershipChangeEvent) LogPositionId() uint16 {
	return 2
}

func (*MembershipChangeEvent) LogPositionSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.LogPositionSinceVersion()
}

func (*MembershipChangeEvent) LogPositionDeprecated() uint16 {
	return 0
}

func (*MembershipChangeEvent) LogPositionMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*MembershipChangeEvent) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*MembershipChangeEvent) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*MembershipChangeEvent) TimestampId() uint16 {
	return 3
}

func (*MembershipChangeEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.TimestampSinceVersion()
}

func (*MembershipChangeEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*MembershipChangeEvent) TimestampMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*MembershipChangeEvent) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*MembershipChangeEvent) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*MembershipChangeEvent) LeaderMemberIdId() uint16 {
	return 4
}

func (*MembershipChangeEvent) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.LeaderMemberIdSinceVersion()
}

func (*MembershipChangeEvent) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*MembershipChangeEvent) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*MembershipChangeEvent) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*MembershipChangeEvent) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*MembershipChangeEvent) ClusterSizeId() uint16 {
	return 5
}

func (*MembershipChangeEvent) ClusterSizeSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) ClusterSizeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.ClusterSizeSinceVersion()
}

func (*MembershipChangeEvent) ClusterSizeDeprecated() uint16 {
	return 0
}

func (*MembershipChangeEvent) ClusterSizeMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) ClusterSizeMinValue() int32 {
	return math.MinInt32 + 1
}

func (*MembershipChangeEvent) ClusterSizeMaxValue() int32 {
	return math.MaxInt32
}

func (*MembershipChangeEvent) ClusterSizeNullValue() int32 {
	return math.MinInt32
}

func (*MembershipChangeEvent) ChangeTypeId() uint16 {
	return 6
}

func (*MembershipChangeEvent) ChangeTypeSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) ChangeTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.ChangeTypeSinceVersion()
}

func (*MembershipChangeEvent) ChangeTypeDeprecated() uint16 {
	return 0
}

func (*MembershipChangeEvent) ChangeTypeMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) MemberIdId() uint16 {
	return 7
}

func (*MembershipChangeEvent) MemberIdSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.MemberIdSinceVersion()
}

func (*MembershipChangeEvent) MemberIdDeprecated() uint16 {
	return 0
}

func (*MembershipChangeEvent) MemberIdMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*MembershipChangeEvent) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*MembershipChangeEvent) MemberIdNullValue() int32 {
	return math.MinInt32
}

func (*MembershipChangeEvent) ClusterMembersMetaAttribute(meta int) string {
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

func (*MembershipChangeEvent) ClusterMembersSinceVersion() uint16 {
	return 0
}

func (m *MembershipChangeEvent) ClusterMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.ClusterMembersSinceVersion()
}

func (*MembershipChangeEvent) ClusterMembersDeprecated() uint16 {
	return 0
}

func (MembershipChangeEvent) ClusterMembersCharacterEncoding() string {
	return "US-ASCII"
}

func (MembershipChangeEvent) ClusterMembersHeaderLength() uint64 {
	return 4
}
