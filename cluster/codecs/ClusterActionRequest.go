// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClusterActionRequest struct {
	LeadershipTermId int64
	LogPosition      int64
	Timestamp        int64
	Action           ClusterActionEnum
}

func (c *ClusterActionRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, c.Timestamp); err != nil {
		return err
	}
	if err := c.Action.Encode(_m, _w); err != nil {
		return err
	}
	return nil
}

func (c *ClusterActionRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !c.TimestampInActingVersion(actingVersion) {
		c.Timestamp = c.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.Timestamp); err != nil {
			return err
		}
	}
	if c.ActionInActingVersion(actingVersion) {
		if err := c.Action.Decode(_m, _r, actingVersion); err != nil {
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

func (c *ClusterActionRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if c.TimestampInActingVersion(actingVersion) {
		if c.Timestamp < c.TimestampMinValue() || c.Timestamp > c.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on c.Timestamp (%v < %v > %v)", c.TimestampMinValue(), c.Timestamp, c.TimestampMaxValue())
		}
	}
	if err := c.Action.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func ClusterActionRequestInit(c *ClusterActionRequest) {
	return
}

func (*ClusterActionRequest) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*ClusterActionRequest) SbeTemplateId() (templateId uint16) {
	return 23
}

func (*ClusterActionRequest) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClusterActionRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterActionRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClusterActionRequest) LeadershipTermIdId() uint16 {
	return 1
}

func (*ClusterActionRequest) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterActionRequest) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LeadershipTermIdSinceVersion()
}

func (*ClusterActionRequest) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*ClusterActionRequest) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*ClusterActionRequest) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterActionRequest) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterActionRequest) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterActionRequest) LogPositionId() uint16 {
	return 2
}

func (*ClusterActionRequest) LogPositionSinceVersion() uint16 {
	return 0
}

func (c *ClusterActionRequest) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogPositionSinceVersion()
}

func (*ClusterActionRequest) LogPositionDeprecated() uint16 {
	return 0
}

func (*ClusterActionRequest) LogPositionMetaAttribute(meta int) string {
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

func (*ClusterActionRequest) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterActionRequest) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterActionRequest) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*ClusterActionRequest) TimestampId() uint16 {
	return 3
}

func (*ClusterActionRequest) TimestampSinceVersion() uint16 {
	return 0
}

func (c *ClusterActionRequest) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.TimestampSinceVersion()
}

func (*ClusterActionRequest) TimestampDeprecated() uint16 {
	return 0
}

func (*ClusterActionRequest) TimestampMetaAttribute(meta int) string {
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

func (*ClusterActionRequest) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterActionRequest) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterActionRequest) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*ClusterActionRequest) ActionId() uint16 {
	return 4
}

func (*ClusterActionRequest) ActionSinceVersion() uint16 {
	return 0
}

func (c *ClusterActionRequest) ActionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ActionSinceVersion()
}

func (*ClusterActionRequest) ActionDeprecated() uint16 {
	return 0
}

func (*ClusterActionRequest) ActionMetaAttribute(meta int) string {
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
