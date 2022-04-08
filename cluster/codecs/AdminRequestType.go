// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type AdminRequestTypeEnum int32
type AdminRequestTypeValues struct {
	SNAPSHOT  AdminRequestTypeEnum
	NullValue AdminRequestTypeEnum
}

var AdminRequestType = AdminRequestTypeValues{0, -2147483648}

func (a AdminRequestTypeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(a)); err != nil {
		return err
	}
	return nil
}

func (a *AdminRequestTypeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(a)); err != nil {
		return err
	}
	return nil
}

func (a AdminRequestTypeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(AdminRequestType)
	for idx := 0; idx < value.NumField(); idx++ {
		if a == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on AdminRequestType, unknown enumeration value %d", a)
}

func (*AdminRequestTypeEnum) EncodedLength() int64 {
	return 4
}

func (*AdminRequestTypeEnum) SNAPSHOTSinceVersion() uint16 {
	return 0
}

func (a *AdminRequestTypeEnum) SNAPSHOTInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.SNAPSHOTSinceVersion()
}

func (*AdminRequestTypeEnum) SNAPSHOTDeprecated() uint16 {
	return 0
}
