package protoactor_test

import (
	"testing"

	"github.com/keecon/protoactor-go/actor"
	"github.com/keecon/protoactor-go/cluster"
	cluster_test_tool "github.com/keecon/protoactor-go/cluster/cluster_test_tool"
	"github.com/keecon/protoactor-go/remote"
	"github.com/keecon/protoactor-go/router"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestGeneratedProtoGoPackageMatchesModulePath(t *testing.T) {
	tests := []struct {
		name       string
		descriptor protoreflect.FileDescriptor
		want       string
	}{
		{name: "actor", descriptor: actor.File_actor_proto, want: "github.com/keecon/protoactor-go/actor"},
		{name: "cluster", descriptor: cluster.File_cluster_proto, want: "/github.com/keecon/protoactor-go/cluster"},
		{name: "gossip", descriptor: cluster.File_gossip_proto, want: "/github.com/keecon/protoactor-go/cluster"},
		{name: "grain", descriptor: cluster.File_grain_proto, want: "/github.com/keecon/protoactor-go/cluster"},
		{name: "pubsub", descriptor: cluster.File_pubsub_proto, want: "/github.com/keecon/protoactor-go/cluster"},
		{name: "pubsub test", descriptor: cluster.File_pubsub_test_proto, want: "/github.com/keecon/protoactor-go/cluster"},
		{name: "cluster test tool", descriptor: cluster_test_tool.File_pubsub_cluster_proto, want: "/github.com/keecon/protoactor-go/cluster/cluster_test_tool"},
		{name: "remote", descriptor: remote.File_remote_proto, want: "/github.com/keecon/protoactor-go/remote"},
		{name: "router", descriptor: router.File_routercontracts_proto, want: "github.com/keecon/protoactor-go/router"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options, ok := test.descriptor.Options().(*descriptorpb.FileOptions)
			if !ok {
				t.Fatalf("unexpected file options type %T", test.descriptor.Options())
			}
			if got := options.GetGoPackage(); got != test.want {
				t.Fatalf("go_package = %q, want %q", got, test.want)
			}
		})
	}
}
