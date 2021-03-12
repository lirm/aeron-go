Design notes
===

The structure of the archive client library is heavily based on the
Java archive library. It's hoped this will aid comprehension, bug fixing,
feature additions etc.

Many design choices are also based upon the golang client library as
the archive library is a layering on top of that.

Finally golang idioms are used where reasonable.

The archive library does not lock and concurrent calls to archive
library calls that invoke the aeron-archive protocol calls should be
externally locked to ensure only one concuurrent access.

Each Archive instance has it's own aeron instance running a proxy
(outgoing/producer) and control (incoming/subscriber) pair to make up
both halves of what is essentially an RPC based protocol.

Archive operations are hence modelled as synchronous calls. Should an
async API be desired then this can be easily achieved by wrapping the
sync call with a channel, see for example aeron.AddSubscription().

The codecs used are generated using Aeron, and generally the archive library uses those types:

 * It was tempting to make the sourceLocation a boolean ```isLocal```
   butkeeping the underlying codec value allows for future
   additions.

 * Where the protocol specifies a BooelanType a golang bool is used
   until encoding.

 * Recording descriptors are returned using the codecs.RecordingDescriptor type.

Questions/Issues
===

Testing:
 * Look for local archive and exec? Test and not run for Travis? Mock? Add jars to repo and fetch?

Are there any Packets bigger than MTU requiring fragment assembly?
 * Errors *could* conceivably be artbitrarily long strings but this is a) unlikely, b) not currently the case.

Logging INFO - I would like henormal operation to be NORMAL so I can
have an intermediate on DEBUG which is voluminous.

FindLatestRecording() currently in basic_replayed_subscriber useful in API?

startRecording return .. arguably this is an aeron-archive issue and relevantID should be recordingId.

stopRecordingBySubscription - seriously?


Backlog
===
 * Have I really understood java's startRecording(), and the relevantId returned from StartRecording exchange.
 * Test failures / reliability
 * Expand testing
 * Review locking decision. Adding lock/unlock in archive is simple.
 * Ephemeral port usage
 * The archive state is bogus
 * 58 FIXMEs
 * AuthConnect, Challenge/Response
 * Improve the Error handling
 * OnAvailableCounter: Not supported yet? This is going to matter I think
  * See also RecordingIdFromCounter
 * Defaults settings and setting
 * RecordingEvent Handler and Recording Admin (detach segment etc)
 * Clean up and initial upstream push

Done
===

 * Close/Disconnect
 * Simplest straight line basic recorded publisher and basic subscriber

