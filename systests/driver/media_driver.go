/*
Copyright 2022 Steven Stern

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/logging"
	"os"
	"os/exec"
	"reflect"
	"syscall"
	"time"
)

// Must match Java's `AERON_DIR_PROP_NAME`
const aeronDirPropName = "aeron.dir"

// Must match Java's `DIR_DELETE_ON_START_PROP_NAME`
const aeronDirDeleteStartPropName = "aeron.dir.delete.on.start"

// Must match Java's `DIR_DELETE_ON_SHUTDOWN_PROP_NAME`
const aeronDirDeleteShutdownPropName = "aeron.dir.delete.on.shutdown"

// Must match Java's `CLIENT_LIVENESS_TIMEOUT_PROP_NAME`, measured in nanos
const aeronClientLivenessTimeoutNs = "aeron.client.liveness.timeout"

// Must match Java's `PUBLICATION_UNBLOCK_TIMEOUT_PROP_NAME`, measured in nanos
const aeronPublicationUnblockTimeoutNs = "aeron.publication.unblock.timeout"

const jarName = "driver/aeron-all-1.39.0.jar"

const mediaDriverClassName = "io.aeron.driver.MediaDriver"

var logger = logging.MustGetLogger("systests")

type MediaDriver struct {
	TempDir string
	cmd     *exec.Cmd
	cxn     *aeron.Aeron
}

func StartMediaDriver() (*MediaDriver, error) {
	tempDir := aeronUniqueTempDir()
	cmd := setupCmd(tempDir)
	setupPdeathsig(cmd)
	if err := cmd.Start(); err != nil {
		logger.Error("couldn't start Media Driver: ", err)
		return nil, err
	}
	cxn, err := waitForStartup(tempDir)
	if err != nil {
		logger.Error("Media Driver timed out during startup: ", err)
		killMediaDriver(cmd)
		return nil, err
	}
	return &MediaDriver{tempDir, cmd, cxn}, nil
}

func (mediaDriver MediaDriver) StopMediaDriver() {
	if err := mediaDriver.cxn.Close(); err != nil {
		logger.Error(err)
	}
	killMediaDriver(mediaDriver.cmd)
	if err := removeTempDir(mediaDriver.TempDir); err != nil {
		logger.Errorf("Failed to clean up cxn directories: %s", err)
	}
}

func removeTempDir(tempDir string) error {
	return os.RemoveAll(tempDir)
}

func killMediaDriver(cmd *exec.Cmd) {
	if err := cmd.Process.Kill(); err != nil {
		logger.Error("Couldn't kill Media Driver:", err)
	}
}

func aeronUniqueTempDir() string {
	id := uuid.New().String()
	return fmt.Sprintf("%s/aeron-%s/%s",
		aeron.DefaultAeronDir,
		aeron.UserName,
		id)
}

func setupCmd(tempDir string) *exec.Cmd {
	cmd := exec.Command(
		"java",
		fmt.Sprintf("-D%s=%s", aeronDirPropName, tempDir),
		fmt.Sprintf("-D%s=true", aeronDirDeleteStartPropName),
		fmt.Sprintf("-D%s=true", aeronDirDeleteShutdownPropName),
		fmt.Sprintf("-D%s=%d", aeronClientLivenessTimeoutNs, time.Minute.Nanoseconds()),
		fmt.Sprintf("-D%s=%d", aeronPublicationUnblockTimeoutNs, 15*time.Minute.Nanoseconds()),
		"-XX:+UnlockDiagnosticVMOptions",
		"-XX:GuaranteedSafepointInterval=300000",
		"-XX:BiasedLockingStartupDelay=0",
		"-cp",
		jarName,
		mediaDriverClassName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// Setting Pdeathsig kills child processes on shutdown, but this only works on linux.
// On other platforms, the media driver could be left stranded if the test panics.
// https://stackoverflow.com/questions/34730941/ensure-executables-called-in-go-process-get-killed-when-process-is-killed
func setupPdeathsig(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	pdeathsig := reflect.ValueOf(cmd.SysProcAttr).Elem().FieldByName("Pdeathsig")
	if pdeathsig.IsValid() {
		pdeathsig.Set(reflect.ValueOf(syscall.SIGTERM))
	}
}

func waitForStartup(tempDir string) (*aeron.Aeron, error) {
	timeout := time.Now().Add(5 * time.Second)
	sleepDuration := 50 * time.Millisecond
	ctx := aeron.NewContext().AeronDir(tempDir).MediaDriverTimeout(10 * time.Second)
	channel := "aeron:ipc"
	streamID := int32(1)
	for time.Now().Before(timeout) {
		cxn, err := aeron.Connect(ctx)
		if err != nil {
			time.Sleep(sleepDuration)
			continue
		}
		pub, err := cxn.AddPublication(channel, streamID)
		if err != nil {
			return nil, err
		}
		if err := pub.Close(); err != nil {
			return nil, err
		}
		return cxn, nil
	}
	return nil, errors.New("timed out waiting for Media Driver connection")
}
