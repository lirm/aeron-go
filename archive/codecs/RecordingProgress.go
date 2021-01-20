// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingProgress struct {
	RecordingId   int64
	StartPosition int64
	Position      int64
}

func (r *RecordingProgress) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.StartPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.Position); err != nil {
		return err
	}
	return nil
}

func (r *RecordingProgress) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.RecordingIdInActingVersion(actingVersion) {
		r.RecordingId = r.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.RecordingId); err != nil {
			return err
		}
	}
	if !r.StartPositionInActingVersion(actingVersion) {
		r.StartPosition = r.StartPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.StartPosition); err != nil {
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

func (r *RecordingProgress) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.RecordingIdInActingVersion(actingVersion) {
		if r.RecordingId < r.RecordingIdMinValue() || r.RecordingId > r.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.RecordingId (%v < %v > %v)", r.RecordingIdMinValue(), r.RecordingId, r.RecordingIdMaxValue())
		}
	}
	if r.StartPositionInActingVersion(actingVersion) {
		if r.StartPosition < r.StartPositionMinValue() || r.StartPosition > r.StartPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.StartPosition (%v < %v > %v)", r.StartPositionMinValue(), r.StartPosition, r.StartPositionMaxValue())
		}
	}
	if r.PositionInActingVersion(actingVersion) {
		if r.Position < r.PositionMinValue() || r.Position > r.PositionMaxValue() {
			return fmt.Errorf("Range check failed on r.Position (%v < %v > %v)", r.PositionMinValue(), r.Position, r.PositionMaxValue())
		}
	}
	return nil
}

func RecordingProgressInit(r *RecordingProgress) {
	return
}

func (*RecordingProgress) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*RecordingProgress) SbeTemplateId() (templateId uint16) {
	return 102
}

func (*RecordingProgress) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingProgress) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*RecordingProgress) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingProgress) RecordingIdId() uint16 {
	return 1
}

func (*RecordingProgress) RecordingIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingProgress) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.RecordingIdSinceVersion()
}

func (*RecordingProgress) RecordingIdDeprecated() uint16 {
	return 0
}

func (*RecordingProgress) RecordingIdMetaAttribute(meta int) string {
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

func (*RecordingProgress) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingProgress) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingProgress) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingProgress) StartPositionId() uint16 {
	return 2
}

func (*RecordingProgress) StartPositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingProgress) StartPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StartPositionSinceVersion()
}

func (*RecordingProgress) StartPositionDeprecated() uint16 {
	return 0
}

func (*RecordingProgress) StartPositionMetaAttribute(meta int) string {
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

func (*RecordingProgress) StartPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingProgress) StartPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingProgress) StartPositionNullValue() int64 {
	return math.MinInt64
}

func (*RecordingProgress) PositionId() uint16 {
	return 3
}

func (*RecordingProgress) PositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingProgress) PositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.PositionSinceVersion()
}

func (*RecordingProgress) PositionDeprecated() uint16 {
	return 0
}

func (*RecordingProgress) PositionMetaAttribute(meta int) string {
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

func (*RecordingProgress) PositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingProgress) PositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingProgress) PositionNullValue() int64 {
	return math.MinInt64
}
