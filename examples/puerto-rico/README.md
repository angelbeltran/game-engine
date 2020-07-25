# Puerto Rico boardgame example


## Generating Protobufs


#### Install Dependencies

`protoc`, `protoc-gen-go`, and `proto-gen-grpc`

#### Build `proto-gen-game`

From the proto-gen-game directory, run `go generate` and `go build`

#### Generate the Protobufs

From this directory, run `go generate`


## Running the Server

After generating the service,

`go run main.go`


## Check Out

#### [actions.proto](actions.proto)

You'll find some interesting service annotations.
The custom options you'll see are how to instruct protoc-gen-game the dynamics, rules, and state of the game.

Also, you'll find a message type, `Response`, with an annotation that indicates to protoc-gen-game that it is used as the output type of each rpc method of the service.
The two fields `state` and `error` are required.
Additional fields are allowed and will likely be populated by plugins to protoc-gen-game.

#### [state.proto](state.proto)

The State message has an annotation indicating it represents the state of the game.
In actions.proto, `state` must be this message.
Precisely one message should have this annotation.
