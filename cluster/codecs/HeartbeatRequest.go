// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type HeartbeatRequest struct {
	CorrelationId      int64
	ResponseStreamId   int32
	ResponseChannel    []uint8
	EncodedCredentials []uint8
}

func (h *HeartbeatRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := h.RangeCheck(h.SbeSchemaVersion(), h.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, h.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, h.ResponseStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(h.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, h.ResponseChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(h.EncodedCredentials))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, h.EncodedCredentials); err != nil {
		return err
	}
	return nil
}

func (h *HeartbeatRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !h.CorrelationIdInActingVersion(actingVersion) {
		h.CorrelationId = h.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &h.CorrelationId); err != nil {
			return err
		}
	}
	if !h.ResponseStreamIdInActingVersion(actingVersion) {
		h.ResponseStreamId = h.ResponseStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &h.ResponseStreamId); err != nil {
			return err
		}
	}
	if actingVersion > h.SbeSchemaVersion() && blockLength > h.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-h.SbeBlockLength()))
	}

	if h.ResponseChannelInActingVersion(actingVersion) {
		var ResponseChannelLength uint32
		if err := _m.ReadUint32(_r, &ResponseChannelLength); err != nil {
			return err
		}
		if cap(h.ResponseChannel) < int(ResponseChannelLength) {
			h.ResponseChannel = make([]uint8, ResponseChannelLength)
		}
		h.ResponseChannel = h.ResponseChannel[:ResponseChannelLength]
		if err := _m.ReadBytes(_r, h.ResponseChannel); err != nil {
			return err
		}
	}

	if h.EncodedCredentialsInActingVersion(actingVersion) {
		var EncodedCredentialsLength uint32
		if err := _m.ReadUint32(_r, &EncodedCredentialsLength); err != nil {
			return err
		}
		if cap(h.EncodedCredentials) < int(EncodedCredentialsLength) {
			h.EncodedCredentials = make([]uint8, EncodedCredentialsLength)
		}
		h.EncodedCredentials = h.EncodedCredentials[:EncodedCredentialsLength]
		if err := _m.ReadBytes(_r, h.EncodedCredentials); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := h.RangeCheck(actingVersion, h.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (h *HeartbeatRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if h.CorrelationIdInActingVersion(actingVersion) {
		if h.CorrelationId < h.CorrelationIdMinValue() || h.CorrelationId > h.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on h.CorrelationId (%v < %v > %v)", h.CorrelationIdMinValue(), h.CorrelationId, h.CorrelationIdMaxValue())
		}
	}
	if h.ResponseStreamIdInActingVersion(actingVersion) {
		if h.ResponseStreamId < h.ResponseStreamIdMinValue() || h.ResponseStreamId > h.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on h.ResponseStreamId (%v < %v > %v)", h.ResponseStreamIdMinValue(), h.ResponseStreamId, h.ResponseStreamIdMaxValue())
		}
	}
	for idx, ch := range h.ResponseChannel {
		if ch > 127 {
			return fmt.Errorf("h.ResponseChannel[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func HeartbeatRequestInit(h *HeartbeatRequest) {
	return
}

func (*HeartbeatRequest) SbeBlockLength() (blockLength uint16) {
	return 12
}

func (*HeartbeatRequest) SbeTemplateId() (templateId uint16) {
	return 79
}

func (*HeartbeatRequest) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*HeartbeatRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 16
}

func (*HeartbeatRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*HeartbeatRequest) SbeSemanticVersion() (semanticVersion string) {
	return "5.4"
}

func (*HeartbeatRequest) CorrelationIdId() uint16 {
	return 1
}

func (*HeartbeatRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (h *HeartbeatRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= h.CorrelationIdSinceVersion()
}

func (*HeartbeatRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*HeartbeatRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*HeartbeatRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*HeartbeatRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*HeartbeatRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*HeartbeatRequest) ResponseStreamIdId() uint16 {
	return 2
}

func (*HeartbeatRequest) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (h *HeartbeatRequest) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= h.ResponseStreamIdSinceVersion()
}

func (*HeartbeatRequest) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*HeartbeatRequest) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*HeartbeatRequest) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*HeartbeatRequest) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*HeartbeatRequest) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*HeartbeatRequest) ResponseChannelMetaAttribute(meta int) string {
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

func (*HeartbeatRequest) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (h *HeartbeatRequest) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= h.ResponseChannelSinceVersion()
}

func (*HeartbeatRequest) ResponseChannelDeprecated() uint16 {
	return 0
}

func (HeartbeatRequest) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (HeartbeatRequest) ResponseChannelHeaderLength() uint64 {
	return 4
}

func (*HeartbeatRequest) EncodedCredentialsMetaAttribute(meta int) string {
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

func (*HeartbeatRequest) EncodedCredentialsSinceVersion() uint16 {
	return 0
}

func (h *HeartbeatRequest) EncodedCredentialsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= h.EncodedCredentialsSinceVersion()
}

func (*HeartbeatRequest) EncodedCredentialsDeprecated() uint16 {
	return 0
}

func (HeartbeatRequest) EncodedCredentialsCharacterEncoding() string {
	return "null"
}

func (HeartbeatRequest) EncodedCredentialsHeaderLength() uint64 {
	return 4
}
