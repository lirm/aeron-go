// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type AdminResponseCodeEnum int32
type AdminResponseCodeValues struct {
	OK                  AdminResponseCodeEnum
	ERROR               AdminResponseCodeEnum
	UNAUTHORISED_ACCESS AdminResponseCodeEnum
	NullValue           AdminResponseCodeEnum
}

var AdminResponseCode = AdminResponseCodeValues{0, 1, 2, -2147483648}

func (a AdminResponseCodeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(a)); err != nil {
		return err
	}
	return nil
}

func (a *AdminResponseCodeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(a)); err != nil {
		return err
	}
	return nil
}

func (a AdminResponseCodeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(AdminResponseCode)
	for idx := 0; idx < value.NumField(); idx++ {
		if a == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on AdminResponseCode, unknown enumeration value %d", a)
}

func (*AdminResponseCodeEnum) EncodedLength() int64 {
	return 4
}

func (*AdminResponseCodeEnum) OKSinceVersion() uint16 {
	return 0
}

func (a *AdminResponseCodeEnum) OKInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.OKSinceVersion()
}

func (*AdminResponseCodeEnum) OKDeprecated() uint16 {
	return 0
}

func (*AdminResponseCodeEnum) ERRORSinceVersion() uint16 {
	return 0
}

func (a *AdminResponseCodeEnum) ERRORInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.ERRORSinceVersion()
}

func (*AdminResponseCodeEnum) ERRORDeprecated() uint16 {
	return 0
}

func (*AdminResponseCodeEnum) UNAUTHORISED_ACCESSSinceVersion() uint16 {
	return 0
}

func (a *AdminResponseCodeEnum) UNAUTHORISED_ACCESSInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.UNAUTHORISED_ACCESSSinceVersion()
}

func (*AdminResponseCodeEnum) UNAUTHORISED_ACCESSDeprecated() uint16 {
	return 0
}
