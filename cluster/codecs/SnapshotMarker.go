// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SnapshotMarker struct {
	TypeId           int64
	LogPosition      int64
	LeadershipTermId int64
	Index            int32
	Mark             SnapshotMarkEnum
	TimeUnit         ClusterTimeUnitEnum
	AppVersion       int32
}

func (s *SnapshotMarker) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.TypeId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.Index); err != nil {
		return err
	}
	if err := s.Mark.Encode(_m, _w); err != nil {
		return err
	}
	if err := s.TimeUnit.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.AppVersion); err != nil {
		return err
	}
	return nil
}

func (s *SnapshotMarker) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.TypeIdInActingVersion(actingVersion) {
		s.TypeId = s.TypeIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.TypeId); err != nil {
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
	if !s.LeadershipTermIdInActingVersion(actingVersion) {
		s.LeadershipTermId = s.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LeadershipTermId); err != nil {
			return err
		}
	}
	if !s.IndexInActingVersion(actingVersion) {
		s.Index = s.IndexNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.Index); err != nil {
			return err
		}
	}
	if s.MarkInActingVersion(actingVersion) {
		if err := s.Mark.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if s.TimeUnitInActingVersion(actingVersion) {
		if err := s.TimeUnit.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !s.AppVersionInActingVersion(actingVersion) {
		s.AppVersion = s.AppVersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.AppVersion); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapshotMarker) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.TypeIdInActingVersion(actingVersion) {
		if s.TypeId < s.TypeIdMinValue() || s.TypeId > s.TypeIdMaxValue() {
			return fmt.Errorf("Range check failed on s.TypeId (%v < %v > %v)", s.TypeIdMinValue(), s.TypeId, s.TypeIdMaxValue())
		}
	}
	if s.LogPositionInActingVersion(actingVersion) {
		if s.LogPosition < s.LogPositionMinValue() || s.LogPosition > s.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on s.LogPosition (%v < %v > %v)", s.LogPositionMinValue(), s.LogPosition, s.LogPositionMaxValue())
		}
	}
	if s.LeadershipTermIdInActingVersion(actingVersion) {
		if s.LeadershipTermId < s.LeadershipTermIdMinValue() || s.LeadershipTermId > s.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeadershipTermId (%v < %v > %v)", s.LeadershipTermIdMinValue(), s.LeadershipTermId, s.LeadershipTermIdMaxValue())
		}
	}
	if s.IndexInActingVersion(actingVersion) {
		if s.Index < s.IndexMinValue() || s.Index > s.IndexMaxValue() {
			return fmt.Errorf("Range check failed on s.Index (%v < %v > %v)", s.IndexMinValue(), s.Index, s.IndexMaxValue())
		}
	}
	if err := s.Mark.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if err := s.TimeUnit.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if s.AppVersionInActingVersion(actingVersion) {
		if s.AppVersion != s.AppVersionNullValue() && (s.AppVersion < s.AppVersionMinValue() || s.AppVersion > s.AppVersionMaxValue()) {
			return fmt.Errorf("Range check failed on s.AppVersion (%v < %v > %v)", s.AppVersionMinValue(), s.AppVersion, s.AppVersionMaxValue())
		}
	}
	return nil
}

func SnapshotMarkerInit(s *SnapshotMarker) {
	s.AppVersion = 0
	return
}

func (*SnapshotMarker) SbeBlockLength() (blockLength uint16) {
	return 40
}

func (*SnapshotMarker) SbeTemplateId() (templateId uint16) {
	return 100
}

func (*SnapshotMarker) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SnapshotMarker) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SnapshotMarker) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SnapshotMarker) TypeIdId() uint16 {
	return 1
}

func (*SnapshotMarker) TypeIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarker) TypeIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TypeIdSinceVersion()
}

func (*SnapshotMarker) TypeIdDeprecated() uint16 {
	return 0
}

func (*SnapshotMarker) TypeIdMetaAttribute(meta int) string {
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

func (*SnapshotMarker) TypeIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotMarker) TypeIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotMarker) TypeIdNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotMarker) LogPositionId() uint16 {
	return 2
}

func (*SnapshotMarker) LogPositionSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarker) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LogPositionSinceVersion()
}

func (*SnapshotMarker) LogPositionDeprecated() uint16 {
	return 0
}

func (*SnapshotMarker) LogPositionMetaAttribute(meta int) string {
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

func (*SnapshotMarker) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotMarker) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotMarker) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotMarker) LeadershipTermIdId() uint16 {
	return 3
}

func (*SnapshotMarker) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarker) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*SnapshotMarker) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*SnapshotMarker) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*SnapshotMarker) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SnapshotMarker) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SnapshotMarker) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*SnapshotMarker) IndexId() uint16 {
	return 4
}

func (*SnapshotMarker) IndexSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarker) IndexInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.IndexSinceVersion()
}

func (*SnapshotMarker) IndexDeprecated() uint16 {
	return 0
}

func (*SnapshotMarker) IndexMetaAttribute(meta int) string {
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

func (*SnapshotMarker) IndexMinValue() int32 {
	return math.MinInt32 + 1
}

func (*SnapshotMarker) IndexMaxValue() int32 {
	return math.MaxInt32
}

func (*SnapshotMarker) IndexNullValue() int32 {
	return math.MinInt32
}

func (*SnapshotMarker) MarkId() uint16 {
	return 5
}

func (*SnapshotMarker) MarkSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarker) MarkInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.MarkSinceVersion()
}

func (*SnapshotMarker) MarkDeprecated() uint16 {
	return 0
}

func (*SnapshotMarker) MarkMetaAttribute(meta int) string {
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

func (*SnapshotMarker) TimeUnitId() uint16 {
	return 6
}

func (*SnapshotMarker) TimeUnitSinceVersion() uint16 {
	return 4
}

func (s *SnapshotMarker) TimeUnitInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TimeUnitSinceVersion()
}

func (*SnapshotMarker) TimeUnitDeprecated() uint16 {
	return 0
}

func (*SnapshotMarker) TimeUnitMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "optional"
	}
	return ""
}

func (*SnapshotMarker) AppVersionId() uint16 {
	return 7
}

func (*SnapshotMarker) AppVersionSinceVersion() uint16 {
	return 4
}

func (s *SnapshotMarker) AppVersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.AppVersionSinceVersion()
}

func (*SnapshotMarker) AppVersionDeprecated() uint16 {
	return 0
}

func (*SnapshotMarker) AppVersionMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "optional"
	}
	return ""
}

func (*SnapshotMarker) AppVersionMinValue() int32 {
	return 1
}

func (*SnapshotMarker) AppVersionMaxValue() int32 {
	return 16777215
}

func (*SnapshotMarker) AppVersionNullValue() int32 {
	return 0
}
