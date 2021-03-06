// Generated by protoc-gen-game. DO NOT EDIT.
package puerto_rico_pb

import (
	"context"
	"fmt"
	"github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"github.com/angelbeltran/game-engine/protoc-gen-game/generate/go_func"
	"google.golang.org/grpc"
	"net"
	"sync"
)

//
// Server and state initialization
//

func NewServer(port uint) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, err
	}

	srv := grpc.NewServer()
	RegisterActionsServer(srv, new(gameEngine))

	return srv, lis, nil
}

type gameEngine struct {
	UnimplementedActionsServer
}

var state = newGameState()

func newGameState() gameState {
	s := gameState{
		State: new(State),
		Mutex: new(sync.Mutex),
	}

	s.State.Players = new(State_Players)
	s.State.Players.Player_1 = new(Player)
	s.State.Players.Player_1.Buildings = new(Player_Buildings)
	s.State.Players.Player_1.Plantations = new(Player_Plantations)
	s.State.Players.Player_1.Goods = new(Player_Goods)
	s.State.Players.Player_2 = new(Player)
	s.State.Players.Player_2.Buildings = new(Player_Buildings)
	s.State.Players.Player_2.Plantations = new(Player_Plantations)
	s.State.Players.Player_2.Goods = new(Player_Goods)
	s.State.Players.Player_3 = new(Player)
	s.State.Players.Player_3.Buildings = new(Player_Buildings)
	s.State.Players.Player_3.Plantations = new(Player_Plantations)
	s.State.Players.Player_3.Goods = new(Player_Goods)
	s.State.Players.Player_4 = new(Player)
	s.State.Players.Player_4.Buildings = new(Player_Buildings)
	s.State.Players.Player_4.Plantations = new(Player_Plantations)
	s.State.Players.Player_4.Goods = new(Player_Goods)
	s.State.Players.Player_5 = new(Player)
	s.State.Players.Player_5.Buildings = new(Player_Buildings)
	s.State.Players.Player_5.Plantations = new(Player_Plantations)
	s.State.Players.Player_5.Goods = new(Player_Goods)
	s.State.Roles = new(State_Roles)
	s.State.Roles.Prospector1 = new(Role)
	s.State.Roles.Prospector2 = new(Role)
	s.State.Roles.Builder = new(Role)
	s.State.Roles.Captain = new(Role)
	s.State.Roles.Craftsman = new(Role)
	s.State.Roles.Mayor = new(Role)
	s.State.Roles.Settler = new(Role)
	s.State.Roles.Trader = new(Role)
	s.State.Plantations = new(State_Plantations)
	s.State.Plantations.Displayed = new(State_Plantations_Displayed)
	s.State.Plantations.Facedown = new(State_Plantations_Counts)
	s.State.Plantations.Discarded = new(State_Plantations_Counts)
	s.State.Goods = new(State_Goods)
	s.State.Buildings = new(State_Buildings)
	s.State.CargoShips = new(State_CargoShips)
	s.State.CargoShips.Ship_4 = new(CargoShip)
	s.State.CargoShips.Ship_5 = new(CargoShip)
	s.State.CargoShips.Ship_6 = new(CargoShip)
	s.State.CargoShips.Ship_7 = new(CargoShip)
	s.State.CargoShips.Ship_8 = new(CargoShip)

	return s
}

func (s *State) Copy() *State {
	c := new(State)
	*c = *s

	c.Players = new(State_Players)
	*c.Players = *(s.Players)
	c.Players.Player_1 = new(Player)
	*c.Players.Player_1 = *(s.Players.Player_1)
	c.Players.Player_1.Buildings = new(Player_Buildings)
	*c.Players.Player_1.Buildings = *(s.Players.Player_1.Buildings)
	c.Players.Player_1.Plantations = new(Player_Plantations)
	*c.Players.Player_1.Plantations = *(s.Players.Player_1.Plantations)
	c.Players.Player_1.Goods = new(Player_Goods)
	*c.Players.Player_1.Goods = *(s.Players.Player_1.Goods)
	c.Players.Player_2 = new(Player)
	*c.Players.Player_2 = *(s.Players.Player_2)
	c.Players.Player_2.Buildings = new(Player_Buildings)
	*c.Players.Player_2.Buildings = *(s.Players.Player_2.Buildings)
	c.Players.Player_2.Plantations = new(Player_Plantations)
	*c.Players.Player_2.Plantations = *(s.Players.Player_2.Plantations)
	c.Players.Player_2.Goods = new(Player_Goods)
	*c.Players.Player_2.Goods = *(s.Players.Player_2.Goods)
	c.Players.Player_3 = new(Player)
	*c.Players.Player_3 = *(s.Players.Player_3)
	c.Players.Player_3.Buildings = new(Player_Buildings)
	*c.Players.Player_3.Buildings = *(s.Players.Player_3.Buildings)
	c.Players.Player_3.Plantations = new(Player_Plantations)
	*c.Players.Player_3.Plantations = *(s.Players.Player_3.Plantations)
	c.Players.Player_3.Goods = new(Player_Goods)
	*c.Players.Player_3.Goods = *(s.Players.Player_3.Goods)
	c.Players.Player_4 = new(Player)
	*c.Players.Player_4 = *(s.Players.Player_4)
	c.Players.Player_4.Buildings = new(Player_Buildings)
	*c.Players.Player_4.Buildings = *(s.Players.Player_4.Buildings)
	c.Players.Player_4.Plantations = new(Player_Plantations)
	*c.Players.Player_4.Plantations = *(s.Players.Player_4.Plantations)
	c.Players.Player_4.Goods = new(Player_Goods)
	*c.Players.Player_4.Goods = *(s.Players.Player_4.Goods)
	c.Players.Player_5 = new(Player)
	*c.Players.Player_5 = *(s.Players.Player_5)
	c.Players.Player_5.Buildings = new(Player_Buildings)
	*c.Players.Player_5.Buildings = *(s.Players.Player_5.Buildings)
	c.Players.Player_5.Plantations = new(Player_Plantations)
	*c.Players.Player_5.Plantations = *(s.Players.Player_5.Plantations)
	c.Players.Player_5.Goods = new(Player_Goods)
	*c.Players.Player_5.Goods = *(s.Players.Player_5.Goods)
	c.Roles = new(State_Roles)
	*c.Roles = *(s.Roles)
	c.Roles.Prospector1 = new(Role)
	*c.Roles.Prospector1 = *(s.Roles.Prospector1)
	c.Roles.Prospector2 = new(Role)
	*c.Roles.Prospector2 = *(s.Roles.Prospector2)
	c.Roles.Builder = new(Role)
	*c.Roles.Builder = *(s.Roles.Builder)
	c.Roles.Captain = new(Role)
	*c.Roles.Captain = *(s.Roles.Captain)
	c.Roles.Craftsman = new(Role)
	*c.Roles.Craftsman = *(s.Roles.Craftsman)
	c.Roles.Mayor = new(Role)
	*c.Roles.Mayor = *(s.Roles.Mayor)
	c.Roles.Settler = new(Role)
	*c.Roles.Settler = *(s.Roles.Settler)
	c.Roles.Trader = new(Role)
	*c.Roles.Trader = *(s.Roles.Trader)
	c.Plantations = new(State_Plantations)
	*c.Plantations = *(s.Plantations)
	c.Plantations.Displayed = new(State_Plantations_Displayed)
	*c.Plantations.Displayed = *(s.Plantations.Displayed)
	c.Plantations.Facedown = new(State_Plantations_Counts)
	*c.Plantations.Facedown = *(s.Plantations.Facedown)
	c.Plantations.Discarded = new(State_Plantations_Counts)
	*c.Plantations.Discarded = *(s.Plantations.Discarded)
	c.Goods = new(State_Goods)
	*c.Goods = *(s.Goods)
	c.Buildings = new(State_Buildings)
	*c.Buildings = *(s.Buildings)
	c.CargoShips = new(State_CargoShips)
	*c.CargoShips = *(s.CargoShips)
	c.CargoShips.Ship_4 = new(CargoShip)
	*c.CargoShips.Ship_4 = *(s.CargoShips.Ship_4)
	c.CargoShips.Ship_5 = new(CargoShip)
	*c.CargoShips.Ship_5 = *(s.CargoShips.Ship_5)
	c.CargoShips.Ship_6 = new(CargoShip)
	*c.CargoShips.Ship_6 = *(s.CargoShips.Ship_6)
	c.CargoShips.Ship_7 = new(CargoShip)
	*c.CargoShips.Ship_7 = *(s.CargoShips.Ship_7)
	c.CargoShips.Ship_8 = new(CargoShip)
	*c.CargoShips.Ship_8 = *(s.CargoShips.Ship_8)

	return c
}

type gameState struct {
	*State
	*sync.Mutex
}

func newResponse() Response {
	var res Response

	res.State = new(State)
	res.State.Players = new(State_Players)
	res.State.Players.Player_1 = new(Player)
	res.State.Players.Player_1.Buildings = new(Player_Buildings)
	res.State.Players.Player_1.Plantations = new(Player_Plantations)
	res.State.Players.Player_1.Goods = new(Player_Goods)
	res.State.Players.Player_2 = new(Player)
	res.State.Players.Player_2.Buildings = new(Player_Buildings)
	res.State.Players.Player_2.Plantations = new(Player_Plantations)
	res.State.Players.Player_2.Goods = new(Player_Goods)
	res.State.Players.Player_3 = new(Player)
	res.State.Players.Player_3.Buildings = new(Player_Buildings)
	res.State.Players.Player_3.Plantations = new(Player_Plantations)
	res.State.Players.Player_3.Goods = new(Player_Goods)
	res.State.Players.Player_4 = new(Player)
	res.State.Players.Player_4.Buildings = new(Player_Buildings)
	res.State.Players.Player_4.Plantations = new(Player_Plantations)
	res.State.Players.Player_4.Goods = new(Player_Goods)
	res.State.Players.Player_5 = new(Player)
	res.State.Players.Player_5.Buildings = new(Player_Buildings)
	res.State.Players.Player_5.Plantations = new(Player_Plantations)
	res.State.Players.Player_5.Goods = new(Player_Goods)
	res.State.Roles = new(State_Roles)
	res.State.Roles.Prospector1 = new(Role)
	res.State.Roles.Prospector2 = new(Role)
	res.State.Roles.Builder = new(Role)
	res.State.Roles.Captain = new(Role)
	res.State.Roles.Craftsman = new(Role)
	res.State.Roles.Mayor = new(Role)
	res.State.Roles.Settler = new(Role)
	res.State.Roles.Trader = new(Role)
	res.State.Plantations = new(State_Plantations)
	res.State.Plantations.Displayed = new(State_Plantations_Displayed)
	res.State.Plantations.Facedown = new(State_Plantations_Counts)
	res.State.Plantations.Discarded = new(State_Plantations_Counts)
	res.State.Goods = new(State_Goods)
	res.State.Buildings = new(State_Buildings)
	res.State.CargoShips = new(State_CargoShips)
	res.State.CargoShips.Ship_4 = new(CargoShip)
	res.State.CargoShips.Ship_5 = new(CargoShip)
	res.State.CargoShips.Ship_6 = new(CargoShip)
	res.State.CargoShips.Ship_7 = new(CargoShip)
	res.State.CargoShips.Ship_8 = new(CargoShip)
	res.Error = new(game_engine_pb.Error)

	return res
}

//
// Service methods
//
func (e *gameEngine) SetPlayers(ctx context.Context, in *Count) (*Response, error) {
	state.Lock()
	defer state.Unlock()

	// Enforce the rules

	allowed, err := go_func.BoolAndBoolToBool_AND(go_func.BoolToBool_NOT(go_func.Bool(bool(state.Started))), go_func.IntAndIntToBool_LTE(go_func.Int(3), go_func.Int(int(in.Count)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprint(err),
			},
		}, nil
	}
	if !allowed {
		return &Response{
			Error: &game_engine_pb.Error{
				Code: "826de622-ad54-4b65-3395-bb4d3828e67b",
				Msg:  "dummy error",
			},
		}, nil
	}

	// Apply any effects

	next := state.State.Copy()
	next.Players.Player_1.Present, err = go_func.IntAndIntToBool_LTE(go_func.Int(3), go_func.Int(int(in.Count))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Players.Player_2.Present, err = go_func.IntAndIntToBool_LTE(go_func.Int(3), go_func.Int(int(in.Count))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Players.Player_3.Present, err = go_func.IntAndIntToBool_LTE(go_func.Int(3), go_func.Int(int(in.Count))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Players.Player_4.Present, err = go_func.IntAndIntToBool_LTE(go_func.Int(4), go_func.Int(int(in.Count))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Players.Player_5.Present, err = go_func.IntAndIntToBool_LTE(go_func.Int(5), go_func.Int(int(in.Count))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}

	state.State = next

	// Construct the response
	res := newResponse()
	res.State.Started = next.Started
	res.State.Players.Player_1.Present = next.Players.Player_1.Present
	res.State.Players.Player_2.Present = next.Players.Player_2.Present
	res.State.Players.Player_3.Present = next.Players.Player_3.Present
	res.State.Players.Player_4.Present = next.Players.Player_4.Present
	res.State.Players.Player_5.Present = next.Players.Player_5.Present

	return &res, nil
}

func (e *gameEngine) Start(ctx context.Context, in *EmptyMsg) (*Response, error) {
	state.Lock()
	defer state.Unlock()

	// Enforce the rules

	allowed, err := go_func.BoolAndBoolToBool_AND(go_func.BoolToBool_NOT(go_func.Bool(bool(state.Started))), go_func.BoolAndBoolToBool_AND(go_func.Bool(bool(state.Players.Player_1.Present)), go_func.BoolAndBoolToBool_AND(go_func.Bool(bool(state.Players.Player_2.Present)), go_func.BoolAndBoolToBool_AND(go_func.Bool(bool(state.Players.Player_3.Present)), go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Players.Player_4.Present)), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Players.Player_5.Present)))))))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprint(err),
			},
		}, nil
	}
	if !allowed {
		return &Response{
			Error: &game_engine_pb.Error{
				Code: "57e0c77b-3e2b-431c-c5ed-ee81a89d5ee6",
				Msg:  "game has already started",
			},
		}, nil
	}

	// Apply any effects

	next := state.State.Copy()
	next.Started, err = go_func.Bool(true).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}

	state.State = next

	// Construct the response
	res := newResponse()
	res.State.Started = next.Started

	return &res, nil
}

func (e *gameEngine) Accept(ctx context.Context, in *RoleChoice) (*Response, error) {
	state.Lock()
	defer state.Unlock()

	// Enforce the rules

	allowed, err := go_func.BoolsToBool_AND(go_func.Bool(bool(state.Started)), go_func.BoolsToBool_OR(go_func.BoolAndBoolToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(1), go_func.Int(int(in.Player))), go_func.IntAndIntToBool_EQ(go_func.Int(0), go_func.Int(int(state.Players.Player_1.Role)))), go_func.BoolAndBoolToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(2), go_func.Int(int(in.Player))), go_func.IntAndIntToBool_EQ(go_func.Int(0), go_func.Int(int(state.Players.Player_2.Role)))), go_func.BoolAndBoolToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(3), go_func.Int(int(in.Player))), go_func.IntAndIntToBool_EQ(go_func.Int(0), go_func.Int(int(state.Players.Player_3.Role)))), go_func.BoolAndBoolToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(4), go_func.Int(int(in.Player))), go_func.IntAndIntToBool_EQ(go_func.Int(0), go_func.Int(int(state.Players.Player_4.Role)))), go_func.BoolAndBoolToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(5), go_func.Int(int(in.Player))), go_func.IntAndIntToBool_EQ(go_func.Int(0), go_func.Int(int(state.Players.Player_5.Role))))), go_func.BoolsToBool_OR(go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(1), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Prospector1.Available)))), go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(2), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Prospector2.Available)))), go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(3), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Builder.Available)))), go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(4), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Captain.Available)))), go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(5), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Craftsman.Available)))), go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(6), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Mayor.Available)))), go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(7), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Settler.Available)))), go_func.BoolsToBool_AND(go_func.IntAndIntToBool_EQ(go_func.Int(8), go_func.Int(int(in.Role))), go_func.BoolToBool_NOT(go_func.Bool(bool(state.Roles.Trader.Available)))))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprint(err),
			},
		}, nil
	}
	if !allowed {
		return &Response{
			Error: &game_engine_pb.Error{
				Code: "446eda13-7e86-40c5-55c0-ac363015f92d",
				Msg:  "dummy error",
			},
		}, nil
	}

	// Apply any effects

	next := state.State.Copy()
	next.Roles.Prospector1.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Prospector1.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(1), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Roles.Prospector2.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Prospector2.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(2), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Roles.Builder.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Builder.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(3), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Roles.Captain.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Captain.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(4), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Roles.Craftsman.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Craftsman.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(5), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Roles.Mayor.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Mayor.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(6), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Roles.Settler.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Settler.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(7), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}
	next.Roles.Trader.Available, err = go_func.BoolAndBoolToBool_OR(go_func.Bool(bool(state.Roles.Trader.Available)), go_func.IntAndIntToBool_EQ(go_func.Int(8), go_func.Int(int(in.Role)))).Value()
	if err != nil {
		return &Response{
			Error: &game_engine_pb.Error{
				Msg: fmt.Sprintf("failed to apply effects: %v", err),
			},
		}, nil
	}

	state.State = next

	// Construct the response
	res := newResponse()

	return &res, nil
}

func (e *gameEngine) Purchase(ctx context.Context, in *BuildingChoice) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) Load(ctx context.Context, in *GoodToShip) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) Craft(ctx context.Context, in *EmptyMsg) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) CraftExtra(ctx context.Context, in *GoodChoice) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) WelcomeColonist(ctx context.Context, in *EmptyMsg) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) WelcomeColonistFromSupply(ctx context.Context, in *EmptyMsg) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) ApplyColonistToBuilding(ctx context.Context, in *BuildingChoice) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) ApplyColonistToPlantation(ctx context.Context, in *PlantationChoice) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) ApplyColonistToQuarry(ctx context.Context, in *EmptyMsg) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) RefillColonistShip(ctx context.Context, in *EmptyMsg) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) Settle(ctx context.Context, in *PlantationChoice) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) ConstructQuarry(ctx context.Context, in *EmptyMsg) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) Trade(ctx context.Context, in *GoodChoice) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

func (e *gameEngine) EndAction(ctx context.Context, in *PlayerChoice) (*Response, error) {

	return &Response{
		Error: &game_engine_pb.Error{
			Msg: "unimplemented",
		},
	}, nil
}

//
// Enum Key Mappings
//

func PuertoRicoPbState_PlayersFieldByPlayerID(msg *State_Players, key PlayerID) (*Player, error) {
	switch key {
	case PlayerID_PLAYER_1:
		return msg.Player_1, nil
	case PlayerID_PLAYER_2:
		return msg.Player_2, nil
	case PlayerID_PLAYER_3:
		return msg.Player_3, nil
	case PlayerID_PLAYER_4:
		return msg.Player_4, nil
	case PlayerID_PLAYER_5:
		return msg.Player_5, nil
	}

	return nil, fmt.Errorf("unexpected PlayerID value: %v", key)
}
