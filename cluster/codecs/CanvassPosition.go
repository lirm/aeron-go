// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CanvassPosition struct {
	LogLeadershipTermId int64
	LogPosition         int64
	LeadershipTermId    int64
	FollowerMemberId    int32
}

func (c *CanvassPosition) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.LogLeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.FollowerMemberId); err != nil {
		return err
	}
	return nil
}

func (c *CanvassPosition) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.LogLeadershipTermIdInActingVersion(actingVersion) {
		c.LogLeadershipTermId = c.LogLeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.LogLeadershipTermId); err != nil {
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
	if !c.LeadershipTermIdInActingVersion(actingVersion) {
		c.LeadershipTermId = c.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.LeadershipTermId); err != nil {
			return err
		}
	}
	if !c.FollowerMemberIdInActingVersion(actingVersion) {
		c.FollowerMemberId = c.FollowerMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.FollowerMemberId); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *CanvassPosition) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.LogLeadershipTermIdInActingVersion(actingVersion) {
		if c.LogLeadershipTermId < c.LogLeadershipTermIdMinValue() || c.LogLeadershipTermId > c.LogLeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on c.LogLeadershipTermId (%v < %v > %v)", c.LogLeadershipTermIdMinValue(), c.LogLeadershipTermId, c.LogLeadershipTermIdMaxValue())
		}
	}
	if c.LogPositionInActingVersion(actingVersion) {
		if c.LogPosition < c.LogPositionMinValue() || c.LogPosition > c.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on c.LogPosition (%v < %v > %v)", c.LogPositionMinValue(), c.LogPosition, c.LogPositionMaxValue())
		}
	}
	if c.LeadershipTermIdInActingVersion(actingVersion) {
		if c.LeadershipTermId < c.LeadershipTermIdMinValue() || c.LeadershipTermId > c.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on c.LeadershipTermId (%v < %v > %v)", c.LeadershipTermIdMinValue(), c.LeadershipTermId, c.LeadershipTermIdMaxValue())
		}
	}
	if c.FollowerMemberIdInActingVersion(actingVersion) {
		if c.FollowerMemberId < c.FollowerMemberIdMinValue() || c.FollowerMemberId > c.FollowerMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.FollowerMemberId (%v < %v > %v)", c.FollowerMemberIdMinValue(), c.FollowerMemberId, c.FollowerMemberIdMaxValue())
		}
	}
	return nil
}

func CanvassPositionInit(c *CanvassPosition) {
	return
}

func (*CanvassPosition) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*CanvassPosition) SbeTemplateId() (templateId uint16) {
	return 50
}

func (*CanvassPosition) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*CanvassPosition) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*CanvassPosition) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CanvassPosition) LogLeadershipTermIdId() uint16 {
	return 1
}

func (*CanvassPosition) LogLeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (c *CanvassPosition) LogLeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogLeadershipTermIdSinceVersion()
}

func (*CanvassPosition) LogLeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*CanvassPosition) LogLeadershipTermIdMetaAttribute(meta int) string {
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

func (*CanvassPosition) LogLeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CanvassPosition) LogLeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CanvassPosition) LogLeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*CanvassPosition) LogPositionId() uint16 {
	return 2
}

func (*CanvassPosition) LogPositionSinceVersion() uint16 {
	return 0
}

func (c *CanvassPosition) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogPositionSinceVersion()
}

func (*CanvassPosition) LogPositionDeprecated() uint16 {
	return 0
}

func (*CanvassPosition) LogPositionMetaAttribute(meta int) string {
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

func (*CanvassPosition) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CanvassPosition) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*CanvassPosition) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*CanvassPosition) LeadershipTermIdId() uint16 {
	return 3
}

func (*CanvassPosition) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (c *CanvassPosition) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeadershipTermIdSinceVersion()
}

func (*CanvassPosition) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*CanvassPosition) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*CanvassPosition) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CanvassPosition) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CanvassPosition) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*CanvassPosition) FollowerMemberIdId() uint16 {
	return 4
}

func (*CanvassPosition) FollowerMemberIdSinceVersion() uint16 {
	return 0
}

func (c *CanvassPosition) FollowerMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.FollowerMemberIdSinceVersion()
}

func (*CanvassPosition) FollowerMemberIdDeprecated() uint16 {
	return 0
}

func (*CanvassPosition) FollowerMemberIdMetaAttribute(meta int) string {
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

func (*CanvassPosition) FollowerMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*CanvassPosition) FollowerMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*CanvassPosition) FollowerMemberIdNullValue() int32 {
	return math.MinInt32
}
