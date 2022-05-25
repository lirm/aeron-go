// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type BackupQuery struct {
	CorrelationId      int64
	ResponseStreamId   int32
	Version            int32
	ResponseChannel    []uint8
	EncodedCredentials []uint8
}

func (b *BackupQuery) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := b.RangeCheck(b.SbeSchemaVersion(), b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, b.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, b.ResponseStreamId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, b.Version); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(b.ResponseChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, b.ResponseChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(b.EncodedCredentials))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, b.EncodedCredentials); err != nil {
		return err
	}
	return nil
}

func (b *BackupQuery) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !b.CorrelationIdInActingVersion(actingVersion) {
		b.CorrelationId = b.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &b.CorrelationId); err != nil {
			return err
		}
	}
	if !b.ResponseStreamIdInActingVersion(actingVersion) {
		b.ResponseStreamId = b.ResponseStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &b.ResponseStreamId); err != nil {
			return err
		}
	}
	if !b.VersionInActingVersion(actingVersion) {
		b.Version = b.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &b.Version); err != nil {
			return err
		}
	}
	if actingVersion > b.SbeSchemaVersion() && blockLength > b.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}

	if b.ResponseChannelInActingVersion(actingVersion) {
		var ResponseChannelLength uint32
		if err := _m.ReadUint32(_r, &ResponseChannelLength); err != nil {
			return err
		}
		if cap(b.ResponseChannel) < int(ResponseChannelLength) {
			b.ResponseChannel = make([]uint8, ResponseChannelLength)
		}
		b.ResponseChannel = b.ResponseChannel[:ResponseChannelLength]
		if err := _m.ReadBytes(_r, b.ResponseChannel); err != nil {
			return err
		}
	}

	if b.EncodedCredentialsInActingVersion(actingVersion) {
		var EncodedCredentialsLength uint32
		if err := _m.ReadUint32(_r, &EncodedCredentialsLength); err != nil {
			return err
		}
		if cap(b.EncodedCredentials) < int(EncodedCredentialsLength) {
			b.EncodedCredentials = make([]uint8, EncodedCredentialsLength)
		}
		b.EncodedCredentials = b.EncodedCredentials[:EncodedCredentialsLength]
		if err := _m.ReadBytes(_r, b.EncodedCredentials); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := b.RangeCheck(actingVersion, b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (b *BackupQuery) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if b.CorrelationIdInActingVersion(actingVersion) {
		if b.CorrelationId < b.CorrelationIdMinValue() || b.CorrelationId > b.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on b.CorrelationId (%v < %v > %v)", b.CorrelationIdMinValue(), b.CorrelationId, b.CorrelationIdMaxValue())
		}
	}
	if b.ResponseStreamIdInActingVersion(actingVersion) {
		if b.ResponseStreamId < b.ResponseStreamIdMinValue() || b.ResponseStreamId > b.ResponseStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on b.ResponseStreamId (%v < %v > %v)", b.ResponseStreamIdMinValue(), b.ResponseStreamId, b.ResponseStreamIdMaxValue())
		}
	}
	if b.VersionInActingVersion(actingVersion) {
		if b.Version != b.VersionNullValue() && (b.Version < b.VersionMinValue() || b.Version > b.VersionMaxValue()) {
			return fmt.Errorf("Range check failed on b.Version (%v < %v > %v)", b.VersionMinValue(), b.Version, b.VersionMaxValue())
		}
	}
	for idx, ch := range b.ResponseChannel {
		if ch > 127 {
			return fmt.Errorf("b.ResponseChannel[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func BackupQueryInit(b *BackupQuery) {
	b.Version = 0
	return
}

func (*BackupQuery) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*BackupQuery) SbeTemplateId() (templateId uint16) {
	return 77
}

func (*BackupQuery) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*BackupQuery) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*BackupQuery) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*BackupQuery) CorrelationIdId() uint16 {
	return 1
}

func (*BackupQuery) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (b *BackupQuery) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.CorrelationIdSinceVersion()
}

func (*BackupQuery) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*BackupQuery) CorrelationIdMetaAttribute(meta int) string {
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

func (*BackupQuery) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*BackupQuery) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*BackupQuery) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*BackupQuery) ResponseStreamIdId() uint16 {
	return 2
}

func (*BackupQuery) ResponseStreamIdSinceVersion() uint16 {
	return 0
}

func (b *BackupQuery) ResponseStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ResponseStreamIdSinceVersion()
}

func (*BackupQuery) ResponseStreamIdDeprecated() uint16 {
	return 0
}

func (*BackupQuery) ResponseStreamIdMetaAttribute(meta int) string {
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

func (*BackupQuery) ResponseStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*BackupQuery) ResponseStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*BackupQuery) ResponseStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*BackupQuery) VersionId() uint16 {
	return 3
}

func (*BackupQuery) VersionSinceVersion() uint16 {
	return 2
}

func (b *BackupQuery) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.VersionSinceVersion()
}

func (*BackupQuery) VersionDeprecated() uint16 {
	return 0
}

func (*BackupQuery) VersionMetaAttribute(meta int) string {
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

func (*BackupQuery) VersionMinValue() int32 {
	return 1
}

func (*BackupQuery) VersionMaxValue() int32 {
	return 16777215
}

func (*BackupQuery) VersionNullValue() int32 {
	return 0
}

func (*BackupQuery) ResponseChannelMetaAttribute(meta int) string {
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

func (*BackupQuery) ResponseChannelSinceVersion() uint16 {
	return 0
}

func (b *BackupQuery) ResponseChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ResponseChannelSinceVersion()
}

func (*BackupQuery) ResponseChannelDeprecated() uint16 {
	return 0
}

func (BackupQuery) ResponseChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (BackupQuery) ResponseChannelHeaderLength() uint64 {
	return 4
}

func (*BackupQuery) EncodedCredentialsMetaAttribute(meta int) string {
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

func (*BackupQuery) EncodedCredentialsSinceVersion() uint16 {
	return 0
}

func (b *BackupQuery) EncodedCredentialsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.EncodedCredentialsSinceVersion()
}

func (*BackupQuery) EncodedCredentialsDeprecated() uint16 {
	return 0
}

func (BackupQuery) EncodedCredentialsCharacterEncoding() string {
	return "null"
}

func (BackupQuery) EncodedCredentialsHeaderLength() uint64 {
	return 4
}
