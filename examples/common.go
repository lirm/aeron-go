package examples

import "flag"

var ExamplesConfig = struct {
	AeronPrefix     string
	ProfilerEnabled bool
	DriverTo        int64
	StreamId        int32
	Channel         string
	Messages        int
}{
	*flag.String("p", "/tmp", "root directory for aeron driver file"),
	*flag.Bool("prof", false, "enable CPU profiling"),
	*flag.Int64("to", 10000, "driver liveliness timeout in ms"),
	int32(*flag.Int("sid", 10, "default streamId to use")),
	*flag.String("chan", "aeron:udp?endpoint=localhost:40123", "default channel to subscribe to"),
	*flag.Int("m", 1000000, "number of messages to send"),
}

var PingPongConfig = struct {
	PongStreamId int32
	PingStreamId int32
	PongChannel  string
	PingChannel  string
}{
	int32(*flag.Int("S", 11, "streamId to use for pong")),
	int32(*flag.Int("s", 10, "streamId to use for ping")),
	*flag.String("C", "aeron:ipc", "pong channel"),
	*flag.String("c", "aeron:ipc", "ping channel"),
}
