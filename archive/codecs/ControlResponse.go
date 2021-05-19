// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ControlResponse struct {
	ControlSessionId int64
	CorrelationId    int64
	RelevantId       int64
	Code             ControlResponseCodeEnum
	Version          int32
	ErrorMessage     []uint8
}

func (c *ControlResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.RelevantId); err != nil {
		return err
	}
	if err := c.Code.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.Version); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.ErrorMessage))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.ErrorMessage); err != nil {
		return err
	}
	return nil
}

func (c *ControlResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.ControlSessionIdInActingVersion(actingVersion) {
		c.ControlSessionId = c.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.ControlSessionId); err != nil {
			return err
		}
	}
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
			return err
		}
	}
	if !c.RelevantIdInActingVersion(actingVersion) {
		c.RelevantId = c.RelevantIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.RelevantId); err != nil {
			return err
		}
	}
	if c.CodeInActingVersion(actingVersion) {
		if err := c.Code.Decode(_m, _r, actingVersion); err != nil {
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

	if c.ErrorMessageInActingVersion(actingVersion) {
		var ErrorMessageLength uint32
		if err := _m.ReadUint32(_r, &ErrorMessageLength); err != nil {
			return err
		}
		if cap(c.ErrorMessage) < int(ErrorMessageLength) {
			c.ErrorMessage = make([]uint8, ErrorMessageLength)
		}
		c.ErrorMessage = c.ErrorMessage[:ErrorMessageLength]
		if err := _m.ReadBytes(_r, c.ErrorMessage); err != nil {
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

func (c *ControlResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.ControlSessionIdInActingVersion(actingVersion) {
		if c.ControlSessionId < c.ControlSessionIdMinValue() || c.ControlSessionId > c.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ControlSessionId (%v < %v > %v)", c.ControlSessionIdMinValue(), c.ControlSessionId, c.ControlSessionIdMaxValue())
		}
	}
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if c.RelevantIdInActingVersion(actingVersion) {
		if c.RelevantId < c.RelevantIdMinValue() || c.RelevantId > c.RelevantIdMaxValue() {
			return fmt.Errorf("Range check failed on c.RelevantId (%v < %v > %v)", c.RelevantIdMinValue(), c.RelevantId, c.RelevantIdMaxValue())
		}
	}
	if err := c.Code.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if c.VersionInActingVersion(actingVersion) {
		if c.Version != c.VersionNullValue() && (c.Version < c.VersionMinValue() || c.Version > c.VersionMaxValue()) {
			return fmt.Errorf("Range check failed on c.Version (%v < %v > %v)", c.VersionMinValue(), c.Version, c.VersionMaxValue())
		}
	}
	return nil
}

func ControlResponseInit(c *ControlResponse) {
	c.Version = 0
	return
}

func (*ControlResponse) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*ControlResponse) SbeTemplateId() (templateId uint16) {
	return 1
}

func (*ControlResponse) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ControlResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ControlResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ControlResponse) ControlSessionIdId() uint16 {
	return 1
}

func (*ControlResponse) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (c *ControlResponse) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ControlSessionIdSinceVersion()
}

func (*ControlResponse) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ControlResponse) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ControlResponse) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ControlResponse) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ControlResponse) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ControlResponse) CorrelationIdId() uint16 {
	return 2
}

func (*ControlResponse) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ControlResponse) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ControlResponse) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ControlResponse) CorrelationIdMetaAttribute(meta int) string {
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

func (*ControlResponse) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ControlResponse) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ControlResponse) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ControlResponse) RelevantIdId() uint16 {
	return 3
}

func (*ControlResponse) RelevantIdSinceVersion() uint16 {
	return 0
}

func (c *ControlResponse) RelevantIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.RelevantIdSinceVersion()
}

func (*ControlResponse) RelevantIdDeprecated() uint16 {
	return 0
}

func (*ControlResponse) RelevantIdMetaAttribute(meta int) string {
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

func (*ControlResponse) RelevantIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ControlResponse) RelevantIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ControlResponse) RelevantIdNullValue() int64 {
	return math.MinInt64
}

func (*ControlResponse) CodeId() uint16 {
	return 4
}

func (*ControlResponse) CodeSinceVersion() uint16 {
	return 0
}

func (c *ControlResponse) CodeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CodeSinceVersion()
}

func (*ControlResponse) CodeDeprecated() uint16 {
	return 0
}

func (*ControlResponse) CodeMetaAttribute(meta int) string {
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

func (*ControlResponse) VersionId() uint16 {
	return 5
}

func (*ControlResponse) VersionSinceVersion() uint16 {
	return 4
}

func (c *ControlResponse) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.VersionSinceVersion()
}

func (*ControlResponse) VersionDeprecated() uint16 {
	return 0
}

func (*ControlResponse) VersionMetaAttribute(meta int) string {
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

func (*ControlResponse) VersionMinValue() int32 {
	return 2
}

func (*ControlResponse) VersionMaxValue() int32 {
	return 16777215
}

func (*ControlResponse) VersionNullValue() int32 {
	return 0
}

func (*ControlResponse) ErrorMessageMetaAttribute(meta int) string {
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

func (*ControlResponse) ErrorMessageSinceVersion() uint16 {
	return 0
}

func (c *ControlResponse) ErrorMessageInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ErrorMessageSinceVersion()
}

func (*ControlResponse) ErrorMessageDeprecated() uint16 {
	return 0
}

func (ControlResponse) ErrorMessageCharacterEncoding() string {
	return "US-ASCII"
}

func (ControlResponse) ErrorMessageHeaderLength() uint64 {
	return 4
}
