// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CloseSessionRequest struct {
	ControlSessionId int64
}

func (c *CloseSessionRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.ControlSessionId); err != nil {
		return err
	}
	return nil
}

func (c *CloseSessionRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.ControlSessionIdInActingVersion(actingVersion) {
		c.ControlSessionId = c.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.ControlSessionId); err != nil {
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

func (c *CloseSessionRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.ControlSessionIdInActingVersion(actingVersion) {
		if c.ControlSessionId < c.ControlSessionIdMinValue() || c.ControlSessionId > c.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ControlSessionId (%v < %v > %v)", c.ControlSessionIdMinValue(), c.ControlSessionId, c.ControlSessionIdMaxValue())
		}
	}
	return nil
}

func CloseSessionRequestInit(c *CloseSessionRequest) {
	return
}

func (*CloseSessionRequest) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*CloseSessionRequest) SbeTemplateId() (templateId uint16) {
	return 3
}

func (*CloseSessionRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*CloseSessionRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*CloseSessionRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CloseSessionRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*CloseSessionRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (c *CloseSessionRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ControlSessionIdSinceVersion()
}

func (*CloseSessionRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*CloseSessionRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*CloseSessionRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CloseSessionRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CloseSessionRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}
