// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClientSession struct {
	ClusterSessionId int64
	ResponseStreamId int32
	ResponseChannel  []uint8
	EncodedPrincipal []uint8
}

func (c *ClientSession) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.ResponseStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ResponseChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.EncodedPrincipal))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.EncodedPrincipal); err != nil {
		return err
	}
	return nil
}

func (c *ClientSession) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.ClusterSessionIdInActingVersion(actingVersion) {
		c.ClusterSessionId = c.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.ClusterSessionId); err != nil {
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

	if c.EncodedPrincipalInActingVersion(actingVersion) {
		var EncodedPrincipalLength uint32
		if err := _m.ReadUint32(_r, &EncodedPrincipalLength); err != nil {
			return err
		}
		if cap(c.EncodedPrincipal) < int(EncodedPrincipalLength) {
			c.EncodedPrincipal = make([]uint8, EncodedPrincipalLength)
		}
		c.EncodedPrincipal = c.EncodedPrincipal[:EncodedPrincipalLength]
		if err := _m.ReadBytes(_r, c.EncodedPrincipal); err != nil {
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

func (c *ClientSession) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.ClusterSessionIdInActingVersion(actingVersion) {
		if c.ClusterSessionId < c.ClusterSessionIdMinValue() || c.ClusterSessionId > c.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ClusterSessionId (%v < %v > %v)", c.ClusterSessionIdMinValue(), c.ClusterSessionId, c.ClusterSessionIdMaxValue())
		}
	}
	if c.ResponseStreamIdInActingVersion(actingVersion) {
		if c.ResponseStreamId < c.ResponseStreamIdMinValue() || c.ResponseStreamId > c.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ResponseStreamId (%v < %v > %v)", c.ResponseStreamIdMinValue(), c.ResponseStreamId, c.ResponseStreamIdMaxValue())
		}
	}
	for idx, ch := range c.ResponseChannel {
		if ch > 127 {
			return fmt.Errorf("c.ResponseChannel[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func ClientSessionInit(c *ClientSession) {
	return
}

func (*ClientSession) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*ClientSession) SbeTemplateId() (templateId uint16) {
	return 102
}

func (*ClientSession) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClientSession) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClientSession) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClientSession) ClusterSessionIdId() uint16 {
	return 1
}

func (*ClientSession) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (c *ClientSession) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ClusterSessionIdSinceVersion()
}

func (*ClientSession) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*ClientSession) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*ClientSession) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClientSession) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClientSession) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ClientSession) ResponseStreamIdId() uint16 {
	return 2
}

func (*ClientSession) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (c *ClientSession) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ResponseStreamIdSinceVersion()
}

func (*ClientSession) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*ClientSession) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*ClientSession) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClientSession) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClientSession) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ClientSession) ResponseChannelMetaAttribute(meta int) string {
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

func (*ClientSession) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (c *ClientSession) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ResponseChannelSinceVersion()
}

func (*ClientSession) ResponseChannelDeprecated() uint16 {
	return 0
}

func (ClientSession) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ClientSession) ResponseChannelHeaderLength() uint64 {
	return 4
}

func (*ClientSession) EncodedPrincipalMetaAttribute(meta int) string {
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

func (*ClientSession) EncodedPrincipalSinceVersion() uint16 {
	return 0
}

func (c *ClientSession) EncodedPrincipalInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.EncodedPrincipalSinceVersion()
}

func (*ClientSession) EncodedPrincipalDeprecated() uint16 {
	return 0
}

func (ClientSession) EncodedPrincipalCharacterEncoding() string {
	return "null"
}

func (ClientSession) EncodedPrincipalHeaderLength() uint64 {
	return 4
}
