package main

import (
	"fmt"

	pb "angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"angelbeltran/game-engine/protoc-gen-game/types"
)

func extractValue(msg *pb.Value) (types.Type, interface{}, error) {
	switch v := msg.GetValue().(type) {
	case *pb.Value_Bool:
		return types.Bool{}, v.Bool, nil
	case *pb.Value_Integer:
		return types.Integer{}, v.Integer, nil
	case *pb.Value_Float:
		return types.Float{}, v.Float, nil
	case *pb.Value_String_:
		return types.String{}, v.String_, nil
	}

	return nil, nil, fmt.Errorf("unexpected value type: %T", msg.GetValue())
}
