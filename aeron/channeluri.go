package aeron

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ChannelUri is a parser for Aeron channel URIs. The format is:
// aeron-uri = "aeron:" media [ "?" param *( "|" param ) ]
// media     = *( "[^?:]" )
// param     = key "=" value
// key       = *( "[^=]" )
// value     = *( "[^|]" )
//
// Multiple params with the same key are allowed, the last value specified takes precedence.
type ChannelUri struct {
	prefix string
	media  string
	params map[string]string
}

const SpyQualifier = "aeron-spy"
const AeronScheme = "aeron"
const AeronPrefix = "aeron:"

const IpcMedia = "ipc"
const UdpMedia = "udp"
const IpcChannel = "aeron:ipc"
const SpyPrefix = "aeron-spy:"
const EndpointParamName = "endpoint"
const InterfaceParamName = "interface"
const InitialTermIdParamName = "init-term-id"
const TermIdParamName = "term-id"
const TermOffsetParamName = "term-offset"
const TermLengthParamName = "term-length"
const MtuLengthParamName = "mtu"
const TtlParamName = "ttl"
const MdcControlParamName = "control"
const MdcControlModeParamName = "control-mode"
const MdcControlModeManual = "manual"
const MdcControlModeDynamic = "dynamic"
const SessionIdParamName = "session-id"
const LingerParamName = "linger"
const ReliableStreamParamName = "reliable"
const TagsParamName = "tags"
const TagPrefix = "tag:"
const SparseParamName = "sparse"
const AliasParamName = "alias"
const EosParamName = "eos"
const TetherParamName = "tether"
const GroupParamName = "group"
const RejoinParamName = "rejoin"
const CongestionControlParamName = "cc"
const FlowControlParamName = "fc"
const GroupTagParamName = "gtag"
const SpiesSimulateConnectionParamName = "ssc"
const SocketSndbufParamName = "so-sndbuf"
const SocketRcvbufParamName = "so-rcvbuf"
const ReceiverWindowLengthParamName = "rcv-wnd"
const MediaRcvTimestampOffsetParamName = "media-rcv-ts-offset"
const ChannelRcvTimestampOffsetParamName = "channel-rcv-ts-offset"
const ChannelSndTimestampOffsetParamName = "channel-snd-ts-offset"

const spyPrefix = SpyQualifier + ":"
const aeronPrefix = AeronScheme + ":"

type parseState int8

const (
	parseStateMedia parseState = iota + 1
	parseStateParamsKey
	parseStateParamsValue
)

// ParseChannelUri parses a string which contains an Aeron URI.
func ParseChannelUri(uriStr string) (uri ChannelUri, err error) {
	uri.params = make(map[string]string)
	if strings.HasPrefix(uriStr, spyPrefix) {
		uri.prefix = SpyQualifier
		uriStr = uriStr[len(spyPrefix):]
	}
	if !strings.HasPrefix(uriStr, aeronPrefix) {
		err = fmt.Errorf("aeron URIs must start with 'aeron:', found '%s'", uriStr)
		return
	}
	uriStr = uriStr[len(aeronPrefix):]

	state := parseStateMedia
	value := ""
	key := ""
	for len(uriStr) > 0 {
		c := uriStr[0]
		switch state {
		case parseStateMedia:
			switch c {
			case '?':
				uri.media = value
				value = ""
				state = parseStateParamsKey
			case ':':
				err = fmt.Errorf("encountered ':' within media definition")
				return
			default:
				value += string(c)
			}
		case parseStateParamsKey:
			if c == '=' {
				key = value
				value = ""
				state = parseStateParamsValue
			} else {
				value += string(c)
			}
		case parseStateParamsValue:
			if c == '|' {
				uri.params[key] = value
				value = ""
				state = parseStateParamsKey
			} else {
				value += string(c)
			}
		default:
			panic("unexpected state")
		}
		uriStr = uriStr[1:]
	}
	switch state {
	case parseStateMedia:
		uri.media = value
	case parseStateParamsValue:
		uri.params[key] = value
	default:
		err = fmt.Errorf("no more input found")
	}
	return
}

// Clone returns a deep copy of a ChannelUri.
func (uri ChannelUri) Clone() (res ChannelUri) {
	res.prefix = uri.prefix
	res.media = uri.media
	res.params = make(map[string]string)
	for k, v := range uri.params {
		res.params[k] = v
	}
	return
}

func (uri ChannelUri) Scheme() string {
	return AeronScheme
}

func (uri ChannelUri) Prefix() string {
	return uri.prefix
}

func (uri ChannelUri) Media() string {
	return uri.media
}

func (uri ChannelUri) IsIpc() bool {
	return uri.media == IpcMedia
}

func (uri ChannelUri) IsUdp() bool {
	return uri.media == UdpMedia
}

func (uri *ChannelUri) SetPrefix(prefix string) {
	uri.prefix = prefix
}

func (uri *ChannelUri) SetMedia(media string) {
	uri.media = media
}

func (uri *ChannelUri) SetControlMode(controlMode string) {
	uri.Set(MdcControlModeParamName, controlMode)
}

func (uri *ChannelUri) SetSessionID(sessionID int32) {
	uri.Set(SessionIdParamName, strconv.Itoa(int(sessionID)))
}

func (uri ChannelUri) Get(key string) string {
	return uri.params[key]
}

func (uri ChannelUri) Set(key string, value string) {
	uri.params[key] = value
}

func (uri ChannelUri) Remove(key string) {
	delete(uri.params, key)
}

func (uri ChannelUri) String() (result string) {
	if uri.prefix != "" {
		result += uri.prefix
		if !strings.HasSuffix(uri.prefix, ":") {
			result += ":"
		}
	}
	result += aeronPrefix
	result += uri.media
	if len(uri.params) > 0 {
		result += "?"
		sortedKeys := make([]string, 0, len(uri.params))
		for k := range uri.params {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
		for _, k := range sortedKeys {
			result += k + "=" + uri.params[k] + "|"
		}
		result = result[:len(result)-1]
	}
	return
}
