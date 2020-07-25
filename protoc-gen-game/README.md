# Protoc Plugin: protoc-gen-game

This package generates a binary, `protoc-gen-game` to be used as a plugin for `protoc`, Google's protobuf compiler.

The plugin is meant to generate a kind of rule engine, one useful for simple game development, eg simple board or card games.

The rule engine is defined by a set of .proto annotations, or options, on a service that will host the rule engine and game server, and on a message that will represent the state of a game.

Check out the [examples](../examples) and [game_engine.proto](game_engine.proto) for an idea of how it all works.
