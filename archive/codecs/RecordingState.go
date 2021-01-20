// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type RecordingStateEnum int32
type RecordingStateValues struct {
	INVALID   RecordingStateEnum
	VALID     RecordingStateEnum
	NullValue RecordingStateEnum
}

var RecordingState = RecordingStateValues{0, 1, -2147483648}

func (r RecordingStateEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(r)); err != nil {
		return err
	}
	return nil
}

func (r *RecordingStateEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(r)); err != nil {
		return err
	}
	return nil
}

func (r RecordingStateEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(RecordingState)
	for idx := 0; idx < value.NumField(); idx++ {
		if r == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on RecordingState, unknown enumeration value %d", r)
}

func (*RecordingStateEnum) EncodedLength() int64 {
	return 4
}

func (*RecordingStateEnum) INVALIDSinceVersion() uint16 {
	return 0
}

func (r *RecordingStateEnum) INVALIDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.INVALIDSinceVersion()
}

func (*RecordingStateEnum) INVALIDDeprecated() uint16 {
	return 0
}

func (*RecordingStateEnum) VALIDSinceVersion() uint16 {
	return 0
}

func (r *RecordingStateEnum) VALIDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.VALIDSinceVersion()
}

func (*RecordingStateEnum) VALIDDeprecated() uint16 {
	return 0
}
