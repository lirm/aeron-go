package aeron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelUri(t *testing.T) {
	assert := assert.New(t)

	assertParseWithMedia(assert, "aeron:udp", "udp")
	assertParseWithMedia(assert, "aeron:ipc", "ipc")
	assertParseWithMedia(assert, "aeron:", "")
	assertParseWithMediaAndPrefix(assert, "aeron-spy:aeron:ipc", "aeron-spy", "ipc")
}

func TestShouldRejectUriWithoutAeronPrefix(t *testing.T) {
	assert := assert.New(t)
	assertInvalid(assert, ":udp")
	assertInvalid(assert, "aeron")
	assertInvalid(assert, "aron:")
	assertInvalid(assert, "eeron:")
}

func TestShouldRejectWithOutOfPlaceColon(t *testing.T) {
	assert := assert.New(t)
	assertInvalid(assert, "aeron:udp:")
}

func TestShouldParseWithMultipleParameters(t *testing.T) {
	assert := assert.New(t)
	assertParseWithParams(
		assert,
		"aeron:udp?endpoint=224.10.9.8|port=4567|interface=199.168.0.3|ttl=16",
		"endpoint", "224.10.9.8",
		"port", "4567",
		"interface", "199.168.0.3",
		"ttl", "16")
}

func TestShouldRoundTripToString(t *testing.T) {
	assert := assert.New(t)

	assertRoundTrip(assert, "aeron:udp?endpoint=224.10.9.8:777")
	assertRoundTrip(assert, "aeron-spy:aeron:udp?control=224.10.9.8:777")
}

func TestCloneShouldRoundTripToString(t *testing.T) {
	assert := assert.New(t)

	assertCloneRoundTrip(assert, "aeron:udp?endpoint=224.10.9.8:777")
	assertCloneRoundTrip(assert, "aeron-spy:aeron:udp?control=224.10.9.8:777")
	assertCloneRoundTrip(assert, "aeron:udp?endpoint=224.10.9.8|interface=199.168.0.3|port=4567|ttl=16")
}

func assertRoundTrip(assert *assert.Assertions, uriString string) {
	uri, err := ParseChannelUri(uriString)
	assert.NoError(err)
	assert.EqualValues(uriString, uri.String())
}

func assertCloneRoundTrip(assert *assert.Assertions, uriString string) {
	uri, err := ParseChannelUri(uriString)
	assert.NoError(err)
	assert.EqualValues(uriString, uri.Clone().String())
}

func assertParseWithParams(assert *assert.Assertions, uriString string, params ...string) {
	uri, err := ParseChannelUri(uriString)
	assert.NoError(err)
	if len(params)%2 != 0 {
		panic("must be key=value params")
	}
	for i := 0; i < len(params); i += 2 {
		assert.EqualValues(uri.Get(params[i]), params[i+1])
	}
}

func assertInvalid(assert *assert.Assertions, uriString string) {
	_, err := ParseChannelUri(uriString)
	assert.Error(err)
}

func assertParseWithMediaAndPrefix(assert *assert.Assertions, uriString string, prefix string, media string) {
	uri, err := ParseChannelUri(uriString)
	assert.NoError(err)
	assert.EqualValues(uri.Scheme(), "aeron")
	assert.EqualValues(uri.Prefix(), prefix)
	assert.EqualValues(uri.Media(), media)
}

func assertParseWithMedia(assert *assert.Assertions, uriString string, media string) {
	assertParseWithMediaAndPrefix(assert, uriString, "", media)
}
