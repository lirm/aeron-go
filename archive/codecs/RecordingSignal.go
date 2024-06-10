// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type RecordingSignalEnum int32
type RecordingSignalValues struct {
	START         RecordingSignalEnum
	STOP          RecordingSignalEnum
	EXTEND        RecordingSignalEnum
	REPLICATE     RecordingSignalEnum
	MERGE         RecordingSignalEnum
	SYNC          RecordingSignalEnum
	DELETE        RecordingSignalEnum
	REPLICATE_END RecordingSignalEnum
	NullValue     RecordingSignalEnum
}

var RecordingSignal = RecordingSignalValues{0, 1, 2, 3, 4, 5, 6, 7, -2147483648}

func (r RecordingSignalEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(r)); err != nil {
		return err
	}
	return nil
}

func (r *RecordingSignalEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(r)); err != nil {
		return err
	}
	return nil
}

func (r RecordingSignalEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(RecordingSignal)
	for idx := 0; idx < value.NumField(); idx++ {
		if r == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on RecordingSignal, unknown enumeration value %d", r)
}

func (*RecordingSignalEnum) EncodedLength() int64 {
	return 4
}

func (*RecordingSignalEnum) STARTSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) STARTInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.STARTSinceVersion()
}

func (*RecordingSignalEnum) STARTDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEnum) STOPSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) STOPInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.STOPSinceVersion()
}

func (*RecordingSignalEnum) STOPDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEnum) EXTENDSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) EXTENDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.EXTENDSinceVersion()
}

func (*RecordingSignalEnum) EXTENDDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEnum) REPLICATESinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) REPLICATEInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.REPLICATESinceVersion()
}

func (*RecordingSignalEnum) REPLICATEDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEnum) MERGESinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) MERGEInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.MERGESinceVersion()
}

func (*RecordingSignalEnum) MERGEDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEnum) SYNCSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) SYNCInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SYNCSinceVersion()
}

func (*RecordingSignalEnum) SYNCDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEnum) DELETESinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) DELETEInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.DELETESinceVersion()
}

func (*RecordingSignalEnum) DELETEDeprecated() uint16 {
	return 0
}

func (*RecordingSignalEnum) REPLICATE_ENDSinceVersion() uint16 {
	return 0
}

func (r *RecordingSignalEnum) REPLICATE_ENDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.REPLICATE_ENDSinceVersion()
}

func (*RecordingSignalEnum) REPLICATE_ENDDeprecated() uint16 {
	return 0
}
