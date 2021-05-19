// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ReplicateRequest struct {
	ControlSessionId   int64
	CorrelationId      int64
	SrcRecordingId     int64
	DstRecordingId     int64
	SrcControlStreamId int32
	SrcControlChannel  []uint8
	LiveDestination    []uint8
}

func (r *ReplicateRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.SrcRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.DstRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.SrcControlStreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.SrcControlChannel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.SrcControlChannel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.LiveDestination))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.LiveDestination); err != nil {
		return err
	}
	return nil
}

func (r *ReplicateRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.ControlSessionIdInActingVersion(actingVersion) {
		r.ControlSessionId = r.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.ControlSessionId); err != nil {
			return err
		}
	}
	if !r.CorrelationIdInActingVersion(actingVersion) {
		r.CorrelationId = r.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.CorrelationId); err != nil {
			return err
		}
	}
	if !r.SrcRecordingIdInActingVersion(actingVersion) {
		r.SrcRecordingId = r.SrcRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.SrcRecordingId); err != nil {
			return err
		}
	}
	if !r.DstRecordingIdInActingVersion(actingVersion) {
		r.DstRecordingId = r.DstRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.DstRecordingId); err != nil {
			return err
		}
	}
	if !r.SrcControlStreamIdInActingVersion(actingVersion) {
		r.SrcControlStreamId = r.SrcControlStreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.SrcControlStreamId); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}

	if r.SrcControlChannelInActingVersion(actingVersion) {
		var SrcControlChannelLength uint32
		if err := _m.ReadUint32(_r, &SrcControlChannelLength); err != nil {
			return err
		}
		if cap(r.SrcControlChannel) < int(SrcControlChannelLength) {
			r.SrcControlChannel = make([]uint8, SrcControlChannelLength)
		}
		r.SrcControlChannel = r.SrcControlChannel[:SrcControlChannelLength]
		if err := _m.ReadBytes(_r, r.SrcControlChannel); err != nil {
			return err
		}
	}

	if r.LiveDestinationInActingVersion(actingVersion) {
		var LiveDestinationLength uint32
		if err := _m.ReadUint32(_r, &LiveDestinationLength); err != nil {
			return err
		}
		if cap(r.LiveDestination) < int(LiveDestinationLength) {
			r.LiveDestination = make([]uint8, LiveDestinationLength)
		}
		r.LiveDestination = r.LiveDestination[:LiveDestinationLength]
		if err := _m.ReadBytes(_r, r.LiveDestination); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := r.RangeCheck(actingVersion, r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReplicateRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.ControlSessionIdInActingVersion(actingVersion) {
		if r.ControlSessionId < r.ControlSessionIdMinValue() || r.ControlSessionId > r.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.ControlSessionId (%v < %v > %v)", r.ControlSessionIdMinValue(), r.ControlSessionId, r.ControlSessionIdMaxValue())
		}
	}
	if r.CorrelationIdInActingVersion(actingVersion) {
		if r.CorrelationId < r.CorrelationIdMinValue() || r.CorrelationId > r.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on r.CorrelationId (%v < %v > %v)", r.CorrelationIdMinValue(), r.CorrelationId, r.CorrelationIdMaxValue())
		}
	}
	if r.SrcRecordingIdInActingVersion(actingVersion) {
		if r.SrcRecordingId < r.SrcRecordingIdMinValue() || r.SrcRecordingId > r.SrcRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SrcRecordingId (%v < %v > %v)", r.SrcRecordingIdMinValue(), r.SrcRecordingId, r.SrcRecordingIdMaxValue())
		}
	}
	if r.DstRecordingIdInActingVersion(actingVersion) {
		if r.DstRecordingId < r.DstRecordingIdMinValue() || r.DstRecordingId > r.DstRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.DstRecordingId (%v < %v > %v)", r.DstRecordingIdMinValue(), r.DstRecordingId, r.DstRecordingIdMaxValue())
		}
	}
	if r.SrcControlStreamIdInActingVersion(actingVersion) {
		if r.SrcControlStreamId < r.SrcControlStreamIdMinValue() || r.SrcControlStreamId > r.SrcControlStreamIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SrcControlStreamId (%v < %v > %v)", r.SrcControlStreamIdMinValue(), r.SrcControlStreamId, r.SrcControlStreamIdMaxValue())
		}
	}
	return nil
}

func ReplicateRequestInit(r *ReplicateRequest) {
	return
}

func (*ReplicateRequest) SbeBlockLength() (blockLength uint16) {
	return 36
}

func (*ReplicateRequest) SbeTemplateId() (templateId uint16) {
	return 50
}

func (*ReplicateRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ReplicateRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ReplicateRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ReplicateRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*ReplicateRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ControlSessionIdSinceVersion()
}

func (*ReplicateRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest) CorrelationIdId() uint16 {
	return 2
}

func (*ReplicateRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CorrelationIdSinceVersion()
}

func (*ReplicateRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest) SrcRecordingIdId() uint16 {
	return 3
}

func (*ReplicateRequest) SrcRecordingIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest) SrcRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SrcRecordingIdSinceVersion()
}

func (*ReplicateRequest) SrcRecordingIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest) SrcRecordingIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest) SrcRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest) SrcRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest) SrcRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest) DstRecordingIdId() uint16 {
	return 4
}

func (*ReplicateRequest) DstRecordingIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest) DstRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.DstRecordingIdSinceVersion()
}

func (*ReplicateRequest) DstRecordingIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest) DstRecordingIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest) DstRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ReplicateRequest) DstRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ReplicateRequest) DstRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ReplicateRequest) SrcControlStreamIdId() uint16 {
	return 5
}

func (*ReplicateRequest) SrcControlStreamIdSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest) SrcControlStreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SrcControlStreamIdSinceVersion()
}

func (*ReplicateRequest) SrcControlStreamIdDeprecated() uint16 {
	return 0
}

func (*ReplicateRequest) SrcControlStreamIdMetaAttribute(meta int) string {
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

func (*ReplicateRequest) SrcControlStreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ReplicateRequest) SrcControlStreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ReplicateRequest) SrcControlStreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ReplicateRequest) SrcControlChannelMetaAttribute(meta int) string {
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

func (*ReplicateRequest) SrcControlChannelSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest) SrcControlChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SrcControlChannelSinceVersion()
}

func (*ReplicateRequest) SrcControlChannelDeprecated() uint16 {
	return 0
}

func (ReplicateRequest) SrcControlChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ReplicateRequest) SrcControlChannelHeaderLength() uint64 {
	return 4
}

func (*ReplicateRequest) LiveDestinationMetaAttribute(meta int) string {
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

func (*ReplicateRequest) LiveDestinationSinceVersion() uint16 {
	return 0
}

func (r *ReplicateRequest) LiveDestinationInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.LiveDestinationSinceVersion()
}

func (*ReplicateRequest) LiveDestinationDeprecated() uint16 {
	return 0
}

func (ReplicateRequest) LiveDestinationCharacterEncoding() string {
	return "US-ASCII"
}

func (ReplicateRequest) LiveDestinationHeaderLength() uint64 {
	return 4
}
