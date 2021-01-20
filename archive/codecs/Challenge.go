// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Challenge struct {
	ControlSessionId int64
	CorrelationId    int64
	Version          int32
	EncodedChallenge []uint8
}

func (c *Challenge) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, c.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.Version); err != nil {
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
	if !c.ControlSessionIdInActingVersion(actingVersion) {
		c.ControlSessionId = c.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.ControlSessionId); err != nil {
			return err
		}
	}
	if !c.CorrelationIdInActingVersion(actingVersion) {
		c.CorrelationId = c.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.CorrelationId); err != nil {
			return err
		}
	}
	if !c.VersionInActingVersion(actingVersion) {
		c.Version = c.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.Version); err != nil {
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
	if c.ControlSessionIdInActingVersion(actingVersion) {
		if c.ControlSessionId < c.ControlSessionIdMinValue() || c.ControlSessionId > c.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on c.ControlSessionId (%v < %v > %v)", c.ControlSessionIdMinValue(), c.ControlSessionId, c.ControlSessionIdMaxValue())
		}
	}
	if c.CorrelationIdInActingVersion(actingVersion) {
		if c.CorrelationId < c.CorrelationIdMinValue() || c.CorrelationId > c.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on c.CorrelationId (%v < %v > %v)", c.CorrelationIdMinValue(), c.CorrelationId, c.CorrelationIdMaxValue())
		}
	}
	if c.VersionInActingVersion(actingVersion) {
		if c.Version != c.VersionNullValue() && (c.Version < c.VersionMinValue() || c.Version > c.VersionMaxValue()) {
			return fmt.Errorf("Range check failed on c.Version (%v < %v > %v)", c.VersionMinValue(), c.Version, c.VersionMaxValue())
		}
	}
	return nil
}

func ChallengeInit(c *Challenge) {
	c.Version = 0
	return
}

func (*Challenge) SbeBlockLength() (blockLength uint16) {
	return 20
}

func (*Challenge) SbeTemplateId() (templateId uint16) {
	return 59
}

func (*Challenge) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*Challenge) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*Challenge) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*Challenge) ControlSessionIdId() uint16 {
	return 1
}

func (*Challenge) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (c *Challenge) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ControlSessionIdSinceVersion()
}

func (*Challenge) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*Challenge) ControlSessionIdMetaAttribute(meta int) string {
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

func (*Challenge) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Challenge) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*Challenge) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*Challenge) CorrelationIdId() uint16 {
	return 2
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

func (*Challenge) VersionId() uint16 {
	return 3
}

func (*Challenge) VersionSinceVersion() uint16 {
	return 0
}

func (c *Challenge) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.VersionSinceVersion()
}

func (*Challenge) VersionDeprecated() uint16 {
	return 0
}

func (*Challenge) VersionMetaAttribute(meta int) string {
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

func (*Challenge) VersionMinValue() int32 {
	return 2
}

func (*Challenge) VersionMaxValue() int32 {
	return 16777215
}

func (*Challenge) VersionNullValue() int32 {
	return 0
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
