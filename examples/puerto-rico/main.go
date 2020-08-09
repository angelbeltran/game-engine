package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"

	pb "github.com/angelbeltran/game-engine/examples/puerto-rico/puerto_rico_pb"
)

//go:generate rm -rf puerto_rico_pb
//go:generate mkdir puerto_rico_pb
//go:generate protoc -I=$GOPATH/src/ --go_out=$GOPATH/src --go-grpc_out=$GOPATH/src --plugin=protoc-gen-game=$GOPATH/src/github.com/angelbeltran/game-engine/protoc-gen-game/protoc-gen-game --game_out=puerto_rico_pb --game_opt=package=puerto_rico_pb $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/actions.proto $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/building.proto $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/cargo_ship.proto $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/good.proto $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/plantation.proto $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/player.proto $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/role.proto $GOPATH/src/github.com/angelbeltran/game-engine/examples/puerto-rico/state.proto

func main() {
	srv, lis, err := pb.NewServer(8080)
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("server shutdown unexpectedly: %v", err)
			os.Exit(1)
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create a client connection: %v", err)
	}

	cli := pb.NewActionsClient(conn)

	if err = testClient(cli); err != nil {
		log.Fatalf("failed test: %v", err)
	}
}

func testClient(cli pb.ActionsClient) error {
	res1, err := cli.SetPlayers(context.Background(), &pb.Count{Count: 3}, grpc.WaitForReady(true))
	if err != nil {
		return fmt.Errorf("failed to perform Count: %v", err)
	}

	fmt.Println("Count:", res1)

	res2, err := cli.Start(context.Background(), &pb.EmptyMsg{}, grpc.WaitForReady(true))
	if err != nil {
		return fmt.Errorf("failed to perform Start: %v", err)
	}

	fmt.Println("Start:", res2)

	return nil
}
