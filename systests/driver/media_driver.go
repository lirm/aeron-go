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

const jarName = "driver/aeron-all-1.39.0.jar"

const mediaDriverClassName = "io.aeron.driver.MediaDriver"

var logger = logging.MustGetLogger("systests")

type MediaDriver struct {
	cmd *exec.Cmd
	cxn *aeron.Aeron
}

func StartMediaDriver() (*MediaDriver, error) {
	// If the previous run crashed and so Media Driver didn't clean up the temp dir, `aeron.Connect()` will try to use
	// it before Media Driver has a chance to execute the "delete on start" code, resulting in a Media Driver crash.
	// This can be avoided by preemptively deleting the directory.

	timeout := time.Now().Add(5 * time.Second)
	sleepDuration := 50 * time.Millisecond
	for removeTempDir() != nil {
		if time.Now().After(timeout) {
			return nil, errors.New("Couldn't delete temp dir.  Zombie process is probably running")
		}
		time.Sleep(sleepDuration)
	}
	cmd := setupCmd()
	setupPdeathsig(cmd)
	if err := cmd.Start(); err != nil {
		logger.Error("couldn't start Media Driver: ", err)
		return nil, err
	}
	cxn, err := waitForStartup()
	if err != nil {
		logger.Error("Media Driver timed out during startup: ", err)
		killMediaDriver(cmd)
		return nil, err
	}
	return &MediaDriver{cmd, cxn}, nil
}

func (mediaDriver MediaDriver) StopMediaDriver() {
	if err := mediaDriver.cxn.Close(); err != nil {
		logger.Error(err)
	}
	killMediaDriver(mediaDriver.cmd)
	if err := removeTempDir(); err != nil {
		logger.Errorf("Failed to clean up Aeron directories: %s", err)
	}
}

func removeTempDir() error {
	return os.RemoveAll(aeronTempDir())
}

func killMediaDriver(cmd *exec.Cmd) {
	if err := cmd.Process.Kill(); err != nil {
		logger.Error("Couldn't kill Media Driver:", err)
	}
}

func aeronTempDir() string {
	return fmt.Sprintf("%s/aeron-%s",
		aeron.DefaultAeronDir,
		aeron.UserName)
}

func setupCmd() *exec.Cmd {
	cmd := exec.Command(
		"java",
		fmt.Sprintf("-D%s=%s", aeronDirPropName, aeronTempDir()),
		fmt.Sprintf("-D%s=true", aeronDirDeleteStartPropName),
		fmt.Sprintf("-D%s=true", aeronDirDeleteShutdownPropName),
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

func waitForStartup() (*aeron.Aeron, error) {
	timeout := time.Now().Add(5 * time.Second)
	sleepDuration := 50 * time.Millisecond
	ctx := aeron.NewContext().AeronDir(aeronTempDir()).MediaDriverTimeout(10 * time.Second)
	channel := "aeron:ipc"
	streamID := int32(1)
	for time.Now().Before(timeout) {
		cxn, err := aeron.Connect(ctx)
		if err != nil {
			time.Sleep(sleepDuration)
			continue
		}
		pubc := cxn.AddPublication(channel, streamID)
		timeout := time.After(5 * time.Second)
		select {
		case <-timeout:
			return nil, errors.New("timed out waiting for Media Driver publication")
		case pub := <-pubc:
			_ = pub.Close()
			err := cxn.Close()
			if err != nil {
				logger.Error("Couldn't close aeron: ", err)
			}
			return cxn, nil
		}
	}
	return nil, errors.New("timed out waiting for Media Driver connection")
}
