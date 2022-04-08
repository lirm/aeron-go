// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type NewLeadershipTermEvent struct {
	LeadershipTermId    int64
	LogPosition         int64
	Timestamp           int64
	TermBaseLogPosition int64
	LeaderMemberId      int32
	LogSessionId        int32
	TimeUnit            ClusterTimeUnitEnum
	AppVersion          int32
}

func (n *NewLeadershipTermEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := n.RangeCheck(n.SbeSchemaVersion(), n.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, n.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.TermBaseLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, n.LeaderMemberId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, n.LogSessionId); err != nil {
		return err
	}
	if err := n.TimeUnit.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, n.AppVersion); err != nil {
		return err
	}
	return nil
}

func (n *NewLeadershipTermEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !n.LeadershipTermIdInActingVersion(actingVersion) {
		n.LeadershipTermId = n.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.LeadershipTermId); err != nil {
			return err
		}
	}
	if !n.LogPositionInActingVersion(actingVersion) {
		n.LogPosition = n.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.LogPosition); err != nil {
			return err
		}
	}
	if !n.TimestampInActingVersion(actingVersion) {
		n.Timestamp = n.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.Timestamp); err != nil {
			return err
		}
	}
	if !n.TermBaseLogPositionInActingVersion(actingVersion) {
		n.TermBaseLogPosition = n.TermBaseLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.TermBaseLogPosition); err != nil {
			return err
		}
	}
	if !n.LeaderMemberIdInActingVersion(actingVersion) {
		n.LeaderMemberId = n.LeaderMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &n.LeaderMemberId); err != nil {
			return err
		}
	}
	if !n.LogSessionIdInActingVersion(actingVersion) {
		n.LogSessionId = n.LogSessionIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &n.LogSessionId); err != nil {
			return err
		}
	}
	if n.TimeUnitInActingVersion(actingVersion) {
		if err := n.TimeUnit.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !n.AppVersionInActingVersion(actingVersion) {
		n.AppVersion = n.AppVersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &n.AppVersion); err != nil {
			return err
		}
	}
	if actingVersion > n.SbeSchemaVersion() && blockLength > n.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-n.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := n.RangeCheck(actingVersion, n.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (n *NewLeadershipTermEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if n.LeadershipTermIdInActingVersion(actingVersion) {
		if n.LeadershipTermId < n.LeadershipTermIdMinValue() || n.LeadershipTermId > n.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LeadershipTermId (%v < %v > %v)", n.LeadershipTermIdMinValue(), n.LeadershipTermId, n.LeadershipTermIdMaxValue())
		}
	}
	if n.LogPositionInActingVersion(actingVersion) {
		if n.LogPosition < n.LogPositionMinValue() || n.LogPosition > n.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on n.LogPosition (%v < %v > %v)", n.LogPositionMinValue(), n.LogPosition, n.LogPositionMaxValue())
		}
	}
	if n.TimestampInActingVersion(actingVersion) {
		if n.Timestamp < n.TimestampMinValue() || n.Timestamp > n.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on n.Timestamp (%v < %v > %v)", n.TimestampMinValue(), n.Timestamp, n.TimestampMaxValue())
		}
	}
	if n.TermBaseLogPositionInActingVersion(actingVersion) {
		if n.TermBaseLogPosition < n.TermBaseLogPositionMinValue() || n.TermBaseLogPosition > n.TermBaseLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on n.TermBaseLogPosition (%v < %v > %v)", n.TermBaseLogPositionMinValue(), n.TermBaseLogPosition, n.TermBaseLogPositionMaxValue())
		}
	}
	if n.LeaderMemberIdInActingVersion(actingVersion) {
		if n.LeaderMemberId < n.LeaderMemberIdMinValue() || n.LeaderMemberId > n.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LeaderMemberId (%v < %v > %v)", n.LeaderMemberIdMinValue(), n.LeaderMemberId, n.LeaderMemberIdMaxValue())
		}
	}
	if n.LogSessionIdInActingVersion(actingVersion) {
		if n.LogSessionId < n.LogSessionIdMinValue() || n.LogSessionId > n.LogSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LogSessionId (%v < %v > %v)", n.LogSessionIdMinValue(), n.LogSessionId, n.LogSessionIdMaxValue())
		}
	}
	if err := n.TimeUnit.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if n.AppVersionInActingVersion(actingVersion) {
		if n.AppVersion != n.AppVersionNullValue() && (n.AppVersion < n.AppVersionMinValue() || n.AppVersion > n.AppVersionMaxValue()) {
			return fmt.Errorf("Range check failed on n.AppVersion (%v < %v > %v)", n.AppVersionMinValue(), n.AppVersion, n.AppVersionMaxValue())
		}
	}
	return nil
}

func NewLeadershipTermEventInit(n *NewLeadershipTermEvent) {
	n.AppVersion = 0
	return
}

func (*NewLeadershipTermEvent) SbeBlockLength() (blockLength uint16) {
	return 48
}

func (*NewLeadershipTermEvent) SbeTemplateId() (templateId uint16) {
	return 24
}

func (*NewLeadershipTermEvent) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*NewLeadershipTermEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*NewLeadershipTermEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*NewLeadershipTermEvent) LeadershipTermIdId() uint16 {
	return 1
}

func (*NewLeadershipTermEvent) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTermEvent) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LeadershipTermIdSinceVersion()
}

func (*NewLeadershipTermEvent) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTermEvent) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTermEvent) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTermEvent) LogPositionId() uint16 {
	return 2
}

func (*NewLeadershipTermEvent) LogPositionSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTermEvent) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LogPositionSinceVersion()
}

func (*NewLeadershipTermEvent) LogPositionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) LogPositionMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTermEvent) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTermEvent) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTermEvent) TimestampId() uint16 {
	return 3
}

func (*NewLeadershipTermEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTermEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.TimestampSinceVersion()
}

func (*NewLeadershipTermEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) TimestampMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTermEvent) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTermEvent) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTermEvent) TermBaseLogPositionId() uint16 {
	return 4
}

func (*NewLeadershipTermEvent) TermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTermEvent) TermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.TermBaseLogPositionSinceVersion()
}

func (*NewLeadershipTermEvent) TermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) TermBaseLogPositionMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) TermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTermEvent) TermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTermEvent) TermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTermEvent) LeaderMemberIdId() uint16 {
	return 5
}

func (*NewLeadershipTermEvent) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTermEvent) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LeaderMemberIdSinceVersion()
}

func (*NewLeadershipTermEvent) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*NewLeadershipTermEvent) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*NewLeadershipTermEvent) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*NewLeadershipTermEvent) LogSessionIdId() uint16 {
	return 6
}

func (*NewLeadershipTermEvent) LogSessionIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTermEvent) LogSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LogSessionIdSinceVersion()
}

func (*NewLeadershipTermEvent) LogSessionIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) LogSessionIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) LogSessionIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*NewLeadershipTermEvent) LogSessionIdMaxValue() int32 {
	return math.MaxInt32
}

func (*NewLeadershipTermEvent) LogSessionIdNullValue() int32 {
	return math.MinInt32
}

func (*NewLeadershipTermEvent) TimeUnitId() uint16 {
	return 7
}

func (*NewLeadershipTermEvent) TimeUnitSinceVersion() uint16 {
	return 4
}

func (n *NewLeadershipTermEvent) TimeUnitInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.TimeUnitSinceVersion()
}

func (*NewLeadershipTermEvent) TimeUnitDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) TimeUnitMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) AppVersionId() uint16 {
	return 8
}

func (*NewLeadershipTermEvent) AppVersionSinceVersion() uint16 {
	return 4
}

func (n *NewLeadershipTermEvent) AppVersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.AppVersionSinceVersion()
}

func (*NewLeadershipTermEvent) AppVersionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTermEvent) AppVersionMetaAttribute(meta int) string {
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

func (*NewLeadershipTermEvent) AppVersionMinValue() int32 {
	return 1
}

func (*NewLeadershipTermEvent) AppVersionMaxValue() int32 {
	return 16777215
}

func (*NewLeadershipTermEvent) AppVersionNullValue() int32 {
	return 0
}
