// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RequestServiceAck struct {
	LogPosition int64
}

func (r *RequestServiceAck) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.LogPosition); err != nil {
		return err
	}
	return nil
}

func (r *RequestServiceAck) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.LogPositionInActingVersion(actingVersion) {
		r.LogPosition = r.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.LogPosition); err != nil {
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

func (r *RequestServiceAck) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.LogPositionInActingVersion(actingVersion) {
		if r.LogPosition < r.LogPositionMinValue() || r.LogPosition > r.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.LogPosition (%v < %v > %v)", r.LogPositionMinValue(), r.LogPosition, r.LogPositionMaxValue())
		}
	}
	return nil
}

func RequestServiceAckInit(r *RequestServiceAck) {
	return
}

func (*RequestServiceAck) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*RequestServiceAck) SbeTemplateId() (templateId uint16) {
	return 108
}

func (*RequestServiceAck) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*RequestServiceAck) SbeSchemaVersion() (schemaVersion uint16) {
	return 16
}

func (*RequestServiceAck) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RequestServiceAck) SbeSemanticVersion() (semanticVersion string) {
	return "5.4"
}

func (*RequestServiceAck) LogPositionId() uint16 {
	return 1
}

func (*RequestServiceAck) LogPositionSinceVersion() uint16 {
	return 0
}

func (r *RequestServiceAck) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.LogPositionSinceVersion()
}

func (*RequestServiceAck) LogPositionDeprecated() uint16 {
	return 0
}

func (*RequestServiceAck) LogPositionMetaAttribute(meta int) string {
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

func (*RequestServiceAck) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RequestServiceAck) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RequestServiceAck) LogPositionNullValue() int64 {
	return math.MinInt64
}
