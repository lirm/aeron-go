// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"math"
)

type VarAsciiEncoding struct {
	Length  uint32
	VarData uint8
}

func (v *VarAsciiEncoding) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteUint32(_w, v.Length); err != nil {
		return err
	}
	if err := _m.WriteUint8(_w, v.VarData); err != nil {
		return err
	}
	return nil
}

func (v *VarAsciiEncoding) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if !v.LengthInActingVersion(actingVersion) {
		v.Length = v.LengthNullValue()
	} else {
		if err := _m.ReadUint32(_r, &v.Length); err != nil {
			return err
		}
	}
	if !v.VarDataInActingVersion(actingVersion) {
		v.VarData = v.VarDataNullValue()
	} else {
		if err := _m.ReadUint8(_r, &v.VarData); err != nil {
			return err
		}
	}
	return nil
}

func (v *VarAsciiEncoding) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if v.LengthInActingVersion(actingVersion) {
		if v.Length < v.LengthMinValue() || v.Length > v.LengthMaxValue() {
			return fmt.Errorf("Range check failed on v.Length (%v < %v > %v)", v.LengthMinValue(), v.Length, v.LengthMaxValue())
		}
	}
	if v.VarDataInActingVersion(actingVersion) {
		if v.VarData < v.VarDataMinValue() || v.VarData > v.VarDataMaxValue() {
			return fmt.Errorf("Range check failed on v.VarData (%v < %v > %v)", v.VarDataMinValue(), v.VarData, v.VarDataMaxValue())
		}
	}
	return nil
}

func VarAsciiEncodingInit(v *VarAsciiEncoding) {
	return
}

func (*VarAsciiEncoding) EncodedLength() int64 {
	return -1
}

func (*VarAsciiEncoding) LengthMinValue() uint32 {
	return 0
}

func (*VarAsciiEncoding) LengthMaxValue() uint32 {
	return 1073741824
}

func (*VarAsciiEncoding) LengthNullValue() uint32 {
	return math.MaxUint32
}

func (*VarAsciiEncoding) LengthSinceVersion() uint16 {
	return 0
}

func (v *VarAsciiEncoding) LengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.LengthSinceVersion()
}

func (*VarAsciiEncoding) LengthDeprecated() uint16 {
	return 0
}

func (*VarAsciiEncoding) VarDataMinValue() uint8 {
	return 0
}

func (*VarAsciiEncoding) VarDataMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*VarAsciiEncoding) VarDataNullValue() uint8 {
	return math.MaxUint8
}

func (*VarAsciiEncoding) VarDataSinceVersion() uint16 {
	return 0
}

func (v *VarAsciiEncoding) VarDataInActingVersion(actingVersion uint16) bool {
	return actingVersion >= v.VarDataSinceVersion()
}

func (*VarAsciiEncoding) VarDataDeprecated() uint16 {
	return 0
}
