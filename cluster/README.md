# aeron-go/cluster

Implementation of [Aeron Cluster](https://github.com/real-logic/Aeron/tree/master/aeron-cluster) service container in Go.
Most structs and functions have near one-to-one parity with the Java classes and
methods on which they are based.

The [Java media driver, archive and consensus module](https://github.com/real-logic/aeron/blob/master/aeron-cluster/src/main/java/io/aeron/cluster/ClusteredMediaDriver.java)
must be used to run a cluster.

The [Aeron Cluster
protocol](http://github.com/real-logic/aeron/blob/master/aeron-cluster/src/main/resources/cluster/aeron-cluster-codecs.xml)
is specified in xml using the [Simple Binary Encoding (SBE)](https://github.com/real-logic/simple-binary-encoding).

## Current State
The implementation is functional and mostly feature complete, including support
for snapshotting, timers, multiple services within the same cluster, sending messages
back to cluster sessions, and service mark files. The [Cluster](cluster.go) interface
lacks of the methods of its [Java equivalent](https://github.com/real-logic/aeron/blob/master/aeron-cluster/src/main/java/io/aeron/cluster/service/Cluster.java),
but these would be trivial additions.

The [Cluster Client](https://github.com/real-logic/aeron/blob/master/aeron-cluster/src/main/java/io/aeron/cluster/client/AeronCluster.java)
has not been implemented yet. One must use the Java version to connect to cluster.

## Examples

[echo_service.go](examples/echo_service.go) implements a basic echo service and can be
used in place of its [Java equivalent](https://github.com/real-logic/aeron/blob/master/aeron-samples/src/main/java/io/aeron/samples/cluster/EchoService.java).

# Backlog
 * godoc improvements
 * testing
 * cluster client
 * cluster session close handling (avoid sending duplicate close requests to consensus module)
 * SBE encoding/decoding improvements
