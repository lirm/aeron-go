# aeron-go/archive

Implementation of [Aeron Archive](https://github.com/real-logic/Aeron/tree/master/aeron-archive) client in Go.

The [Aeron Archive
protocol](http://github.com/real-logic/aeron/blob/master/aeron-archive/src/main/resources/archive/aeron-archive-codecs.xml)
is specified in xml using the [Simple Binary Encoding (SBE)](https://github.com/real-logic/simple-binary-encoding)

## Current State
The implementation is the second beta release. The API may still be changed for bug fixes or significant issues.

# Design

## Guidelines

The structure of the archive client library is heavily based on the
Java archive library. It's hoped this will aid comprehension, bug fixing,
feature additions etc.

Many design choices are also based upon the golang client library as
the archive library is a layering on top of that.

Finally golang idioms are used where reasonable.

The archive library does not lock and concurrent calls on an archive
client to invoke the aeron-archive protocol calls should be externally
locked to ensure only one concurrent access.

### Naming and other choices

Function names used in archive.go which contains the main API are
based on the Java names so that developers can more easily switch
between languages and so that any API documentation is more useful
across implementations. Some differences exist due to capitalization
requirements, lack of polymorphism, etc.

Function names used in encoders.go and proxy.go are based on the
protocol specification. Where the protocol specifies a type that can
be naturally represented in golang, the golang type is used used where
possible until encoding. Examples include the use of `bool` rather than
`BooleanType` and `string` over `[]uint8`

## Structure

The archive protocol is largely an RPC mechanism built on top of
Aeron. Each Archive instance has it's own aeron instance running a
[proxy](proxy.go) (publication/request) and [control](control.go) (subscription/response)
pair. This mirrors the Java implementation. The proxy invokes the
encoders to marshal packets using SBE.

Additionally there are some asynchronous events that can arrive on a
[recordingevents](recordingevents.go) subscription. These
are not enabled by default to avoid using resources when not required.

## Synchronous unlocked API optionally using polling

The implementation provides a synchronous API as the underlying
mechanism is largely an RPC mechanism and archive operations are not
considered high frequency.

Associated with this, the library does not lock and assumes management
of reentrancy is handled by the caller.

If needed it is simple in golang to wrap a synchronous API with a
channel (see for example aeron.AddSubscription(). If overlapping
asynchronous calls are needed then this is where you can add locking.

Some asynchronous events do exist (e.g, recording events) and to be
delivered a polling mechanism is provided. Again this can be easily
wrapped in a goroutine if it's desired but ensure there are no other
operations in progress when polling.

## Examples

Examples are provided for a [basic_recording_publisher](examples/basic_recording_publisher/basic_recording_publisher.go) and [basic_replayed_subscriber](examples/basic_replayed_subscriber/basic_replayed_subscriber.go) that interoperate with the Java examples.

## Security

Enabling security is done via setting the various auth options. [config_test.go](config_test.go) and [archive_test.go](archive_test.go) provide an example.

The actual semantics of the security are dependent upon which authenticator supplier you use and is tested against [secure-logging-archiving-media-driver](secure-logging-archiving-media-driver).

# Backlog
 * godoc improvements
 * more testing
  * archive-media-driver mocking/execution
  * test cleanup in the media driver can be problematic
 * Auth should provide some callout mechanism
 * various FIXMEs

# Bigger picture issues
 * Decided not to do locking in sync api, could subsequently add locks, or just async with locks.
   It may be that the control marshaller should be parameterized for this.
 * Java and C++ poll the counters to determine when a recording has actually started but the counters are not
   available in go. As a result we use delays and 'hope' which isn't ideal.
 * It would be nice to silence the OnAvailableCounter noise
 * Within aeron-go there are cases of Log.Fatalf(), see for example trying to add a publication on a "bogus" channel.

## Release Notes

### 1.0b2
 * Handle different archive clients using same channel/stream pairing
 * Provide Subscription.IsConnected() and IsRecordingEventsConnected()
 * Replace go-logging with zap to avoid reentrancy crashes in logging library
 * Improve the logging by downgrading in severity and message tone some warning level messages
 * Fix a race condition looking up correlationIDs on ControlResponse
 * Fix a return code error in StopRecordingById()
 * Fix unused argumentin StopRecording()
 * Cosmetic improvements for golint and staticcheck
 * Make the Listeners used for async events be per archive client instead of global
 