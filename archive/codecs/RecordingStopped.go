// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingStopped struct {
	RecordingId   int64
	StartPosition int64
	StopPosition  int64
}

func (r *RecordingStopped) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, r.StopPosition); err != nil {
		return err
	}
	return nil
}

func (r *RecordingStopped) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !r.StopPositionInActingVersion(actingVersion) {
		r.StopPosition = r.StopPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.StopPosition); err != nil {
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

func (r *RecordingStopped) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if r.StopPositionInActingVersion(actingVersion) {
		if r.StopPosition < r.StopPositionMinValue() || r.StopPosition > r.StopPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.StopPosition (%v < %v > %v)", r.StopPositionMinValue(), r.StopPosition, r.StopPositionMaxValue())
		}
	}
	return nil
}

func RecordingStoppedInit(r *RecordingStopped) {
	return
}

func (*RecordingStopped) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*RecordingStopped) SbeTemplateId() (templateId uint16) {
	return 103
}

func (*RecordingStopped) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingStopped) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*RecordingStopped) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingStopped) RecordingIdId() uint16 {
	return 1
}

func (*RecordingStopped) RecordingIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingStopped) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.RecordingIdSinceVersion()
}

func (*RecordingStopped) RecordingIdDeprecated() uint16 {
	return 0
}

func (*RecordingStopped) RecordingIdMetaAttribute(meta int) string {
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

func (*RecordingStopped) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingStopped) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingStopped) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingStopped) StartPositionId() uint16 {
	return 2
}

func (*RecordingStopped) StartPositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingStopped) StartPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StartPositionSinceVersion()
}

func (*RecordingStopped) StartPositionDeprecated() uint16 {
	return 0
}

func (*RecordingStopped) StartPositionMetaAttribute(meta int) string {
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

func (*RecordingStopped) StartPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingStopped) StartPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingStopped) StartPositionNullValue() int64 {
	return math.MinInt64
}

func (*RecordingStopped) StopPositionId() uint16 {
	return 3
}

func (*RecordingStopped) StopPositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingStopped) StopPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StopPositionSinceVersion()
}

func (*RecordingStopped) StopPositionDeprecated() uint16 {
	return 0
}

func (*RecordingStopped) StopPositionMetaAttribute(meta int) string {
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

func (*RecordingStopped) StopPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingStopped) StopPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingStopped) StopPositionNullValue() int64 {
	return math.MinInt64
}
