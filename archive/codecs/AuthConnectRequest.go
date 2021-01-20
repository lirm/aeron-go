// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type AuthConnectRequest struct {
	CorrelationId      int64
	ResponseStreamId   int32
	Version            int32
	ResponseChannel    []uint8
	EncodedCredentials []uint8
}

func (a *AuthConnectRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := a.RangeCheck(a.SbeSchemaVersion(), a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, a.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, a.ResponseStreamId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, a.Version); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(a.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, a.ResponseChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(a.EncodedCredentials))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, a.EncodedCredentials); err != nil {
		return err
	}
	return nil
}

func (a *AuthConnectRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !a.CorrelationIdInActingVersion(actingVersion) {
		a.CorrelationId = a.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.CorrelationId); err != nil {
			return err
		}
	}
	if !a.ResponseStreamIdInActingVersion(actingVersion) {
		a.ResponseStreamId = a.ResponseStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &a.ResponseStreamId); err != nil {
			return err
		}
	}
	if !a.VersionInActingVersion(actingVersion) {
		a.Version = a.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &a.Version); err != nil {
			return err
		}
	}
	if actingVersion > a.SbeSchemaVersion() && blockLength > a.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-a.SbeBlockLength()))
	}

	if a.ResponseChannelInActingVersion(actingVersion) {
		var ResponseChannelLength uint32
		if err := _m.ReadUint32(_r, &ResponseChannelLength); err != nil {
			return err
		}
		if cap(a.ResponseChannel) < int(ResponseChannelLength) {
			a.ResponseChannel = make([]uint8, ResponseChannelLength)
		}
		a.ResponseChannel = a.ResponseChannel[:ResponseChannelLength]
		if err := _m.ReadBytes(_r, a.ResponseChannel); err != nil {
			return err
		}
	}

	if a.EncodedCredentialsInActingVersion(actingVersion) {
		var EncodedCredentialsLength uint32
		if err := _m.ReadUint32(_r, &EncodedCredentialsLength); err != nil {
			return err
		}
		if cap(a.EncodedCredentials) < int(EncodedCredentialsLength) {
			a.EncodedCredentials = make([]uint8, EncodedCredentialsLength)
		}
		a.EncodedCredentials = a.EncodedCredentials[:EncodedCredentialsLength]
		if err := _m.ReadBytes(_r, a.EncodedCredentials); err != nil {
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

func (a *AuthConnectRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if a.CorrelationIdInActingVersion(actingVersion) {
		if a.CorrelationId < a.CorrelationIdMinValue() || a.CorrelationId > a.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on a.CorrelationId (%v < %v > %v)", a.CorrelationIdMinValue(), a.CorrelationId, a.CorrelationIdMaxValue())
		}
	}
	if a.ResponseStreamIdInActingVersion(actingVersion) {
		if a.ResponseStreamId < a.ResponseStreamIdMinValue() || a.ResponseStreamId > a.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on a.ResponseStreamId (%v < %v > %v)", a.ResponseStreamIdMinValue(), a.ResponseStreamId, a.ResponseStreamIdMaxValue())
		}
	}
	if a.VersionInActingVersion(actingVersion) {
		if a.Version != a.VersionNullValue() && (a.Version < a.VersionMinValue() || a.Version > a.VersionMaxValue()) {
			return fmt.Errorf("Range check failed on a.Version (%v < %v > %v)", a.VersionMinValue(), a.Version, a.VersionMaxValue())
		}
	}
	return nil
}

func AuthConnectRequestInit(a *AuthConnectRequest) {
	a.Version = 0
	return
}

func (*AuthConnectRequest) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*AuthConnectRequest) SbeTemplateId() (templateId uint16) {
	return 58
}

func (*AuthConnectRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*AuthConnectRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*AuthConnectRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*AuthConnectRequest) CorrelationIdId() uint16 {
	return 1
}

func (*AuthConnectRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (a *AuthConnectRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.CorrelationIdSinceVersion()
}

func (*AuthConnectRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*AuthConnectRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*AuthConnectRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AuthConnectRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AuthConnectRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*AuthConnectRequest) ResponseStreamIdId() uint16 {
	return 2
}

func (*AuthConnectRequest) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (a *AuthConnectRequest) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.ResponseStreamIdSinceVersion()
}

func (*AuthConnectRequest) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*AuthConnectRequest) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*AuthConnectRequest) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*AuthConnectRequest) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*AuthConnectRequest) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*AuthConnectRequest) VersionId() uint16 {
	return 3
}

func (*AuthConnectRequest) VersionSinceVersion() uint16 {
	return 0
}

func (a *AuthConnectRequest) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.VersionSinceVersion()
}

func (*AuthConnectRequest) VersionDeprecated() uint16 {
	return 0
}

func (*AuthConnectRequest) VersionMetaAttribute(meta int) string {
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

func (*AuthConnectRequest) VersionMinValue() int32 {
	return 2
}

func (*AuthConnectRequest) VersionMaxValue() int32 {
	return 16777215
}

func (*AuthConnectRequest) VersionNullValue() int32 {
	return 0
}

func (*AuthConnectRequest) ResponseChannelMetaAttribute(meta int) string {
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

func (*AuthConnectRequest) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (a *AuthConnectRequest) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.ResponseChannelSinceVersion()
}

func (*AuthConnectRequest) ResponseChannelDeprecated() uint16 {
	return 0
}

func (AuthConnectRequest) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (AuthConnectRequest) ResponseChannelHeaderLength() uint64 {
	return 4
}

func (*AuthConnectRequest) EncodedCredentialsMetaAttribute(meta int) string {
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

func (*AuthConnectRequest) EncodedCredentialsSinceVersion() uint16 {
	return 0
}

func (a *AuthConnectRequest) EncodedCredentialsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.EncodedCredentialsSinceVersion()
}

func (*AuthConnectRequest) EncodedCredentialsDeprecated() uint16 {
	return 0
}

func (AuthConnectRequest) EncodedCredentialsCharacterEncoding() string {
	return "null"
}

func (AuthConnectRequest) EncodedCredentialsHeaderLength() uint64 {
	return 4
}
