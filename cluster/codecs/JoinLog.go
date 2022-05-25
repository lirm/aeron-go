// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type JoinLog struct {
	LogPosition    int64
	MaxLogPosition int64
	MemberId       int32
	LogSessionId   int32
	LogStreamId    int32
	IsStartup      BooleanTypeEnum
	Role           int32
	LogChannel     []uint8
}

func (j *JoinLog) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := j.RangeCheck(j.SbeSchemaVersion(), j.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, j.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, j.MaxLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, j.MemberId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, j.LogSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, j.LogStreamId); err != nil {
		return err
	}
	if err := j.IsStartup.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, j.Role); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(j.LogChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, j.LogChannel); err != nil {
		return err
	}
	return nil
}

func (j *JoinLog) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !j.LogPositionInActingVersion(actingVersion) {
		j.LogPosition = j.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &j.LogPosition); err != nil {
			return err
		}
	}
	if !j.MaxLogPositionInActingVersion(actingVersion) {
		j.MaxLogPosition = j.MaxLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &j.MaxLogPosition); err != nil {
			return err
		}
	}
	if !j.MemberIdInActingVersion(actingVersion) {
		j.MemberId = j.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &j.MemberId); err != nil {
			return err
		}
	}
	if !j.LogSessionIdInActingVersion(actingVersion) {
		j.LogSessionId = j.LogSessionIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &j.LogSessionId); err != nil {
			return err
		}
	}
	if !j.LogStreamIdInActingVersion(actingVersion) {
		j.LogStreamId = j.LogStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &j.LogStreamId); err != nil {
			return err
		}
	}
	if j.IsStartupInActingVersion(actingVersion) {
		if err := j.IsStartup.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !j.RoleInActingVersion(actingVersion) {
		j.Role = j.RoleNullValue()
	} else {
		if err := _m.ReadInt32(_r, &j.Role); err != nil {
			return err
		}
	}
	if actingVersion > j.SbeSchemaVersion() && blockLength > j.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-j.SbeBlockLength()))
	}

	if j.LogChannelInActingVersion(actingVersion) {
		var LogChannelLength uint32
		if err := _m.ReadUint32(_r, &LogChannelLength); err != nil {
			return err
		}
		if cap(j.LogChannel) < int(LogChannelLength) {
			j.LogChannel = make([]uint8, LogChannelLength)
		}
		j.LogChannel = j.LogChannel[:LogChannelLength]
		if err := _m.ReadBytes(_r, j.LogChannel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := j.RangeCheck(actingVersion, j.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (j *JoinLog) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if j.LogPositionInActingVersion(actingVersion) {
		if j.LogPosition < j.LogPositionMinValue() || j.LogPosition > j.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on j.LogPosition (%v < %v > %v)", j.LogPositionMinValue(), j.LogPosition, j.LogPositionMaxValue())
		}
	}
	if j.MaxLogPositionInActingVersion(actingVersion) {
		if j.MaxLogPosition < j.MaxLogPositionMinValue() || j.MaxLogPosition > j.MaxLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on j.MaxLogPosition (%v < %v > %v)", j.MaxLogPositionMinValue(), j.MaxLogPosition, j.MaxLogPositionMaxValue())
		}
	}
	if j.MemberIdInActingVersion(actingVersion) {
		if j.MemberId < j.MemberIdMinValue() || j.MemberId > j.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on j.MemberId (%v < %v > %v)", j.MemberIdMinValue(), j.MemberId, j.MemberIdMaxValue())
		}
	}
	if j.LogSessionIdInActingVersion(actingVersion) {
		if j.LogSessionId < j.LogSessionIdMinValue() || j.LogSessionId > j.LogSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on j.LogSessionId (%v < %v > %v)", j.LogSessionIdMinValue(), j.LogSessionId, j.LogSessionIdMaxValue())
		}
	}
	if j.LogStreamIdInActingVersion(actingVersion) {
		if j.LogStreamId < j.LogStreamIdMinValue() || j.LogStreamId > j.LogStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on j.LogStreamId (%v < %v > %v)", j.LogStreamIdMinValue(), j.LogStreamId, j.LogStreamIdMaxValue())
		}
	}
	if err := j.IsStartup.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if j.RoleInActingVersion(actingVersion) {
		if j.Role < j.RoleMinValue() || j.Role > j.RoleMaxValue() {
			return fmt.Errorf("Range check failed on j.Role (%v < %v > %v)", j.RoleMinValue(), j.Role, j.RoleMaxValue())
		}
	}
	for idx, ch := range j.LogChannel {
		if ch > 127 {
			return fmt.Errorf("j.LogChannel[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func JoinLogInit(j *JoinLog) {
	return
}

func (*JoinLog) SbeBlockLength() (blockLength uint16) {
	return 36
}

func (*JoinLog) SbeTemplateId() (templateId uint16) {
	return 40
}

func (*JoinLog) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*JoinLog) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*JoinLog) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*JoinLog) LogPositionId() uint16 {
	return 1
}

func (*JoinLog) LogPositionSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.LogPositionSinceVersion()
}

func (*JoinLog) LogPositionDeprecated() uint16 {
	return 0
}

func (*JoinLog) LogPositionMetaAttribute(meta int) string {
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

func (*JoinLog) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*JoinLog) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*JoinLog) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*JoinLog) MaxLogPositionId() uint16 {
	return 2
}

func (*JoinLog) MaxLogPositionSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) MaxLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.MaxLogPositionSinceVersion()
}

func (*JoinLog) MaxLogPositionDeprecated() uint16 {
	return 0
}

func (*JoinLog) MaxLogPositionMetaAttribute(meta int) string {
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

func (*JoinLog) MaxLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*JoinLog) MaxLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*JoinLog) MaxLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*JoinLog) MemberIdId() uint16 {
	return 3
}

func (*JoinLog) MemberIdSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.MemberIdSinceVersion()
}

func (*JoinLog) MemberIdDeprecated() uint16 {
	return 0
}

func (*JoinLog) MemberIdMetaAttribute(meta int) string {
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

func (*JoinLog) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*JoinLog) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*JoinLog) MemberIdNullValue() int32 {
	return math.MinInt32
}

func (*JoinLog) LogSessionIdId() uint16 {
	return 4
}

func (*JoinLog) LogSessionIdSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) LogSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.LogSessionIdSinceVersion()
}

func (*JoinLog) LogSessionIdDeprecated() uint16 {
	return 0
}

func (*JoinLog) LogSessionIdMetaAttribute(meta int) string {
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

func (*JoinLog) LogSessionIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*JoinLog) LogSessionIdMaxValue() int32 {
	return math.MaxInt32
}

func (*JoinLog) LogSessionIdNullValue() int32 {
	return math.MinInt32
}

func (*JoinLog) LogStreamIdId() uint16 {
	return 5
}

func (*JoinLog) LogStreamIdSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) LogStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.LogStreamIdSinceVersion()
}

func (*JoinLog) LogStreamIdDeprecated() uint16 {
	return 0
}

func (*JoinLog) LogStreamIdMetaAttribute(meta int) string {
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

func (*JoinLog) LogStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*JoinLog) LogStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*JoinLog) LogStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*JoinLog) IsStartupId() uint16 {
	return 6
}

func (*JoinLog) IsStartupSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) IsStartupInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.IsStartupSinceVersion()
}

func (*JoinLog) IsStartupDeprecated() uint16 {
	return 0
}

func (*JoinLog) IsStartupMetaAttribute(meta int) string {
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

func (*JoinLog) RoleId() uint16 {
	return 7
}

func (*JoinLog) RoleSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) RoleInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.RoleSinceVersion()
}

func (*JoinLog) RoleDeprecated() uint16 {
	return 0
}

func (*JoinLog) RoleMetaAttribute(meta int) string {
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

func (*JoinLog) RoleMinValue() int32 {
	return math.MinInt32 + 1
}

func (*JoinLog) RoleMaxValue() int32 {
	return math.MaxInt32
}

func (*JoinLog) RoleNullValue() int32 {
	return math.MinInt32
}

func (*JoinLog) LogChannelMetaAttribute(meta int) string {
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

func (*JoinLog) LogChannelSinceVersion() uint16 {
	return 0
}

func (j *JoinLog) LogChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= j.LogChannelSinceVersion()
}

func (*JoinLog) LogChannelDeprecated() uint16 {
	return 0
}

func (JoinLog) LogChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (JoinLog) LogChannelHeaderLength() uint64 {
	return 4
}
