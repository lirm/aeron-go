package main

import (
	"fmt"
	"os"
	"time"

	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"github.com/corymonroe-coinbase/aeron-go/aeron/idlestrategy"
	"github.com/corymonroe-coinbase/aeron-go/aeron/logbuffer"
	"github.com/corymonroe-coinbase/aeron-go/cluster"
	"github.com/corymonroe-coinbase/aeron-go/cluster/client"
)

type ClusterClient struct {
	cluster      cluster.Cluster
	messageCount int32
}

func (c *ClusterClient) OnConnect(cluster *client.AeronCluster) {
	fmt.Printf("OnConnect - sessionId=%d leaderMemberId=%d leadershipTermId=%d\n",
		cluster.ClusterSessionId(), cluster.LeaderMemberId(), cluster.LeadershipTermId())
}

func (c *ClusterClient) OnDisconnect(cluster *client.AeronCluster, details string) {
	fmt.Printf("OnDisconnect - sessionId=%d (%s)\n", cluster.ClusterSessionId(), details)
}

func (c *ClusterClient) OnMessage(cluster *client.AeronCluster, timestamp int64,
	buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	recvTime := time.Now().UnixNano()
	sendTime := buffer.GetInt64(offset)
	latency := recvTime - sendTime
	fmt.Printf("OnMessage - sessionId=%d timestamp=%d pos=%d length=%d latency=%d\n",
		cluster.ClusterSessionId(), timestamp, header.Position(), length, latency)
	c.messageCount++
}

func (c *ClusterClient) OnNewLeader(cluster *client.AeronCluster, leadershipTermId int64, leaderMemberId int32) {
	fmt.Printf("OnNewLeader - sessionId=%d leaderMemberId=%d leadershipTermId=%d\n",
		cluster.ClusterSessionId(), leaderMemberId, leadershipTermId)
}

func (c *ClusterClient) OnError(cluster *client.AeronCluster, details string) {
	fmt.Printf("OnError - sessionId=%d: %s\n", cluster.ClusterSessionId(), details)
}

func main() {
	ctx := aeron.NewContext()
	if aeronDir := os.Getenv("AERON_DIR"); aeronDir != "" {
		ctx.AeronDir(aeronDir)
		fmt.Println("aeron dir: ", aeronDir)
	} else if _, err := os.Stat("/dev/shm"); err == nil {
		path := fmt.Sprintf("/dev/shm/aeron-%s", aeron.UserName)
		ctx.AeronDir(path)
		fmt.Println("aeron dir: ", path)
	}

	opts := client.NewOptions()
	if idleStr := os.Getenv("NO_OP_IDLE"); idleStr != "" {
		opts.IdleStrategy = &idlestrategy.Busy{}
	}
	opts.IngressChannel = "aeron:udp?alias=cluster-client-ingress|endpoint=localhost:20000"
	opts.IngressEndpoints = "0=localhost:20000,1=localhost:21000,2=localhost:22000"
	opts.EgressChannel = "aeron:udp?alias=cluster-egress|endpoint=localhost:11111"

	listener := &ClusterClient{}
	clusterClient, err := client.NewAeronCluster(ctx, opts, listener)
	if err != nil {
		panic(err)
	}

	sendBuf := atomic.MakeBuffer(make([]byte, 100))
	nextSendTime := time.Now().UnixNano()
	for {
		opts.IdleStrategy.Idle(clusterClient.Poll())
		now := time.Now().UnixNano()
		if clusterClient.IsConnected() && now > nextSendTime {
			nextSendTime = now + time.Second.Nanoseconds()
			sendBuf.PutInt64(0, now)
			if result := clusterClient.Offer(sendBuf, 0, sendBuf.Capacity()); result < 0 {
				fmt.Printf("WARNING: failed to send to cluster, result=%d\n", result)
			}
		}
	}
}
