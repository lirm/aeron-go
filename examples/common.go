/*
Copyright 2016 Stanislav Liberman

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

package examples

import "flag"

var ExamplesConfig = struct {
	AeronPrefix     *string
	ProfilerEnabled *bool
	DriverTo        *int64
	StreamId        *int
	Channel         *string
	Messages        *int
}{
	flag.String("p", "/tmp", "root directory for aeron driver file"),
	flag.Bool("prof", false, "enable CPU profiling"),
	flag.Int64("to", 10000, "driver liveliness timeout in ms"),
	flag.Int("sid", 10, "default streamId to use"),
	flag.String("chan", "aeron:udp?endpoint=localhost:40123", "default channel to subscribe to"),
	flag.Int("m", 1000000, "number of messages to send"),
}

var PingPongConfig = struct {
	PongStreamId *int
	PingStreamId *int
	PongChannel  *string
	PingChannel  *string
}{
	flag.Int("S", 11, "streamId to use for pong"),
	flag.Int("s", 10, "streamId to use for ping"),
	flag.String("C", "aeron:ipc", "pong channel"),
	flag.String("c", "aeron:ipc", "ping channel"),
}
