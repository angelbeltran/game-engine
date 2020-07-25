package main

import (
	"github.com/jhump/goprotoc/plugins"
)

const (
	isActionSetExtensionFieldNumber             = 50000
	actionExtensionFieldNumber                  = 50001
	isGameStateExtensionFieldNumber             = 50002
	isActionServiceResponseExtensionFieldNumber = 50003

	protoPackageName       = "game_engine"
	responseFieldName      = "response"
	responseStateFieldName = "state"
	responseErrorFieldName = "error"
	errorTypeName          = "Error"

	stateVariable = "state"
	inputVariable = "in"
)

var goNames plugins.GoNames
