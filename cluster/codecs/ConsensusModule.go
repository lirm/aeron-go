// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ConsensusModule struct {
	NextSessionId          int64
	NextServiceSessionId   int64
	LogServiceSessionId    int64
	PendingMessageCapacity int32
}

func (c *ConsensusModule) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.NextSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.NextServiceSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.LogServiceSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.PendingMessageCapacity); err != nil {
		return err
	}
	return nil
}

func (c *ConsensusModule) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.NextSessionIdInActingVersion(actingVersion) {
		c.NextSessionId = c.NextSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.NextSessionId); err != nil {
			return err
		}
	}
	if !c.NextServiceSessionIdInActingVersion(actingVersion) {
		c.NextServiceSessionId = c.NextServiceSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.NextServiceSessionId); err != nil {
			return err
		}
	}
	if !c.LogServiceSessionIdInActingVersion(actingVersion) {
		c.LogServiceSessionId = c.LogServiceSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.LogServiceSessionId); err != nil {
			return err
		}
	}
	if !c.PendingMessageCapacityInActingVersion(actingVersion) {
		c.PendingMessageCapacity = c.PendingMessageCapacityNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.PendingMessageCapacity); err != nil {
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

func (c *ConsensusModule) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.NextSessionIdInActingVersion(actingVersion) {
		if c.NextSessionId < c.NextSessionIdMinValue() || c.NextSessionId > c.NextSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.NextSessionId (%v < %v > %v)", c.NextSessionIdMinValue(), c.NextSessionId, c.NextSessionIdMaxValue())
		}
	}
	if c.NextServiceSessionIdInActingVersion(actingVersion) {
		if c.NextServiceSessionId != c.NextServiceSessionIdNullValue() && (c.NextServiceSessionId < c.NextServiceSessionIdMinValue() || c.NextServiceSessionId > c.NextServiceSessionIdMaxValue()) {
			return fmt.Errorf("Range check failed on c.NextServiceSessionId (%v < %v > %v)", c.NextServiceSessionIdMinValue(), c.NextServiceSessionId, c.NextServiceSessionIdMaxValue())
		}
	}
	if c.LogServiceSessionIdInActingVersion(actingVersion) {
		if c.LogServiceSessionId != c.LogServiceSessionIdNullValue() && (c.LogServiceSessionId < c.LogServiceSessionIdMinValue() || c.LogServiceSessionId > c.LogServiceSessionIdMaxValue()) {
			return fmt.Errorf("Range check failed on c.LogServiceSessionId (%v < %v > %v)", c.LogServiceSessionIdMinValue(), c.LogServiceSessionId, c.LogServiceSessionIdMaxValue())
		}
	}
	if c.PendingMessageCapacityInActingVersion(actingVersion) {
		if c.PendingMessageCapacity != c.PendingMessageCapacityNullValue() && (c.PendingMessageCapacity < c.PendingMessageCapacityMinValue() || c.PendingMessageCapacity > c.PendingMessageCapacityMaxValue()) {
			return fmt.Errorf("Range check failed on c.PendingMessageCapacity (%v < %v > %v)", c.PendingMessageCapacityMinValue(), c.PendingMessageCapacity, c.PendingMessageCapacityMaxValue())
		}
	}
	return nil
}

func ConsensusModuleInit(c *ConsensusModule) {
	c.NextServiceSessionId = math.MinInt64
	c.LogServiceSessionId = math.MinInt64
	c.PendingMessageCapacity = 0
	return
}

func (*ConsensusModule) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*ConsensusModule) SbeTemplateId() (templateId uint16) {
	return 105
}

func (*ConsensusModule) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ConsensusModule) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ConsensusModule) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ConsensusModule) NextSessionIdId() uint16 {
	return 1
}

func (*ConsensusModule) NextSessionIdSinceVersion() uint16 {
	return 0
}

func (c *ConsensusModule) NextSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.NextSessionIdSinceVersion()
}

func (*ConsensusModule) NextSessionIdDeprecated() uint16 {
	return 0
}

func (*ConsensusModule) NextSessionIdMetaAttribute(meta int) string {
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

func (*ConsensusModule) NextSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ConsensusModule) NextSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ConsensusModule) NextSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ConsensusModule) NextServiceSessionIdId() uint16 {
	return 2
}

func (*ConsensusModule) NextServiceSessionIdSinceVersion() uint16 {
	return 3
}

func (c *ConsensusModule) NextServiceSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.NextServiceSessionIdSinceVersion()
}

func (*ConsensusModule) NextServiceSessionIdDeprecated() uint16 {
	return 0
}

func (*ConsensusModule) NextServiceSessionIdMetaAttribute(meta int) string {
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

func (*ConsensusModule) NextServiceSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ConsensusModule) NextServiceSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ConsensusModule) NextServiceSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ConsensusModule) LogServiceSessionIdId() uint16 {
	return 3
}

func (*ConsensusModule) LogServiceSessionIdSinceVersion() uint16 {
	return 3
}

func (c *ConsensusModule) LogServiceSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LogServiceSessionIdSinceVersion()
}

func (*ConsensusModule) LogServiceSessionIdDeprecated() uint16 {
	return 0
}

func (*ConsensusModule) LogServiceSessionIdMetaAttribute(meta int) string {
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

func (*ConsensusModule) LogServiceSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ConsensusModule) LogServiceSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ConsensusModule) LogServiceSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ConsensusModule) PendingMessageCapacityId() uint16 {
	return 4
}

func (*ConsensusModule) PendingMessageCapacitySinceVersion() uint16 {
	return 3
}

func (c *ConsensusModule) PendingMessageCapacityInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.PendingMessageCapacitySinceVersion()
}

func (*ConsensusModule) PendingMessageCapacityDeprecated() uint16 {
	return 0
}

func (*ConsensusModule) PendingMessageCapacityMetaAttribute(meta int) string {
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

func (*ConsensusModule) PendingMessageCapacityMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ConsensusModule) PendingMessageCapacityMaxValue() int32 {
	return math.MaxInt32
}

func (*ConsensusModule) PendingMessageCapacityNullValue() int32 {
	return 0
}
