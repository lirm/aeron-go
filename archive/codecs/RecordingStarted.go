// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingStarted struct {
	RecordingId    int64
	StartPosition  int64
	SessionId      int32
	StreamId       int32
	Channel        []uint8
	SourceIdentity []uint8
}

func (r *RecordingStarted) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.StartPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.SessionId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.StreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.Channel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.Channel); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(r.SourceIdentity))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, r.SourceIdentity); err != nil {
		return err
	}
	return nil
}

func (r *RecordingStarted) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.RecordingIdInActingVersion(actingVersion) {
		r.RecordingId = r.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.RecordingId); err != nil {
			return err
		}
	}
	if !r.StartPositionInActingVersion(actingVersion) {
		r.StartPosition = r.StartPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.StartPosition); err != nil {
			return err
		}
	}
	if !r.SessionIdInActingVersion(actingVersion) {
		r.SessionId = r.SessionIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.SessionId); err != nil {
			return err
		}
	}
	if !r.StreamIdInActingVersion(actingVersion) {
		r.StreamId = r.StreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.StreamId); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}

	if r.ChannelInActingVersion(actingVersion) {
		var ChannelLength uint32
		if err := _m.ReadUint32(_r, &ChannelLength); err != nil {
			return err
		}
		if cap(r.Channel) < int(ChannelLength) {
			r.Channel = make([]uint8, ChannelLength)
		}
		r.Channel = r.Channel[:ChannelLength]
		if err := _m.ReadBytes(_r, r.Channel); err != nil {
			return err
		}
	}

	if r.SourceIdentityInActingVersion(actingVersion) {
		var SourceIdentityLength uint32
		if err := _m.ReadUint32(_r, &SourceIdentityLength); err != nil {
			return err
		}
		if cap(r.SourceIdentity) < int(SourceIdentityLength) {
			r.SourceIdentity = make([]uint8, SourceIdentityLength)
		}
		r.SourceIdentity = r.SourceIdentity[:SourceIdentityLength]
		if err := _m.ReadBytes(_r, r.SourceIdentity); err != nil {
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

func (r *RecordingStarted) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.RecordingIdInActingVersion(actingVersion) {
		if r.RecordingId < r.RecordingIdMinValue() || r.RecordingId > r.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.RecordingId (%v < %v > %v)", r.RecordingIdMinValue(), r.RecordingId, r.RecordingIdMaxValue())
		}
	}
	if r.StartPositionInActingVersion(actingVersion) {
		if r.StartPosition < r.StartPositionMinValue() || r.StartPosition > r.StartPositionMaxValue() {
			return fmt.Errorf("Range check failed on r.StartPosition (%v < %v > %v)", r.StartPositionMinValue(), r.StartPosition, r.StartPositionMaxValue())
		}
	}
	if r.SessionIdInActingVersion(actingVersion) {
		if r.SessionId < r.SessionIdMinValue() || r.SessionId > r.SessionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.SessionId (%v < %v > %v)", r.SessionIdMinValue(), r.SessionId, r.SessionIdMaxValue())
		}
	}
	if r.StreamIdInActingVersion(actingVersion) {
		if r.StreamId < r.StreamIdMinValue() || r.StreamId > r.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on r.StreamId (%v < %v > %v)", r.StreamIdMinValue(), r.StreamId, r.StreamIdMaxValue())
		}
	}
	return nil
}

func RecordingStartedInit(r *RecordingStarted) {
	return
}

func (*RecordingStarted) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*RecordingStarted) SbeTemplateId() (templateId uint16) {
	return 101
}

func (*RecordingStarted) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingStarted) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*RecordingStarted) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingStarted) RecordingIdId() uint16 {
	return 1
}

func (*RecordingStarted) RecordingIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingStarted) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.RecordingIdSinceVersion()
}

func (*RecordingStarted) RecordingIdDeprecated() uint16 {
	return 0
}

func (*RecordingStarted) RecordingIdMetaAttribute(meta int) string {
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

func (*RecordingStarted) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingStarted) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingStarted) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingStarted) StartPositionId() uint16 {
	return 2
}

func (*RecordingStarted) StartPositionSinceVersion() uint16 {
	return 0
}

func (r *RecordingStarted) StartPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StartPositionSinceVersion()
}

func (*RecordingStarted) StartPositionDeprecated() uint16 {
	return 0
}

func (*RecordingStarted) StartPositionMetaAttribute(meta int) string {
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

func (*RecordingStarted) StartPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingStarted) StartPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingStarted) StartPositionNullValue() int64 {
	return math.MinInt64
}

func (*RecordingStarted) SessionIdId() uint16 {
	return 3
}

func (*RecordingStarted) SessionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingStarted) SessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SessionIdSinceVersion()
}

func (*RecordingStarted) SessionIdDeprecated() uint16 {
	return 0
}

func (*RecordingStarted) SessionIdMetaAttribute(meta int) string {
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

func (*RecordingStarted) SessionIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingStarted) SessionIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingStarted) SessionIdNullValue() int32 {
	return math.MinInt32
}

func (*RecordingStarted) StreamIdId() uint16 {
	return 4
}

func (*RecordingStarted) StreamIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingStarted) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StreamIdSinceVersion()
}

func (*RecordingStarted) StreamIdDeprecated() uint16 {
	return 0
}

func (*RecordingStarted) StreamIdMetaAttribute(meta int) string {
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

func (*RecordingStarted) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingStarted) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingStarted) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*RecordingStarted) ChannelMetaAttribute(meta int) string {
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

func (*RecordingStarted) ChannelSinceVersion() uint16 {
	return 0
}

func (r *RecordingStarted) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ChannelSinceVersion()
}

func (*RecordingStarted) ChannelDeprecated() uint16 {
	return 0
}

func (RecordingStarted) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (RecordingStarted) ChannelHeaderLength() uint64 {
	return 4
}

func (*RecordingStarted) SourceIdentityMetaAttribute(meta int) string {
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

func (*RecordingStarted) SourceIdentitySinceVersion() uint16 {
	return 0
}

func (r *RecordingStarted) SourceIdentityInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.SourceIdentitySinceVersion()
}

func (*RecordingStarted) SourceIdentityDeprecated() uint16 {
	return 0
}

func (RecordingStarted) SourceIdentityCharacterEncoding() string {
	return "US-ASCII"
}

func (RecordingStarted) SourceIdentityHeaderLength() uint64 {
	return 4
}
