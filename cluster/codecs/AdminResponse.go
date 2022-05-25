// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type AdminResponse struct {
	ClusterSessionId int64
	CorrelationId    int64
	RequestType      AdminRequestTypeEnum
	ResponseCode     AdminResponseCodeEnum
	Message          []uint8
	Payload          []uint8
}

func (a *AdminResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := a.RangeCheck(a.SbeSchemaVersion(), a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, a.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, a.CorrelationId); err != nil {
		return err
	}
	if err := a.RequestType.Encode(_m, _w); err != nil {
		return err
	}
	if err := a.ResponseCode.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(a.Message))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, a.Message); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(a.Payload))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, a.Payload); err != nil {
		return err
	}
	return nil
}

func (a *AdminResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !a.ClusterSessionIdInActingVersion(actingVersion) {
		a.ClusterSessionId = a.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.ClusterSessionId); err != nil {
			return err
		}
	}
	if !a.CorrelationIdInActingVersion(actingVersion) {
		a.CorrelationId = a.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.CorrelationId); err != nil {
			return err
		}
	}
	if a.RequestTypeInActingVersion(actingVersion) {
		if err := a.RequestType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if a.ResponseCodeInActingVersion(actingVersion) {
		if err := a.ResponseCode.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if actingVersion > a.SbeSchemaVersion() && blockLength > a.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-a.SbeBlockLength()))
	}

	if a.MessageInActingVersion(actingVersion) {
		var MessageLength uint32
		if err := _m.ReadUint32(_r, &MessageLength); err != nil {
			return err
		}
		if cap(a.Message) < int(MessageLength) {
			a.Message = make([]uint8, MessageLength)
		}
		a.Message = a.Message[:MessageLength]
		if err := _m.ReadBytes(_r, a.Message); err != nil {
			return err
		}
	}

	if a.PayloadInActingVersion(actingVersion) {
		var PayloadLength uint32
		if err := _m.ReadUint32(_r, &PayloadLength); err != nil {
			return err
		}
		if cap(a.Payload) < int(PayloadLength) {
			a.Payload = make([]uint8, PayloadLength)
		}
		a.Payload = a.Payload[:PayloadLength]
		if err := _m.ReadBytes(_r, a.Payload); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := a.RangeCheck(actingVersion, a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (a *AdminResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if a.ClusterSessionIdInActingVersion(actingVersion) {
		if a.ClusterSessionId < a.ClusterSessionIdMinValue() || a.ClusterSessionId > a.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on a.ClusterSessionId (%v < %v > %v)", a.ClusterSessionIdMinValue(), a.ClusterSessionId, a.ClusterSessionIdMaxValue())
		}
	}
	if a.CorrelationIdInActingVersion(actingVersion) {
		if a.CorrelationId < a.CorrelationIdMinValue() || a.CorrelationId > a.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on a.CorrelationId (%v < %v > %v)", a.CorrelationIdMinValue(), a.CorrelationId, a.CorrelationIdMaxValue())
		}
	}
	if err := a.RequestType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if err := a.ResponseCode.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	for idx, ch := range a.Message {
		if ch > 127 {
			return fmt.Errorf("a.Message[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func AdminResponseInit(a *AdminResponse) {
	return
}

func (*AdminResponse) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*AdminResponse) SbeTemplateId() (templateId uint16) {
	return 27
}

func (*AdminResponse) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*AdminResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*AdminResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*AdminResponse) ClusterSessionIdId() uint16 {
	return 1
}

func (*AdminResponse) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (a *AdminResponse) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.ClusterSessionIdSinceVersion()
}

func (*AdminResponse) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*AdminResponse) ClusterSessionIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*AdminResponse) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AdminResponse) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AdminResponse) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*AdminResponse) CorrelationIdId() uint16 {
	return 2
}

func (*AdminResponse) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (a *AdminResponse) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.CorrelationIdSinceVersion()
}

func (*AdminResponse) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*AdminResponse) CorrelationIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*AdminResponse) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AdminResponse) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AdminResponse) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*AdminResponse) RequestTypeId() uint16 {
	return 3
}

func (*AdminResponse) RequestTypeSinceVersion() uint16 {
	return 0
}

func (a *AdminResponse) RequestTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.RequestTypeSinceVersion()
}

func (*AdminResponse) RequestTypeDeprecated() uint16 {
	return 0
}

func (*AdminResponse) RequestTypeMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*AdminResponse) ResponseCodeId() uint16 {
	return 4
}

func (*AdminResponse) ResponseCodeSinceVersion() uint16 {
	return 0
}

func (a *AdminResponse) ResponseCodeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.ResponseCodeSinceVersion()
}

func (*AdminResponse) ResponseCodeDeprecated() uint16 {
	return 0
}

func (*AdminResponse) ResponseCodeMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*AdminResponse) MessageMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*AdminResponse) MessageSinceVersion() uint16 {
	return 0
}

func (a *AdminResponse) MessageInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.MessageSinceVersion()
}

func (*AdminResponse) MessageDeprecated() uint16 {
	return 0
}

func (AdminResponse) MessageCharacterEncoding() string {
	return "US-ASCII"
}

func (AdminResponse) MessageHeaderLength() uint64 {
	return 4
}

func (*AdminResponse) PayloadMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*AdminResponse) PayloadSinceVersion() uint16 {
	return 0
}

func (a *AdminResponse) PayloadInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.PayloadSinceVersion()
}

func (*AdminResponse) PayloadDeprecated() uint16 {
	return 0
}

func (AdminResponse) PayloadCharacterEncoding() string {
	return "null"
}

func (AdminResponse) PayloadHeaderLength() uint64 {
	return 4
}
