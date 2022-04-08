// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CommitPosition struct {
	LeadershipTermId int64
	LogPosition      int64
	LeaderMemberId   int32
}

func (c *CommitPosition) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.LeaderMemberId); err != nil {
		return err
	}
	return nil
}

func (c *CommitPosition) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !c.LeaderMemberIdInActingVersion(actingVersion) {
		c.LeaderMemberId = c.LeaderMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.LeaderMemberId); err != nil {
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

func (c *CommitPosition) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if c.LeaderMemberIdInActingVersion(actingVersion) {
		if c.LeaderMemberId < c.LeaderMemberIdMinValue() || c.LeaderMemberId > c.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.LeaderMemberId (%v < %v > %v)", c.LeaderMemberIdMinValue(), c.LeaderMemberId, c.LeaderMemberIdMaxValue())
		}
	}
	return nil
}

func CommitPositionInit(c *CommitPosition) {
	return
}

func (*CommitPosition) SbeBlockLength() (blockLength uint16) {
	return 20
}

func (*CommitPosition) SbeTemplateId() (templateId uint16) {
	return 55
}

func (*CommitPosition) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*CommitPosition) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*CommitPosition) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CommitPosition) LeadershipTermIdId() uint16 {
	return 1
}

func (*CommitPosition) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (c *CommitPosition) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeadershipTermIdSinceVersion()
}

func (*CommitPosition) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*CommitPosition) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*CommitPosition) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CommitPosition) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CommitPosition) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*CommitPosition) LogPositionId() uint16 {
	return 2
}

func (*CommitPosition) LogPositionSinceVersion() uint16 {
	return 0
}

func (c *CommitPosition) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogPositionSinceVersion()
}

func (*CommitPosition) LogPositionDeprecated() uint16 {
	return 0
}

func (*CommitPosition) LogPositionMetaAttribute(meta int) string {
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

func (*CommitPosition) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CommitPosition) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*CommitPosition) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*CommitPosition) LeaderMemberIdId() uint16 {
	return 3
}

func (*CommitPosition) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (c *CommitPosition) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeaderMemberIdSinceVersion()
}

func (*CommitPosition) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*CommitPosition) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*CommitPosition) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*CommitPosition) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*CommitPosition) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}
