// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type BackupResponse struct {
	CorrelationId           int64
	LogRecordingId          int64
	LogLeadershipTermId     int64
	LogTermBaseLogPosition  int64
	LastLeadershipTermId    int64
	LastTermBaseLogPosition int64
	CommitPositionCounterId int32
	LeaderMemberId          int32
	Snapshots               []BackupResponseSnapshots
	ClusterMembers          []uint8
}
type BackupResponseSnapshots struct {
	RecordingId         int64
	LeadershipTermId    int64
	TermBaseLogPosition int64
	LogPosition         int64
	Timestamp           int64
	ServiceId           int32
}

func (b *BackupResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := b.RangeCheck(b.SbeSchemaVersion(), b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, b.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.LogRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.LogLeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.LogTermBaseLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.LastLeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.LastTermBaseLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, b.CommitPositionCounterId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, b.LeaderMemberId); err != nil {
		return err
	}
	var SnapshotsBlockLength uint16 = 44
	if err := _m.WriteUint16(_w, SnapshotsBlockLength); err != nil {
		return err
	}
	var SnapshotsNumInGroup uint16 = uint16(len(b.Snapshots))
	if err := _m.WriteUint16(_w, SnapshotsNumInGroup); err != nil {
		return err
	}
	for _, prop := range b.Snapshots {
		if err := prop.Encode(_m, _w); err != nil {
			return err
		}
	}
	if err := _m.WriteUint32(_w, uint32(len(b.ClusterMembers))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, b.ClusterMembers); err != nil {
		return err
	}
	return nil
}

func (b *BackupResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !b.CorrelationIdInActingVersion(actingVersion) {
		b.CorrelationId = b.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.CorrelationId); err != nil {
			return err
		}
	}
	if !b.LogRecordingIdInActingVersion(actingVersion) {
		b.LogRecordingId = b.LogRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.LogRecordingId); err != nil {
			return err
		}
	}
	if !b.LogLeadershipTermIdInActingVersion(actingVersion) {
		b.LogLeadershipTermId = b.LogLeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.LogLeadershipTermId); err != nil {
			return err
		}
	}
	if !b.LogTermBaseLogPositionInActingVersion(actingVersion) {
		b.LogTermBaseLogPosition = b.LogTermBaseLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.LogTermBaseLogPosition); err != nil {
			return err
		}
	}
	if !b.LastLeadershipTermIdInActingVersion(actingVersion) {
		b.LastLeadershipTermId = b.LastLeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.LastLeadershipTermId); err != nil {
			return err
		}
	}
	if !b.LastTermBaseLogPositionInActingVersion(actingVersion) {
		b.LastTermBaseLogPosition = b.LastTermBaseLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.LastTermBaseLogPosition); err != nil {
			return err
		}
	}
	if !b.CommitPositionCounterIdInActingVersion(actingVersion) {
		b.CommitPositionCounterId = b.CommitPositionCounterIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &b.CommitPositionCounterId); err != nil {
			return err
		}
	}
	if !b.LeaderMemberIdInActingVersion(actingVersion) {
		b.LeaderMemberId = b.LeaderMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &b.LeaderMemberId); err != nil {
			return err
		}
	}
	if actingVersion > b.SbeSchemaVersion() && blockLength > b.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}

	if b.SnapshotsInActingVersion(actingVersion) {
		var SnapshotsBlockLength uint16
		if err := _m.ReadUint16(_r, &SnapshotsBlockLength); err != nil {
			return err
		}
		var SnapshotsNumInGroup uint16
		if err := _m.ReadUint16(_r, &SnapshotsNumInGroup); err != nil {
			return err
		}
		if cap(b.Snapshots) < int(SnapshotsNumInGroup) {
			b.Snapshots = make([]BackupResponseSnapshots, SnapshotsNumInGroup)
		}
		b.Snapshots = b.Snapshots[:SnapshotsNumInGroup]
		for i := range b.Snapshots {
			if err := b.Snapshots[i].Decode(_m, _r, actingVersion, uint(SnapshotsBlockLength)); err != nil {
				return err
			}
		}
	}

	if b.ClusterMembersInActingVersion(actingVersion) {
		var ClusterMembersLength uint32
		if err := _m.ReadUint32(_r, &ClusterMembersLength); err != nil {
			return err
		}
		if cap(b.ClusterMembers) < int(ClusterMembersLength) {
			b.ClusterMembers = make([]uint8, ClusterMembersLength)
		}
		b.ClusterMembers = b.ClusterMembers[:ClusterMembersLength]
		if err := _m.ReadBytes(_r, b.ClusterMembers); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := b.RangeCheck(actingVersion, b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (b *BackupResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if b.CorrelationIdInActingVersion(actingVersion) {
		if b.CorrelationId < b.CorrelationIdMinValue() || b.CorrelationId > b.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on b.CorrelationId (%v < %v > %v)", b.CorrelationIdMinValue(), b.CorrelationId, b.CorrelationIdMaxValue())
		}
	}
	if b.LogRecordingIdInActingVersion(actingVersion) {
		if b.LogRecordingId < b.LogRecordingIdMinValue() || b.LogRecordingId > b.LogRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on b.LogRecordingId (%v < %v > %v)", b.LogRecordingIdMinValue(), b.LogRecordingId, b.LogRecordingIdMaxValue())
		}
	}
	if b.LogLeadershipTermIdInActingVersion(actingVersion) {
		if b.LogLeadershipTermId < b.LogLeadershipTermIdMinValue() || b.LogLeadershipTermId > b.LogLeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on b.LogLeadershipTermId (%v < %v > %v)", b.LogLeadershipTermIdMinValue(), b.LogLeadershipTermId, b.LogLeadershipTermIdMaxValue())
		}
	}
	if b.LogTermBaseLogPositionInActingVersion(actingVersion) {
		if b.LogTermBaseLogPosition < b.LogTermBaseLogPositionMinValue() || b.LogTermBaseLogPosition > b.LogTermBaseLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on b.LogTermBaseLogPosition (%v < %v > %v)", b.LogTermBaseLogPositionMinValue(), b.LogTermBaseLogPosition, b.LogTermBaseLogPositionMaxValue())
		}
	}
	if b.LastLeadershipTermIdInActingVersion(actingVersion) {
		if b.LastLeadershipTermId < b.LastLeadershipTermIdMinValue() || b.LastLeadershipTermId > b.LastLeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on b.LastLeadershipTermId (%v < %v > %v)", b.LastLeadershipTermIdMinValue(), b.LastLeadershipTermId, b.LastLeadershipTermIdMaxValue())
		}
	}
	if b.LastTermBaseLogPositionInActingVersion(actingVersion) {
		if b.LastTermBaseLogPosition < b.LastTermBaseLogPositionMinValue() || b.LastTermBaseLogPosition > b.LastTermBaseLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on b.LastTermBaseLogPosition (%v < %v > %v)", b.LastTermBaseLogPositionMinValue(), b.LastTermBaseLogPosition, b.LastTermBaseLogPositionMaxValue())
		}
	}
	if b.CommitPositionCounterIdInActingVersion(actingVersion) {
		if b.CommitPositionCounterId < b.CommitPositionCounterIdMinValue() || b.CommitPositionCounterId > b.CommitPositionCounterIdMaxValue() {
			return fmt.Errorf("Range check failed on b.CommitPositionCounterId (%v < %v > %v)", b.CommitPositionCounterIdMinValue(), b.CommitPositionCounterId, b.CommitPositionCounterIdMaxValue())
		}
	}
	if b.LeaderMemberIdInActingVersion(actingVersion) {
		if b.LeaderMemberId < b.LeaderMemberIdMinValue() || b.LeaderMemberId > b.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on b.LeaderMemberId (%v < %v > %v)", b.LeaderMemberIdMinValue(), b.LeaderMemberId, b.LeaderMemberIdMaxValue())
		}
	}
	for _, prop := range b.Snapshots {
		if err := prop.RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	for idx, ch := range b.ClusterMembers {
		if ch > 127 {
			return fmt.Errorf("b.ClusterMembers[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func BackupResponseInit(b *BackupResponse) {
	return
}

func (b *BackupResponseSnapshots) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt64(_w, b.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.TermBaseLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, b.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, b.ServiceId); err != nil {
		return err
	}
	return nil
}

func (b *BackupResponseSnapshots) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if !b.RecordingIdInActingVersion(actingVersion) {
		b.RecordingId = b.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.RecordingId); err != nil {
			return err
		}
	}
	if !b.LeadershipTermIdInActingVersion(actingVersion) {
		b.LeadershipTermId = b.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.LeadershipTermId); err != nil {
			return err
		}
	}
	if !b.TermBaseLogPositionInActingVersion(actingVersion) {
		b.TermBaseLogPosition = b.TermBaseLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.TermBaseLogPosition); err != nil {
			return err
		}
	}
	if !b.LogPositionInActingVersion(actingVersion) {
		b.LogPosition = b.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.LogPosition); err != nil {
			return err
		}
	}
	if !b.TimestampInActingVersion(actingVersion) {
		b.Timestamp = b.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.Timestamp); err != nil {
			return err
		}
	}
	if !b.ServiceIdInActingVersion(actingVersion) {
		b.ServiceId = b.ServiceIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &b.ServiceId); err != nil {
			return err
		}
	}
	if actingVersion > b.SbeSchemaVersion() && blockLength > b.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}
	return nil
}

func (b *BackupResponseSnapshots) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if b.RecordingIdInActingVersion(actingVersion) {
		if b.RecordingId < b.RecordingIdMinValue() || b.RecordingId > b.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on b.RecordingId (%v < %v > %v)", b.RecordingIdMinValue(), b.RecordingId, b.RecordingIdMaxValue())
		}
	}
	if b.LeadershipTermIdInActingVersion(actingVersion) {
		if b.LeadershipTermId < b.LeadershipTermIdMinValue() || b.LeadershipTermId > b.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on b.LeadershipTermId (%v < %v > %v)", b.LeadershipTermIdMinValue(), b.LeadershipTermId, b.LeadershipTermIdMaxValue())
		}
	}
	if b.TermBaseLogPositionInActingVersion(actingVersion) {
		if b.TermBaseLogPosition < b.TermBaseLogPositionMinValue() || b.TermBaseLogPosition > b.TermBaseLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on b.TermBaseLogPosition (%v < %v > %v)", b.TermBaseLogPositionMinValue(), b.TermBaseLogPosition, b.TermBaseLogPositionMaxValue())
		}
	}
	if b.LogPositionInActingVersion(actingVersion) {
		if b.LogPosition < b.LogPositionMinValue() || b.LogPosition > b.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on b.LogPosition (%v < %v > %v)", b.LogPositionMinValue(), b.LogPosition, b.LogPositionMaxValue())
		}
	}
	if b.TimestampInActingVersion(actingVersion) {
		if b.Timestamp < b.TimestampMinValue() || b.Timestamp > b.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on b.Timestamp (%v < %v > %v)", b.TimestampMinValue(), b.Timestamp, b.TimestampMaxValue())
		}
	}
	if b.ServiceIdInActingVersion(actingVersion) {
		if b.ServiceId < b.ServiceIdMinValue() || b.ServiceId > b.ServiceIdMaxValue() {
			return fmt.Errorf("Range check failed on b.ServiceId (%v < %v > %v)", b.ServiceIdMinValue(), b.ServiceId, b.ServiceIdMaxValue())
		}
	}
	return nil
}

func BackupResponseSnapshotsInit(b *BackupResponseSnapshots) {
	return
}

func (*BackupResponse) SbeBlockLength() (blockLength uint16) {
	return 56
}

func (*BackupResponse) SbeTemplateId() (templateId uint16) {
	return 78
}

func (*BackupResponse) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*BackupResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*BackupResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*BackupResponse) CorrelationIdId() uint16 {
	return 1
}

func (*BackupResponse) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.CorrelationIdSinceVersion()
}

func (*BackupResponse) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*BackupResponse) CorrelationIdMetaAttribute(meta int) string {
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

func (*BackupResponse) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponse) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponse) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponse) LogRecordingIdId() uint16 {
	return 2
}

func (*BackupResponse) LogRecordingIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) LogRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LogRecordingIdSinceVersion()
}

func (*BackupResponse) LogRecordingIdDeprecated() uint16 {
	return 0
}

func (*BackupResponse) LogRecordingIdMetaAttribute(meta int) string {
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

func (*BackupResponse) LogRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponse) LogRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponse) LogRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponse) LogLeadershipTermIdId() uint16 {
	return 3
}

func (*BackupResponse) LogLeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) LogLeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LogLeadershipTermIdSinceVersion()
}

func (*BackupResponse) LogLeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*BackupResponse) LogLeadershipTermIdMetaAttribute(meta int) string {
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

func (*BackupResponse) LogLeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponse) LogLeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponse) LogLeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponse) LogTermBaseLogPositionId() uint16 {
	return 4
}

func (*BackupResponse) LogTermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) LogTermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LogTermBaseLogPositionSinceVersion()
}

func (*BackupResponse) LogTermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*BackupResponse) LogTermBaseLogPositionMetaAttribute(meta int) string {
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

func (*BackupResponse) LogTermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponse) LogTermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponse) LogTermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponse) LastLeadershipTermIdId() uint16 {
	return 5
}

func (*BackupResponse) LastLeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) LastLeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LastLeadershipTermIdSinceVersion()
}

func (*BackupResponse) LastLeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*BackupResponse) LastLeadershipTermIdMetaAttribute(meta int) string {
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

func (*BackupResponse) LastLeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponse) LastLeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponse) LastLeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponse) LastTermBaseLogPositionId() uint16 {
	return 6
}

func (*BackupResponse) LastTermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) LastTermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LastTermBaseLogPositionSinceVersion()
}

func (*BackupResponse) LastTermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*BackupResponse) LastTermBaseLogPositionMetaAttribute(meta int) string {
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

func (*BackupResponse) LastTermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponse) LastTermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponse) LastTermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponse) CommitPositionCounterIdId() uint16 {
	return 7
}

func (*BackupResponse) CommitPositionCounterIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) CommitPositionCounterIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.CommitPositionCounterIdSinceVersion()
}

func (*BackupResponse) CommitPositionCounterIdDeprecated() uint16 {
	return 0
}

func (*BackupResponse) CommitPositionCounterIdMetaAttribute(meta int) string {
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

func (*BackupResponse) CommitPositionCounterIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*BackupResponse) CommitPositionCounterIdMaxValue() int32 {
	return math.MaxInt32
}

func (*BackupResponse) CommitPositionCounterIdNullValue() int32 {
	return math.MinInt32
}

func (*BackupResponse) LeaderMemberIdId() uint16 {
	return 8
}

func (*BackupResponse) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LeaderMemberIdSinceVersion()
}

func (*BackupResponse) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*BackupResponse) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*BackupResponse) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*BackupResponse) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*BackupResponse) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*BackupResponseSnapshots) RecordingIdId() uint16 {
	return 10
}

func (*BackupResponseSnapshots) RecordingIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponseSnapshots) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.RecordingIdSinceVersion()
}

func (*BackupResponseSnapshots) RecordingIdDeprecated() uint16 {
	return 0
}

func (*BackupResponseSnapshots) RecordingIdMetaAttribute(meta int) string {
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

func (*BackupResponseSnapshots) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponseSnapshots) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponseSnapshots) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponseSnapshots) LeadershipTermIdId() uint16 {
	return 11
}

func (*BackupResponseSnapshots) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponseSnapshots) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LeadershipTermIdSinceVersion()
}

func (*BackupResponseSnapshots) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*BackupResponseSnapshots) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*BackupResponseSnapshots) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponseSnapshots) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponseSnapshots) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponseSnapshots) TermBaseLogPositionId() uint16 {
	return 12
}

func (*BackupResponseSnapshots) TermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (b *BackupResponseSnapshots) TermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.TermBaseLogPositionSinceVersion()
}

func (*BackupResponseSnapshots) TermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*BackupResponseSnapshots) TermBaseLogPositionMetaAttribute(meta int) string {
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

func (*BackupResponseSnapshots) TermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponseSnapshots) TermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponseSnapshots) TermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponseSnapshots) LogPositionId() uint16 {
	return 13
}

func (*BackupResponseSnapshots) LogPositionSinceVersion() uint16 {
	return 0
}

func (b *BackupResponseSnapshots) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.LogPositionSinceVersion()
}

func (*BackupResponseSnapshots) LogPositionDeprecated() uint16 {
	return 0
}

func (*BackupResponseSnapshots) LogPositionMetaAttribute(meta int) string {
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

func (*BackupResponseSnapshots) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponseSnapshots) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponseSnapshots) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponseSnapshots) TimestampId() uint16 {
	return 14
}

func (*BackupResponseSnapshots) TimestampSinceVersion() uint16 {
	return 0
}

func (b *BackupResponseSnapshots) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.TimestampSinceVersion()
}

func (*BackupResponseSnapshots) TimestampDeprecated() uint16 {
	return 0
}

func (*BackupResponseSnapshots) TimestampMetaAttribute(meta int) string {
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

func (*BackupResponseSnapshots) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupResponseSnapshots) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupResponseSnapshots) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*BackupResponseSnapshots) ServiceIdId() uint16 {
	return 15
}

func (*BackupResponseSnapshots) ServiceIdSinceVersion() uint16 {
	return 0
}

func (b *BackupResponseSnapshots) ServiceIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ServiceIdSinceVersion()
}

func (*BackupResponseSnapshots) ServiceIdDeprecated() uint16 {
	return 0
}

func (*BackupResponseSnapshots) ServiceIdMetaAttribute(meta int) string {
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

func (*BackupResponseSnapshots) ServiceIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*BackupResponseSnapshots) ServiceIdMaxValue() int32 {
	return math.MaxInt32
}

func (*BackupResponseSnapshots) ServiceIdNullValue() int32 {
	return math.MinInt32
}

func (*BackupResponse) SnapshotsId() uint16 {
	return 9
}

func (*BackupResponse) SnapshotsSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) SnapshotsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.SnapshotsSinceVersion()
}

func (*BackupResponse) SnapshotsDeprecated() uint16 {
	return 0
}

func (*BackupResponseSnapshots) SbeBlockLength() (blockLength uint) {
	return 44
}

func (*BackupResponseSnapshots) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*BackupResponse) ClusterMembersMetaAttribute(meta int) string {
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

func (*BackupResponse) ClusterMembersSinceVersion() uint16 {
	return 0
}

func (b *BackupResponse) ClusterMembersInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ClusterMembersSinceVersion()
}

func (*BackupResponse) ClusterMembersDeprecated() uint16 {
	return 0
}

func (BackupResponse) ClusterMembersCharacterEncoding() string {
	return "US-ASCII"
}

func (BackupResponse) ClusterMembersHeaderLength() uint64 {
	return 4
}
