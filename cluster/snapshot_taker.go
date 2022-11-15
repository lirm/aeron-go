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
	"errors"
	"fmt"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/cluster/codecs"
)

const snapshotTypeId = 2

type snapshotTaker struct {
	marshaller  *codecs.SbeGoMarshaller // currently shared as we're not reentrant (but could be here)
	options     *Options
	publication *aeron.Publication
}

func newSnapshotTaker(
	options *Options,
	publication *aeron.Publication,
) *snapshotTaker {
	return &snapshotTaker{
		marshaller:  codecs.NewSbeGoMarshaller(),
		options:     options,
		publication: publication,
	}
}

func (st *snapshotTaker) markBegin(
	logPosition int64,
	leadershipTermId int64,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) error {
	return st.markSnapshot(logPosition, leadershipTermId, codecs.SnapshotMark.BEGIN, timeUnit, appVersion)
}

func (st *snapshotTaker) markEnd(
	logPosition int64,
	leadershipTermId int64,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) error {
	return st.markSnapshot(logPosition, leadershipTermId, codecs.SnapshotMark.END, timeUnit, appVersion)
}

func (st *snapshotTaker) markSnapshot(
	logPosition int64,
	leadershipTermId int64,
	mark codecs.SnapshotMarkEnum,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) error {
	bytes, err := codecs.SnapshotMarkerPacket(
		st.marshaller,
		st.options.RangeChecking,
		snapshotTypeId,
		logPosition,
		leadershipTermId,
		0,
		mark,
		timeUnit,
		appVersion,
	)
	if err != nil {
		return err
	}
	if _, err := st.offer(bytes); err != nil {
		return fmt.Errorf("snapshotTaker.offer failed: %w", err)
	}
	return nil
}

func (st *snapshotTaker) snapshotSession(session ClientSession) error {
	bytes, err := codecs.ClientSessionPacket(st.marshaller, st.options.RangeChecking,
		session.Id(), session.ResponseStreamId(), []byte(session.ResponseChannel()), session.EncodedPrincipal())
	if err != nil {
		return err
	}
	if _, err := st.offer(bytes); err != nil {
		return fmt.Errorf("snapshotTaker.offer failed: %w", err)
	}
	return nil
}

// Offer to our request publication
func (st *snapshotTaker) offer(bytes []byte) (int64, error) {
	buffer := atomic.MakeBuffer(bytes)
	length := int32(len(bytes))
	start := time.Now()
	var ret int64
	var err error
	for time.Since(start) < st.options.Timeout {
		ret, err = st.publication.Offer(buffer, 0, length, nil)
		if errors.Is(err, aeron.NotConnectedErr) ||
			errors.Is(err, aeron.BackPressuredErr) ||
			errors.Is(err, aeron.AdminActionErr) {
			st.options.IdleStrategy.Idle(0)
		} else {
			return ret, err
		}
	}
	// Give up, returning the last failure
	return ret, err
}
