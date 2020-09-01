// Generated by protoc-gen-game/generation. DO NOT EDIT.
package template

import (
	"fmt"
)

//
// ----- Error constructors -----
//

func failNoFunctionName(funcType string) (interface{}, error) {
	return nil, fmt.Errorf("function name missing for function type %s", funcType)
}

func failUndefinedEffect() (interface{}, error) {
	return nil, fmt.Errorf("undefined effect")
}