// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionConnectRequest struct {
	CorrelationId      int64
	ResponseStreamId   int32
	Version            int32
	ResponseChannel    []uint8
	EncodedCredentials []uint8
}

func (s *SessionConnectRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.ResponseStreamId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.Version); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.ResponseChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.EncodedCredentials))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.EncodedCredentials); err != nil {
		return err
	}
	return nil
}

func (s *SessionConnectRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
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
	if !s.VersionInActingVersion(actingVersion) {
		s.Version = s.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.Version); err != nil {
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

	if s.EncodedCredentialsInActingVersion(actingVersion) {
		var EncodedCredentialsLength uint32
		if err := _m.ReadUint32(_r, &EncodedCredentialsLength); err != nil {
			return err
		}
		if cap(s.EncodedCredentials) < int(EncodedCredentialsLength) {
			s.EncodedCredentials = make([]uint8, EncodedCredentialsLength)
		}
		s.EncodedCredentials = s.EncodedCredentials[:EncodedCredentialsLength]
		if err := _m.ReadBytes(_r, s.EncodedCredentials); err != nil {
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

func (s *SessionConnectRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	if s.ResponseStreamIdInActingVersion(actingVersion) {
		if s.ResponseStreamId < s.ResponseStreamIdMinValue() || s.ResponseStreamId > s.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ResponseStreamId (%v < %v > %v)", s.ResponseStreamIdMinValue(), s.ResponseStreamId, s.ResponseStreamIdMaxValue())
		}
	}
	if s.VersionInActingVersion(actingVersion) {
		if s.Version != s.VersionNullValue() && (s.Version < s.VersionMinValue() || s.Version > s.VersionMaxValue()) {
			return fmt.Errorf("Range check failed on s.Version (%v < %v > %v)", s.VersionMinValue(), s.Version, s.VersionMaxValue())
		}
	}
	for idx, ch := range s.ResponseChannel {
		if ch > 127 {
			return fmt.Errorf("s.ResponseChannel[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func SessionConnectRequestInit(s *SessionConnectRequest) {
	s.Version = 0
	return
}

func (*SessionConnectRequest) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*SessionConnectRequest) SbeTemplateId() (templateId uint16) {
	return 3
}

func (*SessionConnectRequest) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SessionConnectRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SessionConnectRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionConnectRequest) CorrelationIdId() uint16 {
	return 1
}

func (*SessionConnectRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *SessionConnectRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*SessionConnectRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*SessionConnectRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*SessionConnectRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionConnectRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionConnectRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionConnectRequest) ResponseStreamIdId() uint16 {
	return 2
}

func (*SessionConnectRequest) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (s *SessionConnectRequest) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ResponseStreamIdSinceVersion()
}

func (*SessionConnectRequest) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*SessionConnectRequest) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*SessionConnectRequest) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*SessionConnectRequest) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*SessionConnectRequest) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*SessionConnectRequest) VersionId() uint16 {
	return 3
}

func (*SessionConnectRequest) VersionSinceVersion() uint16 {
	return 2
}

func (s *SessionConnectRequest) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.VersionSinceVersion()
}

func (*SessionConnectRequest) VersionDeprecated() uint16 {
	return 0
}

func (*SessionConnectRequest) VersionMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "optional"
	}
	return ""
}

func (*SessionConnectRequest) VersionMinValue() int32 {
	return 1
}

func (*SessionConnectRequest) VersionMaxValue() int32 {
	return 16777215
}

func (*SessionConnectRequest) VersionNullValue() int32 {
	return 0
}

func (*SessionConnectRequest) ResponseChannelMetaAttribute(meta int) string {
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

func (*SessionConnectRequest) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (s *SessionConnectRequest) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ResponseChannelSinceVersion()
}

func (*SessionConnectRequest) ResponseChannelDeprecated() uint16 {
	return 0
}

func (SessionConnectRequest) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (SessionConnectRequest) ResponseChannelHeaderLength() uint64 {
	return 4
}

func (*SessionConnectRequest) EncodedCredentialsMetaAttribute(meta int) string {
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

func (*SessionConnectRequest) EncodedCredentialsSinceVersion() uint16 {
	return 0
}

func (s *SessionConnectRequest) EncodedCredentialsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.EncodedCredentialsSinceVersion()
}

func (*SessionConnectRequest) EncodedCredentialsDeprecated() uint16 {
	return 0
}

func (SessionConnectRequest) EncodedCredentialsCharacterEncoding() string {
	return "null"
}

func (SessionConnectRequest) EncodedCredentialsHeaderLength() uint64 {
	return 4
}
