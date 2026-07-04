// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StandbySnapshot struct {
	CorrelationId      int64
	Version            int32
	ResponseStreamId   int32
	Snapshots          []StandbySnapshotSnapshots
	ResponseChannel    []uint8
	EncodedCredentials []uint8
}
type StandbySnapshotSnapshots struct {
	RecordingId         int64
	LeadershipTermId    int64
	TermBaseLogPosition int64
	LogPosition         int64
	Timestamp           int64
	ServiceId           int32
	ArchiveEndpoint     []uint8
}

func (s *StandbySnapshot) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.Version); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.ResponseStreamId); err != nil {
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
	for i := range s.Snapshots {
		if err := s.Snapshots[i].Encode(_m, _w); err != nil {
			return err
		}
	}
	if err := _m.WriteUint32(_w, uint32(len(s.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.ResponseChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.EncodedCredentials))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.EncodedCredentials); err != nil {
		return err
	}
	return nil
}

func (s *StandbySnapshot) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
			return err
		}
	}
	if !s.VersionInActingVersion(actingVersion) {
		s.Version = s.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.Version); err != nil {
			return err
		}
	}
	if !s.ResponseStreamIdInActingVersion(actingVersion) {
		s.ResponseStreamId = s.ResponseStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.ResponseStreamId); err != nil {
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
			s.Snapshots = make([]StandbySnapshotSnapshots, SnapshotsNumInGroup)
		}
		s.Snapshots = s.Snapshots[:SnapshotsNumInGroup]
		for i := range s.Snapshots {
			if err := s.Snapshots[i].Decode(_m, _r, actingVersion, uint(SnapshotsBlockLength)); err != nil {
				return err
			}
		}
	}

	if s.ResponseChannelInActingVersion(actingVersion) {
		var ResponseChannelLength uint32
		if err := _m.ReadUint32(_r, &ResponseChannelLength); err != nil {
			return err
		}
		if cap(s.ResponseChannel) < int(ResponseChannelLength) {
			s.ResponseChannel = make([]uint8, ResponseChannelLength)
		}
		s.ResponseChannel = s.ResponseChannel[:ResponseChannelLength]
		if err := _m.ReadBytes(_r, s.ResponseChannel); err != nil {
			return err
		}
	}

	if s.EncodedCredentialsInActingVersion(actingVersion) {
		var EncodedCredentialsLength uint32
		if err := _m.ReadUint32(_r, &EncodedCredentialsLength); err != nil {
			return err
		}
		if cap(s.EncodedCredentials) < int(EncodedCredentialsLength) {
			s.EncodedCredentials = make([]uint8, EncodedCredentialsLength)
		}
		s.EncodedCredentials = s.EncodedCredentials[:EncodedCredentialsLength]
		if err := _m.ReadBytes(_r, s.EncodedCredentials); err != nil {
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

func (s *StandbySnapshot) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	if s.VersionInActingVersion(actingVersion) {
		if s.Version < s.VersionMinValue() || s.Version > s.VersionMaxValue() {
			return fmt.Errorf("Range check failed on s.Version (%v < %v > %v)", s.VersionMinValue(), s.Version, s.VersionMaxValue())
		}
	}
	if s.ResponseStreamIdInActingVersion(actingVersion) {
		if s.ResponseStreamId < s.ResponseStreamIdMinValue() || s.ResponseStreamId > s.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ResponseStreamId (%v < %v > %v)", s.ResponseStreamIdMinValue(), s.ResponseStreamId, s.ResponseStreamIdMaxValue())
		}
	}
	for i := range s.Snapshots {
		if err := s.Snapshots[i].RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	for idx, ch := range s.ResponseChannel {
		if ch > 127 {
			return fmt.Errorf("s.ResponseChannel[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func StandbySnapshotInit(s *StandbySnapshot) {
	return
}

func (s *StandbySnapshotSnapshots) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
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
	if err := _m.WriteUint32(_w, uint32(len(s.ArchiveEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.ArchiveEndpoint); err != nil {
		return err
	}
	return nil
}

func (s *StandbySnapshotSnapshots) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
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

	if s.ArchiveEndpointInActingVersion(actingVersion) {
		var ArchiveEndpointLength uint32
		if err := _m.ReadUint32(_r, &ArchiveEndpointLength); err != nil {
			return err
		}
		if cap(s.ArchiveEndpoint) < int(ArchiveEndpointLength) {
			s.ArchiveEndpoint = make([]uint8, ArchiveEndpointLength)
		}
		s.ArchiveEndpoint = s.ArchiveEndpoint[:ArchiveEndpointLength]
		if err := _m.ReadBytes(_r, s.ArchiveEndpoint); err != nil {
			return err
		}
	}
	return nil
}

func (s *StandbySnapshotSnapshots) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	for idx, ch := range s.ArchiveEndpoint {
		if ch > 127 {
			return fmt.Errorf("s.ArchiveEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func StandbySnapshotSnapshotsInit(s *StandbySnapshotSnapshots) {
	return
}

func (*StandbySnapshot) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*StandbySnapshot) SbeTemplateId() (templateId uint16) {
	return 81
}

func (*StandbySnapshot) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*StandbySnapshot) SbeSchemaVersion() (schemaVersion uint16) {
	return 16
}

func (*StandbySnapshot) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StandbySnapshot) SbeSemanticVersion() (semanticVersion string) {
	return "5.4"
}

func (*StandbySnapshot) CorrelationIdId() uint16 {
	return 1
}

func (*StandbySnapshot) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshot) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*StandbySnapshot) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*StandbySnapshot) CorrelationIdMetaAttribute(meta int) string {
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

func (*StandbySnapshot) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StandbySnapshot) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StandbySnapshot) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*StandbySnapshot) VersionId() uint16 {
	return 2
}

func (*StandbySnapshot) VersionSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshot) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.VersionSinceVersion()
}

func (*StandbySnapshot) VersionDeprecated() uint16 {
	return 0
}

func (*StandbySnapshot) VersionMetaAttribute(meta int) string {
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

func (*StandbySnapshot) VersionMinValue() int32 {
	return math.MinInt32 + 1
}

func (*StandbySnapshot) VersionMaxValue() int32 {
	return math.MaxInt32
}

func (*StandbySnapshot) VersionNullValue() int32 {
	return math.MinInt32
}

func (*StandbySnapshot) ResponseStreamIdId() uint16 {
	return 3
}

func (*StandbySnapshot) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshot) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ResponseStreamIdSinceVersion()
}

func (*StandbySnapshot) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*StandbySnapshot) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*StandbySnapshot) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*StandbySnapshot) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*StandbySnapshot) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*StandbySnapshotSnapshots) RecordingIdId() uint16 {
	return 5
}

func (*StandbySnapshotSnapshots) RecordingIdSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshotSnapshots) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.RecordingIdSinceVersion()
}

func (*StandbySnapshotSnapshots) RecordingIdDeprecated() uint16 {
	return 0
}

func (*StandbySnapshotSnapshots) RecordingIdMetaAttribute(meta int) string {
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

func (*StandbySnapshotSnapshots) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StandbySnapshotSnapshots) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StandbySnapshotSnapshots) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*StandbySnapshotSnapshots) LeadershipTermIdId() uint16 {
	return 6
}

func (*StandbySnapshotSnapshots) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshotSnapshots) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*StandbySnapshotSnapshots) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*StandbySnapshotSnapshots) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*StandbySnapshotSnapshots) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StandbySnapshotSnapshots) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StandbySnapshotSnapshots) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*StandbySnapshotSnapshots) TermBaseLogPositionId() uint16 {
	return 7
}

func (*StandbySnapshotSnapshots) TermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshotSnapshots) TermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TermBaseLogPositionSinceVersion()
}

func (*StandbySnapshotSnapshots) TermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*StandbySnapshotSnapshots) TermBaseLogPositionMetaAttribute(meta int) string {
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

func (*StandbySnapshotSnapshots) TermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StandbySnapshotSnapshots) TermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*StandbySnapshotSnapshots) TermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*StandbySnapshotSnapshots) LogPositionId() uint16 {
	return 8
}

func (*StandbySnapshotSnapshots) LogPositionSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshotSnapshots) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LogPositionSinceVersion()
}

func (*StandbySnapshotSnapshots) LogPositionDeprecated() uint16 {
	return 0
}

func (*StandbySnapshotSnapshots) LogPositionMetaAttribute(meta int) string {
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

func (*StandbySnapshotSnapshots) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StandbySnapshotSnapshots) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*StandbySnapshotSnapshots) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*StandbySnapshotSnapshots) TimestampId() uint16 {
	return 9
}

func (*StandbySnapshotSnapshots) TimestampSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshotSnapshots) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TimestampSinceVersion()
}

func (*StandbySnapshotSnapshots) TimestampDeprecated() uint16 {
	return 0
}

func (*StandbySnapshotSnapshots) TimestampMetaAttribute(meta int) string {
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

func (*StandbySnapshotSnapshots) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StandbySnapshotSnapshots) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*StandbySnapshotSnapshots) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*StandbySnapshotSnapshots) ServiceIdId() uint16 {
	return 10
}

func (*StandbySnapshotSnapshots) ServiceIdSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshotSnapshots) ServiceIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ServiceIdSinceVersion()
}

func (*StandbySnapshotSnapshots) ServiceIdDeprecated() uint16 {
	return 0
}

func (*StandbySnapshotSnapshots) ServiceIdMetaAttribute(meta int) string {
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

func (*StandbySnapshotSnapshots) ServiceIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*StandbySnapshotSnapshots) ServiceIdMaxValue() int32 {
	return math.MaxInt32
}

func (*StandbySnapshotSnapshots) ServiceIdNullValue() int32 {
	return math.MinInt32
}

func (*StandbySnapshotSnapshots) ArchiveEndpointMetaAttribute(meta int) string {
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

func (*StandbySnapshotSnapshots) ArchiveEndpointSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshotSnapshots) ArchiveEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ArchiveEndpointSinceVersion()
}

func (*StandbySnapshotSnapshots) ArchiveEndpointDeprecated() uint16 {
	return 0
}

func (StandbySnapshotSnapshots) ArchiveEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (StandbySnapshotSnapshots) ArchiveEndpointHeaderLength() uint64 {
	return 4
}

func (*StandbySnapshot) SnapshotsId() uint16 {
	return 4
}

func (*StandbySnapshot) SnapshotsSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshot) SnapshotsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.SnapshotsSinceVersion()
}

func (*StandbySnapshot) SnapshotsDeprecated() uint16 {
	return 0
}

func (*StandbySnapshotSnapshots) SbeBlockLength() (blockLength uint) {
	return 44
}

func (*StandbySnapshotSnapshots) SbeSchemaVersion() (schemaVersion uint16) {
	return 16
}

func (*StandbySnapshot) ResponseChannelMetaAttribute(meta int) string {
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

func (*StandbySnapshot) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshot) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ResponseChannelSinceVersion()
}

func (*StandbySnapshot) ResponseChannelDeprecated() uint16 {
	return 0
}

func (StandbySnapshot) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (StandbySnapshot) ResponseChannelHeaderLength() uint64 {
	return 4
}

func (*StandbySnapshot) EncodedCredentialsMetaAttribute(meta int) string {
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

func (*StandbySnapshot) EncodedCredentialsSinceVersion() uint16 {
	return 0
}

func (s *StandbySnapshot) EncodedCredentialsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.EncodedCredentialsSinceVersion()
}

func (*StandbySnapshot) EncodedCredentialsDeprecated() uint16 {
	return 0
}

func (StandbySnapshot) EncodedCredentialsCharacterEncoding() string {
	return "null"
}

func (StandbySnapshot) EncodedCredentialsHeaderLength() uint64 {
	return 4
}
