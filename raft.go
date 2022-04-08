package main

import (
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/cluster"
)

func main() {
	ctx := aeron.NewContext()
	opts := &cluster.Options{
		Timeout:       time.Second,
		IdleStrategy:  idlestrategy.NewDefaultBackoffIdleStrategy(),
		RangeChecking: true,
	}
	agent, err := cluster.NewClusteredServiceAgent(ctx, opts)
	if err != nil {
		panic(err)
	}

	agent.OnStart()
}
