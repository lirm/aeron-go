package cluster

import (
	"bytes"
	"fmt"
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/cluster/codecs"
)

type snapshotLoader struct {
	agent      *ClusteredServiceAgent
	img        *aeron.Image
	marshaller *codecs.SbeGoMarshaller
	isDone     bool
	inSnapshot bool
	appVersion int32
	timeUnit   codecs.ClusterTimeUnitEnum
	err        error
}

func newSnapshotLoader(agent *ClusteredServiceAgent, img *aeron.Image) *snapshotLoader {
	return &snapshotLoader{agent: agent, img: img, marshaller: codecs.NewSbeGoMarshaller()}
}

func (loader *snapshotLoader) poll() int {
	if loader.isDone {
		return 0
	}
	return loader.img.Poll(loader.onFragment, 1)
}

func (loader *snapshotLoader) onFragment(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	var hdr codecs.SbeGoMessageHeader
	buf := &bytes.Buffer{}
	buffer.WriteBytes(buf, offset, length)

	if err := hdr.Decode(loader.marshaller, buf); err != nil {
		fmt.Println("header decode error: ", err)
		return
	}

	if hdr.SchemaId != clusterSchemaId {
		return
	}

	switch hdr.TemplateId {
	case snapshotMarkerTemplateId:
		marker := &codecs.SnapshotMarker{}
		if err := marker.Decode(loader.marshaller, buf, hdr.Version, hdr.BlockLength, true); err != nil {
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
		if err := s.Decode(loader.marshaller, buf, hdr.Version, hdr.BlockLength, true); err != nil {
			loader.err = err
			loader.isDone = true
		} else {
			loader.agent.addSessionFromSnapshot(NewContainerClientSession(s.ClusterSessionId,
				s.ResponseStreamId, string(s.ResponseChannel), loader.agent))
		}
	default:
		//fmt.Println("SnapshotLoader: unknown template id: ", hdr.TemplateId)
	}
}
