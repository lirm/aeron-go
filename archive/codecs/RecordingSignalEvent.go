// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingSignalEvent struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	SubscriptionId   int64
	Position         int64
	Signal           RecordingSignalEnum
}

func (r *RecordingSignalEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.SubscriptionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.Position); err != nil {
		return err
	}
	if err := r.Signal.Encode(_m, _w); err != nil {
		return err
	}
	return nil
}

func (r *RecordingSignalEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.ControlSessionIdInActingVersion(actingVersion) {
		r.ControlSessionId = r.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.ControlSessionId); err != nil {
			return err
		}
	}
	if !r.CorrelationIdInActingVersion(actingVersion) {
		r.CorrelationId = r.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.CorrelationId); err != nil {
			return err
		}
	}
	if !r.RecordingIdInActingVersion(actingVersion) {
		r.RecordingId = r.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.RecordingId); err != nil {
			return err
		}
	}
	if !r.SubscriptionIdInActingVersion(actingVersion) {
		r.SubscriptionId = r.SubscriptionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.SubscriptionId); err != nil {
			return err
		}
	}
	if !r.PositionInActingVersion(actingVersion) {
		r.Position = r.PositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.Position); err != nil {
			return err
		}
	}
	if r.SignalInActingVersion(actingVersion) {
		if err := r.Signal.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := r.RangeCheck(actingVersion, r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordingSignalEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.ControlSessionIdInActingVersion(actingVersion) {
		if r.ControlSessionId < r.ControlSessionIdMinValue() || r.ControlSessionId > r.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.ControlSessionId (%v < %v > %v)", r.ControlSessionIdMinValue(), r.ControlSessionId, r.ControlSessionIdMaxValue())
		}
	}
	if r.CorrelationIdInActingVersion(actingVersion) {
		if r.CorrelationId < r.CorrelationIdMinValue() || r.CorrelationId > r.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on r.CorrelationId (%v < %v > %v)", r.CorrelationIdMinValue(), r.CorrelationId, r.CorrelationIdMaxValue())
		}
	}
	if r.RecordingIdInActingVersion(actingVersion) {
		if r.RecordingId < r.RecordingIdMinValue() || r.RecordingId > r.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.RecordingId (%v < %v > %v)", r.RecordingIdMinValue(), r.RecordingId, r.RecordingIdMaxValue())
		}
	}
	if r.SubscriptionIdInActingVersion(actingVersion) {
		if r.SubscriptionId < r.SubscriptionIdMinValue() || r.SubscriptionId > r.SubscriptionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SubscriptionId (%v < %v > %v)", r.SubscriptionIdMinValue(), r.SubscriptionId, r.SubscriptionIdMaxValue())
		}
	}
	if r.PositionInActingVersion(actingVersion) {
		if r.Position < r.PositionMinValue() || r.Position > r.PositionMaxValue() {
			return fmt.Errorf("Range check failed on r.Position (%v < %v > %v)", r.PositionMinValue(), r.Position, r.PositionMaxValue())
		}
	}
	if err := r.Signal.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func RecordingSignalEventInit(r *RecordingSignalEvent) {
	return
}

func (*RecordingSignalEvent) SbeBlockLength() (blockLength uint16) {
	return 44
}

func (*RecordingSignalEvent) SbeTemplateId() (templateId uint16) {
	return 24
}

func (*RecordingSignalEvent) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingSignalEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*RecordingSignalEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingSignalEvent) ControlSessionIdId() uint16 {
	return 1
}

func (*RecordingSignalEvent) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEvent) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ControlSessionIdSinceVersion()
}

func (*RecordingSignalEvent) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEvent) ControlSessionIdMetaAttribute(meta int) string {
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

func (*RecordingSignalEvent) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSignalEvent) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSignalEvent) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSignalEvent) CorrelationIdId() uint16 {
	return 2
}

func (*RecordingSignalEvent) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEvent) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CorrelationIdSinceVersion()
}

func (*RecordingSignalEvent) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEvent) CorrelationIdMetaAttribute(meta int) string {
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

func (*RecordingSignalEvent) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSignalEvent) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSignalEvent) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSignalEvent) RecordingIdId() uint16 {
	return 3
}

func (*RecordingSignalEvent) RecordingIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEvent) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.RecordingIdSinceVersion()
}

func (*RecordingSignalEvent) RecordingIdDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEvent) RecordingIdMetaAttribute(meta int) string {
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

func (*RecordingSignalEvent) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSignalEvent) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSignalEvent) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSignalEvent) SubscriptionIdId() uint16 {
	return 4
}

func (*RecordingSignalEvent) SubscriptionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEvent) SubscriptionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SubscriptionIdSinceVersion()
}

func (*RecordingSignalEvent) SubscriptionIdDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEvent) SubscriptionIdMetaAttribute(meta int) string {
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

func (*RecordingSignalEvent) SubscriptionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSignalEvent) SubscriptionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSignalEvent) SubscriptionIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSignalEvent) PositionId() uint16 {
	return 5
}

func (*RecordingSignalEvent) PositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEvent) PositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.PositionSinceVersion()
}

func (*RecordingSignalEvent) PositionDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEvent) PositionMetaAttribute(meta int) string {
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

func (*RecordingSignalEvent) PositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingSignalEvent) PositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingSignalEvent) PositionNullValue() int64 {
	return math.MinInt64
}

func (*RecordingSignalEvent) SignalId() uint16 {
	return 6
}

func (*RecordingSignalEvent) SignalSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEvent) SignalInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SignalSinceVersion()
}

func (*RecordingSignalEvent) SignalDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEvent) SignalMetaAttribute(meta int) string {
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
