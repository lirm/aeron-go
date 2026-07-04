// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type PendingMessageTracker struct {
	NextServiceSessionId   int64
	LogServiceSessionId    int64
	PendingMessageCapacity int32
	ServiceId              int32
}

func (p *PendingMessageTracker) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := p.RangeCheck(p.SbeSchemaVersion(), p.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, p.NextServiceSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, p.LogServiceSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, p.PendingMessageCapacity); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, p.ServiceId); err != nil {
		return err
	}
	return nil
}

func (p *PendingMessageTracker) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !p.NextServiceSessionIdInActingVersion(actingVersion) {
		p.NextServiceSessionId = p.NextServiceSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &p.NextServiceSessionId); err != nil {
			return err
		}
	}
	if !p.LogServiceSessionIdInActingVersion(actingVersion) {
		p.LogServiceSessionId = p.LogServiceSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &p.LogServiceSessionId); err != nil {
			return err
		}
	}
	if !p.PendingMessageCapacityInActingVersion(actingVersion) {
		p.PendingMessageCapacity = p.PendingMessageCapacityNullValue()
	} else {
		if err := _m.ReadInt32(_r, &p.PendingMessageCapacity); err != nil {
			return err
		}
	}
	if !p.ServiceIdInActingVersion(actingVersion) {
		p.ServiceId = p.ServiceIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &p.ServiceId); err != nil {
			return err
		}
	}
	if actingVersion > p.SbeSchemaVersion() && blockLength > p.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-p.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := p.RangeCheck(actingVersion, p.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (p *PendingMessageTracker) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if p.NextServiceSessionIdInActingVersion(actingVersion) {
		if p.NextServiceSessionId < p.NextServiceSessionIdMinValue() || p.NextServiceSessionId > p.NextServiceSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on p.NextServiceSessionId (%v < %v > %v)", p.NextServiceSessionIdMinValue(), p.NextServiceSessionId, p.NextServiceSessionIdMaxValue())
		}
	}
	if p.LogServiceSessionIdInActingVersion(actingVersion) {
		if p.LogServiceSessionId < p.LogServiceSessionIdMinValue() || p.LogServiceSessionId > p.LogServiceSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on p.LogServiceSessionId (%v < %v > %v)", p.LogServiceSessionIdMinValue(), p.LogServiceSessionId, p.LogServiceSessionIdMaxValue())
		}
	}
	if p.PendingMessageCapacityInActingVersion(actingVersion) {
		if p.PendingMessageCapacity != p.PendingMessageCapacityNullValue() && (p.PendingMessageCapacity < p.PendingMessageCapacityMinValue() || p.PendingMessageCapacity > p.PendingMessageCapacityMaxValue()) {
			return fmt.Errorf("Range check failed on p.PendingMessageCapacity (%v < %v > %v)", p.PendingMessageCapacityMinValue(), p.PendingMessageCapacity, p.PendingMessageCapacityMaxValue())
		}
	}
	if p.ServiceIdInActingVersion(actingVersion) {
		if p.ServiceId < p.ServiceIdMinValue() || p.ServiceId > p.ServiceIdMaxValue() {
			return fmt.Errorf("Range check failed on p.ServiceId (%v < %v > %v)", p.ServiceIdMinValue(), p.ServiceId, p.ServiceIdMaxValue())
		}
	}
	return nil
}

func PendingMessageTrackerInit(p *PendingMessageTracker) {
	p.PendingMessageCapacity = 0
	return
}

func (*PendingMessageTracker) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*PendingMessageTracker) SbeTemplateId() (templateId uint16) {
	return 107
}

func (*PendingMessageTracker) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*PendingMessageTracker) SbeSchemaVersion() (schemaVersion uint16) {
	return 16
}

func (*PendingMessageTracker) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*PendingMessageTracker) SbeSemanticVersion() (semanticVersion string) {
	return "5.4"
}

func (*PendingMessageTracker) NextServiceSessionIdId() uint16 {
	return 1
}

func (*PendingMessageTracker) NextServiceSessionIdSinceVersion() uint16 {
	return 0
}

func (p *PendingMessageTracker) NextServiceSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.NextServiceSessionIdSinceVersion()
}

func (*PendingMessageTracker) NextServiceSessionIdDeprecated() uint16 {
	return 0
}

func (*PendingMessageTracker) NextServiceSessionIdMetaAttribute(meta int) string {
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

func (*PendingMessageTracker) NextServiceSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PendingMessageTracker) NextServiceSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PendingMessageTracker) NextServiceSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*PendingMessageTracker) LogServiceSessionIdId() uint16 {
	return 2
}

func (*PendingMessageTracker) LogServiceSessionIdSinceVersion() uint16 {
	return 0
}

func (p *PendingMessageTracker) LogServiceSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.LogServiceSessionIdSinceVersion()
}

func (*PendingMessageTracker) LogServiceSessionIdDeprecated() uint16 {
	return 0
}

func (*PendingMessageTracker) LogServiceSessionIdMetaAttribute(meta int) string {
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

func (*PendingMessageTracker) LogServiceSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PendingMessageTracker) LogServiceSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PendingMessageTracker) LogServiceSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*PendingMessageTracker) PendingMessageCapacityId() uint16 {
	return 3
}

func (*PendingMessageTracker) PendingMessageCapacitySinceVersion() uint16 {
	return 0
}

func (p *PendingMessageTracker) PendingMessageCapacityInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.PendingMessageCapacitySinceVersion()
}

func (*PendingMessageTracker) PendingMessageCapacityDeprecated() uint16 {
	return 0
}

func (*PendingMessageTracker) PendingMessageCapacityMetaAttribute(meta int) string {
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

func (*PendingMessageTracker) PendingMessageCapacityMinValue() int32 {
	return math.MinInt32 + 1
}

func (*PendingMessageTracker) PendingMessageCapacityMaxValue() int32 {
	return math.MaxInt32
}

func (*PendingMessageTracker) PendingMessageCapacityNullValue() int32 {
	return 0
}

func (*PendingMessageTracker) ServiceIdId() uint16 {
	return 4
}

func (*PendingMessageTracker) ServiceIdSinceVersion() uint16 {
	return 0
}

func (p *PendingMessageTracker) ServiceIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.ServiceIdSinceVersion()
}

func (*PendingMessageTracker) ServiceIdDeprecated() uint16 {
	return 0
}

func (*PendingMessageTracker) ServiceIdMetaAttribute(meta int) string {
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

func (*PendingMessageTracker) ServiceIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*PendingMessageTracker) ServiceIdMaxValue() int32 {
	return math.MaxInt32
}

func (*PendingMessageTracker) ServiceIdNullValue() int32 {
	return math.MinInt32
}
