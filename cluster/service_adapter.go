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

type serviceAdapter struct {
	marshaller   *codecs.SbeGoMarshaller
	agent        *ClusteredServiceAgent
	subscription *aeron.Subscription
}

func (adapter *serviceAdapter) poll() (int, error) {
	if adapter.subscription.IsClosed() {
		return 0, errors.New("subscription closed")
	}
	return adapter.subscription.Poll(adapter.onFragment, 10)
}

func (adapter *serviceAdapter) onFragment(
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

	switch templateId {
	case joinLogTemplateId:
		buf := &bytes.Buffer{}
		buffer.WriteBytes(buf, offset, length)
		joinLog := &codecs.JoinLog{}
		if err := joinLog.Decode(adapter.marshaller, buf, version, blockLength, true); err != nil {
			return err
		}
		adapter.agent.onJoinLog(
			joinLog.LogPosition,
			joinLog.MaxLogPosition,
			joinLog.MemberId,
			joinLog.LogSessionId,
			joinLog.LogStreamId,
			joinLog.IsStartup == codecs.BooleanType.TRUE,
			Role(joinLog.Role),
			string(joinLog.LogChannel),
		)
	case serviceTerminationPosTemplateId:
		logPos := buffer.GetInt64(offset)
		adapter.agent.onServiceTerminationPosition(logPos)
	default:
		// TODO: Is this an error?  Probably, but debug logging gives me pause.
		logger.Debugf("serviceAdapter: unexpected templateId=%d at pos=%d", templateId, header.Position())
	}
	return nil
}
