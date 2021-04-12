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
	SampleStream    *int
	SampleChannel   *string
	ReplayStream    *int
	ReplayChannel   *string
	AeronPrefix     *string
	ProfilerEnabled *bool
	Verbose         *bool
}{
	flag.Int("samplestream", 1001, "default base response control stream to use"),
	flag.String("samplechannel", "aeron:udp?endpoint=localhost:20121", "default response control channel to publish to"),

	flag.Int("replaystream", 1002, "default base replay stream to use"),
	flag.String("replaychannel", "aeron:udp?endpoint=localhost:20121", "default replay to receive from"),

	flag.String("prefix", aeron.DefaultAeronDir+"/aeron-"+aeron.UserName, "root directory for aeron driver file"),
	flag.Bool("profile", false, "enable CPU profiling"),
	flag.Bool("verbose", false, "enable verbose (DEBUG) logging"),
}
