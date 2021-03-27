// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ConnectRequest struct {
	CorrelationId    int64
	ResponseStreamId int32
	Version          int32
	ResponseChannel  []uint8
}

func (c *ConnectRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.ResponseStreamId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.Version); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ResponseChannel); err != nil {
		return err
	}
	return nil
}

func (c *ConnectRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
			return err
		}
	}
	if !c.ResponseStreamIdInActingVersion(actingVersion) {
		c.ResponseStreamId = c.ResponseStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.ResponseStreamId); err != nil {
			return err
		}
	}
	if !c.VersionInActingVersion(actingVersion) {
		c.Version = c.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.Version); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.ResponseChannelInActingVersion(actingVersion) {
		var ResponseChannelLength uint32
		if err := _m.ReadUint32(_r, &ResponseChannelLength); err != nil {
			return err
		}
		if cap(c.ResponseChannel) < int(ResponseChannelLength) {
			c.ResponseChannel = make([]uint8, ResponseChannelLength)
		}
		c.ResponseChannel = c.ResponseChannel[:ResponseChannelLength]
		if err := _m.ReadBytes(_r, c.ResponseChannel); err != nil {
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

func (c *ConnectRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if c.ResponseStreamIdInActingVersion(actingVersion) {
		if c.ResponseStreamId < c.ResponseStreamIdMinValue() || c.ResponseStreamId > c.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ResponseStreamId (%v < %v > %v)", c.ResponseStreamIdMinValue(), c.ResponseStreamId, c.ResponseStreamIdMaxValue())
		}
	}
	if c.VersionInActingVersion(actingVersion) {
		if c.Version != c.VersionNullValue() && (c.Version < c.VersionMinValue() || c.Version > c.VersionMaxValue()) {
			return fmt.Errorf("Range check failed on c.Version (%v < %v > %v)", c.VersionMinValue(), c.Version, c.VersionMaxValue())
		}
	}
	return nil
}

func ConnectRequestInit(c *ConnectRequest) {
	c.Version = 0
	return
}

func (*ConnectRequest) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*ConnectRequest) SbeTemplateId() (templateId uint16) {
	return 2
}

func (*ConnectRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ConnectRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ConnectRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ConnectRequest) CorrelationIdId() uint16 {
	return 1
}

func (*ConnectRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ConnectRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ConnectRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ConnectRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ConnectRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ConnectRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ConnectRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ConnectRequest) ResponseStreamIdId() uint16 {
	return 2
}

func (*ConnectRequest) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (c *ConnectRequest) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ResponseStreamIdSinceVersion()
}

func (*ConnectRequest) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*ConnectRequest) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*ConnectRequest) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ConnectRequest) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ConnectRequest) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ConnectRequest) VersionId() uint16 {
	return 3
}

func (*ConnectRequest) VersionSinceVersion() uint16 {
	return 2
}

func (c *ConnectRequest) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.VersionSinceVersion()
}

func (*ConnectRequest) VersionDeprecated() uint16 {
	return 0
}

func (*ConnectRequest) VersionMetaAttribute(meta int) string {
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

func (*ConnectRequest) VersionMinValue() int32 {
	return 2
}

func (*ConnectRequest) VersionMaxValue() int32 {
	return 16777215
}

func (*ConnectRequest) VersionNullValue() int32 {
	return 0
}

func (*ConnectRequest) ResponseChannelMetaAttribute(meta int) string {
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

func (*ConnectRequest) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (c *ConnectRequest) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ResponseChannelSinceVersion()
}

func (*ConnectRequest) ResponseChannelDeprecated() uint16 {
	return 0
}

func (ConnectRequest) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ConnectRequest) ResponseChannelHeaderLength() uint64 {
	return 4
}
