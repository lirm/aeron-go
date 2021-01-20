// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type ControlResponseCodeEnum int32
type ControlResponseCodeValues struct {
	OK                   ControlResponseCodeEnum
	ERROR                ControlResponseCodeEnum
	RECORDING_UNKNOWN    ControlResponseCodeEnum
	SUBSCRIPTION_UNKNOWN ControlResponseCodeEnum
	NullValue            ControlResponseCodeEnum
}

var ControlResponseCode = ControlResponseCodeValues{0, 1, 2, 3, -2147483648}

func (c ControlResponseCodeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(c)); err != nil {
		return err
	}
	return nil
}

func (c *ControlResponseCodeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(c)); err != nil {
		return err
	}
	return nil
}

func (c ControlResponseCodeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(ControlResponseCode)
	for idx := 0; idx < value.NumField(); idx++ {
		if c == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on ControlResponseCode, unknown enumeration value %d", c)
}

func (*ControlResponseCodeEnum) EncodedLength() int64 {
	return 4
}

func (*ControlResponseCodeEnum) OKSinceVersion() uint16 {
	return 0
}

func (c *ControlResponseCodeEnum) OKInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.OKSinceVersion()
}

func (*ControlResponseCodeEnum) OKDeprecated() uint16 {
	return 0
}

func (*ControlResponseCodeEnum) ERRORSinceVersion() uint16 {
	return 0
}

func (c *ControlResponseCodeEnum) ERRORInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ERRORSinceVersion()
}

func (*ControlResponseCodeEnum) ERRORDeprecated() uint16 {
	return 0
}

func (*ControlResponseCodeEnum) RECORDING_UNKNOWNSinceVersion() uint16 {
	return 0
}

func (c *ControlResponseCodeEnum) RECORDING_UNKNOWNInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.RECORDING_UNKNOWNSinceVersion()
}

func (*ControlResponseCodeEnum) RECORDING_UNKNOWNDeprecated() uint16 {
	return 0
}

func (*ControlResponseCodeEnum) SUBSCRIPTION_UNKNOWNSinceVersion() uint16 {
	return 0
}

func (c *ControlResponseCodeEnum) SUBSCRIPTION_UNKNOWNInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.SUBSCRIPTION_UNKNOWNSinceVersion()
}

func (*ControlResponseCodeEnum) SUBSCRIPTION_UNKNOWNDeprecated() uint16 {
	return 0
}
