// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type SourceLocationEnum int32
type SourceLocationValues struct {
	LOCAL     SourceLocationEnum
	REMOTE    SourceLocationEnum
	NullValue SourceLocationEnum
}

var SourceLocation = SourceLocationValues{0, 1, -2147483648}

func (s SourceLocationEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(s)); err != nil {
		return err
	}
	return nil
}

func (s *SourceLocationEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(s)); err != nil {
		return err
	}
	return nil
}

func (s SourceLocationEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(SourceLocation)
	for idx := 0; idx < value.NumField(); idx++ {
		if s == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on SourceLocation, unknown enumeration value %d", s)
}

func (*SourceLocationEnum) EncodedLength() int64 {
	return 4
}

func (*SourceLocationEnum) LOCALSinceVersion() uint16 {
	return 0
}

func (s *SourceLocationEnum) LOCALInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LOCALSinceVersion()
}

func (*SourceLocationEnum) LOCALDeprecated() uint16 {
	return 0
}

func (*SourceLocationEnum) REMOTESinceVersion() uint16 {
	return 0
}

func (s *SourceLocationEnum) REMOTEInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.REMOTESinceVersion()
}

func (*SourceLocationEnum) REMOTEDeprecated() uint16 {
	return 0
}
