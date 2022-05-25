// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type ClusterActionEnum int32
type ClusterActionValues struct {
	SUSPEND   ClusterActionEnum
	RESUME    ClusterActionEnum
	SNAPSHOT  ClusterActionEnum
	NullValue ClusterActionEnum
}

var ClusterAction = ClusterActionValues{0, 1, 2, -2147483648}

func (c ClusterActionEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(c)); err != nil {
		return err
	}
	return nil
}

func (c *ClusterActionEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(c)); err != nil {
		return err
	}
	return nil
}

func (c ClusterActionEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(ClusterAction)
	for idx := 0; idx < value.NumField(); idx++ {
		if c == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on ClusterAction, unknown enumeration value %d", c)
}

func (*ClusterActionEnum) EncodedLength() int64 {
	return 4
}

func (*ClusterActionEnum) SUSPENDSinceVersion() uint16 {
	return 0
}

func (c *ClusterActionEnum) SUSPENDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.SUSPENDSinceVersion()
}

func (*ClusterActionEnum) SUSPENDDeprecated() uint16 {
	return 0
}

func (*ClusterActionEnum) RESUMESinceVersion() uint16 {
	return 0
}

func (c *ClusterActionEnum) RESUMEInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.RESUMESinceVersion()
}

func (*ClusterActionEnum) RESUMEDeprecated() uint16 {
	return 0
}

func (*ClusterActionEnum) SNAPSHOTSinceVersion() uint16 {
	return 0
}

func (c *ClusterActionEnum) SNAPSHOTInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.SNAPSHOTSinceVersion()
}

func (*ClusterActionEnum) SNAPSHOTDeprecated() uint16 {
	return 0
}
