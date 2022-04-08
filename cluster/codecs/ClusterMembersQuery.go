// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClusterMembersQuery struct {
	CorrelationId int64
	Extended      BooleanTypeEnum
}

func (c *ClusterMembersQuery) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := c.Extended.Encode(_m, _w); err != nil {
		return err
	}
	return nil
}

func (c *ClusterMembersQuery) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
			return err
		}
	}
	if c.ExtendedInActingVersion(actingVersion) {
		if err := c.Extended.Decode(_m, _r, actingVersion); err != nil {
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

func (c *ClusterMembersQuery) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if err := c.Extended.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func ClusterMembersQueryInit(c *ClusterMembersQuery) {
	return
}

func (*ClusterMembersQuery) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*ClusterMembersQuery) SbeTemplateId() (templateId uint16) {
	return 34
}

func (*ClusterMembersQuery) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClusterMembersQuery) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterMembersQuery) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClusterMembersQuery) CorrelationIdId() uint16 {
	return 1
}

func (*ClusterMembersQuery) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterMembersQuery) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ClusterMembersQuery) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ClusterMembersQuery) CorrelationIdMetaAttribute(meta int) string {
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

func (*ClusterMembersQuery) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterMembersQuery) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterMembersQuery) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterMembersQuery) ExtendedId() uint16 {
	return 2
}

func (*ClusterMembersQuery) ExtendedSinceVersion() uint16 {
	return 5
}

func (c *ClusterMembersQuery) ExtendedInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ExtendedSinceVersion()
}

func (*ClusterMembersQuery) ExtendedDeprecated() uint16 {
	return 0
}

func (*ClusterMembersQuery) ExtendedMetaAttribute(meta int) string {
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
