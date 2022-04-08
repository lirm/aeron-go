// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SnapshotRecordings struct {
	CorrelationId   int64
	Snapshots       []SnapshotRecordingsSnapshots
	MemberEndpoints []uint8
}
type SnapshotRecordingsSnapshots struct {
	RecordingId         int64
	LeadershipTermId    int64
	TermBaseLogPosition int64
	LogPosition         int64
	Timestamp           int64
	ServiceId           int32
}

func (s *SnapshotRecordings) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	var SnapshotsBlockLength uint16 = 44
	if err := _m.WriteUint16(_w, SnapshotsBlockLength); err != nil {
		return err
	}
	var SnapshotsNumInGroup uint16 = uint16(len(s.Snapshots))
	if err := _m.WriteUint16(_w, SnapshotsNumInGroup); err != nil {
		return err
	}
	for _, prop := range s.Snapshots {
		if err := prop.Encode(_m, _w); err != nil {
			return err
		}
	}
	if err := _m.WriteUint32(_w, uint32(len(s.MemberEndpoints))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.MemberEndpoints); err != nil {
		return err
	}
	return nil
}

func (s *SnapshotRecordings) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	if s.SnapshotsInActingVersion(actingVersion) {
		var SnapshotsBlockLength uint16
		if err := _m.ReadUint16(_r, &SnapshotsBlockLength); err != nil {
			return err
		}
		var SnapshotsNumInGroup uint16
		if err := _m.ReadUint16(_r, &SnapshotsNumInGroup); err != nil {
			return err
		}
		if cap(s.Snapshots) < int(SnapshotsNumInGroup) {
			s.Snapshots = make([]SnapshotRecordingsSnapshots, SnapshotsNumInGroup)
		}
		s.Snapshots = s.Snapshots[:SnapshotsNumInGroup]
		for i := range s.Snapshots {
			if err := s.Snapshots[i].Decode(_m, _r, actingVersion, uint(SnapshotsBlockLength)); err != nil {
				return err
			}
		}
	}

	if s.MemberEndpointsInActingVersion(actingVersion) {
		var MemberEndpointsLength uint32
		if err := _m.ReadUint32(_r, &MemberEndpointsLength); err != nil {
			return err
		}
		if cap(s.MemberEndpoints) < int(MemberEndpointsLength) {
			s.MemberEndpoints = make([]uint8, MemberEndpointsLength)
		}
		s.MemberEndpoints = s.MemberEndpoints[:MemberEndpointsLength]
		if err := _m.ReadBytes(_r, s.MemberEndpoints); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapshotRecordings) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	for _, prop := range s.Snapshots {
		if err := prop.RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	for idx, ch := range s.MemberEndpoints {
		if ch > 127 {
			return fmt.Errorf("s.MemberEndpoints[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func SnapshotRecordingsInit(s *SnapshotRecordings) {
	return
}

func (s *SnapshotRecordingsSnapshots) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt64(_w, s.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.TermBaseLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.ServiceId); err != nil {
		return err
	}
	return nil
}

func (s *SnapshotRecordingsSnapshots) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if !s.RecordingIdInActingVersion(actingVersion) {
		s.RecordingId = s.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.RecordingId); err != nil {
			return err
		}
	}
	if !s.LeadershipTermIdInActingVersion(actingVersion) {
		s.LeadershipTermId = s.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LeadershipTermId); err != nil {
			return err
		}
	}
	if !s.TermBaseLogPositionInActingVersion(actingVersion) {
		s.TermBaseLogPosition = s.TermBaseLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.TermBaseLogPosition); err != nil {
			return err
		}
	}
	if !s.LogPositionInActingVersion(actingVersion) {
		s.LogPosition = s.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LogPosition); err != nil {
			return err
		}
	}
	if !s.TimestampInActingVersion(actingVersion) {
		s.Timestamp = s.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.Timestamp); err != nil {
			return err
		}
	}
	if !s.ServiceIdInActingVersion(actingVersion) {
		s.ServiceId = s.ServiceIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.ServiceId); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	return nil
}

func (s *SnapshotRecordingsSnapshots) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.RecordingIdInActingVersion(actingVersion) {
		if s.RecordingId < s.RecordingIdMinValue() || s.RecordingId > s.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on s.RecordingId (%v < %v > %v)", s.RecordingIdMinValue(), s.RecordingId, s.RecordingIdMaxValue())
		}
	}
	if s.LeadershipTermIdInActingVersion(actingVersion) {
		if s.LeadershipTermId < s.LeadershipTermIdMinValue() || s.LeadershipTermId > s.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeadershipTermId (%v < %v > %v)", s.LeadershipTermIdMinValue(), s.LeadershipTermId, s.LeadershipTermIdMaxValue())
		}
	}
	if s.TermBaseLogPositionInActingVersion(actingVersion) {
		if s.TermBaseLogPosition < s.TermBaseLogPositionMinValue() || s.TermBaseLogPosition > s.TermBaseLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on s.TermBaseLogPosition (%v < %v > %v)", s.TermBaseLogPositionMinValue(), s.TermBaseLogPosition, s.TermBaseLogPositionMaxValue())
		}
	}
	if s.LogPositionInActingVersion(actingVersion) {
		if s.LogPosition < s.LogPositionMinValue() || s.LogPosition > s.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on s.LogPosition (%v < %v > %v)", s.LogPositionMinValue(), s.LogPosition, s.LogPositionMaxValue())
		}
	}
	if s.TimestampInActingVersion(actingVersion) {
		if s.Timestamp < s.TimestampMinValue() || s.Timestamp > s.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on s.Timestamp (%v < %v > %v)", s.TimestampMinValue(), s.Timestamp, s.TimestampMaxValue())
		}
	}
	if s.ServiceIdInActingVersion(actingVersion) {
		if s.ServiceId < s.ServiceIdMinValue() || s.ServiceId > s.ServiceIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ServiceId (%v < %v > %v)", s.ServiceIdMinValue(), s.ServiceId, s.ServiceIdMaxValue())
		}
	}
	return nil
}

func SnapshotRecordingsSnapshotsInit(s *SnapshotRecordingsSnapshots) {
	return
}

func (*SnapshotRecordings) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*SnapshotRecordings) SbeTemplateId() (templateId uint16) {
	return 73
}

func (*SnapshotRecordings) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SnapshotRecordings) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SnapshotRecordings) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SnapshotRecordings) CorrelationIdId() uint16 {
	return 1
}

func (*SnapshotRecordings) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordings) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*SnapshotRecordings) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordings) CorrelationIdMetaAttribute(meta int) string {
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

func (*SnapshotRecordings) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotRecordings) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotRecordings) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotRecordingsSnapshots) RecordingIdId() uint16 {
	return 4
}

func (*SnapshotRecordingsSnapshots) RecordingIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingsSnapshots) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.RecordingIdSinceVersion()
}

func (*SnapshotRecordingsSnapshots) RecordingIdDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingsSnapshots) RecordingIdMetaAttribute(meta int) string {
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

func (*SnapshotRecordingsSnapshots) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotRecordingsSnapshots) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotRecordingsSnapshots) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotRecordingsSnapshots) LeadershipTermIdId() uint16 {
	return 5
}

func (*SnapshotRecordingsSnapshots) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingsSnapshots) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*SnapshotRecordingsSnapshots) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingsSnapshots) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*SnapshotRecordingsSnapshots) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotRecordingsSnapshots) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotRecordingsSnapshots) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotRecordingsSnapshots) TermBaseLogPositionId() uint16 {
	return 6
}

func (*SnapshotRecordingsSnapshots) TermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingsSnapshots) TermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TermBaseLogPositionSinceVersion()
}

func (*SnapshotRecordingsSnapshots) TermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingsSnapshots) TermBaseLogPositionMetaAttribute(meta int) string {
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

func (*SnapshotRecordingsSnapshots) TermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotRecordingsSnapshots) TermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotRecordingsSnapshots) TermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotRecordingsSnapshots) LogPositionId() uint16 {
	return 7
}

func (*SnapshotRecordingsSnapshots) LogPositionSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingsSnapshots) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LogPositionSinceVersion()
}

func (*SnapshotRecordingsSnapshots) LogPositionDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingsSnapshots) LogPositionMetaAttribute(meta int) string {
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

func (*SnapshotRecordingsSnapshots) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotRecordingsSnapshots) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotRecordingsSnapshots) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotRecordingsSnapshots) TimestampId() uint16 {
	return 8
}

func (*SnapshotRecordingsSnapshots) TimestampSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingsSnapshots) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TimestampSinceVersion()
}

func (*SnapshotRecordingsSnapshots) TimestampDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingsSnapshots) TimestampMetaAttribute(meta int) string {
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

func (*SnapshotRecordingsSnapshots) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotRecordingsSnapshots) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotRecordingsSnapshots) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotRecordingsSnapshots) ServiceIdId() uint16 {
	return 9
}

func (*SnapshotRecordingsSnapshots) ServiceIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordingsSnapshots) ServiceIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ServiceIdSinceVersion()
}

func (*SnapshotRecordingsSnapshots) ServiceIdDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingsSnapshots) ServiceIdMetaAttribute(meta int) string {
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

func (*SnapshotRecordingsSnapshots) ServiceIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*SnapshotRecordingsSnapshots) ServiceIdMaxValue() int32 {
	return math.MaxInt32
}

func (*SnapshotRecordingsSnapshots) ServiceIdNullValue() int32 {
	return math.MinInt32
}

func (*SnapshotRecordings) SnapshotsId() uint16 {
	return 3
}

func (*SnapshotRecordings) SnapshotsSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordings) SnapshotsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.SnapshotsSinceVersion()
}

func (*SnapshotRecordings) SnapshotsDeprecated() uint16 {
	return 0
}

func (*SnapshotRecordingsSnapshots) SbeBlockLength() (blockLength uint) {
	return 44
}

func (*SnapshotRecordingsSnapshots) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SnapshotRecordings) MemberEndpointsMetaAttribute(meta int) string {
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

func (*SnapshotRecordings) MemberEndpointsSinceVersion() uint16 {
	return 0
}

func (s *SnapshotRecordings) MemberEndpointsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.MemberEndpointsSinceVersion()
}

func (*SnapshotRecordings) MemberEndpointsDeprecated() uint16 {
	return 0
}

func (SnapshotRecordings) MemberEndpointsCharacterEncoding() string {
	return "US-ASCII"
}

func (SnapshotRecordings) MemberEndpointsHeaderLength() uint64 {
	return 4
}
