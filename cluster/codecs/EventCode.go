// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type EventCodeEnum int32
type EventCodeValues struct {
	OK                      EventCodeEnum
	ERROR                   EventCodeEnum
	REDIRECT                EventCodeEnum
	AUTHENTICATION_REJECTED EventCodeEnum
	CLOSED                  EventCodeEnum
	NullValue               EventCodeEnum
}

var EventCode = EventCodeValues{0, 1, 2, 3, 4, -2147483648}

func (e EventCodeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(e)); err != nil {
		return err
	}
	return nil
}

func (e *EventCodeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(e)); err != nil {
		return err
	}
	return nil
}

func (e EventCodeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(EventCode)
	for idx := 0; idx < value.NumField(); idx++ {
		if e == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on EventCode, unknown enumeration value %d", e)
}

func (*EventCodeEnum) EncodedLength() int64 {
	return 4
}

func (*EventCodeEnum) OKSinceVersion() uint16 {
	return 0
}

func (e *EventCodeEnum) OKInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.OKSinceVersion()
}

func (*EventCodeEnum) OKDeprecated() uint16 {
	return 0
}

func (*EventCodeEnum) ERRORSinceVersion() uint16 {
	return 0
}

func (e *EventCodeEnum) ERRORInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.ERRORSinceVersion()
}

func (*EventCodeEnum) ERRORDeprecated() uint16 {
	return 0
}

func (*EventCodeEnum) REDIRECTSinceVersion() uint16 {
	return 0
}

func (e *EventCodeEnum) REDIRECTInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.REDIRECTSinceVersion()
}

func (*EventCodeEnum) REDIRECTDeprecated() uint16 {
	return 0
}

func (*EventCodeEnum) AUTHENTICATION_REJECTEDSinceVersion() uint16 {
	return 0
}

func (e *EventCodeEnum) AUTHENTICATION_REJECTEDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.AUTHENTICATION_REJECTEDSinceVersion()
}

func (*EventCodeEnum) AUTHENTICATION_REJECTEDDeprecated() uint16 {
	return 0
}

func (*EventCodeEnum) CLOSEDSinceVersion() uint16 {
	return 0
}

func (e *EventCodeEnum) CLOSEDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.CLOSEDSinceVersion()
}

func (*EventCodeEnum) CLOSEDDeprecated() uint16 {
	return 0
}
