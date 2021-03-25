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

It was tempting to make the sourceLocation a boolean ```isLocal``` but
keeping the underlying codec value allows for future
additionals. However, where the protocol specifies a BooelanType a
bool is used until encoding.

## Structure

The archive protocol is largely an RPC mechanism built on top of
Aeron. Each Archive instance has it's own aeron instance running a
proxy (publication/request) and control (subscription/response)
pair. This mirrors the Java implementation.

Additionally there are some asynchronous events that can arrive on a
RecordingEvents subscription and these are handled by the
recordingevents. These are not enabled by default to avoid using
resources when not required.

# Questions/Issues
Testing:
 * Look for local archive and exec? Test and not run for Travis? Mock? Add jars to repo and fetch?

Are there any Packets bigger than MTU requiring fragment assembly?
 * Errors *could* conceivably be artbitrarily long strings but this is a) unlikely, b) not currently the case.

Logging INFO - I'd rathernormal operations to be NORMAL so I can
have an intermediate on DEBUG which is voluminous.

FindLatestRecording() currently in basic_replayed_subscriber useful in API?

startRecording return .. arguably this is an aeron-archive issue and relevantID should be recordingId? upstream?

startRecording polymorphism ... leading to API choices

# Backlog
 * [S] Improve the Error handling / Error listeners
 * [?] Have I really understood java's startRecording(), and the relevantId returned from StartRecording exchange.
 * [L] Expand testing
  * [S] Test failures
  * [S] Test reliability
  * [?] archive-media-driver mocking/execution
 * [S] Review locking decision. Adding lock/unlock in archive is simple.
 * [?] Ephemeral port usage
 * [S} The archive state is bogus
   * IsOpen()?
 * 39 FIXMEs
 * [?] AuthConnect, Challenge/Response
 * [?] OnAvailableCounter: Not supported yet? This may matter I think
  * See also RecordingIdFromCounter
 * Logging tidy
 * Add remaining archive protocol packets to encoders, control, and API.
 * Is it worth having multiple Idletimers in different places? probably.
 * Encoders might take marshaller as parameter to assist reentrancy
 * [?] Clean up and initial upstream push
 

# Done
 * [x] RecordingEvent Handler and Recording Admin (detach segment etc)
 * [x] Rework the archive context, using it in both the conrol And RecrodingEvents

