package main

import (
	"cluster-basic/shared"
	"fmt"
	console "github.com/asynkron/goconsole"
	"github.com/keecon/protoactor-go/actor"
	"github.com/keecon/protoactor-go/cluster"
	"github.com/keecon/protoactor-go/cluster/clusterproviders/consul"
	"github.com/keecon/protoactor-go/cluster/identitylookup/disthash"
	"github.com/keecon/protoactor-go/remote"
)

func main() {
	c := startNode()

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()
	pid := c.Get("abc", "hello")
	fmt.Printf("Got pid %v", pid)
	res, _ := c.Request("abc", "hello", &shared.HelloRequest{Name: "Roger"})
	fmt.Printf("Got response %v", res)

	fmt.Println()
	console.ReadLine()
	c.Shutdown(true)
}

func startNode() *cluster.Cluster {
	system := actor.NewActorSystem()

	provider, _ := consul.New()
	lookup := disthash.New()
	config := remote.Configure("localhost", 0)
	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config)
	c := cluster.New(system, clusterConfig)
	c.StartMember()

	return c
}
