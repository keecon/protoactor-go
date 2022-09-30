package main

import (
	"fmt"

	"cluster-broadcast/shared"

	console "github.com/asynkron/goconsole"
	"github.com/keecon/protoactor-go/actor"
	"github.com/keecon/protoactor-go/cluster"
	"github.com/keecon/protoactor-go/cluster/clusterproviders/consul"
	"github.com/keecon/protoactor-go/cluster/identitylookup/disthash"
	"github.com/keecon/protoactor-go/remote"
)

func main() {
	cluster := startNode()

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	cluster.Shutdown(true)
}

func startNode() *cluster.Cluster {
	system := actor.NewActorSystem()

	provider, _ := consul.New()
	lookup := disthash.New()
	config := remote.Configure("localhost", 0)

	props := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *actor.Started:
			fmt.Printf("Started %v", msg)
		case *shared.Hello:
			fmt.Printf("Hello %v\n", msg.Name)
		}
	})
	helloKind := cluster.NewKind("hello", props)
	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config, cluster.WithKinds(helloKind))
	c := cluster.New(system, clusterConfig)

	c.StartMember()
	return c
}
