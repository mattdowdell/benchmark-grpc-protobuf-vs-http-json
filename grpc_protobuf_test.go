package benchmarks

import (
	"log"
	"testing"
	"time"

	"github.com/plutov/benchmark-grpc-protobuf-vs-http-json/grpc-protobuf"
	"github.com/plutov/benchmark-grpc-protobuf-vs-http-json/grpc-protobuf/proto"
	"golang.org/x/net/context"
	g "google.golang.org/grpc"
)

var grpcClient proto.APIClient

func init() {
	go grpcprotobuf.Start()
	time.Sleep(time.Second)

	conn, err := g.Dial("127.0.0.1:60000", g.WithInsecure())
	if err != nil {
		log.Fatalf("grpc connection failed: %v", err)
	}

	grpcClient = proto.NewAPIClient(conn)
}

func BenchmarkGRPCProtobuf(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doGRPC(grpcClient, b)
		}
	})
}

func doGRPC(client proto.APIClient, b *testing.B) {
	resp, err := client.CreateUser(context.Background(), &proto.User{
		Email:    "foo@bar.com",
		Name:     "Bench",
		Password: "bench",
	})

	if err != nil {
		b.Fatalf("grpc request failed: %v", err)
	}

	if resp == nil || resp.Code != 200 || resp.User == nil || resp.User.Id != "1000000" {
		b.Fatalf("grpc response is wrong: %v", resp)
	}
}
