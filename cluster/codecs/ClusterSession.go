// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ClusterSession struct {
	ClusterSessionId   int64
	CorrelationId      int64
	OpenedLogPosition  int64
	TimeOfLastActivity int64
	CloseReason        CloseReasonEnum
	ResponseStreamId   int32
	ResponseChannel    []uint8
}

func (c *ClusterSession) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.OpenedLogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.TimeOfLastActivity); err != nil {
		return err
	}
	if err := c.CloseReason.Encode(_m, _w); err != nil {
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
	return nil
}

func (c *ClusterSession) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.ClusterSessionIdInActingVersion(actingVersion) {
		c.ClusterSessionId = c.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.ClusterSessionId); err != nil {
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
	if !c.OpenedLogPositionInActingVersion(actingVersion) {
		c.OpenedLogPosition = c.OpenedLogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.OpenedLogPosition); err != nil {
			return err
		}
	}
	if !c.TimeOfLastActivityInActingVersion(actingVersion) {
		c.TimeOfLastActivity = c.TimeOfLastActivityNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.TimeOfLastActivity); err != nil {
			return err
		}
	}
	if c.CloseReasonInActingVersion(actingVersion) {
		if err := c.CloseReason.Decode(_m, _r, actingVersion); err != nil {
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
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterSession) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.ClusterSessionIdInActingVersion(actingVersion) {
		if c.ClusterSessionId < c.ClusterSessionIdMinValue() || c.ClusterSessionId > c.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ClusterSessionId (%v < %v > %v)", c.ClusterSessionIdMinValue(), c.ClusterSessionId, c.ClusterSessionIdMaxValue())
		}
	}
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if c.OpenedLogPositionInActingVersion(actingVersion) {
		if c.OpenedLogPosition < c.OpenedLogPositionMinValue() || c.OpenedLogPosition > c.OpenedLogPositionMaxValue() {
			return fmt.Errorf("Range check failed on c.OpenedLogPosition (%v < %v > %v)", c.OpenedLogPositionMinValue(), c.OpenedLogPosition, c.OpenedLogPositionMaxValue())
		}
	}
	if c.TimeOfLastActivityInActingVersion(actingVersion) {
		if c.TimeOfLastActivity < c.TimeOfLastActivityMinValue() || c.TimeOfLastActivity > c.TimeOfLastActivityMaxValue() {
			return fmt.Errorf("Range check failed on c.TimeOfLastActivity (%v < %v > %v)", c.TimeOfLastActivityMinValue(), c.TimeOfLastActivity, c.TimeOfLastActivityMaxValue())
		}
	}
	if err := c.CloseReason.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
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

func ClusterSessionInit(c *ClusterSession) {
	return
}

func (*ClusterSession) SbeBlockLength() (blockLength uint16) {
	return 40
}

func (*ClusterSession) SbeTemplateId() (templateId uint16) {
	return 103
}

func (*ClusterSession) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ClusterSession) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ClusterSession) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ClusterSession) ClusterSessionIdId() uint16 {
	return 1
}

func (*ClusterSession) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterSession) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ClusterSessionIdSinceVersion()
}

func (*ClusterSession) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*ClusterSession) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*ClusterSession) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterSession) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterSession) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterSession) CorrelationIdId() uint16 {
	return 2
}

func (*ClusterSession) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterSession) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ClusterSession) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ClusterSession) CorrelationIdMetaAttribute(meta int) string {
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

func (*ClusterSession) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterSession) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterSession) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ClusterSession) OpenedLogPositionId() uint16 {
	return 3
}

func (*ClusterSession) OpenedLogPositionSinceVersion() uint16 {
	return 0
}

func (c *ClusterSession) OpenedLogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.OpenedLogPositionSinceVersion()
}

func (*ClusterSession) OpenedLogPositionDeprecated() uint16 {
	return 0
}

func (*ClusterSession) OpenedLogPositionMetaAttribute(meta int) string {
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

func (*ClusterSession) OpenedLogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterSession) OpenedLogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterSession) OpenedLogPositionNullValue() int64 {
	return math.MinInt64
}

func (*ClusterSession) TimeOfLastActivityId() uint16 {
	return 4
}

func (*ClusterSession) TimeOfLastActivitySinceVersion() uint16 {
	return 0
}

func (c *ClusterSession) TimeOfLastActivityInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.TimeOfLastActivitySinceVersion()
}

func (*ClusterSession) TimeOfLastActivityDeprecated() uint16 {
	return 0
}

func (*ClusterSession) TimeOfLastActivityMetaAttribute(meta int) string {
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

func (*ClusterSession) TimeOfLastActivityMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ClusterSession) TimeOfLastActivityMaxValue() int64 {
	return math.MaxInt64
}

func (*ClusterSession) TimeOfLastActivityNullValue() int64 {
	return math.MinInt64
}

func (*ClusterSession) CloseReasonId() uint16 {
	return 5
}

func (*ClusterSession) CloseReasonSinceVersion() uint16 {
	return 0
}

func (c *ClusterSession) CloseReasonInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CloseReasonSinceVersion()
}

func (*ClusterSession) CloseReasonDeprecated() uint16 {
	return 0
}

func (*ClusterSession) CloseReasonMetaAttribute(meta int) string {
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

func (*ClusterSession) ResponseStreamIdId() uint16 {
	return 6
}

func (*ClusterSession) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (c *ClusterSession) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ResponseStreamIdSinceVersion()
}

func (*ClusterSession) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*ClusterSession) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*ClusterSession) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ClusterSession) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ClusterSession) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ClusterSession) ResponseChannelMetaAttribute(meta int) string {
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

func (*ClusterSession) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (c *ClusterSession) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ResponseChannelSinceVersion()
}

func (*ClusterSession) ResponseChannelDeprecated() uint16 {
	return 0
}

func (ClusterSession) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ClusterSession) ResponseChannelHeaderLength() uint64 {
	return 4
}
