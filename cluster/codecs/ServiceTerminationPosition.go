// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ServiceTerminationPosition struct {
	LogPosition int64
}

func (s *ServiceTerminationPosition) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.LogPosition); err != nil {
		return err
	}
	return nil
}

func (s *ServiceTerminationPosition) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.LogPositionInActingVersion(actingVersion) {
		s.LogPosition = s.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LogPosition); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceTerminationPosition) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.LogPositionInActingVersion(actingVersion) {
		if s.LogPosition < s.LogPositionMinValue() || s.LogPosition > s.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on s.LogPosition (%v < %v > %v)", s.LogPositionMinValue(), s.LogPosition, s.LogPositionMaxValue())
		}
	}
	return nil
}

func ServiceTerminationPositionInit(s *ServiceTerminationPosition) {
	return
}

func (*ServiceTerminationPosition) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*ServiceTerminationPosition) SbeTemplateId() (templateId uint16) {
	return 42
}

func (*ServiceTerminationPosition) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ServiceTerminationPosition) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ServiceTerminationPosition) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ServiceTerminationPosition) LogPositionId() uint16 {
	return 1
}

func (*ServiceTerminationPosition) LogPositionSinceVersion() uint16 {
	return 0
}

func (s *ServiceTerminationPosition) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LogPositionSinceVersion()
}

func (*ServiceTerminationPosition) LogPositionDeprecated() uint16 {
	return 0
}

func (*ServiceTerminationPosition) LogPositionMetaAttribute(meta int) string {
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

func (*ServiceTerminationPosition) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ServiceTerminationPosition) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ServiceTerminationPosition) LogPositionNullValue() int64 {
	return math.MinInt64
}
