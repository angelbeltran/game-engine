module github.com/angelbeltran/game-engine/examples/puerto-rico

go 1.14

require (
	github.com/angelbeltran/game-engine/protoc-gen-game v0.0.2
	github.com/golang/protobuf v1.4.2
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/angelbeltran/game-engine/protoc-gen-game => ../../protoc-gen-game
