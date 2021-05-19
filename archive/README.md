# aeron-go/archive

Implementation of [Aeron Archive](https://github.com/real-logic/Aeron/tree/master/aeron-archive) client in Go.

The [Aeron Archive
protocol](http://github.com/real-logic/aeron/blob/master/aeron-archive/src/main/resources/archive/aeron-archive-codecs.xml)
is specified in xml using the [Simple Binary Encoding (SBE)](https://github.com/real-logic/simple-binary-encoding)

## Current State
The implementation is an alpha release. The API is not yet considered 100% stable.

# Design

## Guidelines

The structure of the archive client library is heavily based on the
Java archive library. It's hoped this will aid comprehension, bug fixing,
feature additions etc.

Many design choices are also based upon the golang client library as
the archive library is a layering on top of that.

Finally golang idioms are used where reasonable.

The archive library does not lock and concurrent calls to archive
library calls that invoke the aeron-archive protocol calls should be
externally locked to ensure only one concuurrent access.

### Naming and other choices

Function names used in archive.go which contains the main API are
based on the Java names so that developers can more easily switch
between langugages and so that any API documentation is more useful
across implementations. Some differences exist due to capitalization
requirements, lack of polymorphism, etc.

Function names used in encoders.go and proxy.go are based on the
protocol specification. Where the protocol specifies a type that cab
ne naturally repreented in golang, the golang type is used used where
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

Examples are provided for a [basic_recording_publisher](examples/basic_recording_publisher/basic_recording_publisher.go) and [basic_replayed_subscriber](examples/basic_replayed_subscriber/basic_replayed_subscriber.go) that interoperate with the Java examples

# Backlog

## Working Set
 * [S] [Bug] RecordingSignalEvents currently throw off the count of
   fragments/records we want. Need a mechanism to adjust for them.
 * [L] Expand testing
  * [M] So many tests to write
  * [?] archive-media-driver mocking/execution
  * test cleanup in the media driver can be problematic
 * [S} The archive state is largely unused. 
   * IsOpen()?
 * 10 FIXMEs
 * [?] Implement AuthConnect, Challenge/Response
 * [?] Add remaining archive protocol packets to proxy, control, archive API, and tests.

## Recently Done
 * Logging at level normal should be mostly quiet if nothing goes wrong
 * Improve the Error handling / Error listeners (mostly)a
 * Ephemeral port usage is dependent upon accessing the counters which is out of scope here and doesn't buy much
 * Error listener
 * Logging tidying
 * Removed the archive context, it was offering little value. Instead,
   the proxy, control, and recrodingevents all have a reference 
 * Made tests a little reliable but cleanup is still a problem

# Bigger picture issues
 * Decided not to do locking in sync api, could subsequently add locks, or just async with locks if desired.
   It may be that the marshaller should be parameterized for this.
 * Java and C++ poll the counters to determine when a recording has actually started but the counters are not
   availabe in go. As a result we use delays and hope which isn't ideal.
 * OnAvailableCounter noise
 * Within aeron-go there are cases of Log.Fatalf(), see for example trying to add a publication on a "bogus" channel.
