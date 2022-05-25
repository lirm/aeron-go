// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"reflect"
)

type SnapshotMarkEnum int32
type SnapshotMarkValues struct {
	BEGIN     SnapshotMarkEnum
	SECTION   SnapshotMarkEnum
	END       SnapshotMarkEnum
	NullValue SnapshotMarkEnum
}

var SnapshotMark = SnapshotMarkValues{0, 1, 2, -2147483648}

func (s SnapshotMarkEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteInt32(_w, int32(s)); err != nil {
		return err
	}
	return nil
}

func (s *SnapshotMarkEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadInt32(_r, (*int32)(s)); err != nil {
		return err
	}
	return nil
}

func (s SnapshotMarkEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(SnapshotMark)
	for idx := 0; idx < value.NumField(); idx++ {
		if s == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on SnapshotMark, unknown enumeration value %d", s)
}

func (*SnapshotMarkEnum) EncodedLength() int64 {
	return 4
}

func (*SnapshotMarkEnum) BEGINSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarkEnum) BEGINInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.BEGINSinceVersion()
}

func (*SnapshotMarkEnum) BEGINDeprecated() uint16 {
	return 0
}

func (*SnapshotMarkEnum) SECTIONSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarkEnum) SECTIONInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.SECTIONSinceVersion()
}

func (*SnapshotMarkEnum) SECTIONDeprecated() uint16 {
	return 0
}

func (*SnapshotMarkEnum) ENDSinceVersion() uint16 {
	return 0
}

func (s *SnapshotMarkEnum) ENDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ENDSinceVersion()
}

func (*SnapshotMarkEnum) ENDDeprecated() uint16 {
	return 0
}
