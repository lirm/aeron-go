// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CancelTimer struct {
	CorrelationId int64
}

func (c *CancelTimer) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	return nil
}

func (c *CancelTimer) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
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

func (c *CancelTimer) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	return nil
}

func CancelTimerInit(c *CancelTimer) {
	return
}

func (*CancelTimer) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*CancelTimer) SbeTemplateId() (templateId uint16) {
	return 32
}

func (*CancelTimer) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*CancelTimer) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*CancelTimer) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CancelTimer) CorrelationIdId() uint16 {
	return 1
}

func (*CancelTimer) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *CancelTimer) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*CancelTimer) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*CancelTimer) CorrelationIdMetaAttribute(meta int) string {
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

func (*CancelTimer) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CancelTimer) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CancelTimer) CorrelationIdNullValue() int64 {
	return math.MinInt64
}
