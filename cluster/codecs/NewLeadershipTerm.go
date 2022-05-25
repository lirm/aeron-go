// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type NewLeadershipTerm struct {
	LogLeadershipTermId     int64
	NextLeadershipTermId    int64
	NextTermBaseLogPosition int64
	NextLogPosition         int64
	LeadershipTermId        int64
	TermBaseLogPosition     int64
	LogPosition             int64
	LeaderRecordingId       int64
	Timestamp               int64
	LeaderMemberId          int32
	LogSessionId            int32
	AppVersion              int32
	IsStartup               BooleanTypeEnum
}

func (n *NewLeadershipTerm) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := n.RangeCheck(n.SbeSchemaVersion(), n.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, n.LogLeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.NextLeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.NextTermBaseLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.NextLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.TermBaseLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.LeaderRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, n.LeaderMemberId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, n.LogSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, n.AppVersion); err != nil {
		return err
	}
	if err := n.IsStartup.Encode(_m, _w); err != nil {
		return err
	}
	return nil
}

func (n *NewLeadershipTerm) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !n.LogLeadershipTermIdInActingVersion(actingVersion) {
		n.LogLeadershipTermId = n.LogLeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.LogLeadershipTermId); err != nil {
			return err
		}
	}
	if !n.NextLeadershipTermIdInActingVersion(actingVersion) {
		n.NextLeadershipTermId = n.NextLeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.NextLeadershipTermId); err != nil {
			return err
		}
	}
	if !n.NextTermBaseLogPositionInActingVersion(actingVersion) {
		n.NextTermBaseLogPosition = n.NextTermBaseLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.NextTermBaseLogPosition); err != nil {
			return err
		}
	}
	if !n.NextLogPositionInActingVersion(actingVersion) {
		n.NextLogPosition = n.NextLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.NextLogPosition); err != nil {
			return err
		}
	}
	if !n.LeadershipTermIdInActingVersion(actingVersion) {
		n.LeadershipTermId = n.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.LeadershipTermId); err != nil {
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
	if !n.LogPositionInActingVersion(actingVersion) {
		n.LogPosition = n.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.LogPosition); err != nil {
			return err
		}
	}
	if !n.LeaderRecordingIdInActingVersion(actingVersion) {
		n.LeaderRecordingId = n.LeaderRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.LeaderRecordingId); err != nil {
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
	if !n.AppVersionInActingVersion(actingVersion) {
		n.AppVersion = n.AppVersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &n.AppVersion); err != nil {
			return err
		}
	}
	if n.IsStartupInActingVersion(actingVersion) {
		if err := n.IsStartup.Decode(_m, _r, actingVersion); err != nil {
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

func (n *NewLeadershipTerm) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if n.LogLeadershipTermIdInActingVersion(actingVersion) {
		if n.LogLeadershipTermId < n.LogLeadershipTermIdMinValue() || n.LogLeadershipTermId > n.LogLeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LogLeadershipTermId (%v < %v > %v)", n.LogLeadershipTermIdMinValue(), n.LogLeadershipTermId, n.LogLeadershipTermIdMaxValue())
		}
	}
	if n.NextLeadershipTermIdInActingVersion(actingVersion) {
		if n.NextLeadershipTermId < n.NextLeadershipTermIdMinValue() || n.NextLeadershipTermId > n.NextLeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on n.NextLeadershipTermId (%v < %v > %v)", n.NextLeadershipTermIdMinValue(), n.NextLeadershipTermId, n.NextLeadershipTermIdMaxValue())
		}
	}
	if n.NextTermBaseLogPositionInActingVersion(actingVersion) {
		if n.NextTermBaseLogPosition < n.NextTermBaseLogPositionMinValue() || n.NextTermBaseLogPosition > n.NextTermBaseLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on n.NextTermBaseLogPosition (%v < %v > %v)", n.NextTermBaseLogPositionMinValue(), n.NextTermBaseLogPosition, n.NextTermBaseLogPositionMaxValue())
		}
	}
	if n.NextLogPositionInActingVersion(actingVersion) {
		if n.NextLogPosition < n.NextLogPositionMinValue() || n.NextLogPosition > n.NextLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on n.NextLogPosition (%v < %v > %v)", n.NextLogPositionMinValue(), n.NextLogPosition, n.NextLogPositionMaxValue())
		}
	}
	if n.LeadershipTermIdInActingVersion(actingVersion) {
		if n.LeadershipTermId < n.LeadershipTermIdMinValue() || n.LeadershipTermId > n.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LeadershipTermId (%v < %v > %v)", n.LeadershipTermIdMinValue(), n.LeadershipTermId, n.LeadershipTermIdMaxValue())
		}
	}
	if n.TermBaseLogPositionInActingVersion(actingVersion) {
		if n.TermBaseLogPosition < n.TermBaseLogPositionMinValue() || n.TermBaseLogPosition > n.TermBaseLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on n.TermBaseLogPosition (%v < %v > %v)", n.TermBaseLogPositionMinValue(), n.TermBaseLogPosition, n.TermBaseLogPositionMaxValue())
		}
	}
	if n.LogPositionInActingVersion(actingVersion) {
		if n.LogPosition < n.LogPositionMinValue() || n.LogPosition > n.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on n.LogPosition (%v < %v > %v)", n.LogPositionMinValue(), n.LogPosition, n.LogPositionMaxValue())
		}
	}
	if n.LeaderRecordingIdInActingVersion(actingVersion) {
		if n.LeaderRecordingId < n.LeaderRecordingIdMinValue() || n.LeaderRecordingId > n.LeaderRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LeaderRecordingId (%v < %v > %v)", n.LeaderRecordingIdMinValue(), n.LeaderRecordingId, n.LeaderRecordingIdMaxValue())
		}
	}
	if n.TimestampInActingVersion(actingVersion) {
		if n.Timestamp < n.TimestampMinValue() || n.Timestamp > n.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on n.Timestamp (%v < %v > %v)", n.TimestampMinValue(), n.Timestamp, n.TimestampMaxValue())
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
	if n.AppVersionInActingVersion(actingVersion) {
		if n.AppVersion != n.AppVersionNullValue() && (n.AppVersion < n.AppVersionMinValue() || n.AppVersion > n.AppVersionMaxValue()) {
			return fmt.Errorf("Range check failed on n.AppVersion (%v < %v > %v)", n.AppVersionMinValue(), n.AppVersion, n.AppVersionMaxValue())
		}
	}
	if err := n.IsStartup.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func NewLeadershipTermInit(n *NewLeadershipTerm) {
	n.AppVersion = 0
	return
}

func (*NewLeadershipTerm) SbeBlockLength() (blockLength uint16) {
	return 88
}

func (*NewLeadershipTerm) SbeTemplateId() (templateId uint16) {
	return 53
}

func (*NewLeadershipTerm) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*NewLeadershipTerm) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*NewLeadershipTerm) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*NewLeadershipTerm) LogLeadershipTermIdId() uint16 {
	return 1
}

func (*NewLeadershipTerm) LogLeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) LogLeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LogLeadershipTermIdSinceVersion()
}

func (*NewLeadershipTerm) LogLeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) LogLeadershipTermIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) LogLeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) LogLeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) LogLeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) NextLeadershipTermIdId() uint16 {
	return 2
}

func (*NewLeadershipTerm) NextLeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) NextLeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.NextLeadershipTermIdSinceVersion()
}

func (*NewLeadershipTerm) NextLeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) NextLeadershipTermIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) NextLeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) NextLeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) NextLeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) NextTermBaseLogPositionId() uint16 {
	return 3
}

func (*NewLeadershipTerm) NextTermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) NextTermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.NextTermBaseLogPositionSinceVersion()
}

func (*NewLeadershipTerm) NextTermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) NextTermBaseLogPositionMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) NextTermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) NextTermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) NextTermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) NextLogPositionId() uint16 {
	return 4
}

func (*NewLeadershipTerm) NextLogPositionSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) NextLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.NextLogPositionSinceVersion()
}

func (*NewLeadershipTerm) NextLogPositionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) NextLogPositionMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) NextLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) NextLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) NextLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) LeadershipTermIdId() uint16 {
	return 5
}

func (*NewLeadershipTerm) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LeadershipTermIdSinceVersion()
}

func (*NewLeadershipTerm) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) TermBaseLogPositionId() uint16 {
	return 6
}

func (*NewLeadershipTerm) TermBaseLogPositionSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) TermBaseLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.TermBaseLogPositionSinceVersion()
}

func (*NewLeadershipTerm) TermBaseLogPositionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) TermBaseLogPositionMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) TermBaseLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) TermBaseLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) TermBaseLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) LogPositionId() uint16 {
	return 7
}

func (*NewLeadershipTerm) LogPositionSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LogPositionSinceVersion()
}

func (*NewLeadershipTerm) LogPositionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) LogPositionMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) LeaderRecordingIdId() uint16 {
	return 8
}

func (*NewLeadershipTerm) LeaderRecordingIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) LeaderRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LeaderRecordingIdSinceVersion()
}

func (*NewLeadershipTerm) LeaderRecordingIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) LeaderRecordingIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) LeaderRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) LeaderRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) LeaderRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) TimestampId() uint16 {
	return 9
}

func (*NewLeadershipTerm) TimestampSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.TimestampSinceVersion()
}

func (*NewLeadershipTerm) TimestampDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) TimestampMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeadershipTerm) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeadershipTerm) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*NewLeadershipTerm) LeaderMemberIdId() uint16 {
	return 10
}

func (*NewLeadershipTerm) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LeaderMemberIdSinceVersion()
}

func (*NewLeadershipTerm) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*NewLeadershipTerm) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*NewLeadershipTerm) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*NewLeadershipTerm) LogSessionIdId() uint16 {
	return 11
}

func (*NewLeadershipTerm) LogSessionIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) LogSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LogSessionIdSinceVersion()
}

func (*NewLeadershipTerm) LogSessionIdDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) LogSessionIdMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) LogSessionIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*NewLeadershipTerm) LogSessionIdMaxValue() int32 {
	return math.MaxInt32
}

func (*NewLeadershipTerm) LogSessionIdNullValue() int32 {
	return math.MinInt32
}

func (*NewLeadershipTerm) AppVersionId() uint16 {
	return 12
}

func (*NewLeadershipTerm) AppVersionSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) AppVersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.AppVersionSinceVersion()
}

func (*NewLeadershipTerm) AppVersionDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) AppVersionMetaAttribute(meta int) string {
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

func (*NewLeadershipTerm) AppVersionMinValue() int32 {
	return 1
}

func (*NewLeadershipTerm) AppVersionMaxValue() int32 {
	return 16777215
}

func (*NewLeadershipTerm) AppVersionNullValue() int32 {
	return 0
}

func (*NewLeadershipTerm) IsStartupId() uint16 {
	return 13
}

func (*NewLeadershipTerm) IsStartupSinceVersion() uint16 {
	return 0
}

func (n *NewLeadershipTerm) IsStartupInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.IsStartupSinceVersion()
}

func (*NewLeadershipTerm) IsStartupDeprecated() uint16 {
	return 0
}

func (*NewLeadershipTerm) IsStartupMetaAttribute(meta int) string {
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
