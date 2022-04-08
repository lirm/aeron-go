// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type NewLeaderEvent struct {
	LeadershipTermId int64
	ClusterSessionId int64
	LeaderMemberId   int32
	IngressEndpoints []uint8
}

func (n *NewLeaderEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := n.RangeCheck(n.SbeSchemaVersion(), n.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, n.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, n.ClusterSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, n.LeaderMemberId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(n.IngressEndpoints))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, n.IngressEndpoints); err != nil {
		return err
	}
	return nil
}

func (n *NewLeaderEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !n.LeadershipTermIdInActingVersion(actingVersion) {
		n.LeadershipTermId = n.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.LeadershipTermId); err != nil {
			return err
		}
	}
	if !n.ClusterSessionIdInActingVersion(actingVersion) {
		n.ClusterSessionId = n.ClusterSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &n.ClusterSessionId); err != nil {
			return err
		}
	}
	if !n.LeaderMemberIdInActingVersion(actingVersion) {
		n.LeaderMemberId = n.LeaderMemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &n.LeaderMemberId); err != nil {
			return err
		}
	}
	if actingVersion > n.SbeSchemaVersion() && blockLength > n.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-n.SbeBlockLength()))
	}

	if n.IngressEndpointsInActingVersion(actingVersion) {
		var IngressEndpointsLength uint32
		if err := _m.ReadUint32(_r, &IngressEndpointsLength); err != nil {
			return err
		}
		if cap(n.IngressEndpoints) < int(IngressEndpointsLength) {
			n.IngressEndpoints = make([]uint8, IngressEndpointsLength)
		}
		n.IngressEndpoints = n.IngressEndpoints[:IngressEndpointsLength]
		if err := _m.ReadBytes(_r, n.IngressEndpoints); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := n.RangeCheck(actingVersion, n.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (n *NewLeaderEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if n.LeadershipTermIdInActingVersion(actingVersion) {
		if n.LeadershipTermId < n.LeadershipTermIdMinValue() || n.LeadershipTermId > n.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LeadershipTermId (%v < %v > %v)", n.LeadershipTermIdMinValue(), n.LeadershipTermId, n.LeadershipTermIdMaxValue())
		}
	}
	if n.ClusterSessionIdInActingVersion(actingVersion) {
		if n.ClusterSessionId < n.ClusterSessionIdMinValue() || n.ClusterSessionId > n.ClusterSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on n.ClusterSessionId (%v < %v > %v)", n.ClusterSessionIdMinValue(), n.ClusterSessionId, n.ClusterSessionIdMaxValue())
		}
	}
	if n.LeaderMemberIdInActingVersion(actingVersion) {
		if n.LeaderMemberId < n.LeaderMemberIdMinValue() || n.LeaderMemberId > n.LeaderMemberIdMaxValue() {
			return fmt.Errorf("Range check failed on n.LeaderMemberId (%v < %v > %v)", n.LeaderMemberIdMinValue(), n.LeaderMemberId, n.LeaderMemberIdMaxValue())
		}
	}
	for idx, ch := range n.IngressEndpoints {
		if ch > 127 {
			return fmt.Errorf("n.IngressEndpoints[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func NewLeaderEventInit(n *NewLeaderEvent) {
	return
}

func (*NewLeaderEvent) SbeBlockLength() (blockLength uint16) {
	return 20
}

func (*NewLeaderEvent) SbeTemplateId() (templateId uint16) {
	return 6
}

func (*NewLeaderEvent) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*NewLeaderEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*NewLeaderEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*NewLeaderEvent) LeadershipTermIdId() uint16 {
	return 1
}

func (*NewLeaderEvent) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeaderEvent) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LeadershipTermIdSinceVersion()
}

func (*NewLeaderEvent) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*NewLeaderEvent) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*NewLeaderEvent) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeaderEvent) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeaderEvent) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*NewLeaderEvent) ClusterSessionIdId() uint16 {
	return 2
}

func (*NewLeaderEvent) ClusterSessionIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeaderEvent) ClusterSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.ClusterSessionIdSinceVersion()
}

func (*NewLeaderEvent) ClusterSessionIdDeprecated() uint16 {
	return 0
}

func (*NewLeaderEvent) ClusterSessionIdMetaAttribute(meta int) string {
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

func (*NewLeaderEvent) ClusterSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*NewLeaderEvent) ClusterSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*NewLeaderEvent) ClusterSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*NewLeaderEvent) LeaderMemberIdId() uint16 {
	return 3
}

func (*NewLeaderEvent) LeaderMemberIdSinceVersion() uint16 {
	return 0
}

func (n *NewLeaderEvent) LeaderMemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.LeaderMemberIdSinceVersion()
}

func (*NewLeaderEvent) LeaderMemberIdDeprecated() uint16 {
	return 0
}

func (*NewLeaderEvent) LeaderMemberIdMetaAttribute(meta int) string {
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

func (*NewLeaderEvent) LeaderMemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*NewLeaderEvent) LeaderMemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*NewLeaderEvent) LeaderMemberIdNullValue() int32 {
	return math.MinInt32
}

func (*NewLeaderEvent) IngressEndpointsMetaAttribute(meta int) string {
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

func (*NewLeaderEvent) IngressEndpointsSinceVersion() uint16 {
	return 0
}

func (n *NewLeaderEvent) IngressEndpointsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= n.IngressEndpointsSinceVersion()
}

func (*NewLeaderEvent) IngressEndpointsDeprecated() uint16 {
	return 0
}

func (NewLeaderEvent) IngressEndpointsCharacterEncoding() string {
	return "US-ASCII"
}

func (NewLeaderEvent) IngressEndpointsHeaderLength() uint64 {
	return 4
}
