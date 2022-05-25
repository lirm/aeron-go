// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type ClusterTimeUnitEnum int32
type ClusterTimeUnitValues struct {
	MILLIS    ClusterTimeUnitEnum
	MICROS    ClusterTimeUnitEnum
	NANOS     ClusterTimeUnitEnum
	NullValue ClusterTimeUnitEnum
}

var ClusterTimeUnit = ClusterTimeUnitValues{0, 1, 2, -2147483648}

func (c ClusterTimeUnitEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(c)); err != nil {
		return err
	}
	return nil
}

func (c *ClusterTimeUnitEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(c)); err != nil {
		return err
	}
	return nil
}

func (c ClusterTimeUnitEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(ClusterTimeUnit)
	for idx := 0; idx < value.NumField(); idx++ {
		if c == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on ClusterTimeUnit, unknown enumeration value %d", c)
}

func (*ClusterTimeUnitEnum) EncodedLength() int64 {
	return 4
}

func (*ClusterTimeUnitEnum) MILLISSinceVersion() uint16 {
	return 0
}

func (c *ClusterTimeUnitEnum) MILLISInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.MILLISSinceVersion()
}

func (*ClusterTimeUnitEnum) MILLISDeprecated() uint16 {
	return 0
}

func (*ClusterTimeUnitEnum) MICROSSinceVersion() uint16 {
	return 0
}

func (c *ClusterTimeUnitEnum) MICROSInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.MICROSSinceVersion()
}

func (*ClusterTimeUnitEnum) MICROSDeprecated() uint16 {
	return 0
}

func (*ClusterTimeUnitEnum) NANOSSinceVersion() uint16 {
	return 0
}

func (c *ClusterTimeUnitEnum) NANOSInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.NANOSSinceVersion()
}

func (*ClusterTimeUnitEnum) NANOSDeprecated() uint16 {
	return 0
}
