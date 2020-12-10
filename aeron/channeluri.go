package aeron

import (
	"fmt"
	"sort"
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

// AeronScheme is a URI Scheme for Aeron channels and destinations.
const AeronScheme = "aeron"

// SpyQualifier is a qualifier for spy subscriptions which spy on outgoing network destined traffic efficiently.
const SpyQualifier = "aeron-spy"

// IpcMedia is the media for IPC.
const IpcMedia = "ipc"

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

func (uri *ChannelUri) SetPrefix(prefix string) {
	uri.prefix = prefix
}

func (uri ChannelUri) SetMedia(media string) {
	uri.media = media
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
