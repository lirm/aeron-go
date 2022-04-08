// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionCloseRequest struct {
	LeadershipTermId int64
	ClusterSessionId int64
}

func (s *SessionCloseRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.ClusterSessionId); err != nil {
		return err
	}
	return nil
}

func (s *SessionCloseRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.LeadershipTermIdInActingVersion(actingVersion) {
		s.LeadershipTermId = s.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LeadershipTermId); err != nil {
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
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionCloseRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.LeadershipTermIdInActingVersion(actingVersion) {
		if s.LeadershipTermId < s.LeadershipTermIdMinValue() || s.LeadershipTermId > s.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on s.LeadershipTermId (%v < %v > %v)", s.LeadershipTermIdMinValue(), s.LeadershipTermId, s.LeadershipTermIdMaxValue())
		}
	}
	if s.ClusterSessionIdInActingVersion(actingVersion) {
		if s.ClusterSessionId < s.ClusterSessionIdMinValue() || s.ClusterSessionId > s.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ClusterSessionId (%v < %v > %v)", s.ClusterSessionIdMinValue(), s.ClusterSessionId, s.ClusterSessionIdMaxValue())
		}
	}
	return nil
}

func SessionCloseRequestInit(s *SessionCloseRequest) {
	return
}

func (*SessionCloseRequest) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*SessionCloseRequest) SbeTemplateId() (templateId uint16) {
	return 4
}

func (*SessionCloseRequest) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*SessionCloseRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*SessionCloseRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionCloseRequest) LeadershipTermIdId() uint16 {
	return 1
}

func (*SessionCloseRequest) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (s *SessionCloseRequest) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LeadershipTermIdSinceVersion()
}

func (*SessionCloseRequest) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*SessionCloseRequest) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*SessionCloseRequest) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionCloseRequest) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionCloseRequest) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*SessionCloseRequest) ClusterSessionIdId() uint16 {
	return 2
}

func (*SessionCloseRequest) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (s *SessionCloseRequest) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ClusterSessionIdSinceVersion()
}

func (*SessionCloseRequest) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*SessionCloseRequest) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*SessionCloseRequest) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*SessionCloseRequest) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*SessionCloseRequest) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}
