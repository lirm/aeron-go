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
based on the Java names so that developers can more easilt switch
between langugages and so that any API documentation is more useful
across implementations. Some differences exist due to capitalization
requirements, lack of polymorphism, etc.

Function names used in encoders.go and proxy.go are based on the
protocol specification.

It was tempting to make the sourceLocation a boolean ```isLocal``` but
keeping the underlying codec value allows for future additions. Where
the protocol specifies a BooelanType a bool is used until encoding.

## Structure

The archive protocol is largely an RPC mechanism built on top of
Aeron. Each Archive instance has it's own aeron instance running a
proxy (publication/request) and control (subscription/response)
pair. This mirrors the Java implementation.

Additionally there are some asynchronous events that can arrive on a
RecordingEvents subscription and these are handled by the
recordingevents. These are not enabled by default to avoid using
resources when not required.

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
wrapped in a goroutine if it's desired.

# Questions/Issues

Logging INFO - I'd rather normal operations to be NORMAL so I can have
an intermediate on DEBUG which is voluminous.

# Backlog
 * [S] Improve the Error handling / Error listeners
 * [L] Expand testing
  * [M] So many tests to write
  * [S] Test reliability
  * [?] archive-media-driver mocking/execution
 * [S] Review locking decision. Adding lock/unlock in archive is simple.
 * [?] Ephemeral port usage
 * [S} The archive state is largely unused. 
   * IsOpen()?
 * 41 FIXMEs
 * [?] AuthConnect, Challenge/Response
 * [?] OnAvailableCounter noise
 * [?] [Bug] RecordingSignalEvents currently throw off the count of
   fragments/records we want. Need a mechanism to adjust for them.
 * [?] Logging tidy and general improvements
 * [?] Add remaining archive protocol packets to encoders, proxy, control, archive API, and tests.
 * [S] Is it worth having multiple Idletimers in different places? probably.
 * [S] Encoders might take marshaller as parameter to assist reentrancy
 * [?] Clean up and initial upstream push
 

# Done
 * [x] RecordingEvent Handler and Recording Admin (detach segment etc)
 * [x] Rework the archive context, using it in both the conrol And RecrodingEvents

