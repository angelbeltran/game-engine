module github.com/angelbeltran/game-engine/examples/puerto-rico

go 1.14

require (
	github.com/angelbeltran/game-engine/protoc-gen-game v0.0.2
	github.com/golang/protobuf v1.4.2
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/angelbeltran/game-engine/protoc-gen-game => ../../protoc-gen-game
