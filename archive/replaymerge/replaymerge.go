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

package replaymerge

import (
	"fmt"
	"strings"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/logbuffer/term"
	"github.com/lirm/aeron-go/archive"
)

const LiveAddMaxWindow = int32(32 * 1024 * 1024)
const ReplayRemoveThreshold = int64(0)
const MergeProgressTimeoutDefaultMs = int64(5 * time.Millisecond)

type State int

const (
	StateResolveReplayPort State = iota
	StateGetRecordingPosition
	StateReplay
	StateCatchup
	StateAttemptLiveJoin
	StateMerged
	StateFailed
	StateClosed
)

func (s State) String() string {
	return [...]string{"ResolveReplayPort", "GetRecordingPosition", "Replay", "Catchup", "AttemptLiveJoin", "Merged", "Failed", "Closed"}[s]
}

// ReplayMerge replays a recorded stream from a starting position and merge with live stream for a full history of a stream.
//
// Once constructed either of Poll or DoWork, interleaved with consumption
// of the Image, should be called in a duty cycle loop until IsMerged is true.
// After which the ReplayMerge can be closed and continued usage can be made of the Image or its
// parent Subscription. If an exception occurs or progress stops, the merge will fail and
// HasFailed will be true.
//
// If the endpoint on the replay destination uses a port of 0, then the OS will assign a port from the ephemeral
// range and this will be added to the replay channel for instructing the archive.
//
// NOTE: Merging is only supported with UDP streams.
type ReplayMerge struct {
	recordingId            int64
	startPosition          int64
	mergeProgressTimeoutMs int64
	replaySessionId        int64
	activeCorrelationId    int64
	nextTargetPosition     int64
	positionOfLastProgress int64
	timeOfLastProgressMs   int64
	isLiveAdded            bool
	isReplayActive         bool
	state                  State
	image                  aeron.Image

	archive           *archive.Archive
	subscription      *aeron.Subscription
	replayDestination string
	liveDestination   string
	replayEndpoint    string
	replayChannelUri  aeron.ChannelUri
}

// NewReplayMerge creates a ReplayMerge to manage the merging of a replayed stream and switching over to live stream as
// appropriate.
//
// Parameters:
//
// subscription           to use for the replay and live stream. Must be a multi-destination subscription.
// archive                to use for the replay.
// replayChannel          to as a template for what the archive will use.
// replayDestination      to send the replay to and the destination added by the Subscription.
// liveDestination        for the live stream and the destination added by the Subscription.
// recordingId            for the replay.
// startPosition          for the replay.
// epochClock             to use for progress checks.
// mergeProgressTimeoutMs to use for progress checks.
func NewReplayMerge(
	subscription *aeron.Subscription,
	archive *archive.Archive,
	replayChannel string,
	replayDestination string,
	liveDestination string,
	recordingId int64,
	startPosition int64,
	mergeProgressTimeoutMs int64) (rm *ReplayMerge, err error) {
	if strings.HasPrefix(subscription.Channel(), aeron.IpcChannel) ||
		strings.HasPrefix(replayChannel, aeron.IpcChannel) ||
		strings.HasPrefix(replayDestination, aeron.IpcChannel) ||
		strings.HasPrefix(liveDestination, aeron.IpcChannel) {
		err = fmt.Errorf("IPC merging is not supported")
		return
	}

	if !strings.Contains(subscription.Channel(), "control-mode=manual") {
		err = fmt.Errorf("Subscription URI must have 'control-mode=manual' uri=%s", subscription.Channel())
		return
	}

	rm = &ReplayMerge{
		archive:                archive,
		subscription:           subscription,
		replayDestination:      replayDestination,
		liveDestination:        liveDestination,
		recordingId:            recordingId,
		startPosition:          startPosition,
		mergeProgressTimeoutMs: mergeProgressTimeoutMs,
		replaySessionId:        aeron.NullValue,
		activeCorrelationId:    aeron.NullValue,
		nextTargetPosition:     aeron.NullValue,
		positionOfLastProgress: aeron.NullValue,
	}

	rm.replayChannelUri, err = aeron.ParseChannelUri(replayChannel)
	if err != nil {
		err = fmt.Errorf("Invalid replay channel '%s'", replayChannel)
		return
	}

	rm.replayChannelUri.Set(aeron.LingerParamName, "0")
	rm.replayChannelUri.Set(aeron.EosParamName, "false")

	var replayDestinationUri aeron.ChannelUri
	replayDestinationUri, err = aeron.ParseChannelUri(replayDestination)
	if err != nil {
		err = fmt.Errorf("Invalid replay destination '%s'", replayDestination)
		return
	}
	rm.replayEndpoint = replayDestinationUri.Get(aeron.EndpointParamName)
	if strings.HasSuffix(rm.replayEndpoint, ":0") {
		rm.state = StateResolveReplayPort
	} else {
		rm.replayChannelUri.Set(aeron.EndpointParamName, rm.replayEndpoint)
		rm.state = StateGetRecordingPosition
	}

	subscription.AddDestination(replayDestination)
	rm.timeOfLastProgressMs = time.Now().UnixMilli()
	return
}

// Close closes and stops any active replay. Will remove the replay destination from the subscription.
// This operation Will NOT remove the live destination if it has been added, so it can be used for live consumption.
func (rm *ReplayMerge) Close() {
	state := rm.state
	if StateClosed != state {
		if !rm.archive.Aeron().IsClosed() {
			if StateMerged != state {
				rm.subscription.RemoveDestination(rm.replayDestination)
			}

			if rm.isReplayActive && rm.archive.Proxy.Publication.IsConnected() {
				rm.stopReplay()
			}
		}

		rm.setState(StateClosed)
	}
}

// Subscription returns the Subscription used to consume the replayed and merged stream.
func (rm *ReplayMerge) Subscription() *aeron.Subscription {
	return rm.subscription
}

// DoWork performs the work of replaying and merging. Should only be used if polling the underlying Image directly.
// Returns indication of work done processing the merge.
func (rm *ReplayMerge) DoWork() (workCount int, err error) {
	nowMs := time.Now().UnixMilli()

	switch rm.state {
	case StateResolveReplayPort:
		workCount, err = rm.resolveReplayPort(nowMs)
		if err != nil {
			rm.setState(StateFailed)
			return
		}
		if err = rm.checkProgress(nowMs); err != nil {
			rm.setState(StateFailed)
			return
		}
	case StateGetRecordingPosition:
		workCount, err = rm.getRecordingPosition(nowMs)
		if err != nil {
			rm.setState(StateFailed)
			return
		}
		if err = rm.checkProgress(nowMs); err != nil {
			rm.setState(StateFailed)
			return
		}
	case StateReplay:
		workCount, err = rm.replay(nowMs)
		if err != nil {
			rm.setState(StateFailed)
			return
		}
		if err = rm.checkProgress(nowMs); err != nil {
			rm.setState(StateFailed)
			return
		}
	case StateCatchup:
		workCount, err = rm.catchup(nowMs)
		if err != nil {
			rm.setState(StateFailed)
			return
		}
		if err = rm.checkProgress(nowMs); err != nil {
			rm.setState(StateFailed)
			return
		}
	case StateAttemptLiveJoin:
		workCount, err = rm.attemptLiveJoin(nowMs)
		if err != nil {
			rm.setState(StateFailed)
			return
		}
		if err = rm.checkProgress(nowMs); err != nil {
			rm.setState(StateFailed)
			return
		}
	}
	return
}

// Poll polls the Image used for replay and merging and live stream. The doWork method
// will be called before the poll so that processing of the merge can be done.
//
// Returns number of fragments processed.
func (rm *ReplayMerge) Poll(fragmentHandler term.FragmentHandler, fragmentLimit int) (workCount int, err error) {
	workCount, err = rm.DoWork()
	if err != nil {
		return
	}
	if rm.image == nil {
		return
	}
	return rm.image.Poll(fragmentHandler, fragmentLimit), nil
}

// IsMerged returns if the live stream merged and the replay stopped?
func (rm *ReplayMerge) IsMerged() bool {
	return rm.state == StateMerged
}

// HasFailed returns if the replay merge failed due to an error?
func (rm *ReplayMerge) HasFailed() bool {
	return rm.state == StateFailed
}

// Image returns the image which is a merge of the replay and live stream.
func (rm *ReplayMerge) Image() aeron.Image {
	return rm.image
}

// IsLiveAdded returns if the live destination added to the subscription.
func (rm *ReplayMerge) IsLiveAdded() bool {
	return rm.isLiveAdded
}

func (rm *ReplayMerge) resolveReplayPort(nowMs int64) (workCount int, err error) {
	resolvedEndpoint := rm.subscription.ResolvedEndpoint()
	if resolvedEndpoint != "" {
		i := strings.LastIndex(resolvedEndpoint, ":")
		rm.replayChannelUri.Set(aeron.EndpointParamName,
			rm.replayEndpoint[0:len(rm.replayEndpoint)-2]+resolvedEndpoint[i:])

		rm.timeOfLastProgressMs = nowMs
		rm.setState(StateGetRecordingPosition)
		workCount += 1
	}

	return
}

func (rm *ReplayMerge) getRecordingPosition(nowMs int64) (workCount int, err error) {
	if aeron.NullValue == rm.activeCorrelationId {
		correlationId := rm.archive.Aeron().NextCorrelationID()

		if rm.archive.Proxy.RecordingPositionRequest(correlationId, rm.recordingId) == nil {
			rm.activeCorrelationId = correlationId
			rm.timeOfLastProgressMs = nowMs
			workCount += 1
		}
		return
	}

	var success bool
	success, err = rm.pollForResponse()
	if err != nil {
		return
	}
	if success {
		rm.nextTargetPosition = rm.polledRelevantId()
		rm.activeCorrelationId = aeron.NullValue

		if archive.RecordingPositionNull == rm.nextTargetPosition {
			correlationId := rm.archive.Aeron().NextCorrelationID()

			if rm.archive.Proxy.StopPositionRequest(correlationId, rm.recordingId) == nil {
				rm.activeCorrelationId = correlationId
				rm.timeOfLastProgressMs = nowMs
				workCount += 1
			}
		} else {
			rm.timeOfLastProgressMs = nowMs
			rm.setState(StateReplay)
		}
	}

	workCount += 1

	return
}

func (rm *ReplayMerge) replay(nowMs int64) (workCount int, err error) {
	if aeron.NullValue == rm.activeCorrelationId {
		correlationId := rm.archive.Aeron().NextCorrelationID()
		if rm.archive.Proxy.ReplayRequest(
			correlationId,
			rm.recordingId,
			rm.startPosition,
			archive.RecordingLengthMax,
			rm.replayChannelUri.String(),
			rm.subscription.StreamID()) == nil {
			rm.activeCorrelationId = correlationId
			rm.timeOfLastProgressMs = nowMs
			workCount += 1
		}
		return
	}

	var success bool
	success, err = rm.pollForResponse()
	if err != nil {
		return
	}
	if success {
		rm.isReplayActive = true
		rm.replaySessionId = rm.polledRelevantId()
		rm.timeOfLastProgressMs = nowMs
		rm.setState(StateCatchup)
		workCount += 1
	}
	return
}

func (rm *ReplayMerge) catchup(nowMs int64) (workCount int, err error) {

	if rm.image == nil && rm.subscription.IsConnected() {
		rm.timeOfLastProgressMs = nowMs
		rm.image = rm.subscription.ImageBySessionID(int32(rm.replaySessionId))
		rm.positionOfLastProgress = aeron.NullValue
		if rm.image != nil {
			rm.positionOfLastProgress = rm.image.Position()
		}
	}

	if rm.image != nil {
		position := rm.image.Position()
		if position >= rm.nextTargetPosition {
			rm.timeOfLastProgressMs = nowMs
			rm.positionOfLastProgress = position
			rm.setState(StateAttemptLiveJoin)
			workCount += 1
		} else if position > rm.positionOfLastProgress {
			rm.timeOfLastProgressMs = nowMs
			rm.positionOfLastProgress = position
		} else if rm.image.IsClosed() {
			err = fmt.Errorf("ReplayMerge Image closed unexpectedly.")
			return
		}
	}
	return
}

func (rm *ReplayMerge) attemptLiveJoin(nowMs int64) (workCount int, err error) {

	if aeron.NullValue == rm.activeCorrelationId {
		correlationId := rm.archive.Aeron().NextCorrelationID()
		if rm.archive.Proxy.RecordingPositionRequest(correlationId, rm.recordingId) == nil {
			rm.activeCorrelationId = correlationId
			workCount += 1
		}
		return
	}

	var success bool
	success, err = rm.pollForResponse()
	if err != nil {
		return
	}
	if success {
		rm.nextTargetPosition = rm.polledRelevantId()
		rm.activeCorrelationId = aeron.NullValue

		if archive.RecordingPositionNull == rm.nextTargetPosition {
			correlationId := rm.archive.Aeron().NextCorrelationID()
			if rm.archive.Proxy.RecordingPositionRequest(correlationId, rm.recordingId) == nil {
				rm.activeCorrelationId = correlationId
			}
		} else {
			nextState := StateCatchup

			if rm.image != nil {
				position := rm.image.Position()

				if rm.shouldAddLiveDestination(position) {
					rm.subscription.AddDestination(rm.liveDestination)
					rm.timeOfLastProgressMs = nowMs
					rm.positionOfLastProgress = position
					rm.isLiveAdded = true
				} else if rm.shouldStopAndRemoveReplay(position) {
					rm.subscription.RemoveDestination(rm.replayDestination)
					rm.stopReplay()
					rm.timeOfLastProgressMs = nowMs
					rm.positionOfLastProgress = position
					nextState = StateMerged
				}
			}

			rm.setState(nextState)
		}

		workCount += 1
	}

	return
}

func (rm *ReplayMerge) stopReplay() {
	correlationId := rm.archive.Aeron().NextCorrelationID()
	if rm.archive.Proxy.StopReplayRequest(correlationId, rm.replaySessionId) == nil {
		rm.isReplayActive = false
	}
}

func (rm *ReplayMerge) setState(newState State) {
	rm.state = newState
	rm.activeCorrelationId = aeron.NullValue
}

func (rm *ReplayMerge) shouldAddLiveDestination(position int64) bool {
	mn := rm.image.TermBufferLength() >> 2
	if mn > LiveAddMaxWindow {
		mn = LiveAddMaxWindow
	}
	return !rm.isLiveAdded &&
		(rm.nextTargetPosition-position) <= int64(mn)
}

func (rm *ReplayMerge) shouldStopAndRemoveReplay(position int64) bool {
	return rm.isLiveAdded &&
		(rm.nextTargetPosition-position) <= ReplayRemoveThreshold &&
		rm.image.ActiveTransportCount() >= 2
}

func (rm *ReplayMerge) checkProgress(nowMs int64) error {
	if nowMs > (rm.timeOfLastProgressMs + rm.mergeProgressTimeoutMs) {
		return fmt.Errorf("ReplayMerge no progress: state=%s", rm.state)
	}
	return nil
}

// Returns whether this succeeded, and what the error is.
func (rm *ReplayMerge) pollForResponse() (bool, error) {
	correlationId := rm.activeCorrelationId
	poller := rm.archive.Control

	if poller.Poll() > 0 && poller.Results.IsPollComplete {
		if poller.Results.ControlResponse.ControlSessionId == rm.archive.SessionID {
			if poller.Results.ErrorResponse != nil {
				err := fmt.Errorf(
					"archive response for correlationId=%d, error=%s",
					correlationId,
					poller.Results.ErrorResponse,
				)
				return false, err
			}
		}
		return poller.Results.CorrelationId == correlationId, nil
	}
	// TODO: (false, nil) is what was here before, but I suspect that (true, nil) is more accurate?
	// Check this when revamping archive code.
	return false, nil
}

func (rm *ReplayMerge) polledRelevantId() int64 {
	poller := rm.archive.Control
	return poller.Results.ControlResponse.RelevantId
}
