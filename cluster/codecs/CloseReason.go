// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type CloseReasonEnum int32
type CloseReasonValues struct {
	CLIENT_ACTION  CloseReasonEnum
	SERVICE_ACTION CloseReasonEnum
	TIMEOUT        CloseReasonEnum
	NullValue      CloseReasonEnum
}

var CloseReason = CloseReasonValues{0, 1, 2, -2147483648}

func (c CloseReasonEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(c)); err != nil {
		return err
	}
	return nil
}

func (c *CloseReasonEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(c)); err != nil {
		return err
	}
	return nil
}

func (c CloseReasonEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(CloseReason)
	for idx := 0; idx < value.NumField(); idx++ {
		if c == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on CloseReason, unknown enumeration value %d", c)
}

func (*CloseReasonEnum) EncodedLength() int64 {
	return 4
}

func (*CloseReasonEnum) CLIENT_ACTIONSinceVersion() uint16 {
	return 0
}

func (c *CloseReasonEnum) CLIENT_ACTIONInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CLIENT_ACTIONSinceVersion()
}

func (*CloseReasonEnum) CLIENT_ACTIONDeprecated() uint16 {
	return 0
}

func (*CloseReasonEnum) SERVICE_ACTIONSinceVersion() uint16 {
	return 0
}

func (c *CloseReasonEnum) SERVICE_ACTIONInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.SERVICE_ACTIONSinceVersion()
}

func (*CloseReasonEnum) SERVICE_ACTIONDeprecated() uint16 {
	return 0
}

func (*CloseReasonEnum) TIMEOUTSinceVersion() uint16 {
	return 0
}

func (c *CloseReasonEnum) TIMEOUTInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.TIMEOUTSinceVersion()
}

func (*CloseReasonEnum) TIMEOUTDeprecated() uint16 {
	return 0
}
