// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CatchupPosition struct {
	LeadershipTermId int64
	LogPosition      int64
	FollowerMemberId int32
	CatchupEndpoint  []uint8
}

func (c *CatchupPosition) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt32(_w, c.FollowerMemberId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.CatchupEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.CatchupEndpoint); err != nil {
		return err
	}
	return nil
}

func (c *CatchupPosition) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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

	if c.CatchupEndpointInActingVersion(actingVersion) {
		var CatchupEndpointLength uint32
		if err := _m.ReadUint32(_r, &CatchupEndpointLength); err != nil {
			return err
		}
		if cap(c.CatchupEndpoint) < int(CatchupEndpointLength) {
			c.CatchupEndpoint = make([]uint8, CatchupEndpointLength)
		}
		c.CatchupEndpoint = c.CatchupEndpoint[:CatchupEndpointLength]
		if err := _m.ReadBytes(_r, c.CatchupEndpoint); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *CatchupPosition) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if c.FollowerMemberIdInActingVersion(actingVersion) {
		if c.FollowerMemberId < c.FollowerMemberIdMinValue() || c.FollowerMemberId > c.FollowerMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on c.FollowerMemberId (%v < %v > %v)", c.FollowerMemberIdMinValue(), c.FollowerMemberId, c.FollowerMemberIdMaxValue())
		}
	}
	for idx, ch := range c.CatchupEndpoint {
		if ch > 127 {
			return fmt.Errorf("c.CatchupEndpoint[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func CatchupPositionInit(c *CatchupPosition) {
	return
}

func (*CatchupPosition) SbeBlockLength() (blockLength uint16) {
	return 20
}

func (*CatchupPosition) SbeTemplateId() (templateId uint16) {
	return 56
}

func (*CatchupPosition) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*CatchupPosition) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*CatchupPosition) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CatchupPosition) LeadershipTermIdId() uint16 {
	return 1
}

func (*CatchupPosition) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (c *CatchupPosition) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeadershipTermIdSinceVersion()
}

func (*CatchupPosition) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*CatchupPosition) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*CatchupPosition) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CatchupPosition) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CatchupPosition) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*CatchupPosition) LogPositionId() uint16 {
	return 2
}

func (*CatchupPosition) LogPositionSinceVersion() uint16 {
	return 0
}

func (c *CatchupPosition) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogPositionSinceVersion()
}

func (*CatchupPosition) LogPositionDeprecated() uint16 {
	return 0
}

func (*CatchupPosition) LogPositionMetaAttribute(meta int) string {
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

func (*CatchupPosition) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CatchupPosition) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*CatchupPosition) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*CatchupPosition) FollowerMemberIdId() uint16 {
	return 3
}

func (*CatchupPosition) FollowerMemberIdSinceVersion() uint16 {
	return 0
}

func (c *CatchupPosition) FollowerMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.FollowerMemberIdSinceVersion()
}

func (*CatchupPosition) FollowerMemberIdDeprecated() uint16 {
	return 0
}

func (*CatchupPosition) FollowerMemberIdMetaAttribute(meta int) string {
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

func (*CatchupPosition) FollowerMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*CatchupPosition) FollowerMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*CatchupPosition) FollowerMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*CatchupPosition) CatchupEndpointMetaAttribute(meta int) string {
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

func (*CatchupPosition) CatchupEndpointSinceVersion() uint16 {
	return 0
}

func (c *CatchupPosition) CatchupEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CatchupEndpointSinceVersion()
}

func (*CatchupPosition) CatchupEndpointDeprecated() uint16 {
	return 0
}

func (CatchupPosition) CatchupEndpointCharacterEncoding() string {
	return "US-ASCII"
}

func (CatchupPosition) CatchupEndpointHeaderLength() uint64 {
	return 4
}
