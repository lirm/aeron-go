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
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lirm/aeron-go/aeron"
	atomicbuf "github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster"
	"github.com/lirm/aeron-go/cluster/client"
	"go.uber.org/zap/zapcore"
)

type egressCollector struct {
	messages chan string
}

func (c *egressCollector) OnConnect(cluster *client.AeronCluster) {}

func (c *egressCollector) OnDisconnect(cluster *client.AeronCluster, details string) {}

func (c *egressCollector) OnMessage(
	cluster *client.AeronCluster,
	timestamp int64,
	buffer *atomicbuf.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	c.messages <- string(buffer.GetBytesArray(offset, length))
}

func (c *egressCollector) OnNewLeader(cluster *client.AeronCluster, leadershipTermId int64, leaderMemberId int32) {
}

func (c *egressCollector) OnError(cluster *client.AeronCluster, details string) {
	logger.Errorf("egress error: %s", details)
}

// serviceRunner drives a clustered service agent's duty cycle on its own
// goroutine until shutdown is requested.
type serviceRunner struct {
	service *echoService
	agent   *cluster.ClusteredServiceAgent
	stop    atomic.Bool
	done    chan struct{}
}

func startServiceAgent(t *testing.T, driver *ClusteredMediaDriver) *serviceRunner {
	t.Helper()
	opts := cluster.NewOptions()
	opts.ClusterDir = driver.ClusterDir
	if os.Getenv("CLUSTER_TEST_DEBUG") != "" {
		opts.Loglevel = zapcore.DebugLevel
	}
	service := newEchoService()
	agent, err := cluster.NewClusteredServiceAgent(aeron.NewContext().AeronDir(driver.AeronDir), opts, service)
	if err != nil {
		t.Fatalf("failed to create service agent: %v", err)
	}

	runner := &serviceRunner{service: service, agent: agent, done: make(chan struct{})}
	started := make(chan error, 1)
	go func() {
		defer close(runner.done)
		// The client library panics when the media driver goes away at
		// shutdown; that must not take the test process down.
		defer func() {
			if r := recover(); r != nil && !runner.stop.Load() {
				t.Errorf("service agent panicked: %v", r)
			}
		}()
		if err := agent.OnStart(); err != nil {
			started <- err
			return
		}
		started <- nil
		for !runner.stop.Load() {
			agent.Idle(agent.DoWork())
		}
	}()

	select {
	case err := <-started:
		if err != nil {
			t.Fatalf("service agent failed to start: %v", err)
		}
	case <-time.After(30 * time.Second):
		dumpClusterState(t, driver)
		t.Fatal("timed out waiting for service agent to join the cluster")
	}
	return runner
}

func (runner *serviceRunner) shutdown() {
	runner.stop.Store(true)
	select {
	case <-runner.done:
	case <-time.After(10 * time.Second):
		logger.Errorf("service agent duty cycle did not stop in time")
	}
	runner.agent.Close()
}

func connectClient(t *testing.T, driver *ClusteredMediaDriver) (*client.AeronCluster, *egressCollector) {
	t.Helper()
	clientOpts := client.NewOptions()
	clientOpts.IngressChannel = "aeron:udp?alias=cluster-client-ingress|endpoint=" + driver.IngressEndpoint
	clientOpts.IngressEndpoints = "0=" + driver.IngressEndpoint
	collector := &egressCollector{messages: make(chan string, 100)}
	clusterClient, err := client.NewAeronCluster(aeron.NewContext().AeronDir(driver.AeronDir), clientOpts, collector)
	if err != nil {
		t.Fatalf("failed to create cluster client: %v", err)
	}
	deadline := time.Now().Add(60 * time.Second)
	for !clusterClient.IsConnected() {
		if time.Now().After(deadline) {
			dumpClusterState(t, driver)
			t.Fatal("timed out waiting for cluster client to connect")
		}
		clientOpts.IdleStrategy.Idle(clusterClient.Poll())
	}
	return clusterClient, collector
}

// dumpClusterState logs the consensus module's view of the cluster and the
// driver process output, for diagnosing timeouts.
func dumpClusterState(t *testing.T, driver *ClusteredMediaDriver) {
	t.Helper()
	t.Logf("driver process exited=%v", driver.Exited())
	if data, err := os.ReadFile(driver.LogPath); err == nil {
		out := string(data)
		if len(out) > 4000 {
			out = out[len(out)-4000:]
		}
		t.Logf("driver log tail:\n%s", out)
	}
	for _, command := range []string{"describe", "errors"} {
		out, err := driver.ClusterTool(command)
		t.Logf("ClusterTool %s (err=%v):\n%s", command, err, out)
	}
}

// sendAndAwaitEchoes offers every payload to the cluster and waits until each
// one has come back on egress.
func sendAndAwaitEchoes(
	t *testing.T,
	clusterClient *client.AeronCluster,
	collector *egressCollector,
	payloads []string,
) {
	t.Helper()
	sendBuf := atomicbuf.MakeBuffer(make([]byte, 256))
	for _, payload := range payloads {
		payloadBytes := []byte(payload)
		sendBuf.PutBytesArray(0, &payloadBytes, 0, int32(len(payloadBytes)))
		deadline := time.Now().Add(10 * time.Second)
		for clusterClient.Offer(sendBuf, 0, int32(len(payloadBytes))) < 0 {
			if time.Now().After(deadline) {
				t.Fatalf("timed out offering %q", payload)
			}
			clusterClient.Poll()
		}
	}

	received := make(map[string]bool)
	deadline := time.After(30 * time.Second)
	for len(received) < len(payloads) {
		clusterClient.Poll()
		select {
		case msg := <-collector.messages:
			received[msg] = true
		case <-deadline:
			t.Fatalf("timed out waiting for echoes - received %d of %d", len(received), len(payloads))
		default:
			time.Sleep(time.Millisecond)
		}
	}
	for _, payload := range payloads {
		if !received[payload] {
			t.Errorf("missing echo for %q", payload)
		}
	}
}

// triggerSnapshot asks the consensus module for a snapshot via ClusterTool
// and waits for the service to take it. The tool only returns success once
// the snapshot recording is acknowledged by the service.
func triggerSnapshot(t *testing.T, driver *ClusteredMediaDriver, service *echoService) {
	t.Helper()
	before := service.snapshots.Load()
	if out, err := driver.ClusterTool("snapshot"); err != nil {
		t.Fatalf("ClusterTool snapshot failed: %v - output: %s", err, out)
	}
	deadline := time.Now().Add(30 * time.Second)
	for service.snapshots.Load() == before {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for service to take a snapshot")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func awaitMessageCount(t *testing.T, service *echoService, expected int32) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for service.messageCount.Load() != expected {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for message count %d, have %d", expected, service.messageCount.Load())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func payloads(prefix string, count int) []string {
	result := make([]string, count)
	for i := range result {
		result[i] = fmt.Sprintf("%s-%d", prefix, i)
	}
	return result
}

func requireDriver(t *testing.T) *ClusteredMediaDriver {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping cluster system test in short mode")
	}
	if jar, ok := JarAvailable(); !ok {
		t.Skipf("aeron-all jar not found at %s - see harness.go for fetch instructions", jar)
	}
	driver, err := StartClusteredMediaDriver()
	if err != nil {
		t.Fatalf("failed to start ClusteredMediaDriver: %v", err)
	}
	return driver
}

// TestClusterEchoAndSnapshot runs the Go service container and cluster
// client against a Java 1.52 ClusteredMediaDriver: client messages must be
// replicated through the log, echoed by the service, and delivered on
// egress; a ClusterTool-triggered snapshot must invoke OnTakeSnapshot and
// be acked (with its recording id) so the tool completes.
func TestClusterEchoAndSnapshot(t *testing.T) {
	driver := requireDriver(t)
	defer driver.Stop()

	runner := startServiceAgent(t, driver)
	defer runner.shutdown()

	clusterClient, collector := connectClient(t, driver)
	defer clusterClient.Close()

	sendAndAwaitEchoes(t, clusterClient, collector, payloads("echo-test-message", 10))
	triggerSnapshot(t, driver, runner.service)

	if got := runner.service.messageCount.Load(); got != 10 {
		t.Errorf("service processed %d messages, expected 10", got)
	}
}

// TestClusterTimers verifies deterministic cluster timers, modeled on
// upstream's ClusterTimerTest: a timer scheduled by the service (from within
// session message processing) must expire through the replicated log and
// reach OnTimerEvent; a cancelled timer must not fire.
func TestClusterTimers(t *testing.T) {
	driver := requireDriver(t)
	defer driver.Stop()

	runner := startServiceAgent(t, driver)
	defer runner.shutdown()

	clusterClient, collector := connectClient(t, driver)
	defer clusterClient.Close()

	// Schedule a short timer and a long one, then cancel the long one.
	sendAndAwaitEchoes(t, clusterClient, collector, []string{
		"timer:100:500",
		"timer:101:60000",
		"cancel-timer:101",
	})

	select {
	case correlationId := <-runner.service.timerEvents:
		if correlationId != 100 {
			t.Errorf("expected timer 100 to fire, got %d", correlationId)
		}
	case <-time.After(15 * time.Second):
		dumpClusterState(t, driver)
		t.Fatal("timed out waiting for timer to fire")
	}

	// The cancelled timer must stay silent; watch briefly for stragglers.
	select {
	case correlationId := <-runner.service.timerEvents:
		t.Errorf("unexpected timer fired: %d", correlationId)
	case <-time.After(2 * time.Second):
	}
}

// TestClusterRestartFromSnapshot verifies the recovery path: after a
// snapshot and more messages, the cluster and the Go service container are
// restarted over the same cluster and archive directories. The service must
// restore its state from the snapshot (replayed through the archive), replay
// the post-snapshot log entries, and then serve new clients.
func TestClusterRestartFromSnapshot(t *testing.T) {
	driver := requireDriver(t)
	currentDriver := driver
	defer func() { currentDriver.Stop() }()

	runner := startServiceAgent(t, driver)
	clusterClient, collector := connectClient(t, driver)

	sendAndAwaitEchoes(t, clusterClient, collector, payloads("pre-snapshot", 10))
	triggerSnapshot(t, driver, runner.service)
	sendAndAwaitEchoes(t, clusterClient, collector, payloads("post-snapshot", 5))

	// Take the whole node down: client, service container, then the driver
	// (gracefully, keeping the cluster and archive directories).
	clusterClient.Close()
	runner.shutdown()
	driver.Shutdown()

	restarted, err := driver.Restart()
	if err != nil {
		t.Fatalf("failed to restart ClusteredMediaDriver: %v", err)
	}
	currentDriver = restarted

	// The new service instance must recover the snapshot state (10) and
	// replay the 5 post-snapshot log entries.
	recoveredRunner := startServiceAgent(t, restarted)
	defer recoveredRunner.shutdown()
	awaitMessageCount(t, recoveredRunner.service, 15)

	// And it must serve new clients.
	newClient, newCollector := connectClient(t, restarted)
	defer newClient.Close()
	sendAndAwaitEchoes(t, newClient, newCollector, payloads("after-restart", 5))
	awaitMessageCount(t, recoveredRunner.service, 20)
}
