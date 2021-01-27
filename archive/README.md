Design notes
===

The structure of the archive client library is heavily based on the
Java archive library. It's hoped this will aid comprehension, bug fixing,
feature additions etc.

Many design choices are also based upon the golang client library as
the archive library is a layering on top of that.

Finally golang idioms are used where reasonable.

The library should be fully reentrant.


Questions
===

Can we assume a 1-1 mapping of archive instance to Control pair? I'm assuming not.

Sync APIs built on top of Async via channel design choice?

Testing:
 * Look for local archive and exec? Test and not run for Travis? Mock? Add jars to repo and fetch?



Backlog
===
Known bugs:
 * Duplicate connection controlresponse

29 FIXMEs

CorrelationID

Simplest straight line basic recorded publisher and basic subscriber

Async Connection establishment

AuthConnect, Challenge/Response

Close/Disconnect

Improve the Error handling

OnAvailableCounter: Not supported yet?

Defaults settings and setting

Testing:
 * Start with connection establishment