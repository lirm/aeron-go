package cluster

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type snapshotLoader struct {
	agent      *ClusteredServiceAgent
	img        aeron.Image
	marshaller *codecs.SbeGoMarshaller
	isDone     bool
	inSnapshot bool
	appVersion int32
	timeUnit   codecs.ClusterTimeUnitEnum
	err        error
}

func newSnapshotLoader(agent *ClusteredServiceAgent, img aeron.Image) *snapshotLoader {
	return &snapshotLoader{agent: agent, img: img, marshaller: codecs.NewSbeGoMarshaller()}
}

func (loader *snapshotLoader) poll() (int, error) {
	if loader.isDone {
		return 0, nil
	}
	return loader.img.Poll(loader.onFragment, 1)
}

func (loader *snapshotLoader) onFragment(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) error {
	if length < SBEHeaderLength {
		return errors.New("fragment is too small to process")
	}
	blockLength := buffer.GetUInt16(offset)
	templateId := buffer.GetUInt16(offset + 2)
	schemaId := buffer.GetUInt16(offset + 4)
	version := buffer.GetUInt16(offset + 6)
	if schemaId != ClusterSchemaId {
		return fmt.Errorf("unexpected fragment with schemaId=%d templateId=%d blockLen=%d version=%d",
			schemaId, templateId, blockLength, version)
	}
	offset += SBEHeaderLength
	length -= SBEHeaderLength

	buf := &bytes.Buffer{}
	buffer.WriteBytes(buf, offset, length)

	// TODO Rework loader.err so the error is returned directly.  No need for this fancy error redirection once I make
	// FragmentHandler and Poll both return errors.
	switch templateId {
	case snapshotMarkerTemplateId:
		marker := &codecs.SnapshotMarker{}
		if err := marker.Decode(loader.marshaller, buf, version, blockLength, true); err != nil {
			loader.err = err
			loader.isDone = true
		} else if marker.TypeId != snapshotTypeId {
			loader.err = fmt.Errorf("unexpected snapshot type: %d", marker.TypeId)
		} else if marker.Mark == codecs.SnapshotMark.BEGIN {
			if loader.inSnapshot {
				loader.err = fmt.Errorf("already in snapshot, pos=%d", header.Position())
				loader.isDone = true
			} else {
				loader.inSnapshot = true
				loader.appVersion = marker.AppVersion
				loader.timeUnit = marker.TimeUnit
			}
		} else if marker.Mark == codecs.SnapshotMark.END {
			if !loader.inSnapshot {
				loader.err = fmt.Errorf("missing beging snapshot, pos=%d", header.Position())
				loader.isDone = true
			} else {
				loader.inSnapshot = false
				loader.isDone = true
			}
		} else {
			loader.err = fmt.Errorf("unexpected snapshot mark, pos=%d inSnapshot=%v mark=%v", header.Position(), loader.inSnapshot, marker.Mark)
			loader.isDone = true
		}
	case clientSessionTemplateId:
		s := codecs.ClientSession{}
		if err := s.Decode(loader.marshaller, buf, version, blockLength, true); err != nil {
			loader.err = err
			loader.isDone = true
		} else {
			loader.agent.addSessionFromSnapshot(newContainerClientSession(s.ClusterSessionId,
				s.ResponseStreamId, string(s.ResponseChannel), s.EncodedPrincipal, loader.agent))
		}
	default:
		logger.Debugf("SnapshotLoader: unknown templateId=%d at pos=%d", templateId, header.Position())
	}
	return nil
}
