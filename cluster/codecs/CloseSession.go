// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CloseSession struct {
	ClusterSessionId int64
}

func (c *CloseSession) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.ClusterSessionId); err != nil {
		return err
	}
	return nil
}

func (c *CloseSession) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.ClusterSessionIdInActingVersion(actingVersion) {
		c.ClusterSessionId = c.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.ClusterSessionId); err != nil {
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

func (c *CloseSession) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.ClusterSessionIdInActingVersion(actingVersion) {
		if c.ClusterSessionId < c.ClusterSessionIdMinValue() || c.ClusterSessionId > c.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ClusterSessionId (%v < %v > %v)", c.ClusterSessionIdMinValue(), c.ClusterSessionId, c.ClusterSessionIdMaxValue())
		}
	}
	return nil
}

func CloseSessionInit(c *CloseSession) {
	return
}

func (*CloseSession) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*CloseSession) SbeTemplateId() (templateId uint16) {
	return 30
}

func (*CloseSession) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*CloseSession) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*CloseSession) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CloseSession) ClusterSessionIdId() uint16 {
	return 1
}

func (*CloseSession) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (c *CloseSession) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ClusterSessionIdSinceVersion()
}

func (*CloseSession) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*CloseSession) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*CloseSession) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CloseSession) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CloseSession) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}
