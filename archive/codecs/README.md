Aeron Codecs
===

The control protocol for the [Java Aeron
Archive](github.com/real-logic/aeron/aeron-archive) is a collection of
[SBE](github.com/real-logic/simple-binary-encoding)
[messages](github.com/real-logic/aeron/blob/master/aeron-archive/src/main/resources/archive/aeron-archive-codecs.xml).

These generated codecs are added to the repository to reduce the build
prerequisites and to better enable configuration management. They were
built using [generate.sh](./generate.sh).

`generate.sh` should be taken more as a encoded description of what
needs to be done as it needs to know the locations of the additional
tools of simple-binary-encoding and agrona.

Should the protocol version change you will need to update the
semantic version in [semanticversion.go](./semanticversion.go)
