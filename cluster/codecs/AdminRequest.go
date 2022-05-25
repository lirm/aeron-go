// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type AdminRequest struct {
	LeadershipTermId int64
	ClusterSessionId int64
	CorrelationId    int64
	RequestType      AdminRequestTypeEnum
	Payload          []uint8
}

func (a *AdminRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := a.RangeCheck(a.SbeSchemaVersion(), a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, a.LeadershipTermId); err != nil {
		return err
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
	if err := _m.WriteUint32(_w, uint32(len(a.Payload))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, a.Payload); err != nil {
		return err
	}
	return nil
}

func (a *AdminRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !a.LeadershipTermIdInActingVersion(actingVersion) {
		a.LeadershipTermId = a.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.LeadershipTermId); err != nil {
			return err
		}
	}
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
	if actingVersion > a.SbeSchemaVersion() && blockLength > a.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-a.SbeBlockLength()))
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

func (a *AdminRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if a.LeadershipTermIdInActingVersion(actingVersion) {
		if a.LeadershipTermId < a.LeadershipTermIdMinValue() || a.LeadershipTermId > a.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on a.LeadershipTermId (%v < %v > %v)", a.LeadershipTermIdMinValue(), a.LeadershipTermId, a.LeadershipTermIdMaxValue())
		}
	}
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
	return nil
}

func AdminRequestInit(a *AdminRequest) {
	return
}

func (*AdminRequest) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*AdminRequest) SbeTemplateId() (templateId uint16) {
	return 26
}

func (*AdminRequest) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*AdminRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*AdminRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*AdminRequest) LeadershipTermIdId() uint16 {
	return 1
}

func (*AdminRequest) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (a *AdminRequest) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.LeadershipTermIdSinceVersion()
}

func (*AdminRequest) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*AdminRequest) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*AdminRequest) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AdminRequest) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AdminRequest) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*AdminRequest) ClusterSessionIdId() uint16 {
	return 2
}

func (*AdminRequest) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (a *AdminRequest) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.ClusterSessionIdSinceVersion()
}

func (*AdminRequest) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*AdminRequest) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*AdminRequest) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AdminRequest) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AdminRequest) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*AdminRequest) CorrelationIdId() uint16 {
	return 3
}

func (*AdminRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (a *AdminRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.CorrelationIdSinceVersion()
}

func (*AdminRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*AdminRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*AdminRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AdminRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AdminRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*AdminRequest) RequestTypeId() uint16 {
	return 4
}

func (*AdminRequest) RequestTypeSinceVersion() uint16 {
	return 0
}

func (a *AdminRequest) RequestTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.RequestTypeSinceVersion()
}

func (*AdminRequest) RequestTypeDeprecated() uint16 {
	return 0
}

func (*AdminRequest) RequestTypeMetaAttribute(meta int) string {
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

func (*AdminRequest) PayloadMetaAttribute(meta int) string {
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

func (*AdminRequest) PayloadSinceVersion() uint16 {
	return 0
}

func (a *AdminRequest) PayloadInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.PayloadSinceVersion()
}

func (*AdminRequest) PayloadDeprecated() uint16 {
	return 0
}

func (AdminRequest) PayloadCharacterEncoding() string {
	return "null"
}

func (AdminRequest) PayloadHeaderLength() uint64 {
	return 4
}
