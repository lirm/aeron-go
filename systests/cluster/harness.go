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

// Package cluster contains integration tests that run the Go clustered
// service container and cluster client against a real Java
// ClusteredMediaDriver (media driver + archive + consensus module).
package clustertests

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/logging"
)

// The jar is not committed; fetch it with:
//
//	curl -fLO --output-dir systests/cluster \
//	  https://repo1.maven.org/maven2/io/aeron/aeron-all/1.52.0/aeron-all-1.52.0.jar
//
// or point AERON_ALL_JAR at an existing copy.
const defaultJarName = "aeron-all-1.52.0.jar"

const clusteredMediaDriverClassName = "io.aeron.cluster.ClusteredMediaDriver"

// Each driver instance gets its own port block, derived from the pid and an
// instance counter: distinct test invocations and distinct drivers within one
// test process never share ports, so a leaked client from a stopped cluster
// cannot send stale traffic into the next one. Ports stay below 32768 so they
// can't collide with the kernel's ephemeral port range (32768-60999 by
// default). Restarts of the same driver deliberately reuse its block.
var driverInstances atomic.Int32

func nextBasePort() int {
	instance := int(driverInstances.Add(1) - 1)
	return 20000 + (os.Getpid()%100)*120 + (instance%6)*20
}

var logger = logging.MustGetLogger("clustertests")

// ClusteredMediaDriver wraps a Java ClusteredMediaDriver child process and
// the temp directories it runs in.
type ClusteredMediaDriver struct {
	AeronDir        string
	ClusterDir      string
	ArchiveDir      string
	IngressEndpoint string
	LogPath         string
	basePort        int
	cmd             *exec.Cmd
	exited          chan struct{}
	exitErr         error
}

func jarPath() string {
	if path := os.Getenv("AERON_ALL_JAR"); path != "" {
		return path
	}
	return defaultJarName
}

// JarAvailable reports whether the aeron-all jar needed to launch the driver
// exists; tests should skip when it does not.
func JarAvailable() (string, bool) {
	path := jarPath()
	_, err := os.Stat(path)
	return path, err == nil
}

func StartClusteredMediaDriver() (*ClusteredMediaDriver, error) {
	jar, ok := JarAvailable()
	if !ok {
		return nil, fmt.Errorf("aeron-all jar not found at %s", jar)
	}

	id := uuid.New().String()
	aeronDir := fmt.Sprintf("%s/aeron-%s/%s/driver", aeron.DefaultAeronDir, aeron.UserName, id)
	baseDir, err := os.MkdirTemp("", "aeron-go-cluster-systest")
	if err != nil {
		return nil, err
	}
	clusterDir := baseDir + "/cluster"
	archiveDir := baseDir + "/archive"
	// The Go service container writes its mark file into the cluster dir,
	// possibly before the consensus module has created it.
	if err := os.MkdirAll(clusterDir, 0o755); err != nil {
		return nil, err
	}
	return launchDriver(aeronDir, clusterDir, archiveDir, nextBasePort())
}

// Restart starts a new driver process over this driver's cluster and archive
// directories, so the cluster recovers from its recording log and snapshots.
// The driver must have been stopped with Shutdown first.
func (driver *ClusteredMediaDriver) Restart() (*ClusteredMediaDriver, error) {
	return launchDriver(driver.AeronDir, driver.ClusterDir, driver.ArchiveDir, driver.basePort)
}

func launchDriver(aeronDir, clusterDir, archiveDir string, basePort int) (*ClusteredMediaDriver, error) {
	jar, ok := JarAvailable()
	if !ok {
		return nil, fmt.Errorf("aeron-all jar not found at %s", jar)
	}

	clusterMembers := fmt.Sprintf("0,localhost:%d,localhost:%d,localhost:%d,localhost:0,localhost:%d",
		basePort, basePort+1, basePort+2, basePort+10)
	archiveControl := fmt.Sprintf("aeron:udp?endpoint=localhost:%d", basePort+10)

	cmd := exec.Command(
		"java",
		"--add-opens=java.base/sun.nio.ch=ALL-UNNAMED",
		"--add-exports=java.base/jdk.internal.misc=ALL-UNNAMED",
		"-XX:+UnlockDiagnosticVMOptions",
		"-XX:GuaranteedSafepointInterval=300000",
		fmt.Sprintf("-Daeron.dir=%s", aeronDir),
		"-Daeron.dir.delete.on.start=true",
		"-Daeron.dir.delete.on.shutdown=true",
		"-Daeron.threading.mode=SHARED",
		fmt.Sprintf("-Daeron.client.liveness.timeout=%d", time.Minute.Nanoseconds()),
		fmt.Sprintf("-Daeron.publication.unblock.timeout=%d", 15*time.Minute.Nanoseconds()),
		fmt.Sprintf("-Daeron.archive.dir=%s", archiveDir),
		fmt.Sprintf("-Daeron.archive.control.channel=%s", archiveControl),
		"-Daeron.archive.replication.channel=aeron:udp?endpoint=localhost:0",
		"-Daeron.archive.threading.mode=SHARED",
		fmt.Sprintf("-Daeron.cluster.dir=%s", clusterDir),
		fmt.Sprintf("-Daeron.cluster.members=%s", clusterMembers),
		"-Daeron.cluster.ingress.channel=aeron:udp?term-length=64k",
		"-Daeron.cluster.replication.channel=aeron:udp?endpoint=localhost:0",
		"-Daeron.cluster.service.count=1",
		"-cp",
		jar,
		clusteredMediaDriverClassName,
	)
	logPath := fmt.Sprintf("%s/driver-%s.log", filepath.Dir(clusterDir), uuid.New().String()[:8])
	logFile, err := os.Create(logPath)
	if err != nil {
		return nil, err
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	setupPdeathsig(cmd)

	driver := &ClusteredMediaDriver{
		AeronDir:        aeronDir,
		ClusterDir:      clusterDir,
		ArchiveDir:      archiveDir,
		IngressEndpoint: fmt.Sprintf("localhost:%d", basePort),
		LogPath:         logPath,
		basePort:        basePort,
		cmd:             cmd,
		exited:          make(chan struct{}),
	}
	logger.Infof("starting ClusteredMediaDriver (log: %s): %s", logPath, cmd)
	// Pdeathsig is delivered when the forking OS thread dies, not the
	// process, so the thread that starts the driver must stay alive for the
	// driver's whole lifetime: keep it locked until the child exits.
	started := make(chan error, 1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		if err := cmd.Start(); err != nil {
			logFile.Close()
			started <- err
			return
		}
		started <- nil
		driver.exitErr = cmd.Wait()
		logFile.Close()
		close(driver.exited)
		logger.Infof("ClusteredMediaDriver pid=%d exited: %v", cmd.Process.Pid, driver.exitErr)
	}()
	if err := <-started; err != nil {
		return nil, err
	}
	if err := driver.awaitMediaDriverReady(); err != nil {
		driver.Stop()
		return nil, err
	}
	return driver, nil
}

// Exited reports whether the driver process has terminated.
func (driver *ClusteredMediaDriver) Exited() bool {
	select {
	case <-driver.exited:
		return true
	default:
		return false
	}
}

// Shutdown stops the driver process gracefully (so the consensus module and
// archive close cleanly) and keeps the cluster and archive directories, ready
// for a Restart.
func (driver *ClusteredMediaDriver) Shutdown() {
	if err := driver.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		logger.Errorf("couldn't signal ClusteredMediaDriver: %v", err)
	}
	select {
	case <-driver.exited:
	case <-time.After(15 * time.Second):
		logger.Errorf("ClusteredMediaDriver did not shut down in time, killing")
		_ = driver.cmd.Process.Kill()
		<-driver.exited
	}
}

// Stop kills the driver process and removes all its directories.
func (driver *ClusteredMediaDriver) Stop() {
	if err := driver.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		logger.Errorf("couldn't kill ClusteredMediaDriver: %v", err)
	}
	<-driver.exited
	for _, dir := range []string{driver.AeronDir, driver.ClusterDir, driver.ArchiveDir} {
		if err := os.RemoveAll(dir); err != nil {
			logger.Errorf("failed to remove %s: %v", dir, err)
		}
	}
}

func (driver *ClusteredMediaDriver) awaitMediaDriverReady() error {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		ctx := aeron.NewContext().AeronDir(driver.AeronDir).MediaDriverTimeout(10 * time.Second)
		cxn, err := aeron.Connect(ctx)
		if err == nil {
			return cxn.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
	return errors.New("timed out waiting for ClusteredMediaDriver to start")
}

// ClusterTool runs an io.aeron.cluster.ClusterTool command (e.g. "snapshot",
// "describe", "errors") against this driver's cluster directory.
func (driver *ClusteredMediaDriver) ClusterTool(command string) (string, error) {
	jar, _ := JarAvailable()
	tool := exec.Command(
		"java",
		"--add-opens=java.base/sun.nio.ch=ALL-UNNAMED",
		"--add-exports=java.base/jdk.internal.misc=ALL-UNNAMED",
		"-cp", jar,
		"io.aeron.cluster.ClusterTool",
		driver.ClusterDir,
		command,
	)
	out, err := tool.CombinedOutput()
	return string(out), err
}

// Setting Pdeathsig kills the child process when the test process dies, but
// this only works on linux; elsewhere a panicking test can strand the driver.
func setupPdeathsig(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	pdeathsig := reflect.ValueOf(cmd.SysProcAttr).Elem().FieldByName("Pdeathsig")
	if pdeathsig.IsValid() {
		pdeathsig.Set(reflect.ValueOf(syscall.SIGTERM))
	}
}
