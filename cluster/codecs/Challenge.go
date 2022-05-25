// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Challenge struct {
	CorrelationId    int64
	ClusterSessionId int64
	EncodedChallenge []uint8
}

func (c *Challenge) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteUint32(_w, uint32(len(c.EncodedChallenge))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.EncodedChallenge); err != nil {
		return err
	}
	return nil
}

func (c *Challenge) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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

	if c.EncodedChallengeInActingVersion(actingVersion) {
		var EncodedChallengeLength uint32
		if err := _m.ReadUint32(_r, &EncodedChallengeLength); err != nil {
			return err
		}
		if cap(c.EncodedChallenge) < int(EncodedChallengeLength) {
			c.EncodedChallenge = make([]uint8, EncodedChallengeLength)
		}
		c.EncodedChallenge = c.EncodedChallenge[:EncodedChallengeLength]
		if err := _m.ReadBytes(_r, c.EncodedChallenge); err != nil {
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

func (c *Challenge) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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

func ChallengeInit(c *Challenge) {
	return
}

func (*Challenge) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*Challenge) SbeTemplateId() (templateId uint16) {
	return 7
}

func (*Challenge) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*Challenge) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*Challenge) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*Challenge) CorrelationIdId() uint16 {
	return 1
}

func (*Challenge) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (c *Challenge) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CorrelationIdSinceVersion()
}

func (*Challenge) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*Challenge) CorrelationIdMetaAttribute(meta int) string {
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

func (*Challenge) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Challenge) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*Challenge) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*Challenge) ClusterSessionIdId() uint16 {
	return 2
}

func (*Challenge) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (c *Challenge) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ClusterSessionIdSinceVersion()
}

func (*Challenge) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*Challenge) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*Challenge) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Challenge) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*Challenge) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*Challenge) EncodedChallengeMetaAttribute(meta int) string {
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

func (*Challenge) EncodedChallengeSinceVersion() uint16 {
	return 0
}

func (c *Challenge) EncodedChallengeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.EncodedChallengeSinceVersion()
}

func (*Challenge) EncodedChallengeDeprecated() uint16 {
	return 0
}

func (Challenge) EncodedChallengeCharacterEncoding() string {
	return "null"
}

func (Challenge) EncodedChallengeHeaderLength() uint64 {
	return 4
}
