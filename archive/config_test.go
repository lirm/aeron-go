// Copyright (C) 2021 Talos, Inc.
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

package archive

import (
	"flag"
	"github.com/lirm/aeron-go/aeron"
	"os/user"
)

func getuser() string {
	user, _ := user.Current()
	return user.Username
}

var TestConfig = struct {
	RequestStream   *int
	RequestChannel  *string
	ResponseStream  *int
	ResponseChannel *string
	AeronPrefix     *string
	ProfilerEnabled *bool
	Timeout         *int64
	Messages        *int
	Payload         *int
	Verbose         *bool
}{
	flag.Int("requeststream", 10, "default request control stream to use"),
	flag.String("requestchannel", "aeron:udp?endpoint=localhost:8010", "default request control channel to publish to"),
	flag.Int("responsestream", 21, "default response control stream to use"),
	flag.String("responsechannel", "aeron:udp?endpoint=localhost:8020", "default response control channel to publish to"),
	flag.String("prefix", "dev/shm//aeron-"+aeron.UserName, "root directory for aeron driver file"),
	flag.Bool("profile", false, "enable CPU profiling"),
	flag.Int64("timeout", 10000, "driver liveliness timeout in ms"),
	flag.Int("messages", 100, "number of messages to send/receive"),
	flag.Int("payload", 256, "messages size"),
	flag.Bool("verbose", false, "enable verbose logging"),
}
