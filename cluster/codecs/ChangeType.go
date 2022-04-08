// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type ChangeTypeEnum int32
type ChangeTypeValues struct {
	JOIN      ChangeTypeEnum
	QUIT      ChangeTypeEnum
	NullValue ChangeTypeEnum
}

var ChangeType = ChangeTypeValues{0, 1, -2147483648}

func (c ChangeTypeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(c)); err != nil {
		return err
	}
	return nil
}

func (c *ChangeTypeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(c)); err != nil {
		return err
	}
	return nil
}

func (c ChangeTypeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(ChangeType)
	for idx := 0; idx < value.NumField(); idx++ {
		if c == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on ChangeType, unknown enumeration value %d", c)
}

func (*ChangeTypeEnum) EncodedLength() int64 {
	return 4
}

func (*ChangeTypeEnum) JOINSinceVersion() uint16 {
	return 0
}

func (c *ChangeTypeEnum) JOINInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.JOINSinceVersion()
}

func (*ChangeTypeEnum) JOINDeprecated() uint16 {
	return 0
}

func (*ChangeTypeEnum) QUITSinceVersion() uint16 {
	return 0
}

func (c *ChangeTypeEnum) QUITInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.QUITSinceVersion()
}

func (*ChangeTypeEnum) QUITDeprecated() uint16 {
	return 0
}
