package main

import (
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/cluster"
)

func main() {
	ctx := aeron.NewContext()
	opts := cluster.NewOptions()
	agent, err := cluster.NewClusteredServiceAgent(ctx, opts)
	if err != nil {
		panic(err)
	}

	agent.OnStart()
	for {
		opts.IdleStrategy.Idle(agent.DoWork())
	}
}
