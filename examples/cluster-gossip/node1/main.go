package main

import (
	"cluster-gossip/shared"
	"fmt"
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/consul"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

func main() {
	c := startNode()

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()
	pid := c.Get("abc", "hello")
	fmt.Printf("Got pid %v\n", pid)
	res, _ := c.Request("abc", "hello", &shared.HelloRequest{Name: "Roger"})
	fmt.Printf("Got response %v\n", res)

	fmt.Println()
	console.ReadLine()
	c.Shutdown(true)
}

func coloredConsoleLogging(system *actor.ActorSystem) *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelError,
		TimeFormat: time.RFC3339,
		AddSource:  true,
	})).With("lib", "Proto.Actor").
		With("system", system.ID)
}

func startNode() *cluster.Cluster {
	system := actor.NewActorSystem(actor.WithLoggerFactory(coloredConsoleLogging))
	system.EventStream.Subscribe(func(evt interface{}) {
		switch msg := evt.(type) {
		case *cluster.ClusterTopology:
			fmt.Printf("\nClusterTopology %v\n\n", msg)
		case *cluster.GossipUpdate:

			heartbeat := &cluster.MemberHeartbeat{}

			fmt.Printf("Member %v\n", msg.MemberID)
			fmt.Printf("Sequence Number %v\n", msg.SeqNumber)

			unpackErr := msg.Value.UnmarshalTo(heartbeat)
			if unpackErr != nil {
				fmt.Printf("Unpack error %v\n", unpackErr)
			} else {
				//loop over as.ActorCount map
				for k, v := range heartbeat.ActorStatistics.ActorCount {
					fmt.Printf("ActorCount %v %v\n", k, v)
				}
			}
		}
	})

	provider, _ := consul.New()
	lookup := disthash.New()
	config := remote.Configure("localhost", 0)
	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config)
	c := cluster.New(system, clusterConfig)
	c.StartMember()

	return c
}
