// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type BooleanTypeEnum int32
type BooleanTypeValues struct {
	FALSE     BooleanTypeEnum
	TRUE      BooleanTypeEnum
	NullValue BooleanTypeEnum
}

var BooleanType = BooleanTypeValues{0, 1, -2147483648}

func (b BooleanTypeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(b)); err != nil {
		return err
	}
	return nil
}

func (b *BooleanTypeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(b)); err != nil {
		return err
	}
	return nil
}

func (b BooleanTypeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(BooleanType)
	for idx := 0; idx < value.NumField(); idx++ {
		if b == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on BooleanType, unknown enumeration value %d", b)
}

func (*BooleanTypeEnum) EncodedLength() int64 {
	return 4
}

func (*BooleanTypeEnum) FALSESinceVersion() uint16 {
	return 0
}

func (b *BooleanTypeEnum) FALSEInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.FALSESinceVersion()
}

func (*BooleanTypeEnum) FALSEDeprecated() uint16 {
	return 0
}

func (*BooleanTypeEnum) TRUESinceVersion() uint16 {
	return 0
}

func (b *BooleanTypeEnum) TRUEInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.TRUESinceVersion()
}

func (*BooleanTypeEnum) TRUEDeprecated() uint16 {
	return 0
}
