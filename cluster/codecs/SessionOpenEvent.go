// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionOpenEvent struct {
	LeadershipTermId int64
	CorrelationId    int64
	ClusterSessionId int64
	Timestamp        int64
	ResponseStreamId int32
	ResponseChannel  []uint8
	EncodedPrincipal []uint8
}

func (s *SessionOpenEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.ResponseStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.ResponseChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.EncodedPrincipal))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.EncodedPrincipal); err != nil {
		return err
	}
	return nil
}

func (s *SessionOpenEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.LeadershipTermIdInActingVersion(actingVersion) {
		s.LeadershipTermId = s.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LeadershipTermId); err != nil {
			return err
		}
	}
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
			return err
		}
	}
	if !s.ClusterSessionIdInActingVersion(actingVersion) {
		s.ClusterSessionId = s.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.ClusterSessionId); err != nil {
			return err
		}
	}
	if !s.TimestampInActingVersion(actingVersion) {
		s.Timestamp = s.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.Timestamp); err != nil {
			return err
		}
	}
	if !s.ResponseStreamIdInActingVersion(actingVersion) {
		s.ResponseStreamId = s.ResponseStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.ResponseStreamId); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	if s.ResponseChannelInActingVersion(actingVersion) {
		var ResponseChannelLength uint32
		if err := _m.ReadUint32(_r, &ResponseChannelLength); err != nil {
			return err
		}
		if cap(s.ResponseChannel) < int(ResponseChannelLength) {
			s.ResponseChannel = make([]uint8, ResponseChannelLength)
		}
		s.ResponseChannel = s.ResponseChannel[:ResponseChannelLength]
		if err := _m.ReadBytes(_r, s.ResponseChannel); err != nil {
			return err
		}
	}

	if s.EncodedPrincipalInActingVersion(actingVersion) {
		var EncodedPrincipalLength uint32
		if err := _m.ReadUint32(_r, &EncodedPrincipalLength); err != nil {
			return err
		}
		if cap(s.EncodedPrincipal) < int(EncodedPrincipalLength) {
			s.EncodedPrincipal = make([]uint8, EncodedPrincipalLength)
		}
		s.EncodedPrincipal = s.EncodedPrincipal[:EncodedPrincipalLength]
		if err := _m.ReadBytes(_r, s.EncodedPrincipal); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionOpenEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.LeadershipTermIdInActingVersion(actingVersion) {
		if s.LeadershipTermId < s.LeadershipTermIdMinValue() || s.LeadershipTermId > s.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeadershipTermId (%v < %v > %v)", s.LeadershipTermIdMinValue(), s.LeadershipTermId, s.LeadershipTermIdMaxValue())
		}
	}
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	if s.ClusterSessionIdInActingVersion(actingVersion) {
		if s.ClusterSessionId < s.ClusterSessionIdMinValue() || s.ClusterSessionId > s.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ClusterSessionId (%v < %v > %v)", s.ClusterSessionIdMinValue(), s.ClusterSessionId, s.ClusterSessionIdMaxValue())
		}
	}
	if s.TimestampInActingVersion(actingVersion) {
		if s.Timestamp < s.TimestampMinValue() || s.Timestamp > s.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on s.Timestamp (%v < %v > %v)", s.TimestampMinValue(), s.Timestamp, s.TimestampMaxValue())
		}
	}
	if s.ResponseStreamIdInActingVersion(actingVersion) {
		if s.ResponseStreamId < s.ResponseStreamIdMinValue() || s.ResponseStreamId > s.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ResponseStreamId (%v < %v > %v)", s.ResponseStreamIdMinValue(), s.ResponseStreamId, s.ResponseStreamIdMaxValue())
		}
	}
	for idx, ch := range s.ResponseChannel {
		if ch > 127 {
			return fmt.Errorf("s.ResponseChannel[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func SessionOpenEventInit(s *SessionOpenEvent) {
	return
}

func (*SessionOpenEvent) SbeBlockLength() (blockLength uint16) {
	return 36
}

func (*SessionOpenEvent) SbeTemplateId() (templateId uint16) {
	return 21
}

func (*SessionOpenEvent) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SessionOpenEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SessionOpenEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionOpenEvent) LeadershipTermIdId() uint16 {
	return 1
}

func (*SessionOpenEvent) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *SessionOpenEvent) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*SessionOpenEvent) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*SessionOpenEvent) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*SessionOpenEvent) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionOpenEvent) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionOpenEvent) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionOpenEvent) CorrelationIdId() uint16 {
	return 2
}

func (*SessionOpenEvent) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *SessionOpenEvent) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*SessionOpenEvent) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*SessionOpenEvent) CorrelationIdMetaAttribute(meta int) string {
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

func (*SessionOpenEvent) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionOpenEvent) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionOpenEvent) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionOpenEvent) ClusterSessionIdId() uint16 {
	return 3
}

func (*SessionOpenEvent) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (s *SessionOpenEvent) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ClusterSessionIdSinceVersion()
}

func (*SessionOpenEvent) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*SessionOpenEvent) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*SessionOpenEvent) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionOpenEvent) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionOpenEvent) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionOpenEvent) TimestampId() uint16 {
	return 4
}

func (*SessionOpenEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (s *SessionOpenEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TimestampSinceVersion()
}

func (*SessionOpenEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*SessionOpenEvent) TimestampMetaAttribute(meta int) string {
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

func (*SessionOpenEvent) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionOpenEvent) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionOpenEvent) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*SessionOpenEvent) ResponseStreamIdId() uint16 {
	return 6
}

func (*SessionOpenEvent) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (s *SessionOpenEvent) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ResponseStreamIdSinceVersion()
}

func (*SessionOpenEvent) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*SessionOpenEvent) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*SessionOpenEvent) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*SessionOpenEvent) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*SessionOpenEvent) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*SessionOpenEvent) ResponseChannelMetaAttribute(meta int) string {
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

func (*SessionOpenEvent) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (s *SessionOpenEvent) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ResponseChannelSinceVersion()
}

func (*SessionOpenEvent) ResponseChannelDeprecated() uint16 {
	return 0
}

func (SessionOpenEvent) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (SessionOpenEvent) ResponseChannelHeaderLength() uint64 {
	return 4
}

func (*SessionOpenEvent) EncodedPrincipalMetaAttribute(meta int) string {
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

func (*SessionOpenEvent) EncodedPrincipalSinceVersion() uint16 {
	return 0
}

func (s *SessionOpenEvent) EncodedPrincipalInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.EncodedPrincipalSinceVersion()
}

func (*SessionOpenEvent) EncodedPrincipalDeprecated() uint16 {
	return 0
}

func (SessionOpenEvent) EncodedPrincipalCharacterEncoding() string {
	return "null"
}

func (SessionOpenEvent) EncodedPrincipalHeaderLength() uint64 {
	return 4
}
