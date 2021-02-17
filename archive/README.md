Design notes
===

The structure of the archive client library is heavily based on the
Java archive library. It's hoped this will aid comprehension, bug fixing,
feature additions etc.

Many design choices are also based upon the golang client library as
the archive library is a layering on top of that.

Finally golang idioms are used where reasonable.

The library should be fully reentrant.

Each Archive instance has it's won aeron instance running a proxy
(outgoing/producer) and control (incoming/subscriber) pair.

Questions/Issues
===

Sync APIs built on top of Async via channel design choice? For now just sync as we get things working.

Testing:
 * Look for local archive and exec? Test and not run for Travis? Mock? Add jars to repo and fetch?


Backlog
===
Ephemeral port usage

Simplest straight line basic recorded publisher and basic subscriber

Now up to 58 FIXMEs

AuthConnect, Challenge/Response

Close/Disconnect

Improve the Error handling

OnAvailableCounter: Not supported yet? 

Defaults settings and setting

