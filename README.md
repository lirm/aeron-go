[![Build Status](https://travis-ci.org/lirm/aeron-go.svg?branch=master)](https://travis-ci.org/lirm/aeron-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lirm/aeron-go)](https://goreportcard.com/report/github.com/lirm/aeron-go)
[![Join the chat at https://gitter.im/aeron-go/Lobby](https://badges.gitter.im/aeron-go/Lobby.svg)](https://gitter.im/aeron-go/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<!--[![Coverage Status](https://coveralls.io/repos/github/lirm/aeron-go/badge.svg?branch=master)](https://coveralls.io/github/lirm/aeron-go?branch=master)-->

# aeron-go

Implementation of [Aeron](https://github.com/real-logic/Aeron) messaging client in Go.

Architecture, design, and protocol of Aeron can be found [here](https://github.com/real-logic/Aeron/wiki)

# Usage

Example subscriber can be found [here](https://github.com/lirm/aeron-go/tree/master/examples/basic_subscriber).

Example publication can be found [here](https://github.com/lirm/aeron-go/tree/master/examples/basic_publisher).

## Common

Instantiate Aeron with Context:
```go
ctx := aeron.NewContext().MediaDriverTimeout(time.Second * 10)

a := aeron.Connect(ctx)
```

## Subscribers

Create subscription:
```go
subscription := <-a.AddSubscription("aeron:ipc", 10)

defer subscription.Close()
```

`aeron.AddSubscription()` returns a channel, so that the user has the choice
of blocking waiting for subscription to register with the driver or do async `select` poll.

Define callback for message processing:
```go
handler := func(buffer *buffers.Atomic, offset int32, length int32, header *logbuffer.Header) {
    bytes := buffer.GetBytesArray(offset, length)

    fmt.Printf("Received a fragment with payload: %s\n", string(bytes))
}
```

Poll for messages:
```go
idleStrategy := idlestrategy.Sleeping{time.Millisecond}

for {
    fragmentsRead := subscription.Poll(handler, 10)
    idleStrategy.Idle(fragmentsRead)
}
```

## Publications

Create publication:
```go
publication := <-a.AddPublication("aeron:ipc", 10)

defer publication.Close()
```

`aeron.AddPublication()` returns a channel, so that the user has the choice
of blocking waiting for publication to register with the driver or do async `select` poll.

Create Aeron buffer to send the message:
```go
message := fmt.Sprintf("this is a message %d", counter)

srcBuffer := buffers.MakeAtomic(([]byte)(message))
```

Optionally make sure that there are connected subscriptions:
```go
for !publication.IsConnected() {
    time.Sleep(time.Millisecond * 10)
}
```

Send the message, by calling `publication.Offer`
```go
ret := publication.Offer(srcBuffer, 0, int32(len(message)), nil)
switch ret {
case aeron.NotConnected:
    log.Print("not connected yet")
case aeron.BackPressured:
    log.Print("back pressured")
default:
    if ret < 0 {
        log.Print("Unrecognized code: %d", ret)
    } else {
        log.Print("success!")
    }
}
```
