// Copyright (C) 2021-2022 Talos, Inc.
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
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/cluster/codecs"
)

// Proxy class for encapsulating encoding and sending of control protocol messages to a cluster
type ConsensusModuleProxy struct {
	marshaller  *codecs.SbeGoMarshaller // currently shared as we're not reentrant (but could be here)
	options     *Options
	publication *aeron.Publication
}

func NewConsensusModuleProxy(
	options *Options,
	publication *aeron.Publication,
) *ConsensusModuleProxy {
	return &ConsensusModuleProxy{
		marshaller:  codecs.NewSbeGoMarshaller(),
		options:     options,
		publication: publication,
	}
}

// Offer to our request publication
func (proxy *ConsensusModuleProxy) Offer(
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	reservedValueSupplier term.ReservedValueSupplier,
) int64 {
	start := time.Now()
	var ret int64
	for time.Since(start) < proxy.options.Timeout {
		ret = proxy.publication.Offer(buffer, offset, length, reservedValueSupplier)
		switch ret {
		// Retry on these
		case aeron.NotConnected, aeron.BackPressured, aeron.AdminAction:
			proxy.options.IdleStrategy.Idle(0)
		// Fail or succeed on other values
		default:
			return ret
		}
	}

	// Give up, returning the last failure
	// logger.Debugf("ConsensusModuleProxy.Offer timing out [%d]", ret)
	return ret
}

// From here we have all the functions that create a data packet and send it on the
// publication. Responses will be processed on the control

// ConnectRequest packet and offer
func (proxy *ConsensusModuleProxy) ServiceAckRequest(
	logPosition int64,
	timestamp int64,
	ackID int64,
	relevantID int64,
	serviceID int32,
) error {
	// Create a packet and send it
	bytes, err := codecs.ServiceAckRequestPacket(
		proxy.marshaller,
		proxy.options.RangeChecking,
		logPosition,
		timestamp,
		ackID,
		relevantID,
		serviceID,
	)

	if err != nil {
		return err
	}

	buffer := atomic.MakeBuffer(bytes, len(bytes))
	claim := &logbuffer.Claim{}
	claim.Wrap(buffer, 0, int32(len(bytes)))
	proxy.publication.TryClaim(int32(len(bytes)), claim)
	claim.Commit()
	/*
		if ret := proxy.Offer(atomic.MakeBuffer(bytes, len(bytes)), 0, int32(len(bytes)), nil); ret < 0 {
			return fmt.Errorf("ConsensusModuleProxy.Offer failed: %d", ret)
		}
	*/

	return nil
}
