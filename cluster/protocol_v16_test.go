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

package cluster

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/lirm/aeron-go/cluster/codecs"
)

// The adapters dispatch on hardcoded template ids; keep them in sync with the
// generated codecs.
func TestTemplateIdsMatchGeneratedCodecs(t *testing.T) {
	cases := []struct {
		name     string
		constant uint16
		codec    uint16
	}{
		{"JoinLog", joinLogTemplateId, (&codecs.JoinLog{}).SbeTemplateId()},
		{"ClusterActionRequest", clusterActionReqTemplateId, (&codecs.ClusterActionRequest{}).SbeTemplateId()},
		{"RequestServiceAck", requestServiceAckTemplateId, (&codecs.RequestServiceAck{}).SbeTemplateId()},
		{"SchemaVersion", ClusterSchemaVersion, (&codecs.JoinLog{}).SbeSchemaVersion()},
		{"SchemaId", ClusterSchemaId, (&codecs.JoinLog{}).SbeSchemaId()},
	}
	for _, c := range cases {
		if c.constant != c.codec {
			t.Errorf("%s: constant=%d codec=%d", c.name, c.constant, c.codec)
		}
	}
}

func TestShouldSnapshot(t *testing.T) {
	agent := &ClusteredServiceAgent{standbySnapshotFlags: clusterActionFlagsDefault}
	if !agent.shouldSnapshot(clusterActionFlagsDefault) {
		t.Error("regular node must snapshot on default flags")
	}
	if agent.shouldSnapshot(clusterActionFlagsStandbySnapshot) {
		t.Error("regular node must not take standby snapshot")
	}

	agent.standbySnapshotFlags = clusterActionFlagsStandbySnapshot
	if !agent.shouldSnapshot(clusterActionFlagsDefault) {
		t.Error("standby node must snapshot on default flags")
	}
	if !agent.shouldSnapshot(clusterActionFlagsStandbySnapshot) {
		t.Error("standby node must take standby snapshot")
	}
}

// A ClusterActionRequest recorded by a pre-v11 cluster has no flags field on
// the wire; the decoder must report the null value so the log adapter can map
// it to the default flags.
func TestClusterActionRequestDecodeV8HasNullFlags(t *testing.T) {
	buf := &bytes.Buffer{}
	for _, v := range []int64{5, 1024, 999} { // leadershipTermId, logPosition, timestamp
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	if err := binary.Write(buf, binary.LittleEndian, int32(codecs.ClusterAction.SNAPSHOT)); err != nil {
		t.Fatal(err)
	}

	m := codecs.NewSbeGoMarshaller()
	req := codecs.ClusterActionRequest{}
	const v8BlockLength = 28
	if err := req.Decode(m, buf, 8, v8BlockLength, true); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if req.LogPosition != 1024 || req.Action != codecs.ClusterAction.SNAPSHOT {
		t.Errorf("bad decode: %+v", req)
	}
	if req.Flags != req.FlagsNullValue() {
		t.Errorf("expected null flags for v8 message, got %d", req.Flags)
	}
}

func TestClusterActionRequestRoundTripV16(t *testing.T) {
	m := codecs.NewSbeGoMarshaller()
	in := codecs.ClusterActionRequest{
		LeadershipTermId: 7,
		LogPosition:      2048,
		Timestamp:        12345,
		Action:           codecs.ClusterAction.SNAPSHOT,
		Flags:            clusterActionFlagsStandbySnapshot,
	}
	buf := &bytes.Buffer{}
	if err := in.Encode(m, buf, true); err != nil {
		t.Fatal(err)
	}
	out := codecs.ClusterActionRequest{}
	if err := out.Decode(m, buf, in.SbeSchemaVersion(), in.SbeBlockLength(), true); err != nil {
		t.Fatal(err)
	}
	if out.Flags != clusterActionFlagsStandbySnapshot {
		t.Errorf("expected flags=%d, got %d", clusterActionFlagsStandbySnapshot, out.Flags)
	}
}

func TestRequestServiceAckRoundTrip(t *testing.T) {
	m := codecs.NewSbeGoMarshaller()
	in := codecs.RequestServiceAck{LogPosition: 4096}
	buf := &bytes.Buffer{}
	if err := in.Encode(m, buf, true); err != nil {
		t.Fatal(err)
	}
	out := codecs.RequestServiceAck{}
	if err := out.Decode(m, buf, in.SbeSchemaVersion(), in.SbeBlockLength(), true); err != nil {
		t.Fatal(err)
	}
	if out.LogPosition != 4096 {
		t.Errorf("expected logPosition=4096, got %d", out.LogPosition)
	}
}

// Mirrors upstream SessionEventCodecCompatibilityTest: an egress SessionEvent
// from a pre-v13 consensus module has no leaderHeartbeatTimeoutNs, and the
// var-length detail must still decode correctly after the shorter block.
func TestSessionEventDecodeV8HasNullHeartbeatTimeout(t *testing.T) {
	buf := &bytes.Buffer{}
	for _, v := range []int64{-4623823, 3583456348756843, 10000000000} { // clusterSessionId, correlationId, leadershipTermId
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	// leaderMemberId, code, version
	for _, v := range []int32{2, int32(codecs.EventCode.REDIRECT), 823} {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	detail := "some very detailed message"
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(detail))); err != nil {
		t.Fatal(err)
	}
	buf.WriteString(detail)

	m := codecs.NewSbeGoMarshaller()
	event := codecs.SessionEvent{}
	const v8BlockLength = 36
	if err := event.Decode(m, buf, 8, v8BlockLength, true); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if event.ClusterSessionId != -4623823 || event.Code != codecs.EventCode.REDIRECT {
		t.Errorf("bad decode: %+v", event)
	}
	if string(event.Detail) != detail {
		t.Errorf("expected detail %q, got %q", detail, string(event.Detail))
	}
	if event.LeaderHeartbeatTimeoutNs != event.LeaderHeartbeatTimeoutNsNullValue() {
		t.Errorf("expected null heartbeat timeout, got %d", event.LeaderHeartbeatTimeoutNs)
	}
}

// A JoinLog from a pre-v16 consensus module carries no isStandby field; it
// must decode as not-standby with the log channel still intact.
func TestJoinLogDecodeV8DefaultsToNotStandby(t *testing.T) {
	buf := &bytes.Buffer{}
	for _, v := range []int64{100, 200} { // logPosition, maxLogPosition
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	// memberId, logSessionId, logStreamId, isStartup, role
	for _, v := range []int32{1, 2, 3, int32(codecs.BooleanType.TRUE), int32(Leader)} {
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	channel := "aeron:ipc"
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(channel))); err != nil {
		t.Fatal(err)
	}
	buf.WriteString(channel)

	m := codecs.NewSbeGoMarshaller()
	joinLog := codecs.JoinLog{}
	const v8BlockLength = 36
	if err := joinLog.Decode(m, buf, 8, v8BlockLength, true); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if joinLog.IsStandby == codecs.BooleanType.TRUE {
		t.Error("v8 JoinLog must not decode as standby")
	}
	if joinLog.IsStartup != codecs.BooleanType.TRUE || string(joinLog.LogChannel) != channel {
		t.Errorf("bad decode: %+v", joinLog)
	}
}
