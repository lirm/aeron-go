package cluster

import (
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer/term"
)

const (
	ClientSessionMockedOffer = 1
)

type ClientSession interface {
	Id() int64
	ResponseStreamId() int32
	ResponseChannel() string
	// EncodedPrinciple() []byte
	Close()
	// TODO: the other close methods are not part of interface.
	// I don't understand the closing bool implementation and why it is needed
	// Leaving out for now unless it is really important
	// IsClosing() bool
	Offer(*atomic.Buffer, int32, int32, term.ReservedValueSupplier) int64
	// TryClaim(...)
}

type ContainerClientSession struct {
	id               int64
	responseStreamId int32
	responseChannel  string
	agent            *ClusteredServiceAgent
	response         *aeron.Publication
}

func NewContainerClientSession(
	id int64,
	responseStreamId int32,
	responseChannel string,
	agent *ClusteredServiceAgent,
) *ContainerClientSession {
	return &ContainerClientSession{
		id:               id,
		responseStreamId: responseStreamId,
		responseChannel:  responseChannel,
		agent:            agent,
		response: <-agent.a.AddPublication(
			responseChannel,
			responseStreamId,
		),
	}
}

func (s *ContainerClientSession) Id() int64 {
	return s.id
}

func (s *ContainerClientSession) ResponseStreamId() int32 {
	return s.responseStreamId
}

func (s *ContainerClientSession) ResponseChannel() string {
	return s.responseChannel
}

func (s *ContainerClientSession) Close() {
	if _, ok := s.agent.getClientSession(s.id); ok {
		s.agent.closeClientSession(s.id)
	}
}

func (s *ContainerClientSession) Offer(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	reservedValueSupplier term.ReservedValueSupplier,
) int64 {
	return s.agent.Offer(
		s.id,
		s.response,
		buffer,
		offset,
		length,
		reservedValueSupplier,
	)
}
