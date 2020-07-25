module github.com/angelbeltran/game-engine/examples/puerto-rico

go 1.14

require (
	github.com/angelbeltran/game-engine/protoc-gen-game v0.0.0-20200725023809-075a31ba5f76
	github.com/golang/protobuf v1.4.2
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
)

//replace github.com/angelbeltran/game-engine@v0.0.0 => /Users/abeltran/go/src/github.com/angelbeltran/game-engine
replace github.com/angelbeltran/game-engine => ../../

replace github.com/angelbeltran/game-engine/game_engine_pb => ../../game_engine_pb

replace github.com/angelbeltran/game-engine/examples/puerto-rico/puerto_rico_pb => ./puerto_rico_pb
