// Copyright (C) 2026 Talos, Inc.
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

package clustertests

import (
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/lirm/aeron-go/aeron"
	atomicbuf "github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster"
	"github.com/lirm/aeron-go/cluster/codecs"
)

// echoService echoes every session message back to its session, mirroring
// upstream's EchoService used by ClusterNodeTest. Messages of the form
// "timer:<id>:<delayMs>" and "cancel-timer:<id>" additionally schedule or
// cancel a cluster timer; expired timers are reported on timerEvents.
type echoService struct {
	cluster      cluster.Cluster
	messageCount atomic.Int32
	snapshots    atomic.Int32
	roleChanges  chan cluster.Role
	timerEvents  chan int64
}

func newEchoService() *echoService {
	return &echoService{
		roleChanges: make(chan cluster.Role, 10),
		timerEvents: make(chan int64, 16),
	}
}

func (s *echoService) handleTimerCommand(payload string) {
	if id, ok := strings.CutPrefix(payload, "cancel-timer:"); ok {
		correlationId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Errorf("echoService: bad cancel-timer command %q: %v", payload, err)
			return
		}
		for !s.cluster.CancelTimer(correlationId) {
			s.cluster.IdleStrategy().Idle(0)
		}
		return
	}
	if args, ok := strings.CutPrefix(payload, "timer:"); ok {
		idStr, delayStr, found := strings.Cut(args, ":")
		if !found {
			logger.Errorf("echoService: bad timer command %q", payload)
			return
		}
		correlationId, err1 := strconv.ParseInt(idStr, 10, 64)
		delayMs, err2 := strconv.ParseInt(delayStr, 10, 64)
		if err1 != nil || err2 != nil {
			logger.Errorf("echoService: bad timer command %q: %v %v", payload, err1, err2)
			return
		}
		for !s.cluster.ScheduleTimer(correlationId, s.cluster.Time()+delayMs) {
			s.cluster.IdleStrategy().Idle(0)
		}
	}
}

func (s *echoService) OnStart(cluster cluster.Cluster, image aeron.Image) {
	s.cluster = cluster
	if image != nil {
		image.Poll(func(buf *atomicbuf.Buffer, offset int32, length int32, hdr *logbuffer.Header) {
			if length == 4 {
				s.messageCount.Store(buf.GetInt32(offset))
			}
		}, 100)
	}
}

func (s *echoService) OnSessionOpen(session cluster.ClientSession, timestamp int64) {}

func (s *echoService) OnSessionClose(
	session cluster.ClientSession,
	timestamp int64,
	closeReason codecs.CloseReasonEnum,
) {
}

func (s *echoService) OnSessionMessage(
	session cluster.ClientSession,
	timestamp int64,
	buffer *atomicbuf.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	s.messageCount.Add(1)
	s.handleTimerCommand(string(buffer.GetBytesArray(offset, length)))
	for {
		result := session.Offer(buffer, offset, length, nil)
		if result >= 0 {
			return
		}
		if result != aeron.BackPressured && result != aeron.AdminAction {
			logger.Errorf("echoService: offer failed - sessionId=%d result=%d", session.Id(), result)
			return
		}
		s.cluster.IdleStrategy().Idle(0)
	}
}

func (s *echoService) OnTimerEvent(correlationId, timestamp int64) {
	select {
	case s.timerEvents <- correlationId:
	default:
		logger.Errorf("echoService: timerEvents channel full, dropping correlationId=%d", correlationId)
	}
}

func (s *echoService) OnTakeSnapshot(publication *aeron.Publication) {
	s.snapshots.Add(1)
	buf := atomicbuf.MakeBuffer(make([]byte, 4))
	buf.PutInt32(0, s.messageCount.Load())
	for {
		result := publication.Offer(buf, 0, buf.Capacity(), nil)
		if result >= 0 {
			return
		}
		if result != aeron.BackPressured && result != aeron.AdminAction {
			logger.Errorf("echoService: snapshot offer failed - result=%d", result)
			return
		}
		s.cluster.IdleStrategy().Idle(0)
	}
}

func (s *echoService) OnRoleChange(role cluster.Role) {
	s.roleChanges <- role
}

func (s *echoService) OnTerminate(cluster cluster.Cluster) {}

func (s *echoService) OnNewLeadershipTermEvent(
	leadershipTermId int64,
	logPosition int64,
	timestamp int64,
	termBaseLogPosition int64,
	leaderMemberId int32,
	logSessionId int32,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) {
}
