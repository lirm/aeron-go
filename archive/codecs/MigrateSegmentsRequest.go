// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type MigrateSegmentsRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	SrcRecordingId   int64
	DstRecordingId   int64
}

func (m *MigrateSegmentsRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := m.RangeCheck(m.SbeSchemaVersion(), m.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, m.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, m.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, m.SrcRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, m.DstRecordingId); err != nil {
		return err
	}
	return nil
}

func (m *MigrateSegmentsRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !m.ControlSessionIdInActingVersion(actingVersion) {
		m.ControlSessionId = m.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &m.ControlSessionId); err != nil {
			return err
		}
	}
	if !m.CorrelationIdInActingVersion(actingVersion) {
		m.CorrelationId = m.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &m.CorrelationId); err != nil {
			return err
		}
	}
	if !m.SrcRecordingIdInActingVersion(actingVersion) {
		m.SrcRecordingId = m.SrcRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &m.SrcRecordingId); err != nil {
			return err
		}
	}
	if !m.DstRecordingIdInActingVersion(actingVersion) {
		m.DstRecordingId = m.DstRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &m.DstRecordingId); err != nil {
			return err
		}
	}
	if actingVersion > m.SbeSchemaVersion() && blockLength > m.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-m.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := m.RangeCheck(actingVersion, m.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (m *MigrateSegmentsRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if m.ControlSessionIdInActingVersion(actingVersion) {
		if m.ControlSessionId < m.ControlSessionIdMinValue() || m.ControlSessionId > m.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on m.ControlSessionId (%v < %v > %v)", m.ControlSessionIdMinValue(), m.ControlSessionId, m.ControlSessionIdMaxValue())
		}
	}
	if m.CorrelationIdInActingVersion(actingVersion) {
		if m.CorrelationId < m.CorrelationIdMinValue() || m.CorrelationId > m.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on m.CorrelationId (%v < %v > %v)", m.CorrelationIdMinValue(), m.CorrelationId, m.CorrelationIdMaxValue())
		}
	}
	if m.SrcRecordingIdInActingVersion(actingVersion) {
		if m.SrcRecordingId < m.SrcRecordingIdMinValue() || m.SrcRecordingId > m.SrcRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on m.SrcRecordingId (%v < %v > %v)", m.SrcRecordingIdMinValue(), m.SrcRecordingId, m.SrcRecordingIdMaxValue())
		}
	}
	if m.DstRecordingIdInActingVersion(actingVersion) {
		if m.DstRecordingId < m.DstRecordingIdMinValue() || m.DstRecordingId > m.DstRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on m.DstRecordingId (%v < %v > %v)", m.DstRecordingIdMinValue(), m.DstRecordingId, m.DstRecordingIdMaxValue())
		}
	}
	return nil
}

func MigrateSegmentsRequestInit(m *MigrateSegmentsRequest) {
	return
}

func (*MigrateSegmentsRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*MigrateSegmentsRequest) SbeTemplateId() (templateId uint16) {
	return 57
}

func (*MigrateSegmentsRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*MigrateSegmentsRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*MigrateSegmentsRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*MigrateSegmentsRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*MigrateSegmentsRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (m *MigrateSegmentsRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.ControlSessionIdSinceVersion()
}

func (*MigrateSegmentsRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*MigrateSegmentsRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*MigrateSegmentsRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*MigrateSegmentsRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*MigrateSegmentsRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*MigrateSegmentsRequest) CorrelationIdId() uint16 {
	return 2
}

func (*MigrateSegmentsRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (m *MigrateSegmentsRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.CorrelationIdSinceVersion()
}

func (*MigrateSegmentsRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*MigrateSegmentsRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*MigrateSegmentsRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*MigrateSegmentsRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*MigrateSegmentsRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*MigrateSegmentsRequest) SrcRecordingIdId() uint16 {
	return 3
}

func (*MigrateSegmentsRequest) SrcRecordingIdSinceVersion() uint16 {
	return 0
}

func (m *MigrateSegmentsRequest) SrcRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.SrcRecordingIdSinceVersion()
}

func (*MigrateSegmentsRequest) SrcRecordingIdDeprecated() uint16 {
	return 0
}

func (*MigrateSegmentsRequest) SrcRecordingIdMetaAttribute(meta int) string {
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

func (*MigrateSegmentsRequest) SrcRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*MigrateSegmentsRequest) SrcRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*MigrateSegmentsRequest) SrcRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*MigrateSegmentsRequest) DstRecordingIdId() uint16 {
	return 4
}

func (*MigrateSegmentsRequest) DstRecordingIdSinceVersion() uint16 {
	return 0
}

func (m *MigrateSegmentsRequest) DstRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.DstRecordingIdSinceVersion()
}

func (*MigrateSegmentsRequest) DstRecordingIdDeprecated() uint16 {
	return 0
}

func (*MigrateSegmentsRequest) DstRecordingIdMetaAttribute(meta int) string {
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

func (*MigrateSegmentsRequest) DstRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*MigrateSegmentsRequest) DstRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*MigrateSegmentsRequest) DstRecordingIdNullValue() int64 {
	return math.MinInt64
}
