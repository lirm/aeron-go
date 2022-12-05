// Copyright 2022 Steven Stern
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logging"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/cluster"
	"github.com/lirm/aeron-go/cluster/codecs"
)

var logger = logging.MustGetLogger("cluster-client")
var marshaller = codecs.NewSbeGoMarshaller()

var TemporaryError = errors.New("temporary error")

type AeronCluster struct {
	opts                 *Options
	aeronClient          *aeron.Aeron
	egressSub            *aeron.Subscription
	ingressChannel       *aeron.ChannelUri
	ingressPub           *aeron.Publication
	clusterSessionId     int64
	leadershipTermId     int64
	leaderMemberId       int32
	memberByIdMap        map[int32]*memberIngress
	fragmentAssembler    *aeron.FragmentAssembler
	egressListener       EgressListener
	sessionMsgHdrBuffer  *atomic.Buffer
	keepAliveBuffer      *atomic.Buffer
	state                clientState
	correlationId        int64
	nextRetryConnectTime int64
	awaitTimeoutTime     int64
}

type memberIngress struct {
	memberId    int32
	endpoint    string
	publication *aeron.Publication
}

type clientState int8

const (
	clientDisconnected clientState = iota
	clientCreatePublications
	clientAwaitPublicationConnected
	clientAwaitConnectReply
	clientConnected
	clientClosed
)

const (
	protocolMajorVersion = 0
	protocolMinorVersion = 2
	protocolPatchVersion = 0
)

var protocolSemanticVersion = util.SemanticVersionCompose(
	protocolMajorVersion, protocolMinorVersion, protocolPatchVersion)

func NewAeronCluster(
	aeronCtx *aeron.Context,
	options *Options,
	egressListener EgressListener,
) (*AeronCluster, error) {
	if egressListener == nil {
		return nil, fmt.Errorf("egressListener is nil")
	}
	ingressChannel, err := aeron.ParseChannelUri(options.IngressChannel)
	if err != nil {
		return nil, err
	}
	if ingressChannel.IsIpc() && options.IngressEndpoints != "" {
		return nil, fmt.Errorf("IngressEndpoints must be empty when using IPC ingress")
	}

	aeronClient, err := aeron.Connect(aeronCtx)
	if err != nil {
		return nil, err
	}

	egressSub, err := aeronClient.AddSubscription(options.EgressChannel, options.EgressStreamId)
	if err != nil {
		return nil, err
	}

	sessionMsgHdrBuf := codecs.MakeClusterMessageBuffer(cluster.SessionMessageHeaderTemplateId, cluster.SessionMessageHdrBlockLength)

	client := &AeronCluster{
		opts:                options,
		aeronClient:         aeronClient,
		egressSub:           egressSub,
		ingressChannel:      &ingressChannel,
		clusterSessionId:    cluster.NullValue,
		leadershipTermId:    cluster.NullValue,
		leaderMemberId:      cluster.NullValue,
		memberByIdMap:       make(map[int32]*memberIngress),
		egressListener:      egressListener,
		state:               clientDisconnected,
		sessionMsgHdrBuffer: sessionMsgHdrBuf,
		keepAliveBuffer:     codecs.MakeClusterMessageBuffer(cluster.SessionKeepAliveTemplateId, 16),
	}
	client.fragmentAssembler = aeron.NewFragmentAssembler(client.onFragment, 0)
	client.updateMemberEndpoints(options.IngressEndpoints)

	return client, nil
}

func (ac *AeronCluster) ClusterSessionId() int64 {
	return ac.clusterSessionId
}

func (ac *AeronCluster) LeadershipTermId() int64 {
	return ac.leadershipTermId
}

func (ac *AeronCluster) LeaderMemberId() int32 {
	return ac.leaderMemberId
}

func (ac *AeronCluster) IsConnected() bool {
	return ac.state == clientConnected
}

func (ac *AeronCluster) IsClosed() bool {
	return ac.state == clientClosed
}

func (ac *AeronCluster) Poll() int {
	switch ac.state {
	case clientDisconnected:
		if time.Now().UnixMilli() > ac.nextRetryConnectTime {
			ac.state = clientCreatePublications
		}
	case clientCreatePublications:
		ret, err := ac.createPublications()
		if err != nil {
			logger.Warningf("error from createPublications %w", err)
		}
		return ret
	case clientAwaitPublicationConnected:
		ret, err := ac.awaitPublicationConnected()
		if err != nil {
			logger.Warningf("error from awaitPublicationConnected %w", err)
		}
		return ret
	case clientAwaitConnectReply:
		now := time.Now().UnixMilli()
		if ac.ingressPub.IsConnected() && now < ac.awaitTimeoutTime {
			return ac.pollEgress(1)
		} else {
			logger.Warningf("timed out waiting for session connect reply")
			ac.state = clientDisconnected
			ac.nextRetryConnectTime = now + (30 * time.Second).Milliseconds()
		}
	case clientConnected:
		if ac.ingressPub.IsConnected() {
			return ac.pollEgress(10)
			// TODO: check if state == closed
		} else {
			ac.egressListener.OnDisconnect(ac, "ingress publication disconnected")
			ac.state = clientCreatePublications
		}
	}
	return 0
}

func (ac *AeronCluster) Offer(buffer *atomic.Buffer, offset, length int32) int64 {
	if ac.state != clientConnected {
		return aeron.NotConnected
	} else {
		hdrBuf := ac.sessionMsgHdrBuffer
		return ac.ingressPub.Offer2(hdrBuf, 0, hdrBuf.Capacity(), buffer, offset, length, nil)
	}
}

func (ac *AeronCluster) SendKeepAlive() bool {
	if !ac.IsConnected() {
		return false
	}
	buf := ac.keepAliveBuffer
	for i := 0; i < 3; i++ {
		if result := ac.ingressPub.Offer(buf, 0, buf.Capacity(), nil); result >= 0 {
			return true
		}
		ac.opts.IdleStrategy.Idle(0)
	}
	return false
}

func (ac *AeronCluster) Close() {
	if ac.IsConnected() && ac.ingressPub.IsConnected() {
		ac.sendCloseSession()
	}
	if ac.ingressPub != nil {
		if err := ac.ingressPub.Close(); err != nil {
			logger.Debugf("error closing ingress publication: %v", err)
		}
		ac.ingressPub = nil
	}
	if ac.egressSub != nil {
		if err := ac.egressSub.Close(); err != nil {
			logger.Debugf("error closing egress subscription: %v", err)
		}
		ac.egressSub = nil
	}
	if err := ac.aeronClient.Close(); err != nil {
		logger.Debugf("error closing aeron client: %v", err)
	}
	ac.state = clientClosed
}

func (ac *AeronCluster) updateMemberEndpoints(endpoints string) {
	if endpoints == "" {
		return
	}
	logger.Debugf("updateMemberEndpoints: %s", endpoints)
	for idx, endpoint := range strings.Split(endpoints, ",") {
		if delim := strings.IndexByte(endpoint, '='); delim > 0 {
			memberId, err := strconv.Atoi(endpoint[:delim])
			if err != nil {
				logger.Warningf("invalid endpoint at idx=%d: %s", idx, endpoint)
				continue
			}
			address := endpoint[delim+1:]
			member := ac.memberByIdMap[int32(memberId)]
			if member == nil {
				member = &memberIngress{
					memberId: int32(memberId),
					endpoint: address,
				}
				ac.memberByIdMap[member.memberId] = member
			} else if address != member.endpoint {
				member.endpoint = address
				if member.publication != nil {
					member.close()
				}
			}
			if member.memberId == ac.leaderMemberId {
				if member.publication == nil {
					pub, err := ac.createIngressPublication(address)
					if err == nil {
						member.publication = pub
					} else {
						logger.Warning(err)
					}
				}
				ac.ingressPub = member.publication
				ac.fragmentAssembler.Clear()
			}
		} else {
			logger.Warningf("endpoint at idx=%d missing '=' separator: %s", idx, endpoint)
		}
	}
}

func (ac *AeronCluster) createPublications() (int, error) {
	if len(ac.memberByIdMap) > 0 {
		for _, member := range ac.memberByIdMap {
			if member.publication == nil {
				pub, err := ac.createIngressPublication(member.endpoint)
				if err != nil {
					return 0, err
				}
				member.publication = pub
			}
		}
	} else if ac.ingressPub == nil {
		pub, err := ac.createIngressPublication(ac.opts.IngressChannel)
		if err != nil {
			return 0, err
		}
		ac.ingressPub = pub
	}
	ac.state = clientAwaitPublicationConnected
	ac.awaitTimeoutTime = time.Now().UnixMilli() + (5 * time.Second).Milliseconds()
	return 1, nil
}

func (ac *AeronCluster) createIngressPublication(endpoint string) (*aeron.Publication, error) {
	if ac.ingressChannel.IsUdp() {
		ac.ingressChannel.Set("endpoint", endpoint)
	}
	channel := ac.ingressChannel.String()
	logger.Debugf("createIngressPublication - endpoint=%s isUdp=%v isExclusive=%v",
		endpoint, ac.ingressChannel.IsUdp(), ac.opts.IsIngressExclusive)
	if ac.opts.IsIngressExclusive {
		return ac.aeronClient.AddExclusivePublication(channel, ac.opts.IngressStreamId)
	} else {
		return ac.aeronClient.AddPublication(channel, ac.opts.IngressStreamId)
	}
}

func (ac *AeronCluster) awaitPublicationConnected() (int, error) {
	responseChannel := ac.egressSub.TryResolveChannelEndpointPort()
	if responseChannel == "" {
		// TODO: Is this an error or success condition?
		return 0, nil
	}
	now := time.Now().UnixMilli()
	if now > ac.awaitTimeoutTime {
		ac.state = clientDisconnected
		// close publications? shouldn't be necessary unless we've hit some bug
		ac.nextRetryConnectTime = now + (30 * time.Second).Milliseconds()
		return 0, errors.New("timed out waiting for connected publication")
	}
	if len(ac.memberByIdMap) > 0 {
		for _, member := range ac.memberByIdMap {
			if member.publication != nil && member.publication.IsConnected() {
				ac.ingressPub = member.publication
				ac.fragmentAssembler.Clear()
				err := ac.sendConnectRequest(responseChannel)
				if err == nil {
					logger.Debugf("sent connect request to memberId=%d correlationId=%d channel=%s",
						member.memberId, ac.correlationId, member.publication.Channel())
					ac.state = clientAwaitConnectReply
					ac.awaitTimeoutTime = now + (3 * time.Second).Milliseconds()
					break
				}
				if !errors.Is(err, TemporaryError) {
					return 0, err
				}
			}
		}
	} else if ac.ingressPub.IsConnected() && ac.sendConnectRequest(responseChannel) == nil {
		ac.state = clientAwaitConnectReply
		ac.awaitTimeoutTime = now + (3 * time.Second).Milliseconds()
	}
	return 0, nil
}

// Returns nil on success, TemporaryError, or any other error is a permanent error.
func (ac *AeronCluster) sendConnectRequest(responseChannel string) error {
	ac.correlationId = ac.aeronClient.NextCorrelationID()
	req := codecs.SessionConnectRequest{
		CorrelationId:    ac.correlationId,
		ResponseStreamId: ac.opts.EgressStreamId,
		Version:          int32(protocolSemanticVersion),
		ResponseChannel:  []byte(responseChannel),
	}
	header := codecs.MessageHeader{
		BlockLength: req.SbeBlockLength(),
		TemplateId:  req.SbeTemplateId(),
		SchemaId:    req.SbeSchemaId(),
		Version:     req.SbeSchemaVersion(),
	}
	writer := new(bytes.Buffer)
	if err := header.Encode(marshaller, writer); err != nil {
		return err
	}
	if err := req.Encode(marshaller, writer, true); err != nil {
		return err
	}
	buffer := atomic.MakeBuffer(writer.Bytes())
	result := ac.ingressPub.Offer(buffer, 0, buffer.Capacity(), nil)
	if result >= 0 {
		return nil
	} else {
		return fmt.Errorf("%w, failed to send connect request, channel=%s result=%d",
			TemporaryError, ac.ingressPub.Channel(), result)
	}
}

func (ac *AeronCluster) pollEgress(fragmentLimit int) int {
	return ac.egressSub.Poll(ac.fragmentAssembler.OnFragment, fragmentLimit)
}

func (ac *AeronCluster) onFragment(buffer *atomic.Buffer, offset, length int32, header *logbuffer.Header) {
	if length < cluster.SBEHeaderLength {
		return
	}
	blockLength := buffer.GetUInt16(offset)
	templateId := buffer.GetUInt16(offset + 2)
	schemaId := buffer.GetUInt16(offset + 4)
	version := buffer.GetUInt16(offset + 6)
	if schemaId != cluster.ClusterSchemaId {
		logger.Errorf("unexpected schemaId=%d templateId=%d blockLen=%d version=%d",
			schemaId, templateId, blockLength, version)
		return
	}
	offset += cluster.SBEHeaderLength
	length -= cluster.SBEHeaderLength

	switch templateId {
	case cluster.SessionMessageHeaderTemplateId:
		ac.onSessionMessage(buffer, offset, length, header)
	case cluster.SessionEventTemplateId:
		ac.onSessionEvent(buffer, offset, length, version, blockLength)
	case cluster.NewLeaderEventTemlateId:
		ac.onNewLeaderEvent(buffer, offset, length, version, blockLength)
	case cluster.ChallengeTemplateId:
		e := codecs.Challenge{}
		buf := bytes.Buffer{}
		buffer.WriteBytes(&buf, offset, length)
		if err := e.Decode(marshaller, &buf, version, blockLength, true); err != nil {
			logger.Errorf("new leader event decode error: %v", err)
		} else {
			logger.Warningf("received challenge, corrId=%d clusterSessionId=%d", e.CorrelationId, e.ClusterSessionId)
		}
	}
}

func (ac *AeronCluster) onNewLeaderEvent(buffer *atomic.Buffer, offset, length int32, version, blockLength uint16) {
	e := codecs.NewLeaderEvent{}
	buf := bytes.Buffer{}
	buffer.WriteBytes(&buf, offset, length)
	if err := e.Decode(marshaller, &buf, version, blockLength, true); err != nil {
		logger.Errorf("new leader event decode error: %v", err)
	} else if ac.state == clientConnected && e.ClusterSessionId == ac.clusterSessionId {
		ac.leadershipTermId = e.LeadershipTermId
		ac.leaderMemberId = e.LeaderMemberId
		ac.sessionMsgHdrBuffer.PutInt64(cluster.SBEHeaderLength, e.LeadershipTermId)
		ac.keepAliveBuffer.PutInt64(cluster.SBEHeaderLength, e.LeadershipTermId)
		if ac.opts.IngressEndpoints != "" {
			if err := ac.ingressPub.Close(); err != nil {
				logger.Warningf("error closing ingress publication: %v", err)
			}
			ac.ingressPub = nil
			ac.updateMemberEndpoints(string(e.IngressEndpoints))
		}
		ac.fragmentAssembler.Clear()
		ac.egressListener.OnNewLeader(ac, e.LeadershipTermId, e.LeaderMemberId)
	} else {
		logger.Debugf("ignored new leader event - state=%v thisSessionId=%d targetSessionId=%d leaderMemberId=%d leaderTermId=%d",
			ac.state, ac.clusterSessionId, e.ClusterSessionId, e.LeaderMemberId, e.LeadershipTermId)
	}
}

func (ac *AeronCluster) onSessionEvent(buffer *atomic.Buffer, offset, length int32, version, blockLength uint16) {
	e := codecs.SessionEvent{}
	buf := bytes.Buffer{}
	buffer.WriteBytes(&buf, offset, length)
	if err := e.Decode(marshaller, &buf, version, blockLength, true); err != nil {
		logger.Errorf("session event decode error: %v", err)
	} else if ac.state == clientAwaitConnectReply && e.CorrelationId == ac.correlationId {
		switch e.Code {
		case codecs.EventCode.OK:
			ac.leadershipTermId = e.LeadershipTermId
			ac.leaderMemberId = e.LeaderMemberId
			ac.clusterSessionId = e.ClusterSessionId
			ac.sessionMsgHdrBuffer.PutInt64(cluster.SBEHeaderLength, e.LeadershipTermId)
			ac.sessionMsgHdrBuffer.PutInt64(cluster.SBEHeaderLength+8, e.ClusterSessionId)
			ac.keepAliveBuffer.PutInt64(cluster.SBEHeaderLength, e.LeadershipTermId)
			ac.keepAliveBuffer.PutInt64(cluster.SBEHeaderLength+8, e.ClusterSessionId)
			ac.state = clientConnected
			ac.closeNonLeaderPublications()
			ac.egressListener.OnConnect(ac)
		case codecs.EventCode.REDIRECT:
			logger.Infof("got redirect - leaderTermId=%d leaderMemberId=%d", e.LeadershipTermId, e.LeaderMemberId)
			ac.leaderMemberId = e.LeaderMemberId
			ac.updateMemberEndpoints(string(e.Detail))
			ac.closeNonLeaderPublications()
			ac.state = clientAwaitPublicationConnected
		case codecs.EventCode.ERROR:
			ac.egressListener.OnError(ac, string(e.Detail))
			ac.state = clientDisconnected
			ac.nextRetryConnectTime = time.Now().UnixMilli() + (5 * time.Second).Milliseconds()
		case codecs.EventCode.AUTHENTICATION_REJECTED:
			ac.egressListener.OnError(ac, fmt.Sprintf("authentication rejected (%s)", string(e.Detail)))
			ac.state = clientDisconnected
			ac.nextRetryConnectTime = time.Now().UnixMilli() + time.Minute.Milliseconds()
		}
	} else if ac.state == clientConnected && e.ClusterSessionId == ac.clusterSessionId {
		if e.Code == codecs.EventCode.CLOSED {
			ac.egressListener.OnDisconnect(ac, string(e.Detail))
			ac.state = clientClosed
		} else if e.Code == codecs.EventCode.ERROR {
			ac.egressListener.OnError(ac, string(e.Detail))
		} else {
			logger.Infof("onSessionEvent - code=%v (%s)", e.Code, string(e.Detail))
		}
	} else {
		logger.Debugf("ignored session event - state=%v thisSessionId=%d targetSessionId=%d code=%d (%s)",
			ac.state, ac.clusterSessionId, e.ClusterSessionId, e.Code, string(e.Detail))
	}
}

func (ac *AeronCluster) closeNonLeaderPublications() {
	for _, member := range ac.memberByIdMap {
		if member.memberId != ac.leaderMemberId {
			member.close()
		}
	}
}

func (ac *AeronCluster) onSessionMessage(buffer *atomic.Buffer, offset, length int32, header *logbuffer.Header) {
	if length < cluster.SessionMessageHeaderLength {
		logger.Errorf("received invalid session message - length: %d", length)
		return
	}
	leadershipTermId := buffer.GetInt64(offset)
	clusterSessionId := buffer.GetInt64(offset + 8)
	timestamp := buffer.GetInt64(offset + 16)
	if ac.state != clientConnected {
		logger.Debugf("received unexpected session message - leadershipTermId=%d clusterSessionId=%d state=%v",
			leadershipTermId, clusterSessionId, ac.state)
	} else if clusterSessionId == ac.clusterSessionId {
		ac.egressListener.OnMessage(ac, timestamp, buffer, offset+cluster.SessionMessageHeaderLength,
			length-cluster.SessionMessageHeaderLength, header)
	} else {
		logger.Debugf("received unexpected session msg - leaderTermId=%d targetSessionId=%d thisSessionId=%d",
			leadershipTermId, clusterSessionId, ac.clusterSessionId)
	}
}

func (ac *AeronCluster) sendCloseSession() {
	buf := codecs.MakeClusterMessageBuffer(cluster.SessionCloseRequestTemplateId, 16)
	buf.PutInt64(cluster.SBEHeaderLength, ac.leadershipTermId)
	buf.PutInt64(cluster.SBEHeaderLength+8, ac.clusterSessionId)
	for i := 0; i < 3; i++ {
		if result := ac.ingressPub.Offer(buf, 0, buf.Capacity(), nil); result >= 0 {
			return
		}
		ac.opts.IdleStrategy.Idle(0)
	}
}

func (member *memberIngress) close() {
	if member.publication != nil {
		if err := member.publication.Close(); err != nil {
			logger.Warningf("error closing member publication, memberId=%d endpoint=%s: %v",
				member.memberId, member.endpoint, err)
		}
		member.publication = nil
	}
}
