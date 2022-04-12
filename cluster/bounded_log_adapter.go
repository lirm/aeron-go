package cluster

import (
	"bytes"
	"fmt"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/codecs"
)

const (
	schemaId                       = 111
	sessionMessageHeaderTemplateId = 1
	timerEventTemplateId           = 20
	sessionOpenTemplateId          = 21
	sessionCloseTemplateId         = 22
	clusterActionReqTemplateId     = 23
	newLeadershipTermTemplateId    = 24
	membershipChangeTemplateId     = 25
)

type BoundedLogAdapter struct {
	marshaller     *codecs.SbeGoMarshaller
	options        *Options
	agent          *ClusteredServiceAgent
	image          *aeron.Image
	maxLogPosition int64
}

func (adapter *BoundedLogAdapter) IsDone() bool {
	return adapter.image.Position() >= adapter.maxLogPosition || adapter.image.IsEndOfStream() || adapter.image.IsClosed()
}

func (adapter *BoundedLogAdapter) Poll() int {
	return adapter.image.BoundedPoll(adapter.onFragment, adapter.maxLogPosition, adapter.options.LogFragmentLimit)
}

func (adapter *BoundedLogAdapter) onFragment(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	var hdr codecs.SbeGoMessageHeader
	buf := &bytes.Buffer{}
	buffer.WriteBytes(buf, offset, length)

	if err := hdr.Decode(adapter.marshaller, buf); err != nil {
		fmt.Println("BoundedLogAdaptor - header decode error: ", err)
	}
	if hdr.SchemaId != schemaId {
		fmt.Println("BoundedLogAdaptor - unexpected schemaId: ", hdr)
		return
	}

	switch hdr.TemplateId {
	case timerEventTemplateId:
		fmt.Println("BoundedLogAdaptor - got timer event")
	case sessionOpenTemplateId:
		open := &codecs.SessionOpenEvent{}
		if err := open.Decode(
			adapter.marshaller,
			buf,
			hdr.Version,
			hdr.BlockLength,
			adapter.options.RangeChecking,
		); err != nil {
			fmt.Println("session open decode error: ", err)
		} else {
			fmt.Println("BoundedLogAdaptor - got session open: ", open)
		}
	case sessionCloseTemplateId:
		ce := &codecs.SessionCloseEvent{}
		if err := ce.Decode(adapter.marshaller, buf, hdr.Version, hdr.BlockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("session close decode error: ", err)
		} else {
			fmt.Println("BoundedLogAdaptor - got session close: ", ce)
		}
	case clusterActionReqTemplateId:
		e := &codecs.ClusterActionRequest{}
		if err := e.Decode(adapter.marshaller, buf, hdr.Version, hdr.BlockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("cluster action request decode error: ", err)
		} else {
			fmt.Println("BoundedLogAdaptor - got cluster action request: ", e)
		}
	case newLeadershipTermTemplateId:
		e := &codecs.NewLeadershipTermEvent{}
		if err := e.Decode(adapter.marshaller, buf, hdr.Version, hdr.BlockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("new leadership term decode error: ", err)
		} else {
			fmt.Println("BoundedLogAdaptor - got new leadership term: ", e)
			adapter.agent.onNewLeadershipTermEvent(e.LeadershipTermId, e.LogPosition, e.Timestamp, e.TermBaseLogPosition,
				e.LeaderMemberId, e.LogSessionId, e.AppVersion)
		}
	case membershipChangeTemplateId:
		e := &codecs.MembershipChangeEvent{}
		if err := e.Decode(adapter.marshaller, buf, hdr.Version, hdr.BlockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("membership change event decode error: ", err)
		} else {
			fmt.Println("BoundedLogAdaptor - got membership change event: ", e)
		}
	case sessionMessageHeaderTemplateId:
		e := &codecs.SessionMessageHeader{}
		if err := e.Decode(adapter.marshaller, buf, hdr.Version, hdr.BlockLength, adapter.options.RangeChecking); err != nil {
			fmt.Println("session message header decode error: ", err)
		} else {
			fmt.Println("BoundedLogAdaptor - got session message: ", e)
		}
	default:
		fmt.Println("unexpected template id: ", hdr.TemplateId)
	}
}

func (adapter *BoundedLogAdapter) Close() error {
	var err error
	if adapter.image != nil {
		err = adapter.image.Close()
		adapter.image = nil
	}
	return err
}
