// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ChallengeResponse struct {
	CorrelationId      int64
	ClusterSessionId   int64
	EncodedCredentials []uint8
}

func (c *ChallengeResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.EncodedCredentials))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.EncodedCredentials); err != nil {
		return err
	}
	return nil
}

func (c *ChallengeResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
			return err
		}
	}
	if !c.ClusterSessionIdInActingVersion(actingVersion) {
		c.ClusterSessionId = c.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.ClusterSessionId); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.EncodedCredentialsInActingVersion(actingVersion) {
		var EncodedCredentialsLength uint32
		if err := _m.ReadUint32(_r, &EncodedCredentialsLength); err != nil {
			return err
		}
		if cap(c.EncodedCredentials) < int(EncodedCredentialsLength) {
			c.EncodedCredentials = make([]uint8, EncodedCredentialsLength)
		}
		c.EncodedCredentials = c.EncodedCredentials[:EncodedCredentialsLength]
		if err := _m.ReadBytes(_r, c.EncodedCredentials); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *ChallengeResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if c.ClusterSessionIdInActingVersion(actingVersion) {
		if c.ClusterSessionId < c.ClusterSessionIdMinValue() || c.ClusterSessionId > c.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ClusterSessionId (%v < %v > %v)", c.ClusterSessionIdMinValue(), c.ClusterSessionId, c.ClusterSessionIdMaxValue())
		}
	}
	return nil
}

func ChallengeResponseInit(c *ChallengeResponse) {
	return
}

func (*ChallengeResponse) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*ChallengeResponse) SbeTemplateId() (templateId uint16) {
	return 8
}

func (*ChallengeResponse) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ChallengeResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ChallengeResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ChallengeResponse) CorrelationIdId() uint16 {
	return 1
}

func (*ChallengeResponse) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *ChallengeResponse) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*ChallengeResponse) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ChallengeResponse) CorrelationIdMetaAttribute(meta int) string {
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

func (*ChallengeResponse) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ChallengeResponse) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ChallengeResponse) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ChallengeResponse) ClusterSessionIdId() uint16 {
	return 2
}

func (*ChallengeResponse) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (c *ChallengeResponse) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ClusterSessionIdSinceVersion()
}

func (*ChallengeResponse) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*ChallengeResponse) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*ChallengeResponse) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ChallengeResponse) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ChallengeResponse) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ChallengeResponse) EncodedCredentialsMetaAttribute(meta int) string {
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

func (*ChallengeResponse) EncodedCredentialsSinceVersion() uint16 {
	return 0
}

func (c *ChallengeResponse) EncodedCredentialsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.EncodedCredentialsSinceVersion()
}

func (*ChallengeResponse) EncodedCredentialsDeprecated() uint16 {
	return 0
}

func (ChallengeResponse) EncodedCredentialsCharacterEncoding() string {
	return "null"
}

func (ChallengeResponse) EncodedCredentialsHeaderLength() uint64 {
	return 4
}
